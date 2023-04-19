package model

import (
	"file_explorer/common"
	"file_explorer/common/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"time"
)

type FileInfo struct {
	Name       string         `yaml:"name,omitempty"`
	Dir        string         `yaml:"dir,omitempty"`
	Type       byte           `yaml:"type,omitempty"`
	Size       int64          `yaml:"size,omitempty"`
	ModifyTime time.Time      `yaml:"modify_time,omitempty"`
	ExtraInfo  *FileExtraInfo `yaml:"extra_info,omitempty"`
}

func (info *FileInfo) GetSize() string {
	if info.IsDir() {
		return ""
	}
	return utils.ConvSize(info.Size)
}

func (info *FileInfo) GetModifyTime() string {
	if info.Type == common.FileTypeDriver {
		return ""
	}
	return info.ModifyTime.Format("2006-01-02 15:04:05")
}

func (info *FileInfo) GetScore() string {
	if info.ExtraInfo != nil && info.ExtraInfo.Score != nil {
		return utils.Conv2Str(*info.ExtraInfo.Score)
	}
	return ""
}

func (info *FileInfo) GetTagsStr() string {
	if info.ExtraInfo != nil && info.ExtraInfo.Tags != nil {
		return utils.Conv2Str(info.ExtraInfo.Tags)
	}
	return ""
}

func (info *FileInfo) GetTagList() []string {
	if info.ExtraInfo != nil && info.ExtraInfo.Tags != nil {
		return info.ExtraInfo.Tags
	}
	return nil
}

func (info *FileInfo) GetUrl() string {
	if info.ExtraInfo != nil && info.ExtraInfo.Url != nil {
		return utils.Conv2Str(*info.ExtraInfo.Url)
	}
	return ""
}
func (info *FileInfo) GetNote() string {
	if info.ExtraInfo != nil && info.ExtraInfo.Note != nil {
		return utils.Conv2Str(*info.ExtraInfo.Note)
	}
	return ""
}

func (info *FileInfo) SetScore(score *int) {
	if info.ExtraInfo == nil {
		info.ExtraInfo = &FileExtraInfo{}
	}
	info.ExtraInfo.Score = score
}

func (info *FileInfo) SetTags(tags []string) {
	if info.ExtraInfo == nil {
		info.ExtraInfo = &FileExtraInfo{}
	}
	info.ExtraInfo.Tags = tags
}
func (info *FileInfo) SetNote(note *string) {
	if info.ExtraInfo == nil {
		info.ExtraInfo = &FileExtraInfo{}
	}
	info.ExtraInfo.Note = note
}
func (info *FileInfo) SetUrl(url *string) {
	if info.ExtraInfo == nil {
		info.ExtraInfo = &FileExtraInfo{}
	}
	info.ExtraInfo.Url = url
}

func (info *FileInfo) GetPath() string {
	return utils.PathJoin(info.Dir, info.Name)
}

func (info *FileInfo) GetIcon() fyne.Resource {
	switch info.Type {
	case common.FileTypeDir:
		return theme.FolderIcon()
	case common.FileTypeRegular:
		return theme.FileIcon()
	case common.FileTypeDriver:
		return theme.ComputerIcon()
	default:
		return nil
	}
}

func (info *FileInfo) FormatExtraInfo() {
	if info.ExtraInfo == nil {
		return
	}

	if info.ExtraInfo.Note != nil && len(*info.ExtraInfo.Note) == 0 {
		info.ExtraInfo.Note = nil
	}
	if len(info.ExtraInfo.Tags) == 0 {
		info.ExtraInfo.Tags = nil
	}
	if info.ExtraInfo.Url != nil && len(*info.ExtraInfo.Url) == 0 {
		info.ExtraInfo.Tags = nil
	}

	if info.ExtraInfo.Note == nil &&
		info.ExtraInfo.Tags == nil &&
		info.ExtraInfo.Score == nil &&
		info.ExtraInfo.Url == nil {
		info.ExtraInfo = nil
	}
}

func (info *FileInfo) IsDir() bool {
	return info.Type == common.FileTypeDir || info.Type == common.FileTypeDriver
}

func (info *FileInfo) Clone() *FileInfo {
	clone := &FileInfo{
		Name:       info.Name,
		Dir:        info.Dir,
		Type:       info.Type,
		Size:       info.Size,
		ModifyTime: info.ModifyTime,
	}
	if info.ExtraInfo != nil {
		clone.ExtraInfo = &FileExtraInfo{}
		if info.ExtraInfo.Note != nil {
			note := *info.ExtraInfo.Note
			clone.ExtraInfo.Note = &note
		}
		if len(info.ExtraInfo.Tags) > 0 {
			clone.ExtraInfo.Tags = make([]string, len(info.ExtraInfo.Tags))
			for i, tag := range info.ExtraInfo.Tags {
				clone.ExtraInfo.Tags[i] = tag
			}
		}
		if info.ExtraInfo.Score != nil {
			score := *info.ExtraInfo.Score
			clone.ExtraInfo.Score = &score
		}
		if info.ExtraInfo.Url != nil {
			url := *info.ExtraInfo.Url
			clone.ExtraInfo.Url = &url
		}
	}
	return clone
}

func FileInfos2Interfaces(fileInfos []*FileInfo) []interface{} {
	interfaces := make([]interface{}, len(fileInfos))
	for i, item := range fileInfos {
		interfaces[i] = item
	}
	return interfaces
}

func Interfaces2FileInfos(interfaces []interface{}) []*FileInfo {
	fileInfos := make([]*FileInfo, len(interfaces))
	for i, item := range interfaces {
		fileInfos[i] = item.(*FileInfo)
	}
	return fileInfos
}
