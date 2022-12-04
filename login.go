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
	// TODO: db操作
	// err := db.Update(func(txn *badger.Txn) error {
	// 	// e := txn.Set([]byte("email"), []byte(email))
	// 	// if e != nil {
	// 	// 	return e
	// 	// }
	// 	// e = txn.Set([]byte("password3"), []byte(password3))
	// 	// if e != nil {
	// 	// 	return e
	// 	// }
	// 	// e = txn.Set([]byte("machineId"), []byte(machineId))
	// 	// if e != nil {
	// 	// 	return e
	// 	// }
	// 	// e = txn.Set([]byte("machineName"), []byte(machineName))
	// 	// if e != nil {
	// 	// 	return e
	// 	// }
	// 	// e = txn.Set([]byte("publicKey"), []byte(publicKey))
	// 	// if e != nil {
	// 	// 	return e
	// 	// }
	// 	return txn.Set([]byte("privateKey"), []byte(privateKey))
	// })
	// if err != nil {
	// 	return err
	// }

	// wb := db.NewWriteBatch()
	// defer wb.Cancel()

	// wb.Set([]byte("privateKey"), []byte(privateKey))
	// wb.Flush()

	// defer db.Close()
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
		return errors.New("GET request failed: " + err.Error())
	}

	encryptedPrivateKeyStr, err := utils.AESCTREncrypt([]byte(encryptedPrivKeyPEMBase64), []byte(rsaPrivateEncryptPassword))
	if err != nil {
		return errors.New("AES encrypt failed: " + err.Error())
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
		return err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return errors.New("HTTP request failed: " + err.Error())
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

	if isNewUser {
		if publicKey != publicKeyPEM.String() {
			return errors.New("publicKey and publicKeyPEM are not equal")
		}

		return registerDevice(email, password[:48], machineId, machineName, publicKey, privateKey)
	}

	privateKeyByte, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		return errors.New("secret decrypt error: " + err.Error())
	}

	privateKeyEncrypted, err := utils.AESCTRDecrypt(privateKeyByte, []byte(rsaPrivateEncryptPassword))

	if err != nil || len(privateKeyEncrypted) == 0 {
		return errors.New("secret decrypt error: " + err.Error())
	}

	decrypted, plaintextBytes, err := utils.AESMACDecryptBytes(privateKeyEncrypted, rsaPrivateKeyPassword)

	if err != nil || !decrypted {
		return errors.New("secret decrypt error: " + err.Error())
	}

	privateKey = string(plaintextBytes)

	return registerDevice(email, password[:48], machineId, machineName, publicKey, privateKey)
}
