package test

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	"testing"
)

func TestWildcard(t *testing.T) {
	str := utils.ReplaceWildcard("https://www.baidu.com/s?wd={name}", map[string]string{
		"name": "zzz",
	})
	logger.Info("%v", str)
}

func TestCmd(t *testing.T) {
	utils.RunCommand(utils.GetExplorerCommand(), "E://")
}
