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
				m.treeCursor = len(m.swarm.SubAgents)
			}
		case tea.KeyDown, tea.KeyType('j'):
			m.treeCursor++
			if m.treeCursor > len(m.swarm.SubAgents) {
				m.treeCursor = 0
			}
		case tea.KeyEnter:
			// Edit Agent
			m.editAgentIdx = m.treeCursor - 1

			var agent AgentConfig
			if m.editAgentIdx == -1 {
				agent = m.swarm.Orchestrator
			} else {
				agent = m.swarm.SubAgents[m.editAgentIdx]
			}

			inputs := make([]textinput.Model, 4)
			inputs[0] = createTextInput("Name")
			inputs[0].SetValue(agent.Name)
			inputs[1] = createTextInput("Focus")
			inputs[1].SetValue(agent.Focus)
			inputs[2] = createTextInput("Model")
			inputs[2].SetValue(agent.Model)
			inputs[3] = createTextInput("Temperature")
			inputs[3].SetValue(agent.Temperature)

			inputs[0].Focus()
			m.editFocus = 0
			m.editInputs = inputs
			m.state = StateEditAgent
		case tea.KeyType('r'): // Run command
			m.state = StateRunning

			m.spinners = make([]spinner.Model, len(m.swarm.SubAgents))
			m.completed = make([]bool, len(m.swarm.SubAgents))
			m.results = make([]AgentResult, len(m.swarm.SubAgents))

			var cmds []tea.Cmd
			for i := range m.spinners {
				m.spinners[i] = createSpinner()
				cmds = append(cmds, m.spinners[i].Tick)
				cmds = append(cmds, runAgent(m.swarm.SubAgents[i], m.swarm.MainTask, i))
			}

			// If no subagents, go straight to orchestrator
			if len(m.swarm.SubAgents) == 0 {
				m.synthesizing = true
				m.synthSpinner = createSpinner()
				return m, tea.Batch(m.synthSpinner.Tick, runSynth(m.swarm.Orchestrator, m.swarm.MainTask, m.results))
			}

			return m, tea.Batch(cmds...)

		case tea.KeyBackspace, tea.KeyType('b'):
			m.state = StateMainMenu
		}
	}
	return m, nil
}

func (m MainModel) viewTree() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Evaluation Tree (Enter to Edit • 'r' to Run)"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Task: %s\n\n", m.swarm.MainTask))

	// Orchestrator Node
	orchNode := fmt.Sprintf("▼ [%s] %s (Model: %s, Temp: %s)", m.swarm.Orchestrator.ID, m.swarm.Orchestrator.Name, m.swarm.Orchestrator.Model, m.swarm.Orchestrator.Temperature)
	if m.treeCursor == 0 {
		b.WriteString(nodeActiveStyle.Render(orchNode))
	} else {
		b.WriteString(nodeNormalStyle.Render(orchNode))
	}
	b.WriteString("\n")

	// Sub-Agents
	for i, spec := range m.swarm.SubAgents {
		prefix := "├──"
		if i == len(m.swarm.SubAgents)-1 {
			prefix = "└──"
		}

		b.WriteString(treeLineStyle.Render(" " + prefix + " "))

		nodeText := fmt.Sprintf("[%s] %s: %s (Model: %s, Temp: %s)", spec.ID, spec.Name, spec.Focus, spec.Model, spec.Temperature)
		if m.treeCursor == i+1 {
			b.WriteString(nodeActiveStyle.Render(nodeText))
		} else {
			b.WriteString(nodeNormalStyle.Render(nodeText))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("up/down: navigate • enter: edit • r: run • b: main menu • esc: quit"))

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
			if m.editAgentIdx == -1 {
				m.swarm.Orchestrator.Name = m.editInputs[0].Value()
				m.swarm.Orchestrator.Focus = m.editInputs[1].Value()
				m.swarm.Orchestrator.Model = m.editInputs[2].Value()
				m.swarm.Orchestrator.Temperature = m.editInputs[3].Value()
			} else {
				m.swarm.SubAgents[m.editAgentIdx].Name = m.editInputs[0].Value()
				m.swarm.SubAgents[m.editAgentIdx].Focus = m.editInputs[1].Value()
				m.swarm.SubAgents[m.editAgentIdx].Model = m.editInputs[2].Value()
				m.swarm.SubAgents[m.editAgentIdx].Temperature = m.editInputs[3].Value()
			}

			m.state = StateTree
			return m, nil
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

	agentId := m.swarm.Orchestrator.ID
	if m.editAgentIdx != -1 {
		agentId = m.swarm.SubAgents[m.editAgentIdx].ID
	}

	b.WriteString(titleStyle.Render(fmt.Sprintf("Editing %s", agentId)))
	b.WriteString("\n\n")

	labels := []string{"Name", "Focus", "Model", "Temperature"}

	for i, input := range m.editInputs {
		if m.editFocus == i {
			b.WriteString(activeInputStyle.Render("► " + labels[i] + ":\n"))
		} else {
			b.WriteString(normalInputStyle.Render("  " + labels[i] + ":\n"))
		}
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("tab/up/down: next • enter: save & back • esc: quit"))

	return boxStyle.Render(b.String())
}
