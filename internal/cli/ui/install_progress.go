// ui/install_progress.go
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	// InstallEvent is the interface for install progress events.
	InstallEvent interface {
		isInstallEvent()
	}

	// InstallStepEvent is sent when a new step is started or completed.
	InstallStepEvent struct {
		Step    int
		Message string
	}

	// InstallErrorEvent is sent when an error occurs.
	InstallErrorEvent struct {
		Message string
	}

	// installProgressModel holds state for our install progress UI.
	installProgressModel struct {
		spinner        spinner.Model
		eventChan      <-chan InstallEvent
		currentMessage string
		completedSteps []string
		done           bool
		errorMessage   string
		cancelled      bool
		cancelFunc     func()
	}
)

// NewInstallProgressModel creates a new model with a spinner.
func NewInstallProgressModel(eventChan <-chan InstallEvent, cancel func()) installProgressModel {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return installProgressModel{
		spinner:        s,
		eventChan:      eventChan,
		completedSteps: []string{},
		cancelFunc:     cancel,
	}
}

// Init starts the spinner and waits for the first event.
func (m *installProgressModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForInstallEvent(m.eventChan),
	)
}

// Update handles incoming messages/events.
func (m *installProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.cancel()
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case InstallStepEvent:
		m.currentMessage = msg.Message
		m.completedSteps = append(m.completedSteps, msg.Message)
		if msg.Step == 8 {
			m.done = true
			return m, tea.Quit
		}

		return m, waitForInstallEvent(m.eventChan)

	case InstallErrorEvent:
		m.errorMessage = msg.Message
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI.
func (m *installProgressModel) View() string {
	if m.errorMessage != "" {
		return errorStyle.Render(fmt.Sprintf("Error: %s!\n", m.errorMessage))
	}

	if m.done {
		return successStyle.Render(m.currentMessage, "\n")
	}

	if m.cancelled {
		return warningStyle.Render("Cancelled!\n")
	}

	stepsView := ""
	for _, step := range m.completedSteps {
		stepsView += stepStyle.Render("✓ "+step) + "\n"
	}

	view := fmt.Sprintf("%s %s\n\n%s\n\n%s",
		m.spinner.View(),
		infoStyle.Render(m.currentMessage),
		stepsView,
		renderLegend(),
	)

	return view
}

func (m *installProgressModel) cancel() {
	m.cancelled = true
	m.cancelFunc()
}

func waitForInstallEvent(eventChan <-chan InstallEvent) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan
	}
}

func (InstallStepEvent) isInstallEvent() {}

func (InstallErrorEvent) isInstallEvent() {}
