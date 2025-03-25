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
		status     status
		spinner    spinner.Model
		eventChan  <-chan InstallEvent
		message    string
		steps      []string
		cancelFunc func()
	}
)

// NewInstallProgressModel creates a new model with a spinner.
func NewInstallProgressModel(eventChan <-chan InstallEvent, cancel func()) installProgressModel {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return installProgressModel{
		status:     statusInProgress,
		spinner:    s,
		eventChan:  eventChan,
		steps:      []string{},
		cancelFunc: cancel,
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
			m.status = statusCancelled
			m.cancelFunc()
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case InstallStepEvent:
		m.message = msg.Message
		m.steps = append(m.steps, msg.Message)
		if msg.Step == 8 {
			m.status = statusDone
			return m, tea.Quit
		}
		return m, waitForInstallEvent(m.eventChan)

	case InstallErrorEvent:
		m.message = msg.Message
		m.status = statusError
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI.
func (m *installProgressModel) View() string {
	switch m.status {
	case statusError:
		return errorStyle.Render(fmt.Sprintf("Error: %s\n", m.message))
	case statusDone:
		return successStyle.Render(fmt.Sprintf("%s\n", m.message))
	case statusCancelled:
		return warningStyle.Render("Cancelled!\n")
	}

	stepsView := ""
	if len(m.steps) > 1 {
		for _, step := range m.steps[:len(m.steps)-1] {
			stepsView += stepStyle.Render("✓ "+step) + "\n"
		}
	}

	view := fmt.Sprintf("%s%s %s\n\n%s",
		stepsView,
		m.spinner.View(),
		infoStyle.Render(m.message),
		renderLegend(),
	)

	return view
}

func waitForInstallEvent(eventChan <-chan InstallEvent) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan
	}
}

func (InstallStepEvent) isInstallEvent() {}

func (InstallErrorEvent) isInstallEvent() {}
