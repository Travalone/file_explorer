package packed_widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
)

func NewNotify(content string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "File Explorer Notify",
		Content: content,
	})
}
func NewNotifyFmt(format string, args ...interface{}) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "File Explorer Notify",
		Content: fmt.Sprintf(format, args...),
	})
}
