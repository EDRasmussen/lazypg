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

	Focus      FocusArea
	ColumnMode bool
	FocusedCol int

	Layout Layout
	Styles Styles
}

func InitialModel(sess *session.Session) Model {
	input := textarea.New()
	input.SetValue("SELECT * FROM users LIMIT 100;")
	input.Focus()
	input.ShowLineNumbers = true
	input.DynamicHeight = true
	input.MinHeight = 1
	input.MaxHeight = 10
	input.MaxContentHeight = 100_000 // should be enough

	tbl := table.New()

	sidebar := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	sidebar.Title = "lazypg"

	return Model{
		Session: sess,
		Sidebar: sidebar,
		Input:   input,
		Table:   tbl,
		Focus:   FocusInput,
		Styles:  NewStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
