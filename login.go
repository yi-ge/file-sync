package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func login(email string, password string, machineName string) error {
	requestURL := apiURL + "/device/add"
	machineId := utils.GetMachineID()
	password = utils.GetSha256Str(password)
	verify := utils.GetSha1Str(password[:16])[8:]
	rsaPrivateKeyPassword := password[32:48]
	rsaPrivateEncryptPassword := password[16:32]
	machineKeyEncryptPassword := password[48:]
	_, encryptedPrivKeyPEMBase64, publicKeyPEM, err := utils.GenerateRSAKeypairPEM(4096, rsaPrivateKeyPassword)
	if err != nil {
		return errors.New("GET request failed: " + err.Error())
	}

	encryptedPrivateKeyStr, err := utils.AESCTREncryptWithBase64([]byte(encryptedPrivKeyPEMBase64), []byte(rsaPrivateEncryptPassword))
	if err != nil {
		return errors.New("AES encrypt failed: " + err.Error())
	}

	encryptedMachineName, err := utils.AESCTREncryptWithBase64([]byte(machineName), []byte(verify))
	if err != nil {
		return errors.New("AES encrypt failed: " + err.Error())
	}
	bodyMap := map[string]string{
		"email":       email,
		"machineId":   machineId,
		"machineName": encryptedMachineName,
		"verify":      verify,
		"publicKey":   publicKeyPEM.String(),
		"privateKey":  encryptedPrivateKeyStr,
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

	if status != 1 && status != 2 {
		msg := jsoniter.Get(body, "msg").ToString()
		return errors.New(msg)
	}

	isNewUser := status == 1
	privateKey := jsoniter.Get(body, "result", "privateKey").ToString()
	publicKey := jsoniter.Get(body, "result", "publicKey").ToString()
	machineKey := jsoniter.Get(body, "result", "machineKey").ToString()

	if publicKey == "" || privateKey == "" {
		return errors.New("publicKey or privateKey is empty")
	}

	publicKey, err = utils.AESCBCDecryptSafety([]byte(verify), publicKey)

	if err != nil {
		return errors.New("secret decrypt error: " + err.Error())
	}

	privateKey, err = utils.AESCBCDecryptSafety([]byte(verify), privateKey)

	if err != nil || privateKey == "" {
		return errors.New("secret decrypt error: " + err.Error())
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
		return errors.New("calculation timeout")
	}

	publicKey = publicKeyCheck[1]
	privateKey = privateKeyCheck[1]

	encryptedMachineKey := utils.AESMACEncryptBytesSafety([]byte(machineKey), machineKeyEncryptPassword)

	if isNewUser {
		if publicKey != publicKeyPEM.String() {
			return errors.New("publicKey and publicKeyPEM are not equal")
		}

		// data := Data{
		// 	Email:                     email,
		// 	Verify:                    verify,
		// 	RsaPrivateKeyPassword:     rsaPrivateKeyPassword,
		// 	RsaPrivateEncryptPassword: rsaPrivateEncryptPassword,
		// 	MachineId:                 machineId,
		// 	MachineName:               machineName,
		// 	EncryptedMachineKey:       string(encryptedMachineKey),
		// }

		// return registerDevice(data, publicKey, privateKey)
	}
	privateKeyEncrypted, err := utils.AESCTRDecryptWithBase64(privateKey, []byte(rsaPrivateEncryptPassword))

	if err != nil || len(privateKeyEncrypted) == 0 {
		return errors.New("secret decrypt error: " + err.Error())
	}

	data := Data{
		Email:                     email,
		Verify:                    verify,
		RsaPrivateKeyPassword:     rsaPrivateKeyPassword,
		RsaPrivateEncryptPassword: rsaPrivateEncryptPassword,
		MachineId:                 machineId,
		MachineName:               machineName,
		EncryptedMachineKey:       encryptedMachineKey,
	}

	return registerDevice(data, publicKey, privateKeyEncrypted)
}
