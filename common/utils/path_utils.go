package utils

import (
	"os"
	"path"
	"strings"
)

func DealWithWindowsPath(path string) string {
	if !IsOsWindows() || path == "/" {
		return path
	}
	if path[0] == '/' {
		path = path[1:]
	}
	if path[len(path)-1] == ':' {
		path = path + "/"
	}
	return path
}

func PathJoin(paths ...string) string {
	return DealWithWindowsPath(path.Join(paths...))
}

func PathSplit(absolutePath string) []string {
	if absolutePath == "/" {
		return []string{absolutePath}
	}
	absolutePath = DealWithWindowsPath(absolutePath)
	dirs := strings.Split(absolutePath, "/")
	if IsOsWindows() {
		dirs = append([]string{"/"}, dirs...)
	} else {
		dirs[0] = "/"
	}
	return dirs
}

func PathExists(path string) bool {
	path = DealWithWindowsPath(path)
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetPrefix(filename string) string {
	index := strings.Index(filename, ".")
	if index <= 0 {
		return filename
	}
	return filename[:index]
}

func AddSuffix(filename string, suffix string) string {
	index := strings.Index(filename, ".")
	if index <= 0 {
		return filename + suffix
	}
	return filename[:index] + suffix + filename[index:]
}
