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

func addConfig(fileId string, fileName string, path string, actionMachineId string, data Data) (jsoniter.Any, error) {
	requestURL := apiURL + "/file/config"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return nil, errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":           utils.GetSha1Str(data.Email),
		"machineId":       machineId,
		"action":          "add",
		"fileId":          fileId,
		"fileName":        fileName,
		"path":            path,
		"actionMachineId": actionMachineId,
		"attribute":       "",
		"timestamp":       strconv.FormatInt(timestamp, 10),
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

	res := jsoniter.Get(body, "result")

	return res, nil
}

func removeConfig(fileId string, actionMachineId string, data Data) error {
	requestURL := apiURL + "/file/config"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":           utils.GetSha1Str(data.Email),
		"machineId":       machineId,
		"fileId":          fileId,
		"action":          "remove",
		"actionMachineId": actionMachineId,
		"attribute":       "",
		"timestamp":       strconv.FormatInt(timestamp, 10),
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
		return err
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		return err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		return errors.New("secret decrypt error: " + err.Error())
	}

	token, err := utils.RsaSignWithSha1HexPkcs1(dataParams, string(plaintextBytes))
	if err != nil {
		return err
	}

	bodyMap["token"] = base64.RawURLEncoding.EncodeToString(token)

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
	// fmt.Println(string(body))

	if err != nil {
		return err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		msg := jsoniter.Get(body, "msg").ToString()
		return errors.New(msg)
	}

	return nil
}

func listConfigs(data Data) (jsoniter.Any, error) {
	requestURL := apiURL + "/file/configs"
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

	configs := jsoniter.Get(body, "result")

	return configs, nil
}

func fileCheck(email string, fileId string, sha256 string) (int, error) {
	requestURL := apiURL + "/file/check"
	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":     utils.GetSha1Str(email),
		"fileId":    fileId,
		"sha256":    sha256,
		"timestamp": strconv.FormatInt(timestamp, 10),
	}

	jsonBody, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		return -1, err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return -1, errors.New("HTTP request failed: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return -1, errors.New("HTTP request failed: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// fmt.Println(string(body))

	if err != nil {
		return -1, err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 0 && status != 1 && status != 2 {
		msg := jsoniter.Get(body, "msg").ToString()
		return -1, errors.New(msg)
	}

	return status, nil
}

func fileUpload(fileId string, fileName string, sha256 string, content string, updateAt int64, data Data) error {
	requestURL := apiURL + "/file/sync"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	// TODO: File content encryption
	content = "File content encryption"

	bodyMap := map[string]string{
		"email":     utils.GetSha1Str(data.Email),
		"machineId": machineId,
		"fileId":    fileId,
		"fileName":  fileName,
		"content":   content,
		"sha256":    sha256,
		"updateAt":  strconv.FormatInt(updateAt, 10),
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

	privateKeyEncrypted, err := getPrivateKey()
	if err != nil {
		return err
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		return err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		return errors.New("secret decrypt error: " + err.Error())
	}

	token, err := utils.RsaSignWithSha1HexPkcs1(dataParams, string(plaintextBytes))
	if err != nil {
		return err
	}

	bodyMap["token"] = base64.RawURLEncoding.EncodeToString(token)

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
	// fmt.Println(string(body))

	if err != nil {
		return err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		msg := jsoniter.Get(body, "msg").ToString()
		return errors.New(msg)
	}

	return nil
}

func fileDownload(fileId string, data Data) (jsoniter.Any, error) {
	requestURL := apiURL + "/file/sync"
	machineId := utils.GetMachineID()

	if machineId != data.MachineId {
		return nil, errors.New("machineId does not match")
	}

	timestamp := time.Now().UnixNano() / 1e6

	bodyMap := map[string]string{
		"email":     utils.GetSha1Str(data.Email),
		"machineId": machineId,
		"fileId":    fileId,
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

	file := jsoniter.Get(body, "result")

	return file, nil
}
