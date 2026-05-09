package tui

import (
	"era/lazypg/internal/session"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusInput
	FocusRows
)

type Model struct {
	Width   int
	Height  int
	Ready   bool
	Loading bool

	Session *session.Session

	Sidebar   list.Model
	Input     textarea.Model
	Table     table.Model
	TableCols []table.Column
	TableRows []table.Row

	Focus      FocusArea
	ColumnMode bool
	FocusedCol int

	Layout Layout
	Styles Styles
}

func InitialModel(sess *session.Session) Model {
	styles := NewStyles()

	input := textarea.New()
	input.SetValue("SELECT * FROM users LIMIT 100;")
	input.SetStyles(styles.InputStyles())
	input.Focus()
	input.ShowLineNumbers = true
	input.DynamicHeight = true
	input.MinHeight = 1
	input.MaxHeight = 10
	input.MaxContentHeight = 100_000 // should be enough

	tbl := table.New()
	tbl.SetStyles(styles.TableStyles())

	delegate := list.NewDefaultDelegate()
	delegate.Styles = styles.SidebarItemStyles()

	sidebar := list.New([]list.Item{}, delegate, 0, 0)
	sidebar.Title = "lazypg"
	sidebar.Styles = styles.SidebarStyles()
	sidebar.SetShowStatusBar(false)
	sidebar.SetShowHelp(false)

	return Model{
		Session: sess,
		Sidebar: sidebar,
		Input:   input,
		Table:   tbl,
		Focus:   FocusInput,
		Styles:  styles,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
