package test

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"testing"
)

func getStrWidth(window fyne.Window, s string) float32 {
	label := widget.NewLabel("")
	window.SetContent(container.NewVBox(label))
	return label.Size().Width - 16
}

func TestUtils(t *testing.T) {
	s := "Coalamode - ビューティフルデイズ (Beautiful Days).mp3"
	runeArr := []rune(s)
	for i := len(runeArr) - 1; i > 0; i-- {
		logger.Info("%v", string(runeArr[i]))
	}
}

func TestPathUtils(t *testing.T) {
	path := "E:/111"
	logger.Info("%v", utils.PathSplit(path))
}
