package extra_info_tab

import (
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fmt"
)

type StatusBar struct {
	*packed_widgets.Label

	tabContext *store.ExtraInfoTabContext
}

func (bar *StatusBar) RefreshStatus() {
	text := fmt.Sprintf("%d项", len(bar.tabContext.FileInfos))
	if bar.tabContext.CheckList.Length() > 0 {
		text += fmt.Sprintf(", %d项被选中", bar.tabContext.CheckList.Length())
	}
	bar.SetText(text)
}

func NewStatusBar(tabContext *store.ExtraInfoTabContext) *StatusBar {
	return &StatusBar{
		Label:      packed_widgets.NewLabel(""),
		tabContext: tabContext,
	}
}
