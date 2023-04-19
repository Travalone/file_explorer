package main_panel

import (
	"file_explorer/common"
	"file_explorer/view/components/main_panel/extra_info_tab"
	"file_explorer/view/components/main_panel/file_create_tab"
	"file_explorer/view/components/main_panel/file_tab"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type FeTab interface {
	GetContainer() *fyne.Container
}

type MainPanel struct {
	*container.DocTabs

	feContext *store.FeContext
}

// newTabItem 新增tab上下文，创建TabItem
func (panel *MainPanel) newTabItem(feContext *store.FeContext, tabContext store.TabContext) *container.TabItem {
	var tab FeTab
	var icon fyne.Resource
	switch tabContext.GetTabType() {
	case common.TabTypeFile:
		icon = theme.FolderOpenIcon()
		tab, _ = file_tab.NewFileTab(tabContext.(*store.FileTabContext), panel.DocTabs.Refresh)
	case common.TabTypeExtraInfo:
		icon = theme.DocumentCreateIcon()
		tab = extra_info_tab.NewEditExtraInfoTab(feContext, tabContext.(*store.ExtraInfoTabContext))
	case common.TabTypeFileCreate:
		icon = theme.FolderNewIcon()
		tab = file_create_tab.NewFileCreateTab(feContext, tabContext.(*store.FileCreateTabContext))
	}
	tabItem := container.NewTabItemWithIcon(tabContext.GetTabLabel(), icon, tab.GetContainer())
	tabContext.SetTabItem(tabItem)
	return tabItem
}

func NewMainPanel(feContext *store.FeContext) *MainPanel {
	mainPanel := &MainPanel{
		DocTabs:   container.NewDocTabs(),
		feContext: feContext,
	}

	// 新增tabContext时同步新增tabItem
	mainPanel.feContext.OnAddTab = func(tabContext store.TabContext) {
		tabItem := mainPanel.newTabItem(feContext, tabContext)
		mainPanel.Append(tabItem)
		mainPanel.Select(tabItem)
	}

	// tabContext删除时关闭tab
	mainPanel.feContext.OnRemoveTab = func(tabItem *container.TabItem) {
		mainPanel.Remove(tabItem)
	}

	// 初始化时新增file tab
	homeFileTabContext := store.NewFileTabContext(feContext.GetFeConfig().Root, mainPanel.feContext)
	mainPanel.feContext.AddTab(homeFileTabContext)

	mainPanel.SetTabLocation(container.TabLocationTop)
	// 手动关闭标签页时删除对应context
	mainPanel.OnClosed = func(tabItem *container.TabItem) {
		feContext.Remove(tabItem)

	}
	return mainPanel
}
