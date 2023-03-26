package main

import (
	"file_explorer/common/logger"
	"file_explorer/resource"
	"file_explorer/service"
	"file_explorer/view"
	"file_explorer/view/packed_widgets"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"os"
)

func testWindow() {
	app := app.New()
	window := app.NewWindow("")
	workDir, _ := os.Getwd()
	label := packed_widgets.NewLabel(workDir)

	config := service.ReadConfig()
	label1 := packed_widgets.NewLabel(config.Root)
	window.SetContent(container.NewVBox(label, label1))
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
