package store

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets/packed_binding"
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"sort"
)

type FileTabContext struct {
	tabItem   *container.TabItem
	Dirs      *packed_binding.BindingList
	FileInfos []*model.FileInfo
	CheckList binding.UntypedList
	FeContext *FeContext
}

func (ctx *FileTabContext) GetTabType() string {
	return common.TabTypeFile
}

func (ctx *FileTabContext) GetTabLabel() string {
	dirs := ctx.GetDirs()
	return fmt.Sprintf("%s: %s", ctx.GetTabType(), dirs[len(dirs)-1])
}

func (ctx *FileTabContext) GetTabItem() *container.TabItem {
	return ctx.tabItem
}

func (ctx *FileTabContext) SetTabItem(tabItem *container.TabItem) {
	ctx.tabItem = tabItem
}

func (ctx *FileTabContext) GetDirs() []string {
	dirs := ctx.Dirs.Get()
	return utils.Interfaces2Strings(dirs)
}

// GetAbsolutePath 获取当前绝对路径
func (ctx *FileTabContext) GetAbsolutePath() string {
	return utils.PathJoin(ctx.GetDirs()...)
}

// GetTagList 根据文件列表获取所有tag
func (ctx *FileTabContext) GetTagList() []string {
	tagList := make([]string, 0)
	for _, fileInfo := range ctx.FileInfos {
		tagList = utils.MergeLists(tagList, fileInfo.GetTagList())
	}
	interfaces := utils.ToInterfaces(tagList)
	sort.Slice(interfaces, func(i, j int) bool {
		a, b := interfaces[i].(string), interfaces[j].(string)
		if len(a) == len(b) {
			return a < b
		}
		return len(a) < len(b)
	})
	return tagList
}

// FilterTags 当前目录文件列表按tags过滤
func (ctx *FileTabContext) FilterTags(filterTags []string, andOr bool) error {
	// 保存选中文件名称 (刷新后创建新对象，不能整个保存)
	checkList, _ := ctx.CheckList.Get()
	checkNameSet := make(map[string]struct{})
	for _, checkItem := range checkList {
		checkNameSet[checkItem.(*model.FileInfo).Name] = struct{}{}
	}

	// 获取当前目录文件列表
	originPath := ctx.GetAbsolutePath()
	dirs, fileInfos := ctx.queryPath(ctx.GetDirs())
	if utils.PathJoin(dirs...) != originPath {
		// 当前目录查询失败，通过刷新返回上级
		ctx.RefreshPath()
		return fmt.Errorf("FilterTags failed, access path err, path=%v", originPath)
	}

	// 过滤
	newCheckList := make([]interface{}, 0, len(checkNameSet))
	filterFileInfos := make([]*model.FileInfo, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		hit := false
		if len(filterTags) == 0 {
			// 没有tag时全部通过
			hit = true
		} else if andOr {
			// 命中所有tag时通过
			for _, filterTag := range filterTags {
				if hit = utils.ListContainsItem(fileInfo.GetTagList(), filterTag); !hit {
					break
				}
			}
		} else {
			// 命中一个tag即通过
			for _, filterTag := range filterTags {
				if hit = utils.ListContainsItem(fileInfo.GetTagList(), filterTag); hit {
					break
				}
			}
		}
		if hit {
			filterFileInfos = append(filterFileInfos, fileInfo)
			// 保存选中文件
			if _, ok := checkNameSet[fileInfo.Name]; ok {
				newCheckList = append(newCheckList, fileInfo)
			}
		}
	}
	ctx.CheckList.Set(newCheckList)

	// 填充extra info
	service.FillFileExtraInfo(originPath, fileInfos)

	ctx.FileInfos = filterFileInfos
	return nil
}

// 查询目录下子文件，查询报错时返回上级再查，返回最终查到的目录和文件列表
func (ctx *FileTabContext) queryPath(dirs []string) ([]string, []*model.FileInfo) {
	for true {
		path := utils.PathJoin(dirs...)
		fileInfos, err := service.QueryFiles(path)
		if err == nil {
			return utils.PathSplit(path), fileInfos
		}
		dirs = dirs[:len(dirs)-1]
	}
	return nil, nil
}

// SortFileInfos 文件列表排序
func (ctx *FileTabContext) SortFileInfos(cmp func(a, b *model.FileInfo) bool) {
	interfaces := model.FileInfos2Interfaces(ctx.FileInfos)
	sort.Slice(interfaces, func(i, j int) bool {
		a := interfaces[i].(*model.FileInfo)
		b := interfaces[j].(*model.FileInfo)
		return cmp(a, b)
	})
	ctx.FileInfos = model.Interfaces2FileInfos(interfaces)
}

// Move 在当前目录打开子目录
func (ctx *FileTabContext) Move(subDir string) {
	dirs := append(ctx.GetDirs(), subDir)
	ctx.setDirs(dirs)
}

// RefreshPath 刷新当前目录
func (ctx *FileTabContext) RefreshPath() {
	ctx.setDirs(ctx.GetDirs())
}

// Back 点击PathBar后退到上级目录
func (ctx *FileTabContext) Back(index int) {
	ctx.setDirs(ctx.GetDirs()[:index+1])
}

// 指定目录，查询文件列表
func (ctx *FileTabContext) setDirs(dirs []string) {
	dirs, fileInfos := ctx.queryPath(dirs)
	service.FillFileExtraInfo(utils.PathJoin(dirs...), fileInfos)
	ctx.Dirs.Set(utils.Strings2Interfaces(dirs))
	ctx.FileInfos = fileInfos
	ctx.CheckList.Set(make([]interface{}, 0))
	if ctx.tabItem != nil {
		ctx.tabItem.Text = ctx.GetTabLabel()
	}
}

func NewFileTabContext(path string, feContext *FeContext) *FileTabContext {
	dirs := utils.PathSplit(path)
	tabContext := &FileTabContext{
		Dirs:      packed_binding.NewBindingList(utils.Strings2Interfaces(dirs)),
		CheckList: binding.NewUntypedList(),
		FeContext: feContext,
	}

	return tabContext
}
