package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AESCTREncrypt AES-CTR Encrypt
func AESCTREncrypt(plainText []byte, key []byte) ([]byte, error) {
	// cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Grouping Mode
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	// Encrypt / Decrypt
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

// AESCTREncryptWithBase64 AES-CTR Encrypt
func AESCTREncryptWithBase64(plainText []byte, key []byte) (string, error) {
	dst, err := AESCTREncrypt(plainText, key)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(dst), nil
}

// AESCTREncrypt AES-CTR Decrypt
func AESCTRDecrypt(encryptData []byte, key []byte) ([]byte, error) {
	return AESCTREncrypt(encryptData, key)
}

// AESCTRDecryptWithBase64 AES-CTR Decrypt
func AESCTRDecryptWithBase64(encryptDataStr string, key []byte) (string, error) {
	encryptData, err := base64.RawURLEncoding.DecodeString(encryptDataStr)
	if err != nil {
		return "", err
	}
	res, err := AESCTRDecrypt(encryptData, key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
