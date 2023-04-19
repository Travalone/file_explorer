package file_create_tab

import (
	"file_explorer/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func newCreateItem() fyne.CanvasObject {
	deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)

	typeSelect := widget.NewSelect([]string{"文件", "目录"}, nil)
	typeSelect.Selected = "文件"

	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "名称"

	scoreBox := widget.NewEntry()
	scoreBox.PlaceHolder = "评分"
	scoreBox.Validator = func(score string) error {
		if len(score) == 0 {
			return nil
		}
		_, err := strconv.Atoi(scoreBox.Text)
		if err != nil {
			return common.NewError("Input score is not number")
		}
		return nil
	}
	scoreEditBox := container.NewBorder(
		nil, nil,
		widget.NewButton("-", func() {
			scoreInt, err := strconv.Atoi(scoreBox.Text)
			if err != nil {
				scoreInt = 0
			}
			scoreBox.SetText(strconv.Itoa(scoreInt - 1))
		}),
		widget.NewButton("+", func() {
			scoreInt, err := strconv.Atoi(scoreBox.Text)
			if err != nil {
				scoreInt = 0
			}
			scoreBox.SetText(strconv.Itoa(scoreInt + 1))
		}),
		scoreBox,
	)

	noteBox := widget.NewEntry()
	noteBox.MultiLine = true
	noteBox.PlaceHolder = "备注"

	urlBox := widget.NewEntry()
	urlBox.PlaceHolder = "链接"

	tagBox := widget.NewEntry()
	tagBox.PlaceHolder = "标签 (逗号隔开)"

	inputBoxes := container.NewGridWithColumns(2,
		container.NewGridWithRows(3,
			container.NewGridWithColumns(2, typeSelect, nameEntry),
			container.NewGridWithColumns(2, scoreEditBox, tagBox),
			urlBox,
		),
		noteBox,
	)
	return container.NewBorder(nil, nil, deleteButton, nil, inputBoxes)
}

func parseCreateItem(obj fyne.CanvasObject) (*widget.Button, *widget.Select, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Entry, *widget.Entry) {
	container := obj.(*fyne.Container)
	inputBoxes, deleteButton := container.Objects[0].(*fyne.Container), container.Objects[1].(*widget.Button)
	leftBoxes, noteBox := inputBoxes.Objects[0].(*fyne.Container), inputBoxes.Objects[1].(*widget.Entry)
	lf1Boxes, lf2Boxes, urlBox := leftBoxes.Objects[0].(*fyne.Container), leftBoxes.Objects[1].(*fyne.Container), leftBoxes.Objects[2].(*widget.Entry)
	typeSelectEntry, nameEntry := lf1Boxes.Objects[0].(*widget.Select), lf1Boxes.Objects[1].(*widget.Entry)
	scoreEditBox, tagsBox := lf2Boxes.Objects[0].(*fyne.Container), lf2Boxes.Objects[1].(*widget.Entry)
	scoreInputBox := scoreEditBox.Objects[0].(*widget.Entry)
	return deleteButton, typeSelectEntry, nameEntry, scoreInputBox, tagsBox, urlBox, noteBox
}
