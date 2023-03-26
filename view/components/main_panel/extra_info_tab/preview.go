package extra_info_tab

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"sort"
	"strconv"
	"strings"
)

type ExtraInfoPreview struct {
	*packed_widgets.Table

	tabContext *store.ExtraInfoTabContext
}

// ResetCheckItems 撤销EditForm的改动，重置为初始状态
func (preview *ExtraInfoPreview) ResetCheckItems() {
	checkIndexList := preview.CheckList.GetCheckIndexList()
	for _, checkIndex := range checkIndexList {
		original := preview.tabContext.FileInfos[checkIndex].Original
		preview.tabContext.FileInfos[checkIndex].New = original.Clone()
	}
	preview.Refresh()
}

// ReloadCheckList 从context回写CheckList状态
func (preview *ExtraInfoPreview) ReloadCheckList() {
	checkList, _ := preview.tabContext.CheckList.Get()
	checkItemSet := make(map[*model.PreviewFileInfo]struct{})
	for _, checkItem := range checkList {
		checkItemSet[checkItem.(*model.PreviewFileInfo)] = struct{}{}
	}
	checkMap := make(map[int]bool)
	for i, fileInfo := range preview.tabContext.FileInfos {
		_, checkMap[i] = checkItemSet[fileInfo]
	}
	preview.CheckList.BatchUpdate(checkMap)
}

// EditForm score更新时回调
func (preview *ExtraInfoPreview) onInputScoreChange() {
	inputScore, _ := preview.tabContext.InputScore.Get()
	checkIndexList := preview.CheckList.GetCheckIndexList()
	for _, index := range checkIndexList {
		fileInfo := preview.tabContext.FileInfos[index].New
		newScore := inputScore
		if newScore == "*" {
			newScore = preview.tabContext.FileInfos[index].Original.GetScore()
		}
		if fileInfo.ExtraInfo == nil {
			fileInfo.ExtraInfo = &model.FileExtraInfo{}
		}
		if newScore == "" {
			fileInfo.ExtraInfo.Score = nil
		} else if s, err := strconv.Atoi(newScore); err == nil {
			fileInfo.ExtraInfo.Score = &s
		}
	}
	preview.Table.Refresh()
}

// EditForm tags更新时回调
func (preview *ExtraInfoPreview) onInputTagsChange() {
	inputTags, _ := preview.tabContext.InputTags.Get()
	inputTags = inputTags[1:]
	checkIndexList := preview.CheckList.GetCheckIndexList()
	for _, index := range checkIndexList {
		fileInfo := preview.tabContext.FileInfos[index].New
		newTags := make([]string, 0)
		for _, inputTag := range inputTags {
			if inputTag == "*" {
				oldTags := preview.tabContext.FileInfos[index].Original.GetTagList()
				newTags = utils.MergeLists(newTags, oldTags)
			} else {
				newTags = utils.MergeLists(newTags, []string{inputTag})
			}
		}
		if fileInfo.ExtraInfo == nil {
			fileInfo.ExtraInfo = &model.FileExtraInfo{}
		}
		sort.Strings(newTags)
		fileInfo.ExtraInfo.Tags = newTags
	}
	preview.Table.Refresh()
}

// EditForm note更新时回调
func (preview *ExtraInfoPreview) onInputNoteChange() {
	inputNote, _ := preview.tabContext.InputNote.Get()
	checkIndexList := preview.CheckList.GetCheckIndexList()
	for _, index := range checkIndexList {
		fileInfo := preview.tabContext.FileInfos[index].New
		newNote := inputNote
		if newNote == "*" {
			newNote = preview.tabContext.FileInfos[index].Original.GetNote()
		}
		if fileInfo.ExtraInfo == nil {
			fileInfo.ExtraInfo = &model.FileExtraInfo{}
		}
		fileInfo.ExtraInfo.Note = &newNote
	}
	preview.Table.Refresh()
}

func NewExtraInfoPreview(tabContext *store.ExtraInfoTabContext) *ExtraInfoPreview {
	preview := &ExtraInfoPreview{
		tabContext: tabContext,
	}

	preview.Table = packed_widgets.NewTable(
		[]*packed_widgets.TableColumn{
			{
				Header: "",
				Width:  30,
				OrderBy: func(desc bool, inited bool) {
					preview.tabContext.SortFileInfos(func(a, b *model.PreviewFileInfo) bool {
						if a.New.Type != b.New.Type {
							if desc {
								return a.New.Type > b.New.Type
							}
							return a.New.Type < b.New.Type
						}
						// 相同时按名称升序
						return strings.ToLower(a.New.Name) < strings.ToLower(b.New.Name)
					})
					if inited {
						preview.ReloadCheckList()
					}
				},
				GetIcon: func(data interface{}) fyne.Resource {
					if data.(*model.FileInfo).Type == common.FILE_TYPE_DIR {
						return theme.FolderIcon()
					}
					if data.(*model.FileInfo).Type == common.FILE_TYPE_DRIVER {
						return theme.ComputerIcon()
					}
					return theme.FileIcon()
				},
			},
			{
				Header:  "名称",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).Name },
				Width:   200,
				OrderBy: func(desc bool, inited bool) {
					preview.tabContext.SortFileInfos(func(a, b *model.PreviewFileInfo) bool {
						res := strings.ToLower(a.New.Name) < strings.ToLower(b.New.Name)
						if desc {
							return !res
						}
						return res
					})
					if inited {
						preview.ReloadCheckList()
					}
				},
			},
			{
				Header:  "评分",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetScore() },
				Width:   60,
				OrderBy: func(desc bool, inited bool) {
					preview.tabContext.SortFileInfos(func(a, b *model.PreviewFileInfo) bool {
						if a.New.GetScore() != b.New.GetScore() {
							if desc {
								return a.New.GetScore() > b.New.GetScore()
							}
							return a.New.GetScore() < b.New.GetScore()
						}
						// 相同时按名称升序
						return strings.ToLower(a.New.Name) < strings.ToLower(b.New.Name)
					})
					if inited {
						preview.ReloadCheckList()
					}
				},
			},
			{
				Header:  "标签",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetTagsStr() },
				Width:   100,
			},
			{
				Header:  "备注",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetNote() },
				Width:   300,
			},
		},
		func(index int) interface{} {
			return tabContext.FileInfos[index].New
		},
		func() int {
			return len(tabContext.FileInfos)
		},
	)

	// 初次刷新，加载CheckList
	preview.Table.RefreshData()

	// vList内check list更新时同步到tab context
	preview.Table.OnCheckChange = func(checkIndexList []int) {
		checkList := make([]interface{}, len(checkIndexList))
		for i, checkIndex := range checkIndexList {
			checkList[i] = preview.tabContext.FileInfos[checkIndex]
		}
		preview.tabContext.CheckList.Set(checkList)
		// 更新EditForm内聚合extra info值的显示
		preview.tabContext.RefreshInputExtraInfo()
	}

	// 单击文件所在行切换选中状态
	preview.Table.OnTapped = func(row int, col int) {
		if row > 0 {
			preview.Table.CheckList.Toggle(row - 1)
		}
	}

	// 预览列表extraInfo实时响应输入变化
	packed_widgets.NewListener(preview.onInputScoreChange).BindData(tabContext.InputScore)
	packed_widgets.NewListener(preview.onInputNoteChange).BindData(tabContext.InputNote)
	packed_widgets.NewListener(preview.onInputTagsChange).BindData(tabContext.InputTags)

	return preview
}
