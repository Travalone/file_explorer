package test

import (
	"file_explorer/common/logger"
	"fyne.io/fyne/v2/data/binding"
	"testing"
	"time"
)

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
