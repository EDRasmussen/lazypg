package tui

import (
	"image/color"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
)

const (
	colorSurface    = "235"
	colorPanel      = "236"
	colorBorder     = "238"
	colorAccent     = "69"
	colorAccentSoft = "111"
	colorText       = "252"
	colorTextMuted  = "245"
	colorTextSubtle = "242"
	colorHeaderText = "230"
	colorHeaderBg   = "60"
)

type Styles struct {
	Sidebar          lipgloss.Style
	Input            lipgloss.Style
	Table            lipgloss.Style
	TableHeader      lipgloss.Style
	Cell             lipgloss.Style
	OddRow           lipgloss.Style
	EvenRow          lipgloss.Style
	TableBorder      lipgloss.Border
	TableBorderColor color.Color
	Status           lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		Sidebar: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorder)).
			Background(lipgloss.Color(colorSurface)).
			Foreground(lipgloss.Color(colorText)).
			Padding(0, 1),

		Input: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorder)).
			Background(lipgloss.Color(colorSurface)).
			Foreground(lipgloss.Color(colorText)).
			Padding(0, 1),

		Table: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorder)).
			Background(lipgloss.Color(colorPanel)).
			Padding(0, 1),

		TableHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorHeaderText)).
			Background(lipgloss.Color(colorHeaderBg)).
			Padding(0, 1),

		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorText)).
			Padding(0, 1),

		OddRow: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorText)),

		EvenRow: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextMuted)),

		TableBorder: lipgloss.RoundedBorder(),

		TableBorderColor: lipgloss.Color(colorBorder),

		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorText)).
			Background(lipgloss.Color(colorSurface)).
			Padding(0, 1),
	}
}

func (s Styles) TableStyles() table.Styles {
	styles := table.DefaultStyles()
	styles.Header = s.TableHeader
	styles.Cell = s.Cell
	styles.Selected = lipgloss.NewStyle()
	return styles
}

func (s Styles) SelectedRowStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorHeaderText))
}

func (s Styles) SelectedCellStyle(active bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorHeaderText))

	if active {
		return style.Background(lipgloss.Color(colorAccent))
	}

	return style.Background(lipgloss.Color(colorHeaderBg))
}

func (s Styles) InputStyles() textarea.Styles {
	styles := textarea.DefaultStyles(true)

	styles.Focused.Base = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorSurface))
	styles.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
	styles.Focused.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextSubtle))
	styles.Focused.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color(colorPanel))
	styles.Focused.CursorLineNumber = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorAccentSoft))
	styles.Focused.EndOfBuffer = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorder))
	styles.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextSubtle))
	styles.Focused.Prompt = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorAccent))

	styles.Blurred.Base = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextMuted)).
		Background(lipgloss.Color(colorSurface))
	styles.Blurred.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted))
	styles.Blurred.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorder))
	styles.Blurred.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color(colorSurface))
	styles.Blurred.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted))
	styles.Blurred.EndOfBuffer = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorder))
	styles.Blurred.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextSubtle))
	styles.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted))

	styles.Cursor.Color = lipgloss.Color(colorAccentSoft)

	return styles
}

func (s Styles) SidebarStyles() list.Styles {
	styles := list.DefaultStyles(true)

	styles.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 0)
	styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorHeaderText)).
		Background(lipgloss.Color(colorAccent)).
		Padding(0, 1)
	styles.Filter.Cursor.Color = lipgloss.Color(colorAccentSoft)
	styles.Filter.Focused.Prompt = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorAccent))
	styles.Filter.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted))
	styles.StatusBar = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextMuted))
	styles.StatusEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextSubtle))
	styles.NoItems = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color(colorTextSubtle))
	styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorTextSubtle))
	styles.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		SetString("•")
	styles.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBorder)).
		SetString("•")

	return styles
}

func (s Styles) SidebarItemStyles() list.DefaultItemStyles {
	styles := list.NewDefaultItemStyles(true)

	styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 0, 0, 1)
	styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextMuted)).
		Padding(0, 0, 0, 1)
	styles.SelectedTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorHeaderText)).
		Background(lipgloss.Color(colorHeaderBg)).
		Padding(0, 1)
	styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorHeaderBg)).
		Padding(0, 1)
	styles.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextSubtle)).
		Padding(0, 0, 0, 1)
	styles.DimmedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBorder)).
		Padding(0, 0, 0, 1)
	styles.FilterMatch = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorAccentSoft))

	return styles
}

func (m Model) FocusedStyle(style lipgloss.Style, focused bool) lipgloss.Style {
	if focused {
		return style.BorderForeground(lipgloss.Color(colorAccent))
	}

	return style.BorderForeground(lipgloss.Color(colorBorder))
}
