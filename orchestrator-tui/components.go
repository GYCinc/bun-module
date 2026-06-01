package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// UI States
type ViewState int

const (
	StateMainMenu ViewState = iota
	StateWizard
	StateTree
	StateEditAgent
	StateRunning
	StateResults
)

// Styling
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	subTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	activeInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EE6FF8"))

	normalInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	treeLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#43BF6D"))

	nodeActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#EE6FF8")).
			Padding(0, 1)

	nodeNormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	wizardDoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#43BF6D")).
			MarginBottom(1)
)

func createTextArea(placeholder string) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.Focus()
	ta.CharLimit = 500
	ta.SetWidth(60)
	ta.SetHeight(5)
	return ta
}

func createTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 100
	ti.Width = 60
	return ti
}

func createSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}
