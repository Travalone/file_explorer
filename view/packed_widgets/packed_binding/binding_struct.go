package packed_binding

import (
	"fyne.io/fyne/v2/data/binding"
	"time"
)

type BindingStruct struct {
	marker binding.Int
	data   interface{}
}

func (list *BindingStruct) AddListener(listener binding.DataListener) {
	panic("pls don't use this")
}
func (list *BindingStruct) RemoveListener(listener binding.DataListener) {
	panic("pls don't use this")
}

func (list *BindingStruct) BindListener(fun func()) {
	NewListener(fun).BindData(list.marker)
}

func (list *BindingStruct) Set(data interface{}) {
	list.data = data
	list.marker.Set(int(time.Now().Unix()))
}

func (list *BindingStruct) Get() interface{} {
	return list.data
}
func NewBindingStruct(data interface{}) *BindingStruct {
	return &BindingStruct{
		marker: binding.NewInt(),
		data:   data,
	}
}
