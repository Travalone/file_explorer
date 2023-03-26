package model

import (
	"file_explorer/common"
	"file_explorer/common/utils"
	"time"
)

type FileInfo struct {
	Name       string         `json:"name"`
	Dir        string         `json:"dir"`
	Type       byte           `json:"type"`
	Size       int64          `json:"size"`
	ModifyTime time.Time      `json:"modify_time"`
	ExtraInfo  *FileExtraInfo `json:"extra_info"`
}

func (info *FileInfo) GetSize() string {
	if info.IsDir() {
		return ""
	}
	return utils.ConvSize(info.Size)
}

func (info *FileInfo) GetModifyTime() string {
	if info.Type == common.FILE_TYPE_DRIVER {
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

func (info *FileInfo) GetNote() string {
	if info.ExtraInfo != nil && info.ExtraInfo.Note != nil {
		return utils.Conv2Str(*info.ExtraInfo.Note)
	}
	return ""
}

func (info *FileInfo) GetPath() string {
	return utils.PathJoin(info.Dir, info.Name)
}

func (info *FileInfo) IsDir() bool {
	return info.Type == common.FILE_TYPE_DIR || info.Type == common.FILE_TYPE_DRIVER
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
