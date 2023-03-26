package packed_widgets

import (
	"file_explorer/common/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Label struct {
	widget.Label   // 不能用指针，否则table滚动时数据会乱
	OnTapped       func()
	OnDoubleTapped func()
}

func GetMinWidth(text string) float32 {
	tmpLabel := widget.NewLabel(text)
	return tmpLabel.MinSize().Width
}

func (label *Label) SetTextWithFixWidth(text interface{}, width float32) {
	label.SetText(text)
	if label.MinSize().Width <= width {
		return
	}
	tmpLabel := widget.NewLabel(label.Text)
	runeArr := []rune(tmpLabel.Text)
	for i := len(runeArr) - 1; i > 0; i-- {
		tmpLabel.SetText(string(runeArr[:i]) + "...")
		if tmpLabel.MinSize().Width <= width {
			break
		}
	}
	label.SetText(tmpLabel.Text)
}

func (label *Label) SetText(text interface{}) {
	t := utils.Conv2Str(text)
	if label.Label.Text != t {
		label.Label.SetText(t)
	}
}

func (label *Label) Tapped(*fyne.PointEvent) {
	if label.OnTapped != nil {
		label.OnTapped()
	}
}

func (label *Label) DoubleTapped(*fyne.PointEvent) {
	if label.OnDoubleTapped != nil {
		label.OnDoubleTapped()
	}
}

func NewLabel(text string) *Label {
	label := &Label{}
	label.SetText(text)
	label.ExtendBaseWidget(label) //不加这行，list滚动时顺序会乱
	return label
}
