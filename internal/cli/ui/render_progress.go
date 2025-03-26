package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type renderProgressModel struct {
	progress        progress.Model
	spinner         spinner.Model
	eventChan       <-chan types.BlenderEvent
	currentFrame    int
	currentMemory   string
	totalFrames     int
	currentSample   int
	totalSamples    int
	startTime       time.Time
	completedFrames []int
	done            bool
	earlyExit       bool
	lastUpdate      time.Time
	errorMessage    string
	cancel          context.CancelFunc
}

var (
	progressStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	currentFrameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
)

func NewRenderProgressModel(totalFrames int, eventChan <-chan types.BlenderEvent, cancel context.CancelFunc) renderProgressModel {
	p := progress.New(
		progress.WithGradient("#4E51D0", "#E06F5A"),
		progress.WithWidth(25),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return renderProgressModel{
		progress:    p,
		spinner:     s,
		eventChan:   eventChan,
		totalFrames: totalFrames,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
		cancel:      cancel,
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
		case "ctrl+c", "esc", "q":
			m.earlyExit = true
			m.cancel()
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 30
		if m.progress.Width < 10 {
			m.progress.Width = 10
		}
		return m, nil

	case types.BlenderEvent:
		if renderEvent, ok := msg.(*types.RenderingEvent); ok {
			return m.handleRenderEvent(renderEvent)
		}

		if savedFileEvent, ok := msg.(*types.SavedFileEvent); ok {
			return m.handleSavedFileEvent(savedFileEvent)
		}

		if genericEvent, ok := msg.(*types.GenericEvent); ok {
			tea.Printf("%s %s\n", checkMark, genericEvent.Message)
		}

		if errorEvent, ok := msg.(*types.ErrorEvent); ok {
			m.errorMessage = errorEvent.Message
			return m, tea.Quit
		}

		return m, waitForBlenderEvent(m.eventChan)

	default:
		return m, nil
	}

	return m, nil
}

func (m *renderProgressModel) View() string {
	if m.done {
		totalTime := time.Since(m.startTime).Truncate(time.Second)
		avgTime := totalTime / time.Duration(len(m.completedFrames))
		return successStyle.Render(fmt.Sprintf(
			"Rendered %d frames in %s (avg: %s per frame).\n",
			m.totalFrames, totalTime, avgTime,
		))
	}

	if m.errorMessage != "" {
		return errorStyle.Render(fmt.Sprintf("Error: %s\n", m.errorMessage))
	}

	if m.earlyExit {
		totalTime := time.Since(m.startTime).Truncate(time.Second)
		return warningStyle.Render(fmt.Sprintf(
			"Rendering cancelled. Completed %d/%d frames in %s.\n",
			len(m.completedFrames), m.totalFrames, totalTime,
		))
	}

	if m.currentFrame == 0 {
		waitingMessage := fmt.Sprintf(
			"%s Waiting for Blender to start rendering...",
			m.spinner.View(),
		)

		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			waitingMessage,
			progressStyle.Render(m.progress.View()),
			renderLegend(),
		)
	}

	visibleFrames := strings.Join(m.getVisibleFrames(), "\n")

	status := m.renderStatus()
	prog := progressStyle.Render(m.progress.View())

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n%s",
		visibleFrames,
		status,
		prog,
		renderLegend(),
	)
}

func (m *renderProgressModel) renderStatus() string {
	elapsed := time.Since(m.startTime).Truncate(time.Second)
	var eta string
	if len(m.completedFrames) > 0 {
		avgTime := elapsed / time.Duration(len(m.completedFrames))
		remainingFrames := m.totalFrames - len(m.completedFrames)
		eta = (time.Duration(remainingFrames) * avgTime).Truncate(time.Second).String()
	} else {
		eta = "Calculating..."
	}

	memory := "N/A"
	if m.currentMemory != "" {
		memory = m.currentMemory
	}

	return fmt.Sprintf(
		"Frame %s/%s | Elapsed: %s | ETA: %s | Memory: %s",
		currentFrameStyle.Render(fmt.Sprintf("%d", m.currentFrame)),
		fmt.Sprintf("%d", m.totalFrames),
		elapsed, eta, memory,
	)
}

func (m *renderProgressModel) getVisibleFrames() []string {
	visibleFrames := []string{}

	for _, frame := range m.completedFrames {
		if frame > 0 {
			visibleFrames = append(visibleFrames, checkMark.String()+" "+fmt.Sprintf("Frame %d", frame))
		}
	}

	frameProgress := progress.New(
		progress.WithGradient("#64B5D9", "#4E51D0"),
		progress.WithWidth(25),
		progress.WithoutPercentage(),
	)

	progressPercent := 0.0
	if m.totalSamples > 0 {
		progressPercent = float64(m.currentSample) / float64(m.totalSamples)
	}

	frameProgressBar := frameProgress.ViewAs(progressPercent)
	currentFrameView := fmt.Sprintf("%s Frame %d (progress): %s", m.spinner.View(), m.currentFrame, frameProgressBar)
	visibleFrames = append(visibleFrames, currentFrameView)

	return visibleFrames
}

func (m *renderProgressModel) handleRenderEvent(e *types.RenderingEvent) (tea.Model, tea.Cmd) {
	m.currentSample = e.Current
	m.totalSamples = e.Total
	m.currentMemory = e.Memory
	m.currentFrame = e.Frame

	totalProgress := (float64(len(m.completedFrames))*float64(m.totalSamples) + float64(m.currentSample)) /
		(float64(m.totalFrames) * float64(m.totalSamples))
	progressCmd := m.progress.SetPercent(totalProgress)

	tea.Printf(
		"%s Frame: %s, Progress: %d/%d, Memory: %s, Peak Memory: %s\n",
		checkMark,
		currentFrameStyle.Render(fmt.Sprintf("%d", e.Frame)),
		e.Current, e.Total,
		e.Memory, e.PeakMemory,
	)

	return m, tea.Batch(progressCmd, waitForBlenderEvent(m.eventChan))
}

func (m *renderProgressModel) handleSavedFileEvent(e *types.SavedFileEvent) (tea.Model, tea.Cmd) {
	// Cycles noise threshold can cause the frame to be completed before the last sample.
	m.currentSample = m.totalSamples
	m.completedFrames = append(m.completedFrames, m.currentFrame)
	tea.Printf("%s Frame %d saved: %s\n", checkMark, m.currentFrame, e.Path)

	// Calculate total progress based on completed frames and current frame's progress
	totalProgress := (float64(len(m.completedFrames)) * float64(m.totalSamples)) /
		(float64(m.totalFrames) * float64(m.totalSamples))
	progressCmd := m.progress.SetPercent(totalProgress)

	if len(m.completedFrames) == m.totalFrames {
		m.done = true
		return m, tea.Quit
	}

	return m, tea.Batch(progressCmd, waitForBlenderEvent(m.eventChan))
}

func waitForBlenderEvent(eventChan <-chan types.BlenderEvent) tea.Cmd {
	return func() tea.Msg {
		return <-eventChan
	}
}
