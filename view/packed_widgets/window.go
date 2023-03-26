package packed_widgets

import "fyne.io/fyne/v2"

func GetAppWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}
