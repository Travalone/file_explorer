package file_tab

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"strings"
)

type FileList struct {
	*packed_widgets.Table

	tabContext *store.FileTabContext
}

// ReloadCheckList 从context回写CheckList状态
func (fileList *FileList) ReloadCheckList() {
	if fileList.tabContext.CheckList.Length() == 0 {
		return
	}
	checkList, _ := fileList.tabContext.CheckList.Get()
	checkItemSet := make(map[*model.FileInfo]struct{})
	for _, checkItem := range checkList {
		checkItemSet[checkItem.(*model.FileInfo)] = struct{}{}
	}
	checkMap := make(map[int]bool)
	for i, fileInfo := range fileList.tabContext.FileInfos {
		_, checkMap[i] = checkItemSet[fileInfo]
	}
	fileList.CheckList.BatchUpdate(checkMap)
}

func NewFileList(tabContext *store.FileTabContext) *FileList {
	fileList := &FileList{tabContext: tabContext}

	fileList.Table = packed_widgets.NewTable(
		[]*packed_widgets.TableColumn{
			{
				Header: "",
				Width:  30,
				OrderBy: func(desc bool, inited bool) {
					fileList.tabContext.SortFileInfos(func(a, b *model.FileInfo) bool {
						if a.Type != b.Type {
							if desc {
								return a.Type > b.Type
							}
							return a.Type < b.Type
						}
						// 相同时按名称升序
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})
					if inited {
						fileList.ReloadCheckList()
					}
				},
				GetIcon: func(data interface{}) fyne.Resource {
					return data.(*model.FileInfo).GetIcon()
				},
			},
			{
				Header:  "名称",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).Name },
				Width:   200,
				OrderBy: func(desc bool, inited bool) {
					fileList.tabContext.SortFileInfos(func(a, b *model.FileInfo) bool {
						// 目录优先
						if a.Type != b.Type {
							return a.Type < b.Type
						}
						res := strings.ToLower(a.Name) < strings.ToLower(b.Name)
						if desc {
							return !res
						}
						return res
					})
					if inited {
						// 未初始化时还没有创建checklist
						// 已经初始化时需要从tabContext reload checklist
						fileList.ReloadCheckList()
					}
				},
				DefaultOrderBy: true,
			},

			{
				Header:  "大小",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetSize() },
				Width:   80,
				OrderBy: func(desc bool, inited bool) {
					fileList.tabContext.SortFileInfos(func(a, b *model.FileInfo) bool {
						// 目录优先
						if a.Type != b.Type {
							return a.Type < b.Type
						}

						if a.Type == common.FileTypeRegular && a.Size != b.Size {
							if desc {
								return a.Size > b.Size
							}
							return a.Size < b.Size
						}
						// 目录、相同大小文件时按名称升序
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})
					if inited {
						fileList.ReloadCheckList()
					}
				},
			},
			{
				Header:  "修改时间",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetModifyTime() },
				Width:   150,
				OrderBy: func(desc bool, inited bool) {
					fileList.tabContext.SortFileInfos(func(a, b *model.FileInfo) bool {
						// 目录优先
						if a.Type != b.Type {
							return a.Type < b.Type
						}
						if !a.ModifyTime.Equal(b.ModifyTime) {
							if desc {
								return a.ModifyTime.Unix() > b.ModifyTime.Unix()
							}
							return a.ModifyTime.Unix() < b.ModifyTime.Unix()
						}
						// 相同时按名称升序
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})
					if inited {
						fileList.ReloadCheckList()
					}
				},
			},
			{
				Header:  "评分",
				GetText: func(data interface{}) string { return data.(*model.FileInfo).GetScore() },
				Width:   60,
				OrderBy: func(desc bool, inited bool) {
					fileList.tabContext.SortFileInfos(func(a, b *model.FileInfo) bool {
						// 目录优先
						if a.Type != b.Type {
							return a.Type < b.Type
						}
						if a.GetScore() != b.GetScore() {
							if desc {
								return a.GetScore() > b.GetScore()
							}
							return a.GetScore() < b.GetScore()
						}
						// 相同时按名称升序
						return strings.ToLower(a.Name) < strings.ToLower(b.Name)
					})
					if inited {
						fileList.ReloadCheckList()
					}
				},
			},
			{
				Header:             "标签",
				GetText:            func(data interface{}) string { return data.(*model.FileInfo).GetTagsStr() },
				Width:              100,
				DoubleTappedExpand: true,
			},
			{
				Header:             "备注",
				GetText:            func(data interface{}) string { return data.(*model.FileInfo).GetNote() },
				Width:              300,
				DoubleTappedExpand: true,
			},
		},
		func(index int) interface{} {
			return tabContext.FileInfos[index]
		},
		func() int {
			return len(tabContext.FileInfos)
		},
	)

	// vList内check list更新时同步到tab context内
	fileList.Table.OnCheckChange = func(checkIndexList []int) {
		checkList := make([]interface{}, len(checkIndexList))
		for i, checkIndex := range checkIndexList {
			checkList[i] = fileList.tabContext.FileInfos[checkIndex]
		}
		fileList.tabContext.CheckList.Set(checkList)
	}

	// 单击文件所在行切换选中状态
	fileList.Table.OnTapped = func(row int, col int) {
		if row > 0 {
			fileList.Table.CheckList.Toggle(row - 1)
		}
	}

	// 双击文件所在行进入目录
	fileList.Table.OnDoubleTapped = func(row int, col int) {
		if row > 0 {
			fileInfo := fileList.tabContext.FileInfos[row-1]
			if fileInfo.IsDir() {
				dir := fileInfo.Name
				fileList.tabContext.Move(dir)
			}
		}
	}

	return fileList
}
