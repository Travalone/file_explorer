package store

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets/packed_binding"
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

	FeConfig    *packed_binding.BindingStruct
	TaskManager *service.FeTaskManager
	OpList      *OpListContext

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
		if tabContext.GetTabType() == common.TabTypeFile {
			fileTabContext := tabContext.(*FileTabContext)
			if path == fileTabContext.GetAbsolutePath() {
				fileTabContext.RefreshPath()
			}
		}
	}
}

func (ctx *FeContext) ClearInvalidFavorites() {
	feConfig := ctx.GetFeConfig()
	validFavorites := make([]*model.Favorites, 0, len(feConfig.Favorites))
	for _, favorites := range ctx.GetFeConfig().Favorites {
		if utils.PathExists(favorites.Path) {
			validFavorites = append(validFavorites, favorites)
		}
	}
	// 没有失效目录
	if len(validFavorites) == len(feConfig.Favorites) {
		return
	}

	feConfig.Favorites = validFavorites
	ctx.RefreshFeConfig()
}

// RefreshFeConfig 刷新fe配置，写入文件
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
		FeConfig:    packed_binding.NewBindingStruct(feConfig),
		TaskManager: service.NewFeTaskManager(),
	}
	feContext.FeConfig.Set(feConfig)
	return feContext
}
