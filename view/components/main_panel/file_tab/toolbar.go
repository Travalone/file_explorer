package file_tab

import (
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FileTabToolbar struct {
	*fyne.Container

	createFileButton      *widget.Button // 新建文件按钮
	addFavoritesButton    *widget.Button // 添加收藏按钮
	deleteFavoritesButton *widget.Button // 取消收藏按钮
	tagButton             *widget.Button // tag过滤按钮
	tagsPopUp             *TagsPopUp     // tag过滤菜单
	urlButton             *widget.Button // url按钮
	editExtraInfoButton   *widget.Button // 编辑ExtraInfo按钮

	fileTab    *FileTab
	tabContext *store.FileTabContext
}

func (toolbar *FileTabToolbar) RefreshData() {
	toolbar.tagsPopUp.RefreshData()
	toolbar.refreshButton()
}

// 刷新按钮显示状态
func (toolbar *FileTabToolbar) refreshButton() {
	curPath := toolbar.tabContext.GetAbsolutePath()

	// 收藏/取消收藏 显示
	if toolbar.tabContext.FeContext.GetFeConfig().FindFavorites(curPath) >= 0 {
		toolbar.addFavoritesButton.Hide()
		toolbar.deleteFavoritesButton.Show()
	} else {
		toolbar.addFavoritesButton.Show()
		toolbar.deleteFavoritesButton.Hide()
	}

	// 编辑ExtraInfo、打开url按钮 显示
	if toolbar.tabContext.CheckList.Length() > 0 {
		toolbar.editExtraInfoButton.Enable()
		toolbar.urlButton.Enable()
	} else {
		toolbar.editExtraInfoButton.Disable()
		toolbar.urlButton.Disable()
	}
}

// 创建新建文件按钮
func (toolbar *FileTabToolbar) newCreateFileButton() {
	toolbar.createFileButton = widget.NewButtonWithIcon("新建", theme.FolderNewIcon(), func() {
		fileCreateTabContext := store.NewFileCreateTabContext(toolbar.tabContext)

		// 添加context，回调创建 tab
		toolbar.tabContext.FeContext.AddTab(fileCreateTabContext)
	})
}

// 创建标签过滤按钮
func (toolbar *FileTabToolbar) newTagButton() {
	toolbar.tagButton = widget.NewButtonWithIcon("标签", theme.DocumentPrintIcon(), nil)
	toolbar.tagsPopUp = NewTagsPopUp(toolbar.tabContext, toolbar.fileTab, toolbar.tagButton)
	toolbar.tagButton.OnTapped = func() {
		toolbar.tagsPopUp.Show()
	}
}

// 创建url按钮
func (toolbar *FileTabToolbar) newUrlButton() {
	toolbar.urlButton = widget.NewButtonWithIcon("打开url", theme.DocumentPrintIcon(), nil)
	toolbar.urlButton.OnTapped = func() {
		checkList, _ := toolbar.tabContext.CheckList.Get()
		urls := make([]string, len(checkList))
		for i, item := range checkList {
			fileInfo := item.(*model.FileInfo)
			urls[i] = fileInfo.GetUrl()
			urls[i] = utils.ReplaceWildcard(urls[i], map[string]string{"name": fileInfo.Name})
		}
		service.OpenUrlsWithDefaultWebExplorer(urls)
	}
}

// 创建添加收藏夹按钮
func (toolbar *FileTabToolbar) newAddFavoritesButton() {
	toolbar.addFavoritesButton = widget.NewButtonWithIcon("收藏", theme.StorageIcon(), func() {
		nameEntry := widget.NewEntry()
		dirs := toolbar.tabContext.GetDirs()
		nameEntry.Text = dirs[len(dirs)-1]
		dialog.ShowForm("收藏", "确认", "取消",
			[]*widget.FormItem{
				widget.NewFormItem("名称", nameEntry),
			},
			func(confirm bool) {
				if !confirm || len(nameEntry.Text) == 0 {
					return
				}
				curPath := toolbar.tabContext.GetAbsolutePath()
				toolbar.tabContext.FeContext.GetFeConfig().AddFavorites(curPath, nameEntry.Text)
				toolbar.tabContext.FeContext.RefreshFeConfig()
				toolbar.refreshButton()
			},
			packed_widgets.GetAppWindow())
	})
}

// 创建取消收藏按钮
func (toolbar *FileTabToolbar) newDeleteFavoritesButton() {
	toolbar.deleteFavoritesButton = widget.NewButtonWithIcon("取消收藏", theme.StorageIcon(), func() {
		curPath := toolbar.tabContext.GetAbsolutePath()
		toolbar.tabContext.FeContext.GetFeConfig().DeleteFavorites(curPath)
		toolbar.tabContext.FeContext.RefreshFeConfig()
		toolbar.refreshButton()
	})
}

// 创建修改ExtraInfo按钮
func (toolbar *FileTabToolbar) newEditExtraInfoButton() {
	toolbar.editExtraInfoButton = widget.NewButtonWithIcon("修改信息", theme.DocumentCreateIcon(), func() {
		// 利用file tab信息(path、选中文件)extra info context
		extraInfoTabContext := store.NewExtraInfoTabContext(toolbar.tabContext)
		if extraInfoTabContext == nil {
			packed_widgets.NewNotify("没有选中文件")
			return
		}
		// 添加extra info context，回调创建extra info tab
		toolbar.tabContext.FeContext.AddTab(extraInfoTabContext)
	})
}

func NewToolbar(tabContext *store.FileTabContext, fileTab *FileTab) *FileTabToolbar {
	toolbar := &FileTabToolbar{
		tabContext: tabContext,
		fileTab:    fileTab,
	}

	toolbar.newCreateFileButton()
	toolbar.newAddFavoritesButton()
	toolbar.newDeleteFavoritesButton()
	toolbar.newTagButton()
	toolbar.newUrlButton()
	toolbar.newEditExtraInfoButton()

	floor1Left := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ComputerIcon(), func() {
			service.OpenPathWithDefaultFileExplorer(toolbar.tabContext.GetAbsolutePath())
		}),
		//widget.NewSeparator(),
		toolbar.createFileButton,
		widget.NewSeparator(),
		toolbar.addFavoritesButton,
		toolbar.deleteFavoritesButton,
	)

	floor2Left := container.NewHBox(
		toolbar.urlButton,
		toolbar.tagButton,
		toolbar.editExtraInfoButton,
	)

	floor2Right := container.NewHBox(
		widget.NewButtonWithIcon("反选", theme.SearchReplaceIcon(), func() {
			if fileTab.FileList != nil && fileTab.FileList.CheckList != nil {
				fileTab.FileList.CheckList.ToggleAll()
			}
		}),
		widget.NewButtonWithIcon("区间选择", theme.ListIcon(), func() {
			if fileTab.FileList != nil && fileTab.FileList.CheckList != nil {
				fileTab.FileList.CheckList.CheckRange()
			}
		}),
	)

	floor1 := container.NewBorder(nil, nil, floor1Left, nil)
	floor2 := container.NewBorder(nil, nil, floor2Left, floor2Right)
	toolbar.Container = container.NewVBox(floor1, floor2)

	return toolbar
}
