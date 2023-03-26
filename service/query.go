package service

import (
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"github.com/shirou/gopsutil/disk"
	"os"
)

// QuerySubDirs 查询子目录
func QuerySubDirs(dir string) ([]string, error) {
	fileInfos, err := QueryPath(dir)
	if err != nil {
		return nil, err
	}

	subDirs := make([]string, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subDirs = append(subDirs, fileInfo.Name)
		}
	}

	return subDirs, nil
}

// QueryPath 查询子文件
func QueryPath(dir string) ([]*model.FileInfo, error) {
	dir = utils.DealWithWindowsPath(dir)
	logger.Debug("QueryPath dir=%s", dir)

	if utils.IsOsWindows() && dir == "/" {
		return getWindowsDrives()
	}

	// 读取目录
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("QueryPath ReadFile failed, dir=%s, err=%s", dir, err)
		return nil, err
	}

	// model转换
	fileInfos := make([]*model.FileInfo, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if fileInfo := buildFileInfo(dir, dirEntry); fileInfo != nil {
			fileInfos = append(fileInfos, fileInfo)
		}
	}

	return fileInfos, nil
}

func getWindowsDrives() ([]*model.FileInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		logger.Error("getWindowsDrives failed, err=%v", err)
		return nil, err
	}
	drives := make([]*model.FileInfo, len(partitions))
	for i, partition := range partitions {
		drives[i] = &model.FileInfo{
			Name: partition.Mountpoint,
			Type: common.FILE_TYPE_DRIVER,
		}
	}
	return drives, nil
}

func buildFileInfo(dir string, entry os.DirEntry) *model.FileInfo {
	rawInfo, err := entry.Info()
	if err != nil {
		logger.Warn("buildFileInfo getInfo failed, err=%v", err)
		return nil
	}
	fileInfo := &model.FileInfo{
		Name:       entry.Name(),
		Dir:        dir,
		Type:       common.FILE_TYPE_REGULAR,
		Size:       rawInfo.Size(),
		ModifyTime: rawInfo.ModTime(),
	}
	if entry.IsDir() {
		fileInfo.Type = common.FILE_TYPE_DIR
	}
	return fileInfo
}
