package main

import (
	"bytes"
	"encoding/base64"
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
	password = utils.GetSha256Str(password)
	verify := utils.GetSha1Str(password[:16])[8:]
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
		"privateKey":  base64.RawURLEncoding.EncodeToString(encryptedPrivateKeyStr),
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

	logger.Error(string(body))

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 && status != 2 {
		logger.Infof("Server out status error: %d\n", status)
		return false, "", "", errors.New("Server out status error: " + strconv.Itoa(status))
	}

	privateKey = jsoniter.Get(body, "result", "privateKey").ToString()
	publicKey = jsoniter.Get(body, "result", "publicKey").ToString()

	if publicKey == "" || privateKey == "" {
		return false, "", "", errors.New("publicKey or privateKey is empty")
	}

	publicKey, err = utils.AESCBCDecryptSafety([]byte(verify), publicKey)

	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	privateKey, err = utils.AESCBCDecryptSafety([]byte(verify), privateKey)

	if err != nil || privateKey == "" {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	privateKeyByte, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	privateKeyEncrypted, err := utils.AESCTRDecrypt(privateKeyByte, []byte(rsaPrivateEncryptPassword))

	if err != nil || len(privateKeyEncrypted) == 0 {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyEncrypted, rsaPrivateKeyPassword)

	if err != nil || !decrypted {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", "", err
	}

	privateKey = string(plaintextBytes)

	return status == 1, publicKey, privateKey, err
}
