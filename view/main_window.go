package view

import (
	"file_explorer/service"
	"file_explorer/view/components/main_panel"
	"file_explorer/view/components/side_panel"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type FeWindow struct {
	fyne.Window
	App       fyne.App
	FeContext *store.FeContext
}

func NewMainWindow() *FeWindow {
	myApp := app.New()
	feWindow := &FeWindow{
		Window:    myApp.NewWindow("File explorer"),
		App:       myApp,
		FeContext: store.NewFeContext(service.ReadConfig()),
	}
	left := side_panel.NewSidePanel(feWindow.FeContext)
	right := main_panel.NewMainPanel(feWindow.FeContext)
	split := container.NewHSplit(left.Scroll, right.DocTabs)
	split.Offset = 0.25
	feWindow.Window.SetContent(container.NewGridWithColumns(1, split))
	feWindow.Window.Resize(fyne.NewSize(800, 800))
	return feWindow
}
