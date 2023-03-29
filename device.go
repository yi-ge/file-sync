package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
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

func listDevices(data Data) (jsoniter.Any, error) {
	requestURL := apiURL + "/device/list"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return nil, errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":     utils.GetSha1Str(data.Email),
		"machineId": machineId,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	var dataParams string
	var keys []string
	for k := range bodyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		dataParams = dataParams + k + "=" + bodyMap[k] + "&"
	}

	dataParams += data.Verify
	// fmt.Println(dataParams)
	// ff := dataParams[0 : len(dataParams)-1]

	privateKeyEncrypted, err := getPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		return nil, err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		return nil, errors.New("secret decrypt error: " + err.Error())
	}

	token, err := utils.RsaSignWithSha1HexPkcs1(dataParams, string(plaintextBytes))
	if err != nil {
		return nil, err
	}

	bodyMap["token"] = base64.RawURLEncoding.EncodeToString(token)

	jsonBody, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.New("HTTP request failed: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP request failed: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if err != nil {
		return nil, err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		msg := jsoniter.Get(body, "msg").ToString()
		return nil, errors.New(msg)
	}

	devices := jsoniter.Get(body, "result")
	// fmt.Printf("%T\n", devices)
	// fmt.Println(devices.ToString())

	return devices, nil
}

func removeDevice(machineKey string, data Data, removeMachineId string) (string, error) {
	requestURL := apiURL + "/device/remove"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return "", errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":           utils.GetSha1Str(data.Email),
		"machineId":       machineId,
		"timestamp":       strconv.FormatInt(timestamp, 10),
		"removeMachineId": removeMachineId,
		"machineKey":      utils.GetSha1Str(machineKey),
	}

	var dataParams string
	var keys []string
	for k := range bodyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		dataParams = dataParams + k + "=" + bodyMap[k] + "&"
	}

	dataParams += data.Verify
	// fmt.Println(dataParams)
	// ff := dataParams[0 : len(dataParams)-1]

	privateKeyEncrypted, err := getPrivateKey()
	if err != nil {
		return "", err
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		return "", err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		return "", errors.New("secret decrypt error: " + err.Error())
	}

	token, err := utils.RsaSignWithSha1HexPkcs1(dataParams, string(plaintextBytes))
	if err != nil {
		return "", err
	}

	bodyMap["token"] = base64.RawURLEncoding.EncodeToString(token)

	jsonBody, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", errors.New("HTTP request failed: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return "", errors.New("HTTP request failed: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		msg := jsoniter.Get(body, "msg").ToString()
		return "", errors.New(msg)
	}

	removedMachineId := jsoniter.Get(body, "result", "removedMachineId").ToString()

	if removedMachineId == machineId {
		err = utils.DeleteRSAKeyPair(workDir + getPathSplitStr())
		if err != nil {
			return "", err
		}

		err = delData()
		if err != nil {
			return "", err
		}

		err = delCache()
		if err != nil {
			return "", err
		}
	}

	return removedMachineId, nil
}
