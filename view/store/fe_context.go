package store

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// TabContext tab接口，每类Tab需要实现
type TabContext interface {
	GetTabType() string
	GetTabLabel() string

	GetTabItem() *container.TabItem
	SetTabItem(*container.TabItem)
}

// FeContext 全局上下文，单例
type FeContext struct {
	TabContexts []TabContext
	Window      fyne.Window

	OnAddTab    func(TabContext)
	OnRemoveTab func(*container.TabItem)

	FeConfig    *packed_widgets.BindingStruct
	TaskManager *service.FeTaskManager

	WorkDir string
}

// AddTab 添加Tab
func (ctx *FeContext) AddTab(tabContext TabContext) {
	ctx.TabContexts = append(ctx.TabContexts, tabContext)
	if ctx.OnAddTab != nil {
		ctx.OnAddTab(tabContext)
	}
}

// Remove 关闭Tab
func (ctx *FeContext) Remove(tabItem *container.TabItem) {
	newTabContexts := make([]TabContext, 0)
	for _, oldTabContext := range ctx.TabContexts {
		if oldTabContext.GetTabItem() != tabItem {
			newTabContexts = append(newTabContexts, oldTabContext)
		}
	}
	ctx.TabContexts = newTabContexts
	ctx.TaskManager.ClearTabTasks(tabItem)
	if ctx.OnRemoveTab != nil {
		ctx.OnRemoveTab(tabItem)
	}
}

// RefreshIfExists 刷新FileTab如果已打开
func (ctx *FeContext) RefreshIfExists(path string) {
	for _, tabContext := range ctx.TabContexts {
		if tabContext.GetTabType() == common.TAB_TYPE_FILE {
			fileTabContext := tabContext.(*FileTabContext)
			if path == fileTabContext.GetAbsolutePath() {
				fileTabContext.RefreshPath()
			}
		}
	}
}

func (ctx *FeContext) RefreshFeConfig() {
	config := ctx.GetFeConfig()
	ctx.FeConfig.Set(config)
	service.WriteConfig(config)
}

func (ctx *FeContext) GetFeConfig() *model.FileExplorerConfig {
	return ctx.FeConfig.Get().(*model.FileExplorerConfig)
}

func NewFeContext(feConfig *model.FileExplorerConfig) *FeContext {
	feContext := &FeContext{
		FeConfig:    packed_widgets.NewBindingStruct(feConfig),
		TaskManager: service.NewFeTaskManager(),
	}
	feContext.FeConfig.Set(feConfig)
	return feContext
}
