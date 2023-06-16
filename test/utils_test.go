package test

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	"sort"
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

func TestSortChs(t *testing.T) {
	a := []string{"八", "堺", "啊"}
	sort.Slice(a, func(i, j int) bool {
		return utils.CmpText(a[i], a[j])
	})

	logger.Info("res=%v", a)

}
