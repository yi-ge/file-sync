package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func login(email string, password string, machineName string) (newUser bool, publicKey string, privateKey string, err error) {
	requestURL := apiURL + "/device/add"
	machineId := utils.GetMachineID()
	verify := utils.GetSha1Str(password[:16])
	rsaPrivateKeyPassword := password[32:48]
	rsaPrivateEncryptPassword := password[16:32]
	_, encryptedPrivKeyPEMBase64, publicKeyPEM, err := utils.GenerateRSAKeypairPEM(4096, rsaPrivateKeyPassword)
	encryptedPrivateKeyStr, err := utils.AESCTREncrypt([]byte(encryptedPrivKeyPEMBase64), []byte(rsaPrivateEncryptPassword))
	if err != nil {
		return false, "", "", err
	}
	bodyMap := map[string]string{
		"email":       email,
		"machineId":   machineId,
		"machineName": machineName,
		"verify":      verify,
		"publicKey":   publicKeyPEM.String(),
		"privateKey":  string(encryptedPrivateKeyStr),
	}
	jsonBody, _ := jsoniter.Marshal(bodyMap)
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return false, "", "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Infof("Body read failed: %s\n", err)
		return false, "", "", err
	}

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 {
		logger.Infof("Server out status error: %s\n", status)
		return false, "", "", errors.New("Server out status error: " + strconv.Itoa(status))
	}

	// result := jsoniter.Get(body, "result")

	privateKey = jsoniter.Get(body, "result", "privateKey").ToString()
	publicKey = jsoniter.Get(body, "result", "publicKey").ToString()

	secretData, err := utils.RsaDecryptWithSha1Base64(privateKey, publicKey)

	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	signBase64 := jsoniter.Get([]byte(secretData), "sign").ToString()

	return true, signBase64, "", err
}
