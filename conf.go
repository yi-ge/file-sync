package main

import (
	"errors"
	"io"
	"os"

	"github.com/yi-ge/file-sync/utils"
)

func setConfig(data string) error {
	machineId := utils.GetMachineIDUseSHA256()
	configPath := getConfigPath()
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encryptData, err := utils.AESCTREncryptWithBase64([]byte(data), []byte(machineId[:32]))
	if err != nil {
		return err
	}

	file.Write([]byte(encryptData))

	return nil
}

func getConfig() string {
	machineId := utils.GetMachineIDUseSHA256()
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.OpenFile(configPath, os.O_RDONLY, 0600)
		if err != nil {
			return ""
		}
		defer file.Close()

		dataBytes, err := io.ReadAll(file)
		if err != nil {
			return ""
		}

		encryptData, err := utils.AESCTRDecryptWithBase64(string(dataBytes), []byte(machineId[:32]))
		if err != nil {
			return ""
		}

		return encryptData
	}

	return ""
}

func delConfig() error {
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New("config file is not exists")
	}
	err := os.Remove(configPath)
	if err != nil {
		return err
	}
	return nil
}
