package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// padding
func padding(src []byte) []byte {
	paddingNum := aes.BlockSize - len(src)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(paddingNum)}, paddingNum)
	return append(src, padText...)
}

// unPadding
func unPadding(src []byte) ([]byte, error) {
	length := len(src)
	unPaddingNum := int(src[length-1])

	if unPaddingNum > length {
		return nil, errors.New("unPadding error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unPaddingNum)], nil
}

func AESCBCEncrypt(key []byte, text string) (string, error) {
	cipherText, err := AESCBCEncryptRaw(key, text)
	if err != nil {
		return "", err
	}

	finalMsg := (base64.StdEncoding.EncodeToString(cipherText))

	return finalMsg, nil
}

func AESCBCEncryptSafety(key []byte, text string) (string, error) {
	cipherText, err := AESCBCEncryptRaw(key, text)
	if err != nil {
		return "", err
	}

	finalMsg := (base64.URLEncoding.EncodeToString(cipherText))

	return finalMsg, nil
}

func AESCBCEncryptRaw(key []byte, text string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := padding([]byte(text))
	cipherText := make([]byte, aes.BlockSize+len(msg))

	// Randomly generated vectors
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], msg)

	return cipherText, nil
}

func AESCBCDecrypt(key []byte, text string) (string, error) {
	decodedMsg, err := base64.StdEncoding.DecodeString((text))

	if err != nil {
		return "", err
	}

	return AESCBCDecryptRaw(key, decodedMsg)
}

func AESCBCDecryptSafety(key []byte, text string) (string, error) {
	decodedMsg, err := base64.URLEncoding.DecodeString((text))

	if err != nil {
		return "", err
	}

	return AESCBCDecryptRaw(key, decodedMsg)
}

func AESCBCDecryptRaw(key []byte, decodedMsg []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multiplier of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(msg, msg)

	unPaddingMsg, err := unPadding(msg)
	if err != nil {
		return "", err
	}

	return string(unPaddingMsg), nil
}
