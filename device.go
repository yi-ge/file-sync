package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
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
	requestURL := apiURL + "/device/list"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":     data.Email,
		"machineId": machineId,
		"timestamp": string(timestamp),
		"token":     "verify",
	}
	jsonBody, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		return err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return errors.New("HTTP request failed: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return errors.New("HTTP request failed: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		msg := jsoniter.Get(body, "msg").ToString()
		return errors.New(msg)
	}

	// privateKey := jsoniter.Get(body, "result").ToString()

	return nil
}
