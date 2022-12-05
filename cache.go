package main

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func cacheInit() error {
	cachePath := getCachePath()
	isNotFirst, err := utils.FileExists(cachePath)
	if err != nil {
		return err
	}

	if !isNotFirst {
		file, err := os.OpenFile(cachePath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		dataMap := make(map[string]interface{})
		jsonBytes, err := jsoniter.Marshal(dataMap)
		if err != nil {
			return err
		}
		file.Write(jsonBytes)
	}

	return nil
}
