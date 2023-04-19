package service

import (
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"os"
)

func CreateFile(fileInfo *model.FileInfo) error {
	path := utils.DealWithWindowsPath(fileInfo.GetPath())

	if fileInfo.Type == common.FileTypeDir {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			logger.Error("CreateFile Mkdir failed, path=%v, err=%v", path, err)
			return err
		}
		return nil
	}

	if fileInfo.Type != common.FileTypeRegular {
		logger.Warn("CreateFile invalid type=%v", fileInfo.Type)
		return nil
	}

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("CreateFile OpenFile failed, path=%v, err=%v", path, err)
		return err
	}
	defer fp.Close()
	return nil
}

func DeleteFile(fileInfos []*model.FileInfo) {
	
}
