package tui

import "charm.land/lipgloss/v2"

type Styles struct {
	Sidebar lipgloss.Style
	Input   lipgloss.Style
	Table   lipgloss.Style
	Status  lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		Sidebar: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		Table: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1),
	}
}

func (m Model) FocusedStyle(style lipgloss.Style, focused bool) lipgloss.Style {
	if focused {
		return style.BorderForeground(lipgloss.Color("63"))
	}

	return style.BorderForeground(lipgloss.Color("240"))
}
