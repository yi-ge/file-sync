package main

import (
	"bytes"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func login(email string, password string, machineName string) (newUser bool, publicKey string, privateKey string) {
	requestURL := apiURL + "/device/add"
	machineId := utils.GetMachineID()
	verify := utils.GetSha1Str(password[:64])
	publicKey = ""
	privateKey = ""
	bodyMap := map[string]string{
		"email":       email,
		"machineId":   machineId,
		"machineName": machineName,
		"verify":      verify,
		"publicKey":   publicKey,
		"privateKey":  privateKey,
	}
	jsonBody, _ := jsoniter.Marshal(bodyMap)
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return false, "", ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Infof("Body read failed: %s\n", err)
		return false, "", ""
	}

	status := jsoniter.Get(body, "status").ToString()

	if status != "1" {
		logger.Infof("Server out status error: %s\n", status)
		return false, "", ""
	}

	secretData, err := utils.RsaDecryptWithSha1Base64(jsoniter.Get(body, "result").ToString(), publicKey)

	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", ""
	}

	signBase64 := jsoniter.Get([]byte(secretData), "sign").ToString()

	return true, signBase64, ""
}
