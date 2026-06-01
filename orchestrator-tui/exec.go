package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type AgentCompleteMsg struct {
	Index  int
	Result AgentResult
}

type SynthCompleteMsg struct {
	Result string
}

func runAgent(spec SubAgentSpec, task string, index int) tea.Cmd {
	return func() tea.Msg {
		res := MockExecuteSubagent(spec, task, 2*time.Second)
		return AgentCompleteMsg{Index: index, Result: res}
	}
}

func runSynth(task string, results []AgentResult) tea.Cmd {
	return func() tea.Msg {
		res := MockSynthesize(task, results, 3*time.Second)
		return SynthCompleteMsg{Result: res}
	}
}

func (m *MainModel) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmds []tea.Cmd
		if m.synthesizing {
			var cmd tea.Cmd
			m.synthSpinner, cmd = m.synthSpinner.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			for i := range m.spinners {
				if !m.completed[i] {
					var cmd tea.Cmd
					m.spinners[i], cmd = m.spinners[i].Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}
		return m, tea.Batch(cmds...)

	case AgentCompleteMsg:
		m.completed[msg.Index] = true
		m.results[msg.Index] = msg.Result

		// Check if all done
		allDone := true
		for _, done := range m.completed {
			if !done {
				allDone = false
				break
			}
		}

		if allDone {
			m.synthesizing = true
			m.synthSpinner = createSpinner()
			return m, tea.Batch(m.synthSpinner.Tick, runSynth(m.taskInput.Value(), m.results))
		}

		return m, nil

	case SynthCompleteMsg:
		m.finalResult = msg.Result
		m.state = StateResults
		return m, nil
	}

	return m, nil
}

func (m MainModel) viewRunning() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Execution Running"))
	b.WriteString("\n\n")

	if !m.synthesizing {
		b.WriteString(subTitleStyle.Render("Sub-Agents Working in Parallel:"))
		b.WriteString("\n\n")

		for i, spec := range m.specs {
			if m.completed[i] {
				b.WriteString(fmt.Sprintf("%s [%s] %s: Done ✓\n", treeLineStyle.Render("├──"), spec.ID, spec.Name))
			} else {
				b.WriteString(fmt.Sprintf("%s [%s] %s: %s\n", treeLineStyle.Render("├──"), spec.ID, spec.Name, m.spinners[i].View()))
			}
		}
	} else {
		b.WriteString(subTitleStyle.Render("All Agents Completed."))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s Orchestrator Synthesizing %s\n", treeLineStyle.Render("└──"), m.synthSpinner.View()))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("esc: quit"))

	return boxStyle.Render(b.String())
}

func (m *MainModel) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyBackspace, tea.KeyType('b'):
			m.state = StateSetup
		}
	}
	return m, nil
}

func (m MainModel) viewResults() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Final Synthesized Results"))
	b.WriteString("\n\n")
	b.WriteString(m.finalResult)
	b.WriteString("\n\n")

	b.WriteString(subTitleStyle.Render("Individual Outputs:"))
	b.WriteString("\n")
	for _, res := range m.results {
		b.WriteString(fmt.Sprintf("• [%s] %s\n", res.Name, res.Output))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("b: back to start • esc: quit"))

	return boxStyle.Render(b.String())
}
