package main

import (
	"fmt"
	"os/user"
	"path/filepath"

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

	if configPath == "" {
		homeDir := u.HomeDir

		if homeDir == "" {
			return fmt.Errorf("could not find home directory")
		}

		workDir = filepath.Join(homeDir, ".file-sync")
	} else {
		workDir = configPath
	}

	utils.MakeDirIfNotExist(workDir)

	return nil
}

func getConfigPath() string {
	return filepath.Join(workDir, "config")
}

func getCachePath() string {
	return filepath.Join(workDir, "cache.json")
}

func getDataPath() string {
	return filepath.Join(workDir, "data.json")
}

func getPathSplitStr() string {
	return string(filepath.Separator)
}
