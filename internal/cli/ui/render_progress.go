package ui

import (
	"fmt"
	"strconv"
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
	totalFrames     int
	currentSample   int
	totalSamples    int
	startTime       time.Time
	done            bool
	earlyExit       bool
	completedFrames []int
	displayLimit    int
}

var (
	progressStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	currentFrameStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle          = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.Color("42"))
	earlyExitStyle     = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.Color("9"))
	checkMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
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
		progress:     p,
		spinner:      s,
		eventChan:    eventChan,
		totalFrames:  totalFrames,
		startTime:    time.Now(),
		displayLimit: 5, // Show 5 frames: 1 + ellipsis + 3 recent + current
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
		m.progress.Width = msg.Width - 10 // Dynamically adjust progress bar width
		return m, nil

	case types.BlenderEvent:
		if renderEvent, ok := msg.(*types.RenderEvent); ok {
			return m.handleRenderEvent(renderEvent)
		}
		if genericEvent, ok := msg.(*types.GenericEvent); ok {
			tea.Printf("%s %s\n", checkMark, genericEvent.Message)
		}
		return m, waitForBlenderEvent(m.eventChan)

	default:
		return m, nil
	}

	return m, nil
}

func (m *renderProgressModel) View() string {
	if m.done {
		return doneStyle.Render(fmt.Sprintf("Rendering complete! Rendered %d frames.\n", m.totalFrames))
	}

	if m.earlyExit {
		return earlyExitStyle.Render(fmt.Sprintf("Rendering stopped early. Completed %d/%d frames.\n", len(m.completedFrames), m.totalFrames))
	}

	spin := m.spinner.View()
	prog := progressStyle.Render(m.progress.View())

	status := fmt.Sprintf(
		"Rendering frame %s/%s | Elapsed: %s",
		currentFrameStyle.Render(fmt.Sprintf("%d", m.currentFrame)),
		fmt.Sprintf("%d", m.totalFrames),
		time.Since(m.startTime).Truncate(time.Second),
	)

	visibleFrames := strings.Join(m.getVisibleFrames(), "\n")

	return fmt.Sprintf(
		"%s %s\n%s\n\n%s",
		spin, statusMessageStyle.Render(status), prog, visibleFrames,
	)
}

func (m *renderProgressModel) getVisibleFrames() []string {
	visibleFrames := []string{}

	if len(m.completedFrames) > 0 {
		visibleFrames = append(visibleFrames, checkMark.String()+" "+fmt.Sprintf("Frame %d", m.completedFrames[0]))
	}

	if len(m.completedFrames) > m.displayLimit {
		visibleFrames = append(visibleFrames, "...")
	}

	start := len(m.completedFrames) - (m.displayLimit - 1)
	if start < 1 { // Ensure we don't include the first frame again
		start = 1
	}

	if start < len(m.completedFrames) { // Ensure start is within bounds
		for _, frame := range m.completedFrames[start:] {
			visibleFrames = append(visibleFrames, checkMark.String()+" "+fmt.Sprintf("Frame %d", frame))
		}
	}

	frameProgress := progress.New(
		progress.WithGradient("yellow", "green"),
		progress.WithWidth(30),
		progress.WithoutPercentage(),
	)

	progressPercent := 0.0
	if m.totalSamples > 0 {
		progressPercent = float64(m.currentSample) / float64(m.totalSamples)
	}
	frameProgressBar := frameProgress.ViewAs(progressPercent)
	visibleFrames = append(visibleFrames, fmt.Sprintf("Frame %d (in progress): %s", m.currentFrame, frameProgressBar))

	return visibleFrames
}

func (m *renderProgressModel) handleRenderEvent(e *types.RenderEvent) (tea.Model, tea.Cmd) {
	if e.Frame != m.currentFrame {
		m.currentFrame = e.Frame
		m.completedFrames = append(m.completedFrames, e.Frame) // Mark frame as completed

		if len(m.completedFrames) > m.displayLimit {
			m.completedFrames = m.completedFrames[1:]
		}
	}

	currentSample, _ := strconv.Atoi(e.Data["current"])
	totalSamples, _ := strconv.Atoi(e.Data["total"])
	m.currentSample = currentSample
	m.totalSamples = totalSamples

	frameProgress := 0.0
	if totalSamples > 0 {
		frameProgress = float64(currentSample) / float64(totalSamples)
	}

	totalProgress := (float64(m.currentFrame-1) + frameProgress) / float64(m.totalFrames)
	progressCmd := m.progress.SetPercent(totalProgress)

	tea.Printf(
		"%s Frame: %s, Progress: %d/%d, Memory: %s, Peak Memory: %s\n",
		checkMark,
		currentFrameStyle.Render(fmt.Sprintf("%d", m.currentFrame)),
		currentSample, totalSamples,
		e.Memory, e.PeakMemory,
	)

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
