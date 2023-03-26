package packed_widgets

import "fyne.io/fyne/v2"

func NewNotify(content string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "File Explorer Notify",
		Content: content,
	})
}
