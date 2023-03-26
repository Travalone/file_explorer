package test

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	"file_explorer/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	path := "/Users/bytedance/Documents"

	fileInfos, _ := service.QueryPath(path)

	bindingData := make([]binding.Struct, len(fileInfos))
	for index, _ := range fileInfos {
		bindingData[index] = binding.BindStruct(fileInfos[index])
	}

	myApp := app.New()
	w := myApp.NewWindow("Two Way")

	list := widget.NewTable(
		func() (int, int) {
			return 1, 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			data, _ := bindingData[0].GetValue("Size")
			logger.Info("%+v", utils.Conv2Str(data))
			value, err := data.(binding.Int).Get()
			valueStr := ""
			if err != nil {
				valueStr, _ = data.(binding.String).Get()
			} else {
				valueStr = fmt.Sprintf("%d", value)
			}
			logger.Info("value=%v", valueStr)
			o.(*widget.Label).SetText(valueStr)
		},
	)

	w.SetContent(list)

	w.ShowAndRun()
}

func testListener(content string) (binding.String, binding.DataListener) {
	data := binding.NewString()
	data.Set(content)
	l := binding.NewDataListener(func() {
		str, _ := data.Get()
		logger.Info(str)
	})
	data.AddListener(l)
	return data, l
}

func TestBindData(t *testing.T) {

	d := binding.NewInt()
	d.Set(1)
	d.AddListener(binding.NewDataListener(func() {
		s, _ := d.Get()
		logger.Info("%v", s)
	}))
	time.Sleep(time.Millisecond * 1000)
	d.Set(1)

	for true {
	}
}
