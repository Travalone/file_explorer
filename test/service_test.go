package test

import (
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/service"
	"sort"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	a := []string{}
	a = nil
	logger.Info("%v", a)
}

func TestQueryPath(t *testing.T) {
	path := "/Users/bytedance"

	fileInfos, _ := service.QueryPath(path)

	interfaces := model.FileInfos2Interfaces(fileInfos)

	sort.Slice(interfaces, func(i, j int) bool {
		a := interfaces[i].(*model.FileInfo)
		b := interfaces[j].(*model.FileInfo)
		if a.Type < b.Type {
			return true
		} else if a.Type > b.Type {
			return false
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	fileInfos = model.Interfaces2FileInfos(interfaces)

	for _, fileInfo := range fileInfos {
		logger.Info("%v", fileInfo)
	}

}

func TestQuerySubDirs(t *testing.T) {
	path := "/Users/bytedance/Documents"

	subDirs, _ := service.QuerySubDirs(path)

	for _, subDir := range subDirs {
		logger.Info("%v", subDir)
	}

}

func TestWalkPath(t *testing.T) {
	//path := "/Users/bytedance/Documents"
	//ok := false
	//walker := service.NewPathWalker(path, 1000, func(s string, b bool) {
	//	logger.Info("%v", s)
	//	ok = b
	//})
	//walker.Start()
	//for !ok {
	//	time.Sleep(time.Second)
	//}

	//rCount, dCount, size, err := walker.GetResult()
	//logger.Info("%v %v %v %v", rCount, dCount, size, err)

}
