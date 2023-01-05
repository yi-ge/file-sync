package main

import (
	"os"
)

func getPublicKey() ([]byte, error) {
	data, err := os.ReadFile(workDir + getPathSplitStr() + ".pub.pem")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getPrivateKey() ([]byte, error) {
	data, err := os.ReadFile(workDir + getPathSplitStr() + ".priv.pem")
	if err != nil {
		return nil, err
	}

	return data, nil
}
