package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) updateTree(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp, tea.KeyType('k'):
			m.treeCursor--
			if m.treeCursor < 0 {
				m.treeCursor = len(m.specs)
			}
		case tea.KeyDown, tea.KeyType('j'):
			m.treeCursor++
			if m.treeCursor > len(m.specs) {
				m.treeCursor = 0
			}
		case tea.KeyEnter:
			if m.treeCursor == 0 {
				// Execute
				m.state = StateRunning

				// Initialize Spinners
				m.spinners = make([]spinner.Model, len(m.specs))
				m.completed = make([]bool, len(m.specs))
				m.results = make([]AgentResult, len(m.specs))

				var cmds []tea.Cmd
				for i := range m.spinners {
					m.spinners[i] = createSpinner()
					cmds = append(cmds, m.spinners[i].Tick)
					cmds = append(cmds, runAgent(m.specs[i], m.taskInput.Value(), i))
				}
				return m, tea.Batch(cmds...)
			} else {
				// Edit Agent
				m.editAgentIdx = m.treeCursor - 1
				agent := m.specs[m.editAgentIdx]

				inputs := make([]textinput.Model, 4)
				inputs[0] = createTextInput("Name")
				inputs[0].SetValue(agent.Name)
				inputs[1] = createTextInput("Focus")
				inputs[1].SetValue(agent.Focus)
				inputs[2] = createTextInput("Instructions")
				inputs[2].SetValue(agent.Instructions)
				inputs[3] = createTextInput("Expectations")
				inputs[3].SetValue(agent.OutputExpectations)

				inputs[0].Focus()
				m.editFocus = 0
				m.editInputs = inputs
				m.state = StateEditAgent
			}
		case tea.KeyBackspace, tea.KeyType('b'):
			m.state = StateSetup
		}
	}
	return m, nil
}

func (m MainModel) viewTree() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Agent Hierarchy (Press Enter to Execute/Edit)"))
	b.WriteString("\n\n")

	// Orchestrator Node
	if m.treeCursor == 0 {
		b.WriteString(nodeActiveStyle.Render("▼ [Orchestrator] (Run All)"))
	} else {
		b.WriteString(nodeNormalStyle.Render("▼ [Orchestrator]"))
	}
	b.WriteString("\n")

	// Sub-Agents
	for i, spec := range m.specs {
		prefix := "├──"
		if i == len(m.specs)-1 {
			prefix = "└──"
		}

		b.WriteString(treeLineStyle.Render(" " + prefix + " "))

		nodeText := fmt.Sprintf("[%s] %s: %s", spec.ID, spec.Name, spec.Focus)
		if m.treeCursor == i+1 {
			b.WriteString(nodeActiveStyle.Render(nodeText))
		} else {
			b.WriteString(nodeNormalStyle.Render(nodeText))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("up/down: navigate • enter: edit/run • b: back • esc: quit"))

	return boxStyle.Render(b.String())
}

// --- Edit View ---

func (m *MainModel) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown:
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.editFocus--
				if m.editFocus < 0 {
					m.editFocus = len(m.editInputs) - 1
				}
			} else {
				m.editFocus++
				if m.editFocus > len(m.editInputs)-1 {
					m.editFocus = 0
				}
			}

			for i := range m.editInputs {
				if i == m.editFocus {
					m.editInputs[i].Focus()
				} else {
					m.editInputs[i].Blur()
				}
			}
			return m, nil

		case tea.KeyEnter:
			// Save back to spec
			m.specs[m.editAgentIdx].Name = m.editInputs[0].Value()
			m.specs[m.editAgentIdx].Focus = m.editInputs[1].Value()
			m.specs[m.editAgentIdx].Instructions = m.editInputs[2].Value()
			m.specs[m.editAgentIdx].OutputExpectations = m.editInputs[3].Value()

			m.state = StateTree
			return m, nil

		case tea.KeyBackspace:
			// Only go back on backspace if input is empty or if it's explicitly a command,
			// but textinput handles backspace. So we'll use a specific binding or just Esc/Save.
			// Actually, let textinput handle backspace.
		}
	}

	for i := range m.editInputs {
		m.editInputs[i], cmd = m.editInputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) viewEdit() string {
	var b strings.Builder

	agent := m.specs[m.editAgentIdx]
	b.WriteString(titleStyle.Render(fmt.Sprintf("Editing %s", agent.ID)))
	b.WriteString("\n\n")

	labels := []string{"Name", "Focus", "Instructions", "Expectations"}

	for i, input := range m.editInputs {
		if m.editFocus == i {
			b.WriteString(activeInputStyle.Render("► " + labels[i] + ":\n"))
		} else {
			b.WriteString(normalInputStyle.Render("  " + labels[i] + ":\n"))
		}
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("tab: next • enter: save & back • esc: quit"))

	return boxStyle.Render(b.String())
}
