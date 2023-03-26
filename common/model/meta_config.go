package model

type FileExtraInfo struct {
	Note  *string  `json:"note" yaml:"note,omitempty"`
	Tags  []string `json:"tags" yaml:"tags,omitempty"`
	Score *int     `json:"score" yaml:"score,omitempty"`
}

type MetaData struct {
	FileExtraInfos map[string]*FileExtraInfo `yaml:"file_extra_infos,omitempty"`
}

func (meta *MetaData) SetFileExtraInfo(filename string, info *FileExtraInfo) {
	if meta.FileExtraInfos == nil {
		meta.FileExtraInfos = make(map[string]*FileExtraInfo)
	}
	meta.FileExtraInfos[filename] = info
	if info == nil || (info.Score == nil && len(info.Tags) == 0 && (info.Note == nil || len(*info.Note) == 0)) {
		delete(meta.FileExtraInfos, filename)
		return
	}

	if len(info.Tags) == 0 {
		meta.FileExtraInfos[filename].Tags = nil
	}
	if info.Note != nil && len(*info.Note) == 0 {
		meta.FileExtraInfos[filename].Note = nil
	}
}
