package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) View() tea.View {
	if !m.Ready {
		return tea.NewView("loading...")
	}

	sidebarStyle := m.FocusedStyle(m.Styles.Sidebar, m.Focus == FocusSidebar).
		Width(m.Layout.SidebarWidth).
		Height(m.Height)

	inputStyle := m.FocusedStyle(m.Styles.Input, m.Focus == FocusInput).
		Width(m.Layout.MainWidth).
		Height(m.Layout.InputHeight)

	tableStyle := m.FocusedStyle(m.Styles.Table, m.Focus == FocusRows).
		Width(m.Layout.MainWidth).
		Height(m.Layout.TableHeight)

	input := inputStyle.Render(m.Input.View())
	table := tableStyle.Render(m.Table.View())
	status := m.Styles.Status.
		Width(m.Layout.MainWidth).
		Render(m.statusText(m.Layout.MainWidth))

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		input,
		table,
		status,
	)

	if m.Layout.SidebarWidth == 0 {
		return tea.NewView(right)
	}

	sidebar := sidebarStyle.Render(m.Sidebar.View())

	return tea.NewView(lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		right,
	))
}

func (m Model) statusText(width int) string {
	focus := "input"

	switch m.Focus {
	case FocusSidebar:
		focus = "sidebar"
	case FocusInput:
		focus = "input"
	case FocusRows:
		focus = "rows"
	}

	contentWidth := max(1, width-m.Styles.Status.GetHorizontalFrameSize())
	if m.Focus == FocusRows && m.ColumnMode {
		return strings.Join([]string{
			ansi.Truncate(fmt.Sprintf("focus: %s • mode: columns %d/%d • left/right or h/l: move • esc: rows • tab: switch • q: quit • shift+enter/f5: execute query", focus, m.FocusedCol+1, len(m.TableCols)), contentWidth, "…"),
			ansi.Truncate("column: "+m.focusedColumnTitle(), contentWidth, "…"),
			ansi.Truncate("value: "+m.focusedCellValue(), contentWidth, "…"),
		}, "\n")
	}

	parts := []string{fmt.Sprintf("focus: %s", focus)}

	if m.Focus == FocusRows {
		if m.ColumnMode {
			parts = append(parts, fmt.Sprintf("mode: columns %d/%d", m.FocusedCol+1, len(m.TableCols)))
			if title := m.focusedColumnTitle(); title != "" {
				parts = append(parts, fmt.Sprintf("column: %s", title))
			}
			if value := m.focusedCellValue(); value != "" {
				parts = append(parts, fmt.Sprintf("value: %s", value))
			}
			parts = append(parts, "left/right or h/l: move", "esc: rows")
		} else {
			parts = append(parts, "enter: columns")
		}
	}

	parts = append(parts, "tab: switch", "q: quit", "shift+enter/f5: execute query")
	text := strings.Join(parts, " • ")
	return ansi.Truncate(text, contentWidth, "…")
}

func (m Model) focusedColumnTitle() string {
	if m.FocusedCol < 0 || m.FocusedCol >= len(m.TableCols) {
		return ""
	}

	return m.TableCols[m.FocusedCol].Title
}

func (m Model) focusedCellValue() string {
	row := m.Table.SelectedRow()
	if m.FocusedCol < 0 || m.FocusedCol >= len(row) {
		return ""
	}

	return row[m.FocusedCol]
}
