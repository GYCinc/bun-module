package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *MainModel) updateWizard(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.wizardInput.Value()

			switch m.wizardStep {
			case 0: // Main Task
				m.swarm.MainTask = val
				if m.swarm.MainTask == "" {
					m.swarm.MainTask = "Default task"
				}
				m.wizardStep++
				m.wizardInput.SetValue("")
				m.wizardInput.Placeholder = "e.g., Tech Lead, Architect"
			case 1: // Orchestrator Name
				m.swarm.Orchestrator.Name = val
				m.swarm.Orchestrator.ID = "orch_1"
				m.wizardStep++
				m.wizardInput.SetValue("gpt-4o")
				m.wizardInput.Placeholder = "Model (e.g. gpt-4o)"
			case 2: // Orchestrator Model
				m.swarm.Orchestrator.Model = val
				m.wizardStep++
				m.wizardInput.SetValue("0.7")
				m.wizardInput.Placeholder = "Temperature (0.0 - 1.0)"
			case 3: // Orchestrator Temp
				m.swarm.Orchestrator.Temperature = val
				m.wizardStep++ // Move to asking for Sub-Agents
				m.wizardInput.SetValue("")
				m.wizardInput.Placeholder = "(y/n)"
			case 4: // Add Sub-agent?
				if strings.ToLower(val) == "y" || strings.ToLower(val) == "yes" {
					m.wizardStep = 5 // Start sub-agent flow
					m.wizardInput.SetValue("")
					m.wizardInput.Placeholder = "e.g., Frontend Dev, Researcher"
				} else {
					// Done with wizard!
					m.state = StateTree
					m.treeCursor = 0
					return m, nil
				}
			case 5: // Sub-agent Name
				newAgent := AgentConfig{ID: fmt.Sprintf("sub_%d", len(m.swarm.SubAgents)+1), Name: val}
				m.swarm.SubAgents = append(m.swarm.SubAgents, newAgent)
				m.wizardStep++
				m.wizardInput.SetValue("")
				m.wizardInput.Placeholder = "e.g., Write React components"
			case 6: // Sub-agent Focus
				m.swarm.SubAgents[len(m.swarm.SubAgents)-1].Focus = val
				m.wizardStep++
				m.wizardInput.SetValue("claude-3-opus")
				m.wizardInput.Placeholder = "Model (e.g. claude-3-opus)"
			case 7: // Sub-agent Model
				m.swarm.SubAgents[len(m.swarm.SubAgents)-1].Model = val
				m.wizardStep++
				m.wizardInput.SetValue("0.5")
				m.wizardInput.Placeholder = "Temperature (0.0 - 1.0)"
			case 8: // Sub-agent Temp
				m.swarm.SubAgents[len(m.swarm.SubAgents)-1].Temperature = val
				m.wizardStep = 4 // Ask to add another
				m.wizardInput.SetValue("")
				m.wizardInput.Placeholder = "(y/n)"
			}
			return m, nil
		}
	}

	m.wizardInput, cmd = m.wizardInput.Update(msg)
	return m, cmd
}

func (m MainModel) viewWizard() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("The Rabbit Hole - Swarm Builder"))
	b.WriteString("\n\n")

	// Render collapsed history
	if m.wizardStep > 0 {
		b.WriteString(wizardDoneStyle.Render(fmt.Sprintf("✓ Main Task: %s", m.swarm.MainTask)))
		b.WriteString("\n")
	}
	if m.wizardStep > 3 {
		b.WriteString(wizardDoneStyle.Render(fmt.Sprintf("✓ Orchestrator: %s (Model: %s, Temp: %s)",
			m.swarm.Orchestrator.Name, m.swarm.Orchestrator.Model, m.swarm.Orchestrator.Temperature)))
		b.WriteString("\n")
	}
	if m.wizardStep > 4 && len(m.swarm.SubAgents) > 0 {
		b.WriteString(wizardDoneStyle.Render(fmt.Sprintf("✓ Sub-agents configured: %d", len(m.swarm.SubAgents))))
		b.WriteString("\n")
		// Show most recently added sub-agent if we are back at step 4
		if m.wizardStep == 4 {
			last := m.swarm.SubAgents[len(m.swarm.SubAgents)-1]
			b.WriteString(wizardDoneStyle.Render(fmt.Sprintf("  └─ Last added: %s (%s)", last.Name, last.Focus)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Current prompt
	prompt := ""
	switch m.wizardStep {
	case 0:
		prompt = "What is the primary objective of this swarm?"
	case 1:
		prompt = "Who is leading this? Name the Orchestrator."
	case 2:
		prompt = "Which model should the Orchestrator use?"
	case 3:
		prompt = "What temperature should the Orchestrator use?"
	case 4:
		prompt = "Would you like to add a specialized sub-agent to the team? (y/n)"
	case 5:
		prompt = "What is the name/role of this sub-agent?"
	case 6:
		prompt = "What is their specific focus or skill?"
	case 7:
		prompt = "Which model should this sub-agent use?"
	case 8:
		prompt = "What temperature should this sub-agent use?"
	}

	b.WriteString(activeInputStyle.Render(prompt))
	b.WriteString("\n")
	b.WriteString(m.wizardInput.View())
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("enter: submit • esc: quit"))

	return boxStyle.Render(b.String())
}
