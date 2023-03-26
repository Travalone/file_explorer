package packed_widgets

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/atomic"
	"time"
)

type CheckList struct {
	checkStatus   []binding.Bool
	batchChanging bool
	checkChanging *atomic.Bool

	OnCheckChange func([]int)
}

func (checkList *CheckList) Bind(index int, checkBox *widget.Check) {
	if index < 0 || index >= len(checkList.checkStatus) {
		return
	}
	checkBox.Bind(checkList.checkStatus[index])
}

func (checkList *CheckList) GetCheckIndexList() []int {
	checkIndexList := make([]int, 0, len(checkList.checkStatus))
	for i := range checkList.checkStatus {
		if checkList.checkStatus[i] == nil {
			continue
		}
		if b, _ := checkList.checkStatus[i].Get(); b {
			checkIndexList = append(checkIndexList, i)
		}
	}
	return checkIndexList
}

func (checkList *CheckList) Toggle(index int) {
	if index < 0 || index >= len(checkList.checkStatus) {
		return
	}
	b, _ := checkList.checkStatus[index].Get()
	checkList.checkStatus[index].Set(!b)
}

func (checkList *CheckList) BatchUpdate(checkMap map[int]bool) {
	count := len(checkMap)
	if count == 0 {
		return
	}
	checkList.batchChanging = true
	for index, b := range checkMap {
		if count == 1 {
			checkList.batchChanging = false
		}
		checkList.checkStatus[index].Set(b)
		count--
	}
}

func (checkList *CheckList) CheckAll(b bool) {
	checkMap := make(map[int]bool)
	for i := range checkList.checkStatus {
		checkMap[i] = b
	}
	checkList.BatchUpdate(checkMap)
}

func (checkList *CheckList) ToggleAll() {
	checkMap := make(map[int]bool)
	for i := range checkList.checkStatus {
		b, _ := checkList.checkStatus[i].Get()
		checkMap[i] = !b
	}
	checkList.BatchUpdate(checkMap)
}

func (checkList *CheckList) CheckRange() {
	start, end := 0, 0
	for i := range checkList.checkStatus {
		b, _ := checkList.checkStatus[i].Get()
		if b {
			start = i
			break
		}
	}
	for i := len(checkList.checkStatus) - 1; i > 0; i-- {
		b, _ := checkList.checkStatus[i].Get()
		if b {
			end = i
			break
		}
	}
	checkMap := make(map[int]bool)
	for i := start + 1; i < end; i++ {
		checkMap[i] = true
	}
	checkList.BatchUpdate(checkMap)
}

func (checkList *CheckList) addChangeCallback(index int) {
	NewListener(func() {
		if checkList.batchChanging {
			//logger.Warn("batch changing")
			return
		}
		go func() {
			// listeners异步执行，执行前获取锁
			if !checkList.checkChanging.CAS(false, true) {
				//logger.Warn("on changing")
				return
			}
			if checkList.OnCheckChange != nil {
				// 上级回调
				checkList.OnCheckChange(checkList.GetCheckIndexList())
				// 稍微延个时，避免无用地扫描多次，只要不大于连续点鼠标两次的时间间隔就行
				time.Sleep(time.Millisecond * 50)
			}
			// 执行结束后解锁
			checkList.checkChanging.CAS(true, false)
		}()
	}).BindData(checkList.checkStatus[index])
}

func NewCheckList(length int) *CheckList {
	checkList := &CheckList{
		checkStatus:   make([]binding.Bool, length),
		checkChanging: atomic.NewBool(false),
	}

	// 任一单选状态更新时回调
	for i := range checkList.checkStatus {
		checkList.checkStatus[i] = binding.NewBool()
		checkList.addChangeCallback(i)
	}

	return checkList
}
