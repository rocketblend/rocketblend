package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	// ProgressEvent is the interface for progress events.
	ProgressEvent interface {
		isProgressEvent()
	}

	// StepEvent is sent when a new step is started or completed.
	StepEvent struct {
		Step    int
		Message string
	}

	// ErrorEvent is sent when an error occurs.
	ErrorEvent struct {
		Message string
	}

	// CompletionEvent is sent when progress completes successfully.
	CompletionEvent struct {
		Message string
	}

	// progressModel holds state for our generic progress UI.
	progressModel struct {
		status     status
		spinner    spinner.Model
		eventChan  <-chan ProgressEvent
		message    string
		steps      []string
		cancelFunc func()
	}
)

// Run starts the progress UI.
func Run(ctx context.Context, work func(ctx context.Context, eventChan chan<- ProgressEvent) error) error {
	eventChan := make(chan ProgressEvent, 10)
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(eventChan)
		if err := work(ctxWithCancel, eventChan); err != nil {
			eventChan <- ErrorEvent{Message: err.Error()}
		}
	}()

	m := newProgressModel(eventChan, cancel)
	program := tea.NewProgram(&m, tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	return nil
}

// newProgressModel creates a new model with a spinner.
func newProgressModel(eventChan <-chan ProgressEvent, cancel func()) progressModel {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return progressModel{
		status:     statusInProgress,
		spinner:    s,
		eventChan:  eventChan,
		steps:      []string{},
		cancelFunc: cancel,
	}
}

// Init starts the spinner and waits for the first event.
func (m *progressModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForProgressEvent(m.eventChan),
	)
}

// Update handles incoming messages/events.
func (m *progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case StepEvent:
		m.message = msg.Message
		m.steps = append(m.steps, msg.Message)
		return m, waitForProgressEvent(m.eventChan)

	case CompletionEvent:
		m.message = msg.Message
		m.status = statusDone
		return m, tea.Quit

	case ErrorEvent:
		m.message = msg.Message
		m.status = statusError
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI.
func (m *progressModel) View() string {
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

func waitForProgressEvent(eventChan <-chan ProgressEvent) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan
	}
}

func (StepEvent) isProgressEvent() {}

func (ErrorEvent) isProgressEvent() {}

func (CompletionEvent) isProgressEvent() {}
