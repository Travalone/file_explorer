package side_panel

import (
	"file_explorer/common/utils"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2/widget"
)

func NewSearchBox(feContext *store.FeContext) *widget.Entry {
	searchBox := widget.NewEntry()
	searchBox.PlaceHolder = "Path"
	searchBox.OnSubmitted = func(path string) {
		if utils.PathExists(path) {
			fileTabContext := store.NewFileTabContext(path, feContext)
			feContext.AddTab(fileTabContext)
			return
		}
		packed_widgets.NewNotifyFmt("Invalid path: %s", path)
	}

	return searchBox
}
