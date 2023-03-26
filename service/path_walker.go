package service

import (
	"errors"
	"file_explorer/common"
	"file_explorer/common/logger"
	"file_explorer/common/model"
	"file_explorer/common/utils"
	"io/fs"
	"path/filepath"
	"time"
)

var mockErr = errors.New("stop")

type PathWalker struct {
	refreshTime int64
	fileInfos   []*model.FileInfo

	walking  func(string)
	finished func(int, int, int64, bool, error)

	cancel bool
}

func (walker *PathWalker) Cancel() {
	walker.cancel = true
}

func (walker *PathWalker) Start() {
	go func() {
		regularCount, dirCount, size, ok := 0, 0, int64(0), true
		var err error

		for _, fileInfo := range walker.fileInfos {
			if fileInfo.Type == common.FILE_TYPE_REGULAR {
				// 文件直接累加
				regularCount += 1
				size += fileInfo.Size
				continue
			}

			// 目录dfs遍历
			lastUpdateTime := int64(0)
			root := utils.PathJoin(fileInfo.Dir, fileInfo.Name)
			err = filepath.WalkDir(root, func(path string, info fs.DirEntry, e error) error {
				if walker.cancel {
					// 手动停止
					logger.Debug("PathWalker cancelled, root=%v, path=%v", root, path)
					return mockErr
				}
				// 跳过root
				if path == root {
					return nil
				}
				// 每隔refreshTime (单位: ms)，回调更新状态
				if time.Now().UnixMilli()-lastUpdateTime >= walker.refreshTime {
					process := path[len(root):]
					walker.walking(process)
					lastUpdateTime = time.Now().UnixMilli()
				}
				// 统计文件、目录数量，占用空间
				if info.IsDir() {
					dirCount += 1
				} else {
					regularCount += 1
					fileInfo, e := info.Info()
					if e != nil {
						logger.Error("PathWalker getInfo failed, path=%v, err=%v", path, e)
						return e
					} else {
						size += fileInfo.Size()
					}
				}
				return nil
			})
			if err == mockErr {
				// 手动停止
				err = nil
				walker.cancel = false
				ok = false
				break
			}
		}

		// 遍历结束后回调
		walker.finished(regularCount, dirCount, size, ok, err)
	}()
}

func NewPathWalker(refreshTime int64, fileInfos []*model.FileInfo, walking func(string), finished func(int, int, int64, bool, error)) *PathWalker {
	return &PathWalker{
		refreshTime: refreshTime,
		fileInfos:   fileInfos,
		walking:     walking,
		finished:    finished,
	}
}
