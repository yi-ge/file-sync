package main

import (
	"io/ioutil"
)

func getPublicKey() ([]byte, error) {
	data, err := ioutil.ReadFile(workDir + getPathSplitStr() + ".pub.pem")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getPrivateKey() ([]byte, error) {
	data, err := ioutil.ReadFile(workDir + getPathSplitStr() + ".priv.pem")
	if err != nil {
		return nil, err
	}

	return data, nil
}
