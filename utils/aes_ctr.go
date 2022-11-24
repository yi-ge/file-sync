package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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

// AESCTREncrypt AES-CTR Decrypt
func AESCTRDecrypt(encryptData []byte, key []byte) ([]byte, error) {
	return AESCTREncrypt(encryptData, key)
}
