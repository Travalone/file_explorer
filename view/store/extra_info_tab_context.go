package store

import (
	"file_explorer/common"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"sort"
)

var (
	DefaultTagList = []string{""}
)

type ExtraInfoTabContext struct {
	tabItem *container.TabItem

	Path      string
	FileInfos []*model.PreviewFileInfo
	CheckList binding.UntypedList

	InputScore binding.String
	InputTags  binding.StringList
	InputNote  binding.String
	InputUrl   binding.String
}

func (ctx *ExtraInfoTabContext) GetTabType() string {
	return common.TabTypeExtraInfo
}

func (ctx *ExtraInfoTabContext) GetTabLabel() string {
	dirs := utils.PathSplit(ctx.Path)
	return fmt.Sprintf("%s: %s", ctx.GetTabType(), dirs[len(dirs)-1])
}

func (ctx *ExtraInfoTabContext) GetTabItem() *container.TabItem {
	return ctx.tabItem
}

func (ctx *ExtraInfoTabContext) SetTabItem(tabItem *container.TabItem) {
	ctx.tabItem = tabItem
}

// RefreshInputExtraInfo 根据preview内选中项目，刷新EditForm内显示的聚合数据
func (ctx *ExtraInfoTabContext) RefreshInputExtraInfo() {
	checkList, _ := ctx.CheckList.Get()
	var score, note, url *string = nil, nil, nil
	var tags []string = nil
	for _, checkItem := range checkList {
		fileInfo := checkItem.(*model.PreviewFileInfo)
		checkItemScore := fileInfo.New.GetScore()
		checkItemTags := fileInfo.New.GetTagList()
		checkItemNote := fileInfo.New.GetNote()
		checkItemUrl := fileInfo.New.GetUrl()

		// extraInfos单一值时正常显示，多个值时显示*
		if score == nil {
			score = &checkItemScore
		} else if *score != checkItemScore {
			score = utils.PtrStr("*")
		}

		if tags == nil {
			tags = checkItemTags
		} else if !utils.CheckListsEqual(tags, checkItemTags) {
			tags = []string{"*"}
		}

		if note == nil {
			note = &checkItemNote
		} else if *note != checkItemNote {
			note = utils.PtrStr("*")
		}

		if url == nil {
			url = &checkItemUrl
		} else if *url != checkItemUrl {
			url = utils.PtrStr("*")
		}
	}

	// 空值处理
	if score == nil {
		score = utils.PtrStr("")
	}
	tags = append(DefaultTagList, tags...)
	if note == nil {
		note = utils.PtrStr("")
	}
	if url == nil {
		url = utils.PtrStr("")
	}

	ctx.InputScore.Set(*score)
	ctx.InputTags.Set(tags)
	ctx.InputNote.Set(*note)
	ctx.InputUrl.Set(*url)
}

// SortFileInfos 文件列表排序
func (ctx *ExtraInfoTabContext) SortFileInfos(cmp func(a, b *model.PreviewFileInfo) bool) {
	interfaces := model.PreviewFileInfos2Interfaces(ctx.FileInfos)
	sort.Slice(interfaces, func(i, j int) bool {
		a := interfaces[i].(*model.PreviewFileInfo)
		b := interfaces[j].(*model.PreviewFileInfo)
		return cmp(a, b)
	})
	ctx.FileInfos = model.Interfaces2PreviewFileInfos(interfaces)
}

func NewExtraInfoTabContext(fileTabContext *FileTabContext) *ExtraInfoTabContext {
	if fileTabContext.CheckList.Length() == 0 {
		return nil
	}

	ctx := &ExtraInfoTabContext{
		Path:       fileTabContext.GetAbsolutePath(),
		FileInfos:  make([]*model.PreviewFileInfo, fileTabContext.CheckList.Length()),
		CheckList:  binding.NewUntypedList(),
		InputScore: binding.NewString(),
		InputTags:  binding.NewStringList(),
		InputNote:  binding.NewString(),
		InputUrl:   binding.NewString(),
	}

	// 删除最后一个item后tag编辑栏hList消失，但label还在，所以这里设置一个空item占位
	ctx.InputTags.Set(DefaultTagList)

	// 设置preview列表
	checkList, _ := fileTabContext.CheckList.Get()
	for i, checkItem := range checkList {
		ctx.FileInfos[i] = model.NewPreviewFileInfo(checkItem.(*model.FileInfo))
	}

	return ctx
}
