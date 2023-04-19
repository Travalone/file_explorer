package model

type FileExtraInfo struct {
	Note  *string  `yaml:"note,omitempty"`
	Tags  []string `yaml:"tags,omitempty"`
	Score *int     `yaml:"score,omitempty"`
	Url   *string  `yaml:"url,omitempty"`
}

type FileExtraInfoMap struct {
	FileExtraInfos map[string]*FileExtraInfo `yaml:"file_extra_infos,omitempty"`
}

func (extraInfoMap *FileExtraInfoMap) SetFileExtraInfo(fileInfo *FileInfo) {
	if extraInfoMap.FileExtraInfos == nil {
		extraInfoMap.FileExtraInfos = make(map[string]*FileExtraInfo)
	}

	fileInfo.FormatExtraInfo()

	extraInfoMap.FileExtraInfos[fileInfo.Name] = fileInfo.ExtraInfo
	if fileInfo.ExtraInfo == nil {
		delete(extraInfoMap.FileExtraInfos, fileInfo.Name)
		return
	}
}
