package ui

import "github.com/charmbracelet/lipgloss"

var (
	stepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	legendStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true)
	checkMark    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func renderLegend() string {
	return legendStyle.Render("Press 'q' or 'esc' to cancel.")
}
