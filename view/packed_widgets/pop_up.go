package packed_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PopUp struct {
	*widget.PopUp
	parent fyne.CanvasObject
}

func (popUp *PopUp) Show() {
	parentPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(popUp.parent)
	pos := fyne.NewPos(parentPos.X, parentPos.Y+theme.Padding()+popUp.parent.Size().Height)
	popUp.ShowAtPosition(pos)
}

func NewPopUp(parent fyne.CanvasObject, content fyne.CanvasObject) *PopUp {
	popUp := &PopUp{
		PopUp: widget.NewPopUp(
			content,
			GetAppWindow().Canvas(),
		),
		parent: parent,
	}
	return popUp
}
