package file_tab

import (
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	TagsAnd = "and"
	TagsOr  = "or"
)

type TagsPopUp struct {
	*packed_widgets.PopUp
	andOr     *widget.RadioGroup
	fileTab   *FileTab
	checkList *packed_widgets.CheckList

	tabContext *store.FileTabContext
}

func (popUp *TagsPopUp) refreshFileList(tagList []string) {
	checkIndexList := popUp.checkList.GetCheckIndexList()
	filterTags := make([]string, len(checkIndexList))
	for i, checkIndex := range checkIndexList {
		filterTags[i] = tagList[checkIndex]
	}
	// 更新数据
	err := popUp.tabContext.FilterTags(filterTags, popUp.andOr.Selected == TagsAnd)
	if err != nil {
		return
	}
	// 刷新视图
	popUp.fileTab.FileList.RefreshData()
	popUp.fileTab.FileList.ReloadCheckList()
}

func (popUp *TagsPopUp) RefreshData() {
	// 获取当前文件列表tag list
	tagList := popUp.tabContext.GetTagList()
	popUp.checkList = packed_widgets.NewCheckList(len(tagList))

	// 构造tagList
	tagCheckBoxes := container.NewVBox()
	for i, tag := range tagList {
		checkBox := widget.NewCheck(tag, nil)
		popUp.checkList.Bind(i, checkBox)
		tagCheckBoxes.Add(checkBox)
	}
	popUp.Content.(*fyne.Container).Objects[0].(*container.Scroll).Content = tagCheckBoxes

	// 过滤选项变化时通过刷新文件列表
	popUp.checkList.OnCheckChange = func(checkIndexList []int) {
		popUp.refreshFileList(tagList)
	}
	popUp.andOr.OnChanged = func(s string) {
		popUp.refreshFileList(tagList)
	}
}

func (popUp *TagsPopUp) GetCheckTags() ([]string, bool) {
	if popUp.checkList == nil {
		return nil, false
	}
	tagList := popUp.tabContext.GetTagList()
	checkIndexList := popUp.checkList.GetCheckIndexList()

	checkTags := make([]string, len(checkIndexList))
	for i, checkIndex := range checkIndexList {
		checkTags[i] = tagList[checkIndex]
	}

	return checkTags, popUp.andOr.Selected == TagsAnd
}

func NewTagsPopUp(tabContext *store.FileTabContext, fileTab *FileTab, parent fyne.CanvasObject) *TagsPopUp {
	tagsPopUp := &TagsPopUp{
		tabContext: tabContext,
		fileTab:    fileTab,
		andOr:      widget.NewRadioGroup([]string{TagsOr, TagsAnd}, nil),
	}
	tagsPopUp.andOr.Horizontal = true
	tagsPopUp.andOr.Selected = TagsOr
	tagsPopUp.PopUp = packed_widgets.NewPopUp(parent,
		container.NewBorder(
			tagsPopUp.andOr, nil, nil, nil,
			container.NewVScroll(container.NewVBox()),
		),
	)
	tagsPopUp.Resize(fyne.NewSize(200, 400))

	return tagsPopUp
}
