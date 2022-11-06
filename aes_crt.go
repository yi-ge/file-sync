package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// 加密
func aesCtrEncrypt(plainText []byte, key []byte) ([]byte, error) {
	//TODO
	//aes包，go内置标准库

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

func aes_crt_test() {
	src := "你好"
	fmt.Println("原文：", src)
	//16byte密钥
	key := []byte("1234567887654321")
	//调用加密函数
	encryptData, err := aesCtrEncrypt([]byte(src), key)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("密文: %x\n", encryptData)

	//调用解密函数
	plainText, err := aesCtrDecrypt(encryptData, key)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("解密后明文: %s\n", plainText)
}
