package main

import (
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func login(email string, password string, hostname string) (newUser bool, publicKey string, privateKey string) {
	requestURL := apiURL + "/device/add"
	resp, err := http.Post(requestURL, "application/json", nil)

	if err != nil {
		logger.Infof("GET request failed: %s\n", err)
		return false, "", ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Infof("Body read failed: %s\n", err)
		return false, "", ""
	}

	status := jsoniter.Get(body, "status").ToString()

	if status != "1" {
		logger.Infof("Server out status error: %s\n", status)
		return false, "", ""
	}

	secretData, err := utils.PublicDecryptWithBase64(jsoniter.Get(body, "result").ToString(), publicKey)

	if err != nil {
		logger.Errorf("Secret decrypt error: %s\n", err)
		return false, "", ""
	}

	signBase64 := jsoniter.Get(secretData, "sign").ToString()

	return true, signBase64, ""
}
