package extra_info_tab

import (
	"file_explorer/common"
	"file_explorer/view/packed_widgets/packed_binding"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ExtraInfoTab struct {
	*fyne.Container

	editForm  *ExtraInfoEditForm
	toolbar   *fyne.Container
	preview   *ExtraInfoPreview
	statusBar *StatusBar

	tabContext *store.ExtraInfoTabContext
}

func (tab *ExtraInfoTab) GetTabLabel() string {
	return common.TabTypeExtraInfo
}

func (tab *ExtraInfoTab) GetContainer() *fyne.Container {
	return tab.Container
}

func NewEditExtraInfoTab(feContext *store.FeContext, extraInfoTabContext *store.ExtraInfoTabContext) *ExtraInfoTab {
	tab := &ExtraInfoTab{
		tabContext: extraInfoTabContext,
	}

	// 创建子组件
	tab.preview = NewExtraInfoPreview(tab.tabContext)
	tab.editForm = NewExtraInfoEditForm(tab.tabContext, tab.preview)
	tab.toolbar = NewToolbar(tab.preview)
	tab.statusBar = NewStatusBar(extraInfoTabContext)

	// 改动提交成功后回调
	tab.editForm.OnSubmit = func() {
		// 关闭tab
		feContext.Remove(extraInfoTabContext.GetTabItem())
		// 如果目录打开则刷新
		feContext.RefreshIfExists(extraInfoTabContext.Path)
	}

	// preview内选中文件改变时回调
	packed_binding.NewListener(func() {
		// 选中项更新时刷新状态栏
		tab.statusBar.RefreshStatus()
		// EditForm 聚合值刷新
		tab.editForm.TagList.Refresh()
	}).BindData(tab.tabContext.CheckList)

	// 默认全选
	tab.preview.Table.CheckList.CheckAll(true)
	// 初始状态刷新
	tab.preview.ResetCheckItems()
	tab.tabContext.RefreshInputExtraInfo()

	tab.Container = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("1. \"*\"表示修改前内容; 2. 双击移除标签；3. 选中项可修改，提交时包含所有项"),
			tab.editForm.Form,
			tab.toolbar,
		), tab.statusBar.Label, nil, nil, tab.preview.Table.Table,
	)

	return tab
}
