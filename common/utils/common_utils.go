package utils

import (
	"file_explorer/common/logger"
	"os/exec"
	"runtime"
)

func IsOsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsOsMac() bool {
	return runtime.GOOS == "darwin"
}

func IsOsLinux() bool {
	return runtime.GOOS == "linux"
}

func GetExplorerCommand() string {
	if IsOsWindows() {
		return "start"
	}
	if IsOsMac() {
		return "open"
	}
	if IsOsLinux() {
		return "xdg-open"
	}
	return ""
}

func RunCommand(name string, args ...string) error {
	if IsOsWindows() {
		args = append([]string{"/C", name}, args...)
		name = "cmd"
	}
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		logger.Error("RunCommand failed, name=%v, args=%v, err=%v", name, args, err)
		return err
	}
	return nil
}

