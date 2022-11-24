package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// 加密
func aesCtrEncrypt(plainText []byte, key []byte) ([]byte, error) {
	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

// 解密
func aesCtrDecrypt(encryptData []byte, key []byte) ([]byte, error) {
	return aesCtrEncrypt(encryptData, key)
}
