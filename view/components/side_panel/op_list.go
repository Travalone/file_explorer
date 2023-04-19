package side_panel

import (
	"file_explorer/common/model"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fmt"
	"fyne.io/fyne/v2"
)

type OpList struct {
	Toolbar    *fyne.Container
	OpFileList *packed_widgets.Table
	Context    *store.OpListContext
}

func NewOpList() *OpList {
	columns := []*packed_widgets.TableColumn{
		{
			Header: "名称",
			Width:  115,
			OrderBy: func(desc bool, inited bool) {

			},
			GetText: func(data interface{}) string {
				return data.(*model.FileInfo).Name
			},
		},
		{
			Header: "ㄨ",
			Width:  30,
			OrderBy: func(desc bool, inited bool) {

			},
			GetText: func(data interface{}) string {
				return "ㄨ"
			},
		},
	}
	opList := packed_widgets.NewTable(columns,
		func(i int) interface{} {
			return &model.FileInfo{
				Name: fmt.Sprintf("%d", i+1),
				Type: byte(i + 1),
			}
		},
		func() int {
			return 2
		})

	opList.RefreshData()

	return nil
}
