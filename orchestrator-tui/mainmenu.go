package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp, tea.KeyType('k'):
			m.menuCursor--
			if m.menuCursor < 0 {
				m.menuCursor = 1
			}
		case tea.KeyDown, tea.KeyType('j'):
			m.menuCursor++
			if m.menuCursor > 1 {
				m.menuCursor = 0
			}
		case tea.KeyEnter:
			if m.menuCursor == 0 {
				// Create New Swarm (Rabbit Hole)
				m.state = StateWizard
				m.wizardStep = 0
				m.swarm = Swarm{}
				m.wizardInput = createTextInput("Enter Main Task...")
				m.wizardInput.Focus()
				return m, nil
			} else {
				// Load pre-existing (For now, just load a mock one)
				m.swarm = Swarm{
					MainTask:     "Build a generic web app",
					Orchestrator: AgentConfig{ID: "orch_1", Name: "Tech Lead", Focus: "Oversee dev", Model: "gpt-4o", Temperature: "0.7"},
					SubAgents: []AgentConfig{
						{ID: "sub_1", Name: "Frontend", Focus: "React", Model: "claude-3-opus", Temperature: "0.5"},
						{ID: "sub_2", Name: "Backend", Focus: "Go API", Model: "gpt-4o", Temperature: "0.7"},
					},
				}
				m.state = StateTree
				m.treeCursor = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m MainModel) viewMainMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Welcome to the Swarm Builder"))
	b.WriteString("\n\n")

	choices := []string{"Create a New Swarm (The Rabbit Hole)", "Load Pre-existing Swarm"}

	for i, choice := range choices {
		cursor := "  "
		if m.menuCursor == i {
			cursor = activeInputStyle.Render("► ")
			b.WriteString(cursor + activeInputStyle.Render(choice) + "\n")
		} else {
			b.WriteString(cursor + normalInputStyle.Render(choice) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("up/down: move • enter: select • esc: quit"))

	return boxStyle.Render(b.String())
}
