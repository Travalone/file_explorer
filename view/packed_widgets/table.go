package packed_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TableColumn struct {
	Header             string                          // 列头文本内容
	GetText            func(interface{}) string        // 每行的对应列显示的文本内容
	Width              int                             // 每列宽度，表格每行内容不固定，只能固定宽度
	OrderBy            func(desc bool, inited bool)    // reload用于区分文件列表已经加载时排序和初次打开目录时排序 (有无CheckList)
	GetIcon            func(interface{}) fyne.Resource // 图标
	DefaultOrderBy     bool                            // 默认排序列
	DoubleTappedExpand bool                            // 双击放大，会覆盖Table.OnDoubleTapped
}

type Table struct {
	*widget.Table

	columns       []*TableColumn
	getData       func(index int) interface{}
	getDataLength func() int
	CheckList     *CheckList

	orderColumn int
	orderDesc   bool

	OnCheckChange  func([]int)
	OnTapped       func(row int, col int)
	OnDoubleTapped func(row int, col int)
}

func (table *Table) RefreshData() {
	if table.columns[table.orderColumn].OrderBy != nil {
		table.columns[table.orderColumn].OrderBy(table.orderDesc, false)
	}
	table.Table.Refresh()
	table.ScrollToTop()

	table.CheckList = NewCheckList(table.getDataLength())
	table.CheckList.OnCheckChange = func(indexList []int) {
		if table.OnCheckChange != nil {
			table.OnCheckChange(indexList)
		}
	}
}

func (table *Table) GetContent(row int, col int) string {
	if col == 0 {
		return ""
	}
	if row == 0 {
		header := table.columns[col-1].Header
		if table.orderColumn == col-1 {
			if table.orderDesc {
				header += " v"
			} else {
				header += " ^"
			}
		}
		return header
	}
	if table.columns[col-1].GetText == nil {
		return ""
	}
	return table.columns[col-1].GetText(table.getData(row - 1))
}

func NewTable(columns []*TableColumn, getData func(index int) interface{}, getDataLength func() int) *Table {
	table := &Table{
		columns:       columns,
		getData:       getData,
		getDataLength: getDataLength,
	}
	table.Table = widget.NewTable(
		func() (int, int) {
			return 1 + table.getDataLength(), 1 + len(table.columns)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewIcon(nil),
				NewLabel(""),
			)
		},
		func(cell widget.TableCellID, obj fyne.CanvasObject) {
			if table.CheckList == nil {
				return
			}

			cellContainer := obj.(*fyne.Container)
			checkBox := cellContainer.Objects[0].(*widget.Check)
			icon := cellContainer.Objects[1].(*widget.Icon)
			label := cellContainer.Objects[2].(*Label)

			if cell.Col == 0 {
				// checkBox
				checkBox.Show()
				icon.Hide()
				label.Hide()
				if cell.Row == 0 {
					checkBox.OnChanged = table.CheckList.CheckAll
				} else {
					table.CheckList.Bind(cell.Row-1, checkBox)
				}
			} else {
				// Label
				checkBox.Hide()
				icon.Hide()
				label.Show()

				text := table.GetContent(cell.Row, cell.Col)
				width := float32(table.columns[cell.Col-1].Width)
				label.SetTextWithFixWidth(text, width)
				if label.Size().Width != width {
					label.Resize(fyne.NewSize(width, label.Size().Height))
				}

				if cell.Row > 0 && table.columns[cell.Col-1].GetIcon != nil {
					icon.Resource = table.columns[cell.Col-1].GetIcon(table.getData(cell.Row - 1))
					icon.Show()
				}

				if table.OnTapped != nil {
					label.OnTapped = func() {
						// 点击header时排序
						if cell.Row == 0 && cell.Col > 0 && table.columns[cell.Col-1].OrderBy != nil {
							if table.orderColumn != cell.Col-1 {
								// 切换排序列
								table.orderColumn = cell.Col - 1
								table.orderDesc = false
							} else {
								// 切换顺序
								table.orderDesc = !table.orderDesc
							}
							// 调用newTable时定义的排序方法
							table.columns[cell.Col-1].OrderBy(table.orderDesc, true)
							table.Refresh()
						}

						table.OnTapped(cell.Row, cell.Col)
					}
				}

				if table.columns[cell.Col-1].DoubleTappedExpand {
					label.OnDoubleTapped = func() {
						textEntry := widget.NewEntry()
						textEntry.Text = text
						textEntry.MultiLine = true
						textEntry.Wrapping = fyne.TextWrapBreak //自动换行

						popUp := NewPopUp(obj, textEntry)
						popUp.Resize(fyne.NewSize(400, 200))
						popUp.Show()
					}
				} else if table.OnDoubleTapped != nil {
					label.OnDoubleTapped = func() {
						table.OnDoubleTapped(cell.Row, cell.Col)
					}
				}

			}
		},
	)

	table.SetColumnWidth(0, 35)
	for i, column := range columns {
		table.SetColumnWidth(i+1, float32(column.Width))
		if column.DefaultOrderBy {
			table.orderColumn = i
		}
	}

	return table
}
