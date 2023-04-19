package store

import (
	"file_explorer/common/model"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/packed_widgets/packed_binding"
)

type OpListContext struct {
	OpFileInfos *packed_binding.BindingList
	CheckList   *packed_widgets.CheckList
	FeContext   *FeContext
}

func (context *OpListContext) GetOpFileInfos() []*model.FileInfo {
	return model.Interfaces2FileInfos(context.OpFileInfos.Get())
}

func NewOpListContext(feContext *FeContext) *OpListContext {
	context := &OpListContext{
		OpFileInfos: packed_binding.NewBindingList(nil),
		FeContext:   feContext,
	}

	return context
}
