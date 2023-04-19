package service

import (
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sort"
	"strings"
	"time"
)

// FillFileExtraInfo 指定目录，填充fileInfo
func FillFileExtraInfo(dir string, fileInfos []*model.FileInfo) {
	extraInfo, err := ReadLatestExtraInfo(dir)
	if err != nil {
		logger.Debug("FillFileExtraInfo ReadExtraInfo failed, dir=%v, err=%v", dir, err)
		return
	}

	if extraInfo == nil || extraInfo.FileExtraInfos == nil {
		logger.Debug("FillFileExtraInfo no file_extra_infos, dir=%v", dir)
		return
	}

	for _, fileInfo := range fileInfos {
		if fileExtraInfo, ok := extraInfo.FileExtraInfos[fileInfo.Name]; ok {
			fileInfo.ExtraInfo = fileExtraInfo
		}
	}
}

// GetExtraInfoFiles 获取指定目录下的ExtraInfo files
func GetExtraInfoFiles(dir string) []*model.FileInfo {
	metaDir := utils.PathJoin(dir, common.MetaDateDirName)
	dirEntries, err := os.ReadDir(metaDir)
	if err != nil || len(dirEntries) == 0 {
		return nil
	}

	extraInfoFileInterfaces := make([]interface{}, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), utils.GetPrefix(common.ExtraInfoFilename)) {
			extraInfoFileInterfaces = append(extraInfoFileInterfaces, buildFileInfo(metaDir, dirEntry))
		}
	}

	if len(extraInfoFileInterfaces) == 0 {
		return nil
	}

	// 按修改时间倒序排列
	sort.Slice(extraInfoFileInterfaces, func(i, j int) bool {
		a := extraInfoFileInterfaces[i].(*model.FileInfo)
		b := extraInfoFileInterfaces[j].(*model.FileInfo)
		return a.ModifyTime.Unix() > b.ModifyTime.Unix()
	})

	extraInfoFiles := make([]*model.FileInfo, len(dirEntries))
	for i, extraInfoFileInterface := range extraInfoFileInterfaces {
		extraInfoFiles[i] = extraInfoFileInterface.(*model.FileInfo)
	}

	return extraInfoFiles
}

// ReadLatestExtraInfo 指定目录，读取最新的extra info
func ReadLatestExtraInfo(dir string) (*model.FileExtraInfoMap, error) {
	// 获取最新的ExtraInfo file
	extraInfoFiles := GetExtraInfoFiles(dir)
	if len(extraInfoFiles) == 0 {
		return nil, nil
	}

	return ReadExtraInfo(extraInfoFiles[0].GetPath())
}

// ReadExtraInfo 指定绝对路径，读取extra info
func ReadExtraInfo(path string) (*model.FileExtraInfoMap, error) {
	// 读取ExtraInfo文件
	extraInfoFileData, err := os.ReadFile(path)
	if err != nil {
		logger.Debug("ReadExtraInfo ReadFile failed, path=%v, err=%v", path, err)
		return nil, err
	}

	// 反序列化
	extraInfo := &model.FileExtraInfoMap{}
	err = yaml.Unmarshal(extraInfoFileData, extraInfo)
	if err != nil {
		logger.Error("ReadExtraInfo Unmarshal failed, err=%v", err)
		return nil, err
	}

	return extraInfo, nil
}

// CreateMetaDirIfNotExist 在指定目录下创建meta目录(如果不存在)
func CreateMetaDirIfNotExist(dir string) error {
	metaDirPath := utils.PathJoin(dir, common.MetaDateDirName)
	if !utils.PathExists(metaDirPath) {
		err := os.Mkdir(metaDirPath, os.ModePerm)
		if err != nil {
			logger.Error("WriteExtraInfoFile create dir failed, path=%v, err=%v", metaDirPath, err)
			return err
		}
	}
	return nil
}

// WriteExtraInfoFile 指定目录，写入ExtraInfo
func WriteExtraInfoFile(dir string, extraInfo *model.FileExtraInfoMap) error {
	logger.Info("write extra info file %v", extraInfo)
	// ExtraInfo配置序列化
	bytes, err := yaml.Marshal(extraInfo)
	if err != nil {
		logger.Error("WriteExtraInfoFile Marshal failed, extraInfoMap=%v, err=%v", extraInfo, err)
		return err
	}

	err = CreateMetaDirIfNotExist(dir)
	if err != nil {
		logger.Error("CreateMetaDirIfNotExist failed, dir=%v, err=%v", dir, err)
		return err
	}

	extraInfoFilename := utils.AddSuffix(common.ExtraInfoFilename, fmt.Sprintf(".%d", time.Now().Unix()))

	// 写入ExtraInfo文件
	fp, err := os.OpenFile(utils.PathJoin(dir, common.MetaDateDirName, extraInfoFilename), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("WriteExtraInfoFile WriteFile failed, dir=%v, err=%v", dir, err)
		return err
	}
	defer fp.Close()
	fp.Write(bytes)

	return nil
}
