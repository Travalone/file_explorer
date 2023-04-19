package main

import (
	"file_explorer/common/logger"
	"file_explorer/resource"
	"file_explorer/view"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func testWindow() {
	app := app.New()
	window := app.NewWindow("")

	entry := widget.NewEntry()
	entry.MultiLine = true

	window.SetContent(entry)
	window.Resize(fyne.NewSize(600, 800))
	window.ShowAndRun()
}

func main() {
	logger.SetLogLevel(logger.LevelDebug)
	//testWindow()

	mainWindow := view.NewMainWindow()
	mainWindow.App.Settings().SetTheme(&resource.FeTheme{})
	mainWindow.ShowAndRun()
}
