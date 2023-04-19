package file_create_tab

import (
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/packed_widgets/packed_binding"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

type FileCreateTab struct {
	*fyne.Container

	createList   *widget.List
	submitButton *widget.Button

	createFileInfos *packed_binding.BindingList

	feContext  *store.FeContext
	tabContext *store.FileCreateTabContext
}

func (tab *FileCreateTab) GetTabLabel() string {
	return common.TabTypeFileCreate
}

func (tab *FileCreateTab) GetContainer() *fyne.Container {
	return tab.Container
}

func (tab *FileCreateTab) RefreshData() {
	tab.createList.Refresh()
}

func transformType(str string) byte {
	if str == "目录" {
		return common.FileTypeDir
	}
	return common.FileTypeRegular
}

func (tab *FileCreateTab) createListUpdateItem(index widget.ListItemID, obj fyne.CanvasObject) {
	fileInfo := tab.createFileInfos.GetIndex(index).(*model.FileInfo)
	deleteButton, typeSelectEntry, nameEntry, scoreInputBox, tagsBox, urlBox, noteBox := parseCreateItem(obj)

	if deleteButton.OnTapped == nil {
		deleteButton.OnTapped = func() {
			tab.createFileInfos.Remove(index)
			tab.RefreshData()
		}
	}

	if typeSelectEntry.OnChanged == nil {
		typeSelectEntry.OnChanged = func(itemType string) {
			fileInfo.Type = transformType(itemType)
		}
	}

	if nameEntry.Validator == nil {
		nameEntry.Validator = func(name string) error {
			name = strings.TrimSpace(name)
			path := utils.PathJoin(tab.tabContext.Path, name)
			if utils.PathExists(path) {
				return common.NewError("file already exists")
			}
			return nil
		}
	}
	if nameEntry.OnChanged == nil {
		nameEntry.OnChanged = func(name string) {
			name = strings.TrimSpace(name)
			fileInfo.Name = name
		}
	}

	if scoreInputBox.OnChanged == nil {
		scoreInputBox.OnChanged = func(scoreStr string) {
			scoreStr = strings.TrimSpace(scoreStr)
			score, err := strconv.Atoi(scoreStr)
			if err == nil {
				fileInfo.SetScore(&score)
			} else if len(scoreStr) == 0 {
				fileInfo.SetScore(nil)
			}
		}
	}

	if tagsBox.OnChanged == nil {
		tagsBox.OnChanged = func(tagsStr string) {
			rawTags := strings.Split(tagsStr, ",")
			tags := make([]string, 0, len(rawTags))
			for _, rawTag := range rawTags {
				rawTag = strings.TrimSpace(rawTag)
				if len(rawTag) > 0 {
					tags = append(tags, rawTag)
				}
			}
			fileInfo.SetTags(tags)
		}
	}

	if urlBox.OnChanged == nil {
		urlBox.OnChanged = func(url string) {
			url = strings.TrimSpace(url)
			fileInfo.SetUrl(&url)
		}
	}

	if noteBox.OnChanged == nil {
		noteBox.OnChanged = func(note string) {
			note = strings.TrimSpace(note)
			fileInfo.SetNote(&note)
		}
	}
}

func (tab *FileCreateTab) submitCreate() error {
	hasExtraInfo := false
	for _, item := range tab.createFileInfos.Get() {
		fileInfo := item.(*model.FileInfo)
		if len(fileInfo.Name) == 0 {
			continue
		}

		fileInfo.FormatExtraInfo()
		if fileInfo.ExtraInfo != nil {
			hasExtraInfo = true
		}

		err := service.CreateFile(fileInfo)
		if err != nil {
			logger.Error("submitCreate createFile failed, err=%v", err)
			return err
		}
	}

	if !hasExtraInfo {
		logger.Debug("submitCreate no extra info")
		return nil
	}

	// 更新extraInfo
	extraInfo, _ := service.ReadLatestExtraInfo(tab.tabContext.Path)
	if extraInfo == nil {
		extraInfo = &model.FileExtraInfoMap{}
	}

	for _, item := range tab.createFileInfos.Get() {
		extraInfo.SetFileExtraInfo(item.(*model.FileInfo))
	}

	err := service.WriteExtraInfoFile(tab.tabContext.Path, extraInfo)
	if err != nil {
		logger.Warn("submitCreate WriteExtraInfoFile failed, err=%v", err)
	}

	return nil
}

func NewFileCreateTab(feContext *store.FeContext, tabContext *store.FileCreateTabContext) *FileCreateTab {
	tab := &FileCreateTab{
		tabContext:      tabContext,
		feContext:       feContext,
		createFileInfos: packed_binding.NewBindingList(nil),
	}

	// 初始状态添加一个item占位
	tab.createFileInfos.Append(&model.FileInfo{
		Dir:  tab.tabContext.Path,
		Type: common.FileTypeRegular,
	})

	// 添加项目按钮
	addItemButton := widget.NewButtonWithIcon("添加", theme.ContentAddIcon(), func() {
		tab.createFileInfos.Append(&model.FileInfo{
			Dir:  tab.tabContext.Path,
			Type: common.FileTypeRegular,
		})
		tab.RefreshData()
	})

	// 创建submit按钮
	tab.submitButton = widget.NewButtonWithIcon("提交", theme.ConfirmIcon(), func() {
		err := tab.submitCreate()
		if err != nil {
			packed_widgets.NewNotify(err.Error())
			return
		}
		tab.feContext.Remove(tab.tabContext.GetTabItem())
		tab.feContext.RefreshIfExists(tab.tabContext.Path)
	})

	// 工具栏
	toolbar := container.NewBorder(nil, nil,
		addItemButton,
		tab.submitButton,
	)

	// 创建项目列表
	tab.createList = widget.NewList(
		func() int {
			return tab.createFileInfos.Length()
		},
		newCreateItem,
		tab.createListUpdateItem,
	)

	tab.Container = container.NewBorder(toolbar, nil, nil, nil, tab.createList)

	return tab
}
