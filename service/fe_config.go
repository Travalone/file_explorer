package service

import (
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
)

func getWorkDir() string {
	if utils.IsOsWindows() {
		return "."
	}
	current, _ := user.Current()
	workDir := utils.PathJoin(current.HomeDir, ".fe_config")
	if !utils.PathExists(workDir) {
		err := os.Mkdir(workDir, os.ModePerm)
		if err != nil {
			logger.Error("Create work dir failed, workDir=%v, err=%v", workDir, err)
			return ""
		}
	}
	return workDir
}

func ReadConfig() *model.FileExplorerConfig {
	configData, err := os.ReadFile(utils.PathJoin(getWorkDir(), "conf.yml"))
	if err != nil {
		logger.Error("ReadConfig failed err=%v", err)
		return &model.FileExplorerConfig{Root: "/"}
	}

	// 反序列化
	config := &model.FileExplorerConfig{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		logger.Error("ReadConfig unmarshal failed, err=%v", err)
		return &model.FileExplorerConfig{Root: "/"}
	}
	return config
}

func WriteConfig(config *model.FileExplorerConfig) error {
	logger.Debug("write config %v", config)
	// 序列化
	bytes, err := yaml.Marshal(config)
	if err != nil {
		logger.Error("WriteConfig Marshal failed, config=%v, err=%v", config, err)
		return err
	}

	// 写入
	fp, err := os.OpenFile(utils.PathJoin(getWorkDir(), "conf.yml"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("WriteConfig WriteFile failed, err=%v", err)
		return err
	}
	defer fp.Close()
	fp.Write(bytes)

	return nil
}
