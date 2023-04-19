package packed_binding

import "fyne.io/fyne/v2/data/binding"

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
