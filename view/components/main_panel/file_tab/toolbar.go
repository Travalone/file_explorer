package file_tab

import (
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

	addFavoritesButton    *widget.Button
	deleteFavoritesButton *widget.Button
	tagButton             *widget.Button
	tagsPopUp             *TagsPopUp
	editExtraInfoButton   *widget.Button

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

	// 是否可修改ExtraInfo
	if toolbar.tabContext.CheckList.Length() > 0 {
		toolbar.editExtraInfoButton.Enable()
	} else {
		toolbar.editExtraInfoButton.Disable()
	}
}

// 创建标签过滤按钮
func (toolbar *FileTabToolbar) newTagButton() {
	toolbar.tagButton = widget.NewButtonWithIcon("标签", theme.DocumentPrintIcon(), nil)
	toolbar.tagsPopUp = NewTagsPopUp(toolbar.tabContext, toolbar.fileTab, toolbar.tagButton)
	toolbar.tagButton.OnTapped = func() {
		toolbar.tagsPopUp.Show()
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

	toolbar.newAddFavoritesButton()
	toolbar.newDeleteFavoritesButton()
	toolbar.newTagButton()
	toolbar.newEditExtraInfoButton()

	left := container.NewHBox(
		toolbar.addFavoritesButton,
		toolbar.deleteFavoritesButton,

		widget.NewSeparator(),

		toolbar.tagButton,
		toolbar.editExtraInfoButton,
	)

	right := container.NewHBox(
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

	toolbar.Container = container.NewBorder(nil, nil, left, right)

	return toolbar
}
