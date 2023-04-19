package model

type PreviewFileInfo struct {
	Original *FileInfo
	New      *FileInfo
}

func NewPreviewFileInfo(fileInfo *FileInfo) *PreviewFileInfo {
	return &PreviewFileInfo{
		Original: fileInfo,
		New:      fileInfo.Clone(),
	}
}

func PreviewFileInfos2Interfaces(fileInfos []*PreviewFileInfo) []interface{} {
	interfaces := make([]interface{}, len(fileInfos))
	for i, item := range fileInfos {
		interfaces[i] = item
	}
	return interfaces
}

func Interfaces2PreviewFileInfos(interfaces []interface{}) []*PreviewFileInfo {
	fileInfos := make([]*PreviewFileInfo, len(interfaces))
	for i, item := range interfaces {
		fileInfos[i] = item.(*PreviewFileInfo)
	}
	return fileInfos
}
