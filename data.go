package main

import (
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

type Data struct {
	Name                      string
	Email                     string
	Verify                    string
	RsaPrivateKeyPassword     string
	RsaPrivateEncryptPassword string
	MachineId                 string
	MachineName               string
	EncryptedMachineKey       string
}

func dataInit() error {
	dataPath := getDataPath()
	isNotFirst, err := utils.FileExists(dataPath)
	if err != nil {
		return err
	}

	if !isNotFirst {
		file, err := os.OpenFile(dataPath, os.O_CREATE|os.O_WRONLY, 0600)
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

func getData() ([]Data, error) {
	dataPath := getDataPath()
	file, err := os.OpenFile(dataPath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []Data
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

func setData(data Data) error {
	dataPath := getDataPath()
	file, err := os.OpenFile(dataPath, os.O_CREATE|os.O_WRONLY, 0600)
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
