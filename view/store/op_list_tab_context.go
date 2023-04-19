package store

import (
	"file_explorer/common"
	"fyne.io/fyne/v2/container"
)

type OpListTabContext struct {
	tabItem   *container.TabItem
	FeContext *FeContext
}

func (ctx *OpListTabContext) GetTabType() string {
	return common.TabTypeOpList
}

func (ctx *OpListTabContext) GetTabLabel() string {
	return "操作列表"
}

func (ctx *OpListTabContext) GetTabItem() *container.TabItem {
	return ctx.tabItem
}

func (ctx *OpListTabContext) SetTabItem(tabItem *container.TabItem) {
	ctx.tabItem = tabItem
}
