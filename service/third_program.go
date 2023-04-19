package service

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	net "net/url"
	"strings"
)

func OpenPathWithDefaultFileExplorer(path string) {
	path = utils.DealWithWindowsPath(path)

	if utils.IsOsWindows() && path == "/" {
		utils.RunCommand("explorer.exe")
		return
	}

	utils.RunCommand(utils.GetExplorerCommand(), path)
}

func OpenUrlsWithDefaultWebExplorer(urls []string) {
	for _, url := range urls {
		if len(url) == 0 {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(url), "http") {
			url = "https://" + url
		}
		if _, err := net.ParseRequestURI(url); err != nil {
			logger.Debug("ParseRequestURI failed, url=%v, err=%v", url, err)
			continue
		}
		utils.RunCommand(utils.GetExplorerCommand(), url)
	}
}
