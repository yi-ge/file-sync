package main

import (
	"errors"
	"io"
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

func getCache() (jsoniter.Any, error) {
	cachePath := getCachePath()
	file, err := os.OpenFile(cachePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data jsoniter.Any
	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func setCache(data jsoniter.Any) error {
	cachePath := getCachePath()
	file, err := os.OpenFile(cachePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonBytes, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	file.Write(jsonBytes)
	return nil
}

func delCache() error {
	cachePath := getCachePath()
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return errors.New("cache file is not exists")
	}
	err := os.Remove(cachePath)
	if err != nil {
		return err
	}
	return nil
}
