package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func registerDevice(
	email string,
	password3 string,
	machineId string,
	machineName string,
	publicKey string,
	privateKey string) error {
	return nil
}

func login(email string, password string, machineName string) error {
	requestURL := apiURL + "/device/add"
	machineId := utils.GetMachineID()
	password = utils.GetSha256Str(password)
	verify := utils.GetSha1Str(password[:16])[8:]
	rsaPrivateKeyPassword := password[32:48]
	rsaPrivateEncryptPassword := password[16:32]
	_, encryptedPrivKeyPEMBase64, publicKeyPEM, err := utils.GenerateRSAKeypairPEM(4096, rsaPrivateKeyPassword)
	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return err
	}

	encryptedPrivateKeyStr, err := utils.AESCTREncrypt([]byte(encryptedPrivKeyPEMBase64), []byte(rsaPrivateEncryptPassword))
	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return err
	}
	bodyMap := map[string]string{
		"email":       email,
		"machineId":   machineId,
		"machineName": machineName,
		"verify":      verify,
		"publicKey":   publicKeyPEM.String(),
		"privateKey":  base64.RawURLEncoding.EncodeToString(encryptedPrivateKeyStr),
	}
	jsonBody, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Infof("Body read failed: %s\n", err)
		return err
	}

	logger.Error(string(body))

	status := jsoniter.Get(body, "status").ToInt()

	if status != 1 && status != 2 {
		logger.Infof("Server out status error: %d\n", status)
		return errors.New("Server out status error: " + strconv.Itoa(status))
	}

	isNewUser := status == 1
	privateKey := jsoniter.Get(body, "result", "privateKey").ToString()
	publicKey := jsoniter.Get(body, "result", "publicKey").ToString()

	if publicKey == "" || privateKey == "" {
		return errors.New("publicKey or privateKey is empty")
	}

	publicKey, err = utils.AESCBCDecryptSafety([]byte(verify), publicKey)

	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return err
	}

	publicKeyCheck := strings.Split(publicKey, "@")
	privateKeyCheck := strings.Split(privateKey, "@")

	publicKeyTimestamp, err := strconv.ParseInt(publicKeyCheck[0], 10, 64)
	if err != nil {
		return err
	}

	privateKeyTimestamp, err := strconv.ParseInt(privateKeyCheck[0], 10, 64)
	if err != nil {
		return err
	}

	nowTimestamp := time.Now().UnixNano() / 1e6

	if nowTimestamp-publicKeyTimestamp > 2000 || nowTimestamp-privateKeyTimestamp > 2120 {
		return errors.New("Calculation timeout.")
	}

	publicKey = publicKeyCheck[1]
	privateKey = privateKeyCheck[1]

	privateKey, err = utils.AESCBCDecryptSafety([]byte(verify), privateKey)

	if err != nil || privateKey == "" {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return err
	}

	if isNewUser {
		if publicKey != publicKeyPEM.String() {
			logger.Errorf("publicKey and publicKeyPEM are not equal")
			return errors.New("publicKey and publicKeyPEM are not equal")
		}

		return registerDevice(email, password[:48], machineId, machineName, publicKey, privateKey)
	}

	privateKeyByte, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return err
	}

	privateKeyEncrypted, err := utils.AESCTRDecrypt(privateKeyByte, []byte(rsaPrivateEncryptPassword))

	if err != nil || len(privateKeyEncrypted) == 0 {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return err
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyEncrypted, rsaPrivateKeyPassword)

	if err != nil || !decrypted {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return err
	}

	privateKey = string(plaintextBytes)

	return registerDevice(email, password[:48], machineId, machineName, publicKey, privateKey)
}
