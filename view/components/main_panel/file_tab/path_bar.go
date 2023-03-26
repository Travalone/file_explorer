package file_tab

import (
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
)

type PathBar struct {
	*packed_widgets.HList

	tabContext *store.FileTabContext
}

func NewPathBar(tabContext *store.FileTabContext) *PathBar {
	pathBar := &PathBar{
		HList: packed_widgets.NewHList(
			0,
			func(index int) string {
				tags := tabContext.Dirs.Get()
				return tags[index]
			},
			func() int {
				return tabContext.Dirs.Length()
			}),
		tabContext: tabContext,
	}

	// 地址栏单击后退
	pathBar.OnTapped = func(col int) {
		pathBar.tabContext.Back(col)
	}

	return pathBar
}
