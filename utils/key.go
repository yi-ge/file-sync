package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// WriteKeyFile writes a public or private key file depending on the permissions, 644 for public, 400 for private
func WriteKeyFile(pem *bytes.Buffer, path string, permission int) (bool, error) {
	pemByte, _ := io.ReadAll(pem)
	keyFile, err := WriteByteFile(path, pemByte, permission, false)
	if err != nil {
		return false, err
	}
	return keyFile, nil
}

// WriteRSAKeyPair creates key pair files
func WriteRSAKeyPair(privKey *bytes.Buffer, pubKey *bytes.Buffer, path string) (bool, bool, error) {
	privKeyFile, err := WriteKeyFile(privKey, path+".priv.pem", 0400)
	if err != nil {
		return false, false, err
	}

	pubKeyFile, err := WriteKeyFile(pubKey, path+".pub.pem", 0644)
	if err != nil {
		return privKeyFile, false, err
	}
	return privKeyFile, pubKeyFile, nil
}

// DeleteRSAKeyPair delete key pair files
func DeleteRSAKeyPair(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("key pair files is not exists")
	}

	err := os.Remove(path + ".priv.pem")
	if err != nil {
		return err
	}

	err = os.Remove(path + ".pub.pem")
	if err != nil {
		return err
	}

	return nil
}
