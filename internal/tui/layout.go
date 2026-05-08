package tui

type Layout struct {
	SidebarWidth int
	MainWidth    int
	InputHeight  int
	TableHeight  int
	StatusHeight int
}

func NewLayout(sidebarWidth int, mainWidth int, inputHeight int, tableHeight int, statusHeight int) Layout {
	return Layout{
		SidebarWidth: sidebarWidth,
		MainWidth:    mainWidth,
		InputHeight:  inputHeight,
		TableHeight:  tableHeight,
		StatusHeight: statusHeight,
	}
}
