package main

import (
	"bytes"
	"errors"

	"github.com/yi-ge/file-sync/utils"
)

func registerDevice(
	data Data, publicKey string, privateKey string) error {
	err := setData(data)
	if err != nil {
		return err
	}

	privateKeyBuff := bytes.NewBufferString(privateKey)
	publicKeyBuff := bytes.NewBufferString(publicKey)

	privKeyFile, pubKeyFile, err := utils.WriteRSAKeyPair(privateKeyBuff, publicKeyBuff, workDir+getPathSplitStr())
	if err != nil {
		return err
	}

	if !privKeyFile {
		return errors.New("private key write failure")
	}

	if !pubKeyFile {
		return errors.New("public key write failure")
	}
	return nil
}

func removeDevice(machineKey string, data Data) error {
	return nil
}

func listDevices(data Data) error {
	return nil
}
