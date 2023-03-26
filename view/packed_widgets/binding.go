package packed_widgets

import (
	"fyne.io/fyne/v2/data/binding"
	"time"
)

type Listener struct {
	binding.DataListener
	inited bool
}

func (listener *Listener) BindData(data binding.DataItem) {
	data.AddListener(listener.DataListener)
}

func NewListener(fun func()) *Listener {
	l := Listener{}
	// 所有binding data的listener在同一个后台线程中按添加顺序执行
	l.DataListener = binding.NewDataListener(func() {
		// 为binding data绑定listener时会执行一次，这里给欺骗掉，为了preview表默认全选能set上，不会被初始执行抢占锁
		if l.inited {
			fun()
		}
		l.inited = true
	})
	return &l
}

// BindingStringList binding自带的类型有些问题导致没法直接用，set的list长度不变时就不能触发listener，这对目录刷新不太友好，所以这里自己实现一个
// PS：bindingInt set时值不变也不触发listener
type BindingStringList struct {
	marker binding.Int
	data   []string
}

func (list *BindingStringList) AddListener(listener binding.DataListener) {
	panic("pls don't use this")
}
func (list *BindingStringList) RemoveListener(listener binding.DataListener) {
	panic("pls don't use this")
}

func (list *BindingStringList) BindListener(fun func()) {
	NewListener(fun).BindData(list.marker)
}

func (list *BindingStringList) Set(strList []string) {
	list.data = strList
	list.marker.Set(int(time.Now().Unix()))
}
func (list *BindingStringList) Append(str string) {
	list.data = append(list.data, str)
	list.marker.Set(int(time.Now().Unix()))
}

func (list *BindingStringList) Get() []string {
	return list.data
}
func (list *BindingStringList) Length() int {
	return len(list.data)
}

func NewBindingStringList(strList []string) *BindingStringList {
	return &BindingStringList{
		marker: binding.NewInt(),
		data:   strList,
	}
}

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
