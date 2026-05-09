package tui

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true
		m.resize()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {

		case "shift+enter", "f5":
			if err := m.executeQuery(); err != nil {
				m.printErr(err)
			}
			return m, nil

		case "enter":
			if m.Focus == FocusRows && len(m.TableCols) > 0 && !m.ColumnMode {
				m.ColumnMode = true
				m.resize()
				m.refreshTableColumns()
				m.refreshTableRows()
				return m, nil
			}

		case "esc":
			if m.Focus == FocusRows && m.ColumnMode {
				m.ColumnMode = false
				m.resize()
				m.refreshTableColumns()
				m.refreshTableRows()
				return m, nil
			}

		case "left", "h":
			if m.Focus == FocusRows && m.ColumnMode {
				m.moveFocusedColumn(-1)
				return m, nil
			}

		case "right", "l":
			if m.Focus == FocusRows && m.ColumnMode {
				m.moveFocusedColumn(1)
				return m, nil
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			m.Focus = nextFocus(m.Focus)
			m.syncFocus()
			return m, nil
		}
	}

	switch m.Focus {
	case FocusSidebar:
		var cmd tea.Cmd
		m.Sidebar, cmd = m.Sidebar.Update(msg)
		cmds = append(cmds, cmd)

	case FocusInput:
		var cmd tea.Cmd
		prevHeight := m.Input.Height()
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
		if m.Ready && m.Input.Height() != prevHeight {
			m.resize()
		}

	case FocusRows:
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
		m.refreshTableRows()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateTable(cols []table.Column, rows []table.Row) {
	m.TableCols = append(m.TableCols[:0], cols...)
	m.TableRows = cloneRows(rows)
	if len(m.TableCols) == 0 {
		m.ColumnMode = false
		m.FocusedCol = 0
	} else {
		m.FocusedCol = clampInt(m.FocusedCol, 0, len(m.TableCols)-1)
	}

	m.refreshTableColumns()
	m.refreshTableRows()
}

func (m *Model) executeQuery() error {
	query := strings.TrimSpace(m.Input.Value())
	if query == "" {
		return nil
	}

	cols, rows, err := m.Session.ExecuteQuery(context.Background(), query)
	if err != nil {
		return err
	}

	m.updateTable(cols, rows)
	return nil
}

func (m *Model) printErr(err error) {
	m.updateTable(
		[]table.Column{
			{
				Title: "Error",
				Width: m.Table.Width(),
			},
		},
		[]table.Row{
			{err.Error()},
		},
	)

}

func (m *Model) resize() {
	sidebarWidth := 28
	if m.Width < 100 {
		sidebarWidth = 22
	}
	if m.Width < 70 {
		sidebarWidth = 0
	}

	statusHeight := m.statusHeight()

	mainWidth := m.Width - sidebarWidth
	inputWidth := max(1, mainWidth-m.Styles.Input.GetHorizontalFrameSize())
	m.Input.SetWidth(inputWidth)

	inputHeight := m.Input.Height() + m.Styles.Input.GetVerticalFrameSize()
	tableHeight := m.Height - inputHeight - statusHeight

	if mainWidth < 1 {
		mainWidth = 1
	}

	if tableHeight < 1 {
		tableHeight = 1
	}

	m.Layout = NewLayout(sidebarWidth, mainWidth, inputHeight, tableHeight, statusHeight)
	m.resizeChildren()
}

func (m *Model) resizeChildren() {
	sidebarStyle := m.Styles.Sidebar.
		Width(m.Layout.SidebarWidth).
		Height(m.Height)

	inputStyle := m.Styles.Input.
		Width(m.Layout.MainWidth).
		Height(m.Layout.InputHeight)

	tableStyle := m.Styles.Table.
		Width(m.Layout.MainWidth).
		Height(m.Layout.TableHeight)

	m.Sidebar.SetSize(
		max(1, m.Layout.SidebarWidth-sidebarStyle.GetHorizontalFrameSize()),
		max(1, m.Height-sidebarStyle.GetVerticalFrameSize()),
	)

	m.Input.SetWidth(max(1, m.Layout.MainWidth-inputStyle.GetHorizontalFrameSize()))
	if !m.Input.DynamicHeight {
		m.Input.SetHeight(max(1, m.Layout.InputHeight-inputStyle.GetVerticalFrameSize()))
	}

	m.Table.SetWidth(max(1, m.Layout.MainWidth-tableStyle.GetHorizontalFrameSize()))
	m.Table.SetHeight(max(1, m.Layout.TableHeight-tableStyle.GetVerticalFrameSize()))

	if len(m.TableCols) > 0 {
		m.refreshTableColumns()
	}
}

func fitColumns(cols []table.Column, rows []table.Row, availableWidth int, focusedCol int, columnMode bool, cellPaddingWidth int) []table.Column {
	if len(cols) == 0 {
		return nil
	}

	fitted := make([]table.Column, len(cols))
	ideal := make([]int, len(cols))
	widths := make([]int, len(cols))

	contentWidth := max(len(cols), availableWidth-len(cols)*cellPaddingWidth)

	used := 0
	for i, col := range cols {
		fitted[i] = col

		ideal[i] = max(1, ansi.StringWidth(col.Title))
		for _, row := range rows {
			if i >= len(row) {
				continue
			}
			ideal[i] = max(ideal[i], ansi.StringWidth(row[i]))
		}

		widths[i] = min(ideal[i], max(3, ansi.StringWidth(col.Title)))
		used += widths[i]
	}

	for used > contentWidth && shrinkWidths(widths, ideal, 3) {
		used--
	}

	for used > contentWidth && shrinkWidths(widths, ideal, 1) {
		used--
	}

	for used < contentWidth {
		grew := false
		for i := range widths {
			if widths[i] >= ideal[i] {
				continue
			}
			widths[i]++
			used++
			grew = true
			if used == contentWidth {
				break
			}
		}
		if !grew {
			break
		}
	}

	for i := 0; used < contentWidth; i = (i + 1) % len(widths) {
		widths[i]++
		used++
	}

	for i := range fitted {
		fitted[i].Width = widths[i]
	}

	if columnMode {
		rebalanceFocusedWidth(fitted, ideal, focusedCol)
	}

	return fitted
}

func rebalanceFocusedWidth(cols []table.Column, ideal []int, focusedCol int) {
	if focusedCol < 0 || focusedCol >= len(cols) {
		return
	}

	current := cols[focusedCol].Width
	target := min(ideal[focusedCol], max(current, 36))
	for current < target {
		donor := widestDonor(cols, ideal, focusedCol)
		if donor == -1 {
			break
		}

		cols[donor].Width--
		cols[focusedCol].Width++
		current++
	}
}

func widestDonor(cols []table.Column, ideal []int, focusedCol int) int {
	donor := -1
	for i := range cols {
		if i == focusedCol {
			continue
		}

		minWidth := 3
		if ideal[i] < minWidth {
			minWidth = ideal[i]
		}
		if minWidth < 1 {
			minWidth = 1
		}

		if cols[i].Width <= minWidth {
			continue
		}

		if donor == -1 || cols[i].Width > cols[donor].Width {
			donor = i
		}
	}

	return donor
}

func (m Model) statusHeight() int {
	if m.Focus == FocusRows && m.ColumnMode {
		return 3
	}

	return 1
}

func clampInt(v, low, high int) int {
	if high < low {
		return low
	}

	return min(max(v, low), high)
}

func shrinkWidths(widths []int, ideal []int, floor int) bool {
	for i := range widths {
		minWidth := floor
		if ideal[i] < minWidth {
			minWidth = ideal[i]
		}
		if minWidth < 1 {
			minWidth = 1
		}

		if widths[i] <= minWidth {
			continue
		}

		widths[i]--
		return true
	}

	return false
}

func tableCellPaddingWidth(styles Styles) int {
	return max(styles.TableHeader.GetHorizontalFrameSize(), styles.Cell.GetHorizontalFrameSize())
}

func cloneRows(rows []table.Row) []table.Row {
	cloned := make([]table.Row, len(rows))
	for i, row := range rows {
		cloned[i] = append(table.Row(nil), row...)
	}

	return cloned
}

func (m *Model) refreshTableRows() {
	if len(m.TableRows) == 0 {
		m.Table.SetRows(nil)
		return
	}

	rows := cloneRows(m.TableRows)
	cursor := m.Table.Cursor()
	highlightCell := m.Focus == FocusRows && m.ColumnMode

	for i, row := range rows {
		baseStyle := m.Styles.OddRow
		if i%2 == 1 {
			baseStyle = m.Styles.EvenRow
		}

		for j, value := range row {
			style := baseStyle
			if i == cursor {
				style = lipgloss.NewStyle().Inherit(style).Inherit(m.Styles.SelectedRowStyle())
			}
			if highlightCell && i == cursor && j == m.FocusedCol {
				style = lipgloss.NewStyle().Inherit(style).Inherit(m.Styles.SelectedCellStyle(m.Focus == FocusRows))
			}

			rows[i][j] = style.Render(value)
		}
	}

	m.Table.SetRows(rows)
}

func (m *Model) refreshTableColumns() {
	if len(m.TableCols) == 0 {
		m.Table.SetColumns(nil)
		return
	}

	cols := append([]table.Column(nil), m.TableCols...)
	if m.ColumnMode && m.FocusedCol >= 0 && m.FocusedCol < len(cols) {
		cols[m.FocusedCol].Title = "[" + cols[m.FocusedCol].Title + "]"
	}

	m.Table.SetColumns(fitColumns(cols, m.TableRows, m.Table.Width(), m.FocusedCol, m.ColumnMode, tableCellPaddingWidth(m.Styles)))
}

func (m *Model) moveFocusedColumn(delta int) {
	if len(m.TableCols) == 0 {
		return
	}

	next := clampInt(m.FocusedCol+delta, 0, len(m.TableCols)-1)
	if next == m.FocusedCol {
		return
	}

	m.FocusedCol = next
	m.refreshTableColumns()
	m.refreshTableRows()
}

func (m *Model) syncFocus() {
	resized := false
	if m.Focus != FocusRows && m.ColumnMode {
		m.ColumnMode = false
		resized = true
		m.refreshTableColumns()
	}

	if resized && m.Ready {
		m.resize()
	}

	if m.Focus == FocusInput {
		m.Input.Focus()
	} else {
		m.Input.Blur()
	}

	if m.Focus == FocusRows {
		m.Table.Focus()
	} else {
		m.Table.Blur()
	}

	m.refreshTableRows()
}

func nextFocus(f FocusArea) FocusArea {
	switch f {
	case FocusSidebar:
		return FocusInput
	case FocusInput:
		return FocusRows
	default:
		return FocusSidebar
	}
}
