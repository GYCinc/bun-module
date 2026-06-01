package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	state ViewState

	// Setup View
	taskInput   textarea.Model
	purpose     string
	purposeIdx  int
	purposes    []string
	maxAgents   int
	setupInputs []textinput.Model
	setupFocus  int // 0 = task, 1 = purpose, 2 = maxAgents

	// Tree View
	specs      []SubAgentSpec
	treeCursor int // 0 = orchestrator, 1..N = subagents

	// Edit View
	editInputs   []textinput.Model
	editFocus    int
	editAgentId  string
	editAgentIdx int

	// Running View
	spinners     []spinner.Model
	completed    []bool
	results      []AgentResult
	synthesizing bool
	synthSpinner spinner.Model
	finalResult  string
}

func InitialModel() MainModel {
	ta := createTextArea("Enter the main task...")

	tiPurpose := createTextInput("")
	tiPurpose.SetValue("researcher") // Default
	tiPurpose.Blur()

	tiMax := createTextInput("")
	tiMax.SetValue("5") // Default
	tiMax.Blur()

	return MainModel{
		state:       StateSetup,
		taskInput:   ta,
		purposeIdx:  0,
		purposes:    []string{"researcher", "coder", "analyst", "custom"},
		purpose:     "researcher",
		maxAgents:   5,
		setupInputs: []textinput.Model{tiPurpose, tiMax},
		setupFocus:  0,
	}
}

func (m MainModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateSetup:
		return m.updateSetup(msg)
	case StateTree:
		return m.updateTree(msg)
	case StateEditAgent:
		return m.updateEdit(msg)
	case StateRunning:
		return m.updateRunning(msg)
	case StateResults:
		return m.updateResults(msg)
	}
	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case StateSetup:
		return m.viewSetup()
	case StateTree:
		return m.viewTree()
	case StateEditAgent:
		return m.viewEdit()
	case StateRunning:
		return m.viewRunning()
	case StateResults:
		return m.viewResults()
	}
	return ""
}

// --- Setup View ---

func (m *MainModel) updateSetup(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.setupFocus--
				if m.setupFocus < 0 {
					m.setupFocus = 2
				}
			} else {
				m.setupFocus++
				if m.setupFocus > 2 {
					m.setupFocus = 0
				}
			}

			if m.setupFocus == 0 {
				m.taskInput.Focus()
				m.setupInputs[0].Blur()
				m.setupInputs[1].Blur()
			} else if m.setupFocus == 1 {
				m.taskInput.Blur()
				m.setupInputs[0].Focus()
				m.setupInputs[1].Blur()
			} else {
				m.taskInput.Blur()
				m.setupInputs[0].Blur()
				m.setupInputs[1].Focus()
			}

			return m, nil

		case tea.KeyEnter:
			if m.setupFocus == 1 {
				m.purposeIdx++
				if m.purposeIdx >= len(m.purposes) {
					m.purposeIdx = 0
				}
				m.purpose = m.purposes[m.purposeIdx]
				m.setupInputs[0].SetValue(m.purpose)
				return m, nil
			}

			if m.setupFocus == 2 {
				// Parse and submit
				task := m.taskInput.Value()
				if task == "" {
					task = "Default Task"
				}

				maxA, err := strconv.Atoi(m.setupInputs[1].Value())
				if err != nil || maxA < 1 || maxA > 5 {
					maxA = 5 // default fallback
				}

				blueprint := NewRoleBlueprint(m.purpose, maxA)
				m.specs = blueprint.BuildSpecs(task)

				m.state = StateTree
				m.treeCursor = 0
				return m, nil
			}
		}
	}

	if m.setupFocus == 0 {
		m.taskInput, cmd = m.taskInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.setupFocus == 1 {
		// Read-only toggle, but handled above
	} else {
		m.setupInputs[1], cmd = m.setupInputs[1].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) viewSetup() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("1-Orchestrator + 5-Subagent Runtime"))
	b.WriteString("\n\n")

	// Task Input
	if m.setupFocus == 0 {
		b.WriteString(activeInputStyle.Render("► Main Task:"))
	} else {
		b.WriteString(normalInputStyle.Render("  Main Task:"))
	}
	b.WriteString("\n")
	b.WriteString(m.taskInput.View())
	b.WriteString("\n\n")

	// Purpose
	if m.setupFocus == 1 {
		b.WriteString(activeInputStyle.Render(fmt.Sprintf("► Purpose (Press Enter to toggle): %s", m.purpose)))
	} else {
		b.WriteString(normalInputStyle.Render(fmt.Sprintf("  Purpose: %s", m.purpose)))
	}
	b.WriteString("\n\n")

	// Max Agents
	if m.setupFocus == 2 {
		b.WriteString(activeInputStyle.Render("► Max Agents (1-5): "))
	} else {
		b.WriteString(normalInputStyle.Render("  Max Agents (1-5): "))
	}
	b.WriteString(m.setupInputs[1].View())
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("tab/shift+tab: move • enter: submit • ctrl+c: quit"))

	return boxStyle.Render(b.String())
}
