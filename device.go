package main

import (
	"bytes"
	"errors"
	"fmt"
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
		fmt.Println("key:", k, "Value:", bodyMap[k])
		dataParams = dataParams + k + "=" + bodyMap[k] + "&"
	}

	dataParams += data.Verify
	fmt.Println(dataParams)
	// ff := dataParams[0 : len(dataParams)-1]
	// fmt.Println("去掉最后一个&：", ff)

	//对字符串进行sha1哈希
	// h := sha1.New()
	// h.Write([]byte(dataParams))
	// bs := h.Sum(nil)
	// sign := hex.EncodeToString(bs)

	// fmt.Println(sign)

	prvKey, err := getPrivateKey()
	if err != nil {
		return err
	}
	token, err := utils.RsaSignWithSha1Hex(dataParams, string(prvKey))
	if err != nil {
		return err
	}

	bodyMap["token"] = utils.Base64SafetyEncode(token)

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
	fmt.Println(string(body))

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
