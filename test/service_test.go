package test

import (
	"file_explorer/common/logger"
	"os"
	"os/exec"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	path := "/Users/bytedance/Downloads/test.txt"
	os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	//err := os.Remove(path)
	//if err != nil {
	//	return
	//}
}

func TestOpenUrl(t *testing.T) {
	//service.OpenUrlsWithDefaultWebExplorer([]string{"baidu.com"})
	//service.OpenPathWithDefaultFileExplorer("E://1 2")
	cmd := exec.Command("cmd", "/C", "start", "", "explorer")
	err := cmd.Run()
	logger.Debug(cmd.String())
	if err != nil {
		cmd.String()
		logger.Debug("cmd err, cmd=%v, err=%v", cmd.String(), err)
		return
	}
}
