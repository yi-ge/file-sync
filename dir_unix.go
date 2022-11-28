//go:build !windows

package main

import (
	"os/user"

	"github.com/yi-ge/file-sync/utils"
)

func getDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := u.HomeDir

	dir := homeDir + "/.file-sync"

	utils.MakeDirIfNotExist(dir)

	return dir, nil
}

func getDBPath() string {
	dir, err := getDir()
	if err != nil {
		log.Fatal(err)
	}

	return dir + "/db"
}
