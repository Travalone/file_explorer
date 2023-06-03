package main

import (
	"file_explorer/common/logger"
	"file_explorer/resource"
	"file_explorer/view"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func testWindow() {
	app := app.New()
	window := app.NewWindow("")
	//bindStr := binding.NewString()
	//entry.MultiLine = true
	////entry.Bind(bindStr)
	//entry.OnChanged = func(s string) {
	//	logger.Debug("%s", s)
	//}
	//bindStr.AddListener(binding.NewDataListener(func() {
	//	str, _ := bindStr.Get()
	//	logger.Debug("%s", str)
	//}))
	window.SetContent(
		container.NewHSplit(container.NewWithoutLayout(widget.NewEntry()), widget.NewLabel("2")),
		)

	window.Resize(fyne.NewSize(600, 800))
	app.Settings().SetTheme(&resource.FeTheme{})
	window.ShowAndRun()
}

func main() {
	logger.SetLogLevel(logger.LevelDebug)
	//testWindow()
	//
	mainWindow := view.NewMainWindow()
	mainWindow.App.Settings().SetTheme(&resource.FeTheme{})
	mainWindow.ShowAndRun()

}
