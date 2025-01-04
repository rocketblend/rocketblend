package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	renderProgressModel struct {
		progress     progress.Model
		spinner      spinner.Model
		totalFrames  int
		currentFrame int
		eventChan    <-chan types.BlenderEvent // Read-only channel for events
		done         bool
	}
)

func NewRenderProgressModel(totalFrames int, eventChan <-chan types.BlenderEvent) renderProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return renderProgressModel{
		progress:    p,
		spinner:     s,
		totalFrames: totalFrames,
		eventChan:   eventChan,
	}
}

func (m renderProgressModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForBlenderEvent(m.eventChan),
	)
}

func (m *renderProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case types.BlenderEvent:
		if renderEvent, ok := msg.(*types.RenderEvent); ok {
			return m.handleRenderEvent(renderEvent)
		}

		return m, waitForBlenderEvent(m.eventChan)
	}

	return m, nil
}

func (m *renderProgressModel) View() string {
	if m.done {
		return fmt.Sprintf("Rendering complete! Rendered %d frames.\n", m.totalFrames)
	}

	spin := m.spinner.View()
	prog := m.progress.View()
	return fmt.Sprintf("%s Rendering frame %d/%d\n%s", spin, m.currentFrame, m.totalFrames, prog)
}

func (m *renderProgressModel) handleRenderEvent(e *types.RenderEvent) (tea.Model, tea.Cmd) {
	if e.Frame != m.currentFrame {
		m.currentFrame = e.Frame
	}

	currentSample, _ := strconv.Atoi(e.Data["current"])
	totalSamples, _ := strconv.Atoi(e.Data["total"])
	frameProgress := 0.0
	if totalSamples > 0 {
		frameProgress = float64(currentSample) / float64(totalSamples)
	}

	totalProgress := (float64(m.currentFrame-1) + frameProgress) / float64(m.totalFrames)
	progressCmd := m.progress.SetPercent(totalProgress)

	if m.currentFrame >= m.totalFrames && frameProgress >= 1.0 {
		m.done = true
		return m, tea.Quit
	}

	return m, tea.Batch(progressCmd, waitForBlenderEvent(m.eventChan))
}

func waitForBlenderEvent(eventChan <-chan types.BlenderEvent) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan // Receive an event from the channel
	}
}
