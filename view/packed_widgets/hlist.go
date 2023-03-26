package packed_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HList struct {
	*widget.Table
	FixedWidth int // >0时固定长度，否则自适应

	getData   func(index int) string
	getLength func() int

	OnTapped       func(col int)
	OnDoubleTapped func(col int)
}

func (hList *HList) RefreshData() {
	// 调整每列宽度
	for i := 0; i < hList.getLength(); i++ {
		width := float32(hList.FixedWidth)
		if width <= 0 {
			width = GetMinWidth(hList.getData(i))
		}
		hList.SetColumnWidth(i, width)
	}

	hList.Table.Refresh()
}

func NewHList(width int, getData func(index int) string, getDataLength func() int) *HList {
	hList := &HList{
		FixedWidth: width,
		getData:    getData,
		getLength:  getDataLength,
	}
	hList.Table = widget.NewTable(
		func() (int, int) {
			return 1, hList.getLength()
		},
		func() fyne.CanvasObject {
			return NewLabel("")
		},
		func(cell widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*Label)
			text := hList.getData(cell.Col)
			if hList.FixedWidth > 0 {
				label.SetTextWithFixWidth(text, float32(hList.FixedWidth))
			} else {
				label.SetText(text)
			}

			if hList.OnTapped != nil {
				label.OnTapped = func() {
					hList.OnTapped(cell.Col)
				}
			}
			if hList.OnDoubleTapped != nil {
				label.OnDoubleTapped = func() {
					hList.OnDoubleTapped(cell.Col)
				}
			}
		},
	)

	return hList
}
