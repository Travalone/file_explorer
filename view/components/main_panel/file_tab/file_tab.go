package file_tab

import (
	"file_explorer/common"
	"file_explorer/view/packed_widgets/packed_binding"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type FileTab struct {
	*fyne.Container
	TabLabel  binding.String
	ToolBar   *FileTabToolbar
	PathBar   *PathBar
	FileList  *FileList
	StatusBar *StatusBar

	tabContext *store.FileTabContext
}

func (tab *FileTab) GetTabLabel() string {
	return common.TabTypeFile
}

func (tab *FileTab) GetContainer() *fyne.Container {
	return tab.Container
}

func (tab *FileTab) OpenDir(dirName string) {
	tab.tabContext.Move(dirName)
}

func NewFileTab(fileTabContext *store.FileTabContext, tabRefresh func()) (tab *FileTab, err error) {
	tab = &FileTab{
		tabContext: fileTabContext,
	}

	/* 创建子组件 */
	// 工具栏
	tab.ToolBar = NewToolbar(tab.tabContext, tab)
	// 地址栏
	tab.PathBar = NewPathBar(tab.tabContext)
	// 文件列表
	tab.FileList = NewFileList(tab.tabContext)
	// 状态栏
	tab.StatusBar = NewStatusBar(tab.tabContext)

	/* 监听context目录更新时，刷新组件 */
	// 目录改变时，刷新地址栏、文件列表、状态栏
	tab.tabContext.Dirs.BindListener(func() {
		// 刷新视图
		tab.PathBar.RefreshData()
		tab.FileList.RefreshData()
		tab.ToolBar.RefreshData()
		tab.StatusBar.RefreshStatus()
		tabRefresh() // 刷新tab标题，mainPanel整个传入太重了，也没其他比较合适的办法，就传个方法了
	})

	// 初始状态通过刷新当前目录触发上面的listeners (但不知道为什么，不跳转仅刷新时，目录listener会被调用两次)
	tab.tabContext.RefreshPath()

	// 文件列表选中状态更新时，刷新状态栏、工具栏
	packed_binding.NewListener(func() {
		tab.StatusBar.RefreshStatus()
		tab.ToolBar.refreshButton()
	}).BindData(tab.tabContext.CheckList)

	top := container.NewVBox(tab.ToolBar.Container, tab.PathBar.Table)
	tab.Container = container.NewBorder(top, tab.StatusBar.Scroll, nil, nil, tab.FileList.Table.Table)

	return tab, nil
}
