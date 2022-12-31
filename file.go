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
