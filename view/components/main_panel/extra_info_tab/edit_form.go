package extra_info_tab

import (
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ExtraInfoEditForm struct {
	*widget.Form

	NoteBox *widget.Entry
	TagList *packed_widgets.HList

	OnSubmit   func()
	tabContext *store.ExtraInfoTabContext
}

// 评分输入栏
func (form *ExtraInfoEditForm) newScoreEditBox() fyne.CanvasObject {
	scoreBox := widget.NewEntryWithData(form.tabContext.InputScore)
	scoreBox.Validator = func(score string) error {
		if score == "*" || score == "" {
			return nil
		}
		if _, e := strconv.Atoi(score); e == nil {
			return nil
		}
		return common.NewError("Input score is not number")
	}
	scoreEditBox := container.NewHBox(
		widget.NewButton("-", func() {
			scoreInt, _ := strconv.Atoi(scoreBox.Text)
			scoreBox.SetText(strconv.Itoa(scoreInt - 1))
			form.tabContext.InputScore.Set(scoreBox.Text)
		}),
		scoreBox,
		widget.NewButton("+", func() {
			scoreInt, _ := strconv.Atoi(scoreBox.Text)
			scoreBox.SetText(strconv.Itoa(scoreInt + 1))
			form.tabContext.InputScore.Set(scoreBox.Text)
		}),
	)

	return scoreEditBox
}

// 标签输入栏
func (form *ExtraInfoEditForm) newTagInputBox() *widget.Entry {
	inputBox := widget.NewEntry()
	inputBox.OnSubmitted = func(tag string) {
		tag = strings.TrimSpace(tag)
		if len(tag) == 0 {
			return
		}
		// 被选中项合并tag，同步到context
		oldInputTags, _ := form.tabContext.InputTags.Get()
		newInputTags := utils.MergeLists(oldInputTags[1:], []string{tag})
		sort.Slice(newInputTags, func(i, j int) bool {
			if len(newInputTags[i]) == len(newInputTags[j]) {
				return utils.CmpText(newInputTags[i], newInputTags[j])
			}
			return len(newInputTags[i]) < len(newInputTags[j])
		})
		newInputTags = append(store.DefaultTagList, newInputTags...)
		form.tabContext.InputTags.Set(newInputTags)
		// 清除输入
		inputBox.SetText("")
		// 刷新视图
		form.TagList.RefreshData()
	}
	return inputBox
}

// 标签展示栏
func (form *ExtraInfoEditForm) newTagsEditBox() fyne.CanvasObject {
	// tag展示列表
	form.TagList = packed_widgets.NewHList(
		0,
		func(index int) string {
			tag, _ := form.tabContext.InputTags.GetValue(index)
			return tag
		},
		func() int {
			return form.tabContext.InputTags.Length()
		},
	)
	// 双击tag删除
	form.TagList.OnDoubleTapped = func(col int) {
		oldTags, _ := form.tabContext.InputTags.Get()
		if col >= len(oldTags) || col == 0 {
			return
		}

		newTags := oldTags[:col]
		newTags = append(newTags, oldTags[col+1:]...)
		form.tabContext.InputTags.Set(newTags)
	}

	// tag输入框
	tagInputBox := form.newTagInputBox()
	tagsBox := container.NewBorder(nil, nil, tagInputBox, nil, form.TagList.Table)
	return tagsBox
}

// 备注输入栏
func (form *ExtraInfoEditForm) newNoteEditBox() fyne.CanvasObject {
	form.NoteBox = widget.NewEntry()
	// entry直接绑定bindData有bug，输入一串中文可能会吞字或溢出
	//noteBox.Bind(form.tabContext.InputNote)
	form.NoteBox.OnChanged = func(note string) {
		form.tabContext.InputNote.Set(note)
	}
	form.NoteBox.Validator = func(note string) error {
		note = strings.TrimSpace(note)
		// note不超过300字符
		if utf8.RuneCountInString(note) > 300 {
			return common.NewError("over length")
		}
		return nil
	}
	form.NoteBox.MultiLine = true
	form.NoteBox.SetMinRowsVisible(4)
	form.NoteBox.Wrapping = fyne.TextWrapBreak //自动换行
	return form.NoteBox
}

// 链接输入栏
func (form *ExtraInfoEditForm) newUrlEditBox() fyne.CanvasObject {
	urlBox := widget.NewEntry()
	urlBox.Bind(form.tabContext.InputUrl)
	return urlBox
}

// SubmitExtraInfo 提交ExtraInfo改动
func (form *ExtraInfoEditForm) SubmitExtraInfo() error {
	extraInfo, _ := service.ReadLatestExtraInfo(form.tabContext.Path)
	if extraInfo == nil {
		extraInfo = &model.FileExtraInfoMap{}
	}

	for _, fileInfo := range form.tabContext.FileInfos {
		extraInfo.SetFileExtraInfo(fileInfo.New)
	}

	return service.WriteExtraInfoFile(form.tabContext.Path, extraInfo)
}

func NewExtraInfoEditForm(tabContext *store.ExtraInfoTabContext, preview *ExtraInfoPreview) *ExtraInfoEditForm {
	form := &ExtraInfoEditForm{
		tabContext: tabContext,
		Form: &widget.Form{
			Items:    nil,
			OnCancel: func() {},
			OnSubmit: func() {},
		},
	}

	// 表单Cancel：被选中元素重置成初始值
	form.Form.OnCancel = func() {
		preview.ResetCheckItems()
		// 刷新EditForm内显示的聚合值
		tabContext.RefreshInputExtraInfo()
		form.TagList.Refresh()
	}

	// 表单提交
	form.Form.OnSubmit = func() {
		err := form.SubmitExtraInfo()
		if err != nil {
			logger.Error("SubmitExtraInfo failed, err=%v", err)
			packed_widgets.NewNotify(err.Error())
		}
		if form.OnSubmit != nil {
			// 回调上级
			form.OnSubmit()
		}
	}

	form.Append("评分", form.newScoreEditBox())
	form.Append("标签", form.newTagsEditBox())
	form.Append("备注", form.newNoteEditBox())
	form.Append("链接", form.newUrlEditBox())

	return form
}
