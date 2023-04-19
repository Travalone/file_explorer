package store

import (
	"file_explorer/common"
	"file_explorer/common/utils"
	"fmt"
	"fyne.io/fyne/v2/container"
)

type FileCreateTabContext struct {
	tabItem *container.TabItem

	Path string
}

func (ctx *FileCreateTabContext) GetTabType() string {
	return common.TabTypeFileCreate
}

func (ctx *FileCreateTabContext) GetTabLabel() string {
	dirs := utils.PathSplit(ctx.Path)
	return fmt.Sprintf("%s: %s", ctx.GetTabType(), dirs[len(dirs)-1])
}

func (ctx *FileCreateTabContext) GetTabItem() *container.TabItem {
	return ctx.tabItem
}

func (ctx *FileCreateTabContext) SetTabItem(tabItem *container.TabItem) {
	ctx.tabItem = tabItem
}

func NewFileCreateTabContext(fileTabContext *FileTabContext) *FileCreateTabContext {
	ctx := &FileCreateTabContext{
		Path: fileTabContext.GetAbsolutePath(),
	}

	return ctx
}
