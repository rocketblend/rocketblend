package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	message   string
	confirmed bool
	quitting  bool
}

func Confirm(message string) (bool, error) {
	m := confirmModel{
		message: message,
	}

	p := tea.NewProgram(&m)
	finalModel, err := p.Run()
	if err != nil {
		return false, fmt.Errorf("failed to run confirmation UI: %w", err)
	}

	result := finalModel.(*confirmModel)
	return result.confirmed, nil
}

func (m *confirmModel) Init() tea.Cmd {
	return nil
}

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "q", "esc", "ctrl+c":
			m.confirmed = false
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *confirmModel) View() string {
	if m.quitting {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	return fmt.Sprintf("%s\n\n%s",
		style.Render(m.message),
		promptStyle.Render("Continue? (y/n)"),
	)
}
