package file_tab

import (
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	*container.Scroll
	Label      *packed_widgets.Label
	WalkButton *widget.Button
	StopButton *widget.Button

	tabContext *store.FileTabContext
}

func (bar *StatusBar) StopWalker() {
	if bar.StopButton.OnTapped != nil {
		// stop按钮点击事件回调不为空时，说明存在walker正在运行，强制停止、隐藏按钮
		bar.StopButton.OnTapped()
		bar.StopButton.OnTapped = nil
		bar.StopButton.Hide()
	}
}

func (bar *StatusBar) queryStatus() {
	totalRegularCount, totalDirCount, totalSize := 0, 0, int64(0)
	checkList, _ := bar.tabContext.CheckList.Get()
	fileInfos := model.Interfaces2FileInfos(checkList)
	pathWalker := service.NewPathWalker(
		500,
		fileInfos,
		func(process string) {
			// 扫描进度更新回调：更新label
			bar.Label.SetText(process)
		},
		func(regularCount int, dirCount int, size int64, ok bool, err error) {
			if err != nil {
				// 异常结束，显示错误信息、隐藏stop按钮
				bar.Label.SetText(err.Error())
				bar.StopButton.Hide()
				return
			}
			if !ok {
				// 手动结束，还原选中状态、隐藏stop按钮、清除点击事件
				bar.RefreshStatus()
				bar.StopButton.Hide()
				bar.StopButton.OnTapped = nil
				return
			}
			// 成功结束，累计占用
			totalRegularCount += regularCount
			totalDirCount += dirCount
			totalSize += size
			// 更新label
			bar.Label.SetText(fmt.Sprintf("选中%d项, 子文件%d项, 子目录%d项, 占用%s",
				1, totalRegularCount, totalDirCount, utils.ConvSize(totalSize)))
			// 隐藏stop按钮、清除点击事件
			bar.StopButton.Hide()
			bar.StopButton.OnTapped = nil
		},
	)

	// 设置stop按钮点击事件：停止遍历
	bar.StopButton.OnTapped = func() {
		bar.tabContext.FeContext.TaskManager.CancelTabTask(bar.tabContext.GetTabItem(), pathWalker)
	}
	// 开始遍历
	bar.tabContext.FeContext.TaskManager.SubmitTabTask(bar.tabContext.GetTabItem(), pathWalker)
	pathWalker.Start()
	// 显示stop按钮
	bar.StopButton.Show()
}

func (bar *StatusBar) RefreshStatus() {
	bar.WalkButton.Hide()
	bar.StopWalker()

	text := fmt.Sprintf("%d项", len(bar.tabContext.FileInfos))
	checkList, _ := bar.tabContext.CheckList.Get()
	if len(checkList) == 0 {
		bar.Label.SetText(text)
		return
	}

	text += fmt.Sprintf(", %d项被选中", len(checkList))

	// 获取项目是否全为文件
	allRegular := true
	totalSize := int64(0)
	for _, item := range checkList {
		fileInfo := item.(*model.FileInfo)
		if fileInfo.IsDir() {
			allRegular = false
			break
		}
		totalSize += fileInfo.Size
	}

	if allRegular {
		// 选中项全为文件时显示占用
		text += fmt.Sprintf(", 占用%s", utils.ConvSize(totalSize))
		bar.Label.SetText(text)
		bar.WalkButton.Hide()
		return
	}

	// 选中项存在目录，显示计算按钮
	bar.Label.SetText(text)
	bar.WalkButton.Show()
}

func NewStatusBar(tabContext *store.FileTabContext) *StatusBar {
	statusBar := &StatusBar{
		Label:      packed_widgets.NewLabel(""),
		WalkButton: widget.NewButton("计算空间", nil),
		StopButton: widget.NewButtonWithIcon("", theme.CancelIcon(), nil),

		tabContext: tabContext,
	}

	// walk按钮点击事件：切换按钮、查询状态
	statusBar.WalkButton.OnTapped = func() {
		statusBar.WalkButton.Hide()
		statusBar.StopButton.Show()
		statusBar.queryStatus()
	}
	// 初始状态下隐藏walk、stop按钮
	statusBar.WalkButton.Hide()
	statusBar.StopButton.Hide()
	statusBar.Scroll = container.NewHScroll(container.NewHBox(statusBar.StopButton, statusBar.Label, statusBar.WalkButton))
	return statusBar
}
