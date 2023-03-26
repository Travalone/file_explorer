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
	"time"
)

// FillFileExtraInfo 指定目录，填充fileInfo
func FillFileExtraInfo(dir string, fileInfos []*model.FileInfo) {
	metaData, err := ReadLatestMetaFile(dir)
	if err != nil {
		logger.Debug("FillFileExtraInfo ReadMetaFile failed, dir=%v, err=%v", dir, err)
		return
	}

	if metaData == nil || metaData.FileExtraInfos == nil {
		logger.Debug("FillFileExtraInfo no file_extra_infos, dir=%v", dir)
		return
	}

	for _, fileInfo := range fileInfos {
		if fileExtraInfo, ok := metaData.FileExtraInfos[fileInfo.Name]; ok {
			fileInfo.ExtraInfo = fileExtraInfo
		}
	}
}

// GetMetaFiles 获取指定目录下的meta files
func GetMetaFiles(path string) []*model.FileInfo {
	metaDir := utils.PathJoin(path, common.META_DIR_NAME)
	dirEntries, err := os.ReadDir(metaDir)
	if err != nil {
		return nil
	}
	metaFileInterfaces := make([]interface{}, len(dirEntries))
	for i, dirEntry := range dirEntries {
		metaFileInterfaces[i] = buildFileInfo(metaDir, dirEntry)
	}
	// 按修改时间倒序排列
	sort.Slice(metaFileInterfaces, func(i, j int) bool {
		a := metaFileInterfaces[i].(*model.FileInfo)
		b := metaFileInterfaces[j].(*model.FileInfo)
		return a.ModifyTime.Unix() > b.ModifyTime.Unix()
	})

	metaFiles := make([]*model.FileInfo, len(dirEntries))
	for i, metaFileInterface := range metaFileInterfaces {
		metaFiles[i] = metaFileInterface.(*model.FileInfo)
	}

	return metaFiles
}

// ReadLatestMetaFile 指定目录，读取最新的meta data
func ReadLatestMetaFile(dir string) (*model.MetaData, error) {
	// 获取最新的meta file
	metaFilenames := GetMetaFiles(dir)
	if len(metaFilenames) == 0 {
		return nil, nil
	}

	return ReadMetaFile(metaFilenames[0].GetPath())
}

// ReadMetaFile 指定绝对路径，读取meta data
func ReadMetaFile(path string) (*model.MetaData, error) {
	// 读取meta文件
	metaFileData, err := os.ReadFile(path)
	if err != nil {
		logger.Debug("ReadMetaFile ReadFile failed, path=%s, err=%s", path, err)
		return nil, err
	}

	// 反序列化
	metaData := &model.MetaData{}
	err = yaml.Unmarshal(metaFileData, metaData)
	if err != nil {
		logger.Error("ReadMetaFile Unmarshal failed, err=%s", err)
		return nil, err
	}

	return metaData, nil
}

// WriteMetaFile 指定目录，写入metaInfo
func WriteMetaFile(path string, metaData *model.MetaData) error {
	logger.Info("write meta file %v", metaData)
	// meta配置序列化
	bytes, err := yaml.Marshal(metaData)
	if err != nil {
		logger.Error("WriteMetaFile Marshal failed, metaData=%s, err=%s", metaData, err)
		return err
	}

	metaDirPath := utils.PathJoin(path, common.META_DIR_NAME)
	if !utils.PathExists(metaDirPath) {
		err := os.Mkdir(metaDirPath, os.ModePerm)
		if err != nil {
			logger.Error("WriteMetaFile create dir failed, path=%s, err=%s", metaDirPath, err)
			return err
		}
	}

	// meta file名称，添加时间戳后缀
	metaFilename := utils.AddSuffix(common.EXTRA_INFO_FILE_NAME, fmt.Sprintf(".%d", time.Now().Unix()))

	// 写入meta文件
	fp, err := os.OpenFile(utils.PathJoin(metaDirPath, metaFilename), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("WriteMetaFile WriteFile failed, path=%s, err=%s", path, err)
		return err
	}
	defer fp.Close()
	fp.Write(bytes)

	return nil
}
