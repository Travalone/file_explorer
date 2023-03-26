package service

import (
	"fyne.io/fyne/v2/container"
)

type FeTask interface {
	Start()
	Cancel()
}

type FeTaskManager struct {
	tabTaskMap  map[*container.TabItem][]FeTask // tab任务，随tabContext删除而取消
	globalTasks []FeTask                        // 全局任务
}

func (manager *FeTaskManager) SubmitGlobalTask(task FeTask) {
	manager.globalTasks = append(manager.globalTasks, task)
	task.Start()
}

func (manager *FeTaskManager) CancelGlobalTask(task FeTask) {
	for index, item := range manager.globalTasks {
		if item == task {
			task.Cancel()
			manager.globalTasks = append(manager.globalTasks[:index], manager.globalTasks[index+1:]...)
			return
		}
	}
}

func (manager *FeTaskManager) SubmitTabTask(tabItem *container.TabItem, task FeTask) {
	if manager.tabTaskMap[tabItem] == nil {
		manager.tabTaskMap[tabItem] = make([]FeTask, 0)
	}
	manager.tabTaskMap[tabItem] = append(manager.tabTaskMap[tabItem], task)
	task.Start()
}

func (manager *FeTaskManager) CancelTabTask(tabItem *container.TabItem, task FeTask) {
	if manager.tabTaskMap[tabItem] == nil {
		return
	}
	for index, item := range manager.tabTaskMap[tabItem] {
		if task == item {
			task.Cancel()
			tasks := manager.tabTaskMap[tabItem]
			tasks = append(tasks[:index], tasks[index+1:]...)
			return
		}
	}
}
func (manager *FeTaskManager) ClearTabTasks(tabItem *container.TabItem) {
	if manager.tabTaskMap[tabItem] == nil {
		return
	}
	for _, item := range manager.tabTaskMap[tabItem] {
		item.Cancel()
	}
	manager.tabTaskMap[tabItem] = nil
}

func NewFeTaskManager() *FeTaskManager {
	return &FeTaskManager{
		tabTaskMap:  make(map[*container.TabItem][]FeTask),
		globalTasks: make([]FeTask, 0),
	}
}
