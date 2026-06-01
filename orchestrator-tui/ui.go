package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	state ViewState

	// Main Menu
	menuCursor int

	// Swarm Data
	swarm Swarm

	// Wizard
	wizardStep  int
	wizardInput textinput.Model

	// Tree View
	treeCursor int

	// Edit View
	editInputs   []textinput.Model
	editFocus    int
	editAgentIdx int // -1 for orchestrator, 0+ for subagents

	// Running View
	spinners     []spinner.Model
	completed    []bool
	results      []AgentResult
	synthesizing bool
	synthSpinner spinner.Model
	finalResult  string
}

func InitialModel() MainModel {
	return MainModel{
		state:      StateMainMenu,
		menuCursor: 0,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateMainMenu:
		return m.updateMainMenu(msg)
	case StateWizard:
		return m.updateWizard(msg)
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
	case StateMainMenu:
		return m.viewMainMenu()
	case StateWizard:
		return m.viewWizard()
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
