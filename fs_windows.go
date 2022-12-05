//go:build windows

package main

import (
	"os/user"

	"github.com/yi-ge/file-sync/utils"
)

var (
	workDir string
)

func fsInit() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := u.HomeDir

	workDir = homeDir + "\\.file-sync"

	utils.MakeDirIfNotExist(workDir)

	return nil
}

func getCachePath() string {
	return workDir + "\\cache.json"
}

func getDataPath() string {
	return workDir + "\\data.json"
}

func getPathSplitStr() string {
	return "\\"
}
