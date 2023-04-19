package packed_binding

import (
	"fyne.io/fyne/v2/data/binding"
	"time"
)

// BindingList binding自带的类型有些问题导致没法直接用，set的list长度不变时就不能触发listener，这对目录刷新不太友好，所以这里自己实现一个
// PS：bindingInt set时值不变也不触发listener
type BindingList struct {
	marker binding.Int
	data   []interface{}
}

func (list *BindingList) AddListener(listener binding.DataListener) {
	panic("pls don't use this")
}
func (list *BindingList) RemoveListener(listener binding.DataListener) {
	panic("pls don't use this")
}

func (list *BindingList) BindListener(fun func()) {
	NewListener(fun).BindData(list.marker)
}

func (list *BindingList) Set(data []interface{}) {
	list.data = data
	list.marker.Set(int(time.Now().Unix()))
}
func (list *BindingList) Append(item interface{}) {
	list.data = append(list.data, item)
	list.marker.Set(int(time.Now().Unix()))
}

func (list *BindingList) Get() []interface{} {
	return list.data
}

func (list *BindingList) GetIndex(index int) interface{} {
	return list.data[index]
}

func (list *BindingList) Remove(index int) {
	list.data = append(list.data[:index], list.data[index+1:]...)
	list.marker.Set(int(time.Now().Unix()))
}

func (list *BindingList) Length() int {
	return len(list.data)
}

func NewBindingList(data []interface{}) *BindingList {
	return &BindingList{
		marker: binding.NewInt(),
		data:   data,
	}
}
