package extra_info_tab

import (
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newLoadHistoryButton(preview *ExtraInfoPreview) *widget.Button {
	historyButton := widget.NewButtonWithIcon("还原", theme.HistoryIcon(), nil)
	history := container.NewVBox()
	popUp := packed_widgets.NewPopUp(historyButton, container.NewVScroll(history))
	popUp.Resize(fyne.NewSize(200, 300))

	metaFileInfos := service.GetMetaFiles(preview.tabContext.Path)
	for _, item := range metaFileInfos {
		metaFileInfo := item
		button := widget.NewButton(metaFileInfo.GetModifyTime(), nil)
		button.OnTapped = func() {
			// 读取指定metaFile
			metaPath := metaFileInfo.GetPath()
			metaData, err := service.ReadMetaFile(metaPath)
			if err != nil {
				logger.Error("newLoadHistoryButton ReadMetaFile failed, path=%v", metaPath)
				return
			}
			logger.Debug("load extra info history, time=%v, path=%v, metaData=%v", metaFileInfo.GetModifyTime(), metaPath, metaData)

			if metaData.FileExtraInfos == nil {
				metaData.FileExtraInfos = make(map[string]*model.FileExtraInfo)
			}

			// 加载extraInfo
			for _, fileInfo := range preview.tabContext.FileInfos {
				fileInfo.New.ExtraInfo = &model.FileExtraInfo{}
				if metaData.FileExtraInfos[fileInfo.New.Name] != nil {
					fileInfo.New.ExtraInfo = metaData.FileExtraInfos[fileInfo.New.Name]
				}
			}

			// 刷新视图
			preview.Refresh()
			preview.tabContext.RefreshInputExtraInfo()
			popUp.Hide()
		}
		history.Add(button)
	}

	historyButton.OnTapped = func() {
		popUp.Show()
	}
	return historyButton
}

func moveCheckItems(preview *ExtraInfoPreview, offset int) {
	checkIndexList := preview.CheckList.GetCheckIndexList()
	fileInfos := make([]interface{}, len(preview.tabContext.FileInfos))
	for i, fileInfo := range preview.tabContext.FileInfos {
		fileInfos[i] = fileInfo
	}
	utils.MoveListItems(fileInfos, checkIndexList, offset)
	preview.tabContext.FileInfos = model.Interfaces2PreviewFileInfos(fileInfos)
	preview.Table.Refresh()
	preview.ReloadCheckList()
}

func NewToolbar(preview *ExtraInfoPreview) *fyne.Container {
	left := container.NewHBox(
		widget.NewButtonWithIcon("上移", theme.MoveUpIcon(), func() {
			moveCheckItems(preview, -1)
		}),
		widget.NewButtonWithIcon("下移", theme.MoveDownIcon(), func() {
			moveCheckItems(preview, 1)
		}),
		widget.NewButtonWithIcon("置顶", theme.UploadIcon(), func() {
			moveCheckItems(preview, -len(preview.tabContext.FileInfos))
		}),
		widget.NewButtonWithIcon("置底", theme.DownloadIcon(), func() {
			moveCheckItems(preview, len(preview.tabContext.FileInfos))

		}),

		widget.NewSeparator(),

		newLoadHistoryButton(preview),
	)

	right := container.NewHBox(
		widget.NewButtonWithIcon("反选", theme.SearchReplaceIcon(), func() {
			if preview.CheckList != nil {
				preview.CheckList.ToggleAll()
			}
		}),
		widget.NewButtonWithIcon("区间选择", theme.ListIcon(), func() {
			if preview.CheckList != nil {
				preview.CheckList.CheckRange()
			}
		}),
	)

	toolbar := container.NewBorder(nil, nil, left, right)

	return toolbar
}
