package utils

import (
	"bytes"
	"io/ioutil"
)

// writeRSAKeyPair creates key pair files
func writeRSAKeyPair(privKey *bytes.Buffer, pubKey *bytes.Buffer, path string) (bool, bool, error) {
	privKeyFile, err := writeKeyFile(privKey, path+".priv.pem", 0400)
	if err != nil {
		return false, false, err
	}

	pubKeyFile, err := writeKeyFile(pubKey, path+".pub.pem", 0644)
	if err != nil {
		return privKeyFile, false, err
	}
	return privKeyFile, pubKeyFile, nil
}

// writeKeyFile writes a public or private key file depending on the permissions, 644 for public, 400 for private
func writeKeyFile(pem *bytes.Buffer, path string, permission int) (bool, error) {
	pemByte, _ := ioutil.ReadAll(pem)
	keyFile, err := WriteByteFile(path, pemByte, permission, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}
