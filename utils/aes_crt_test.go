package utils

import (
	"testing"
)

func TestAesCrt(t *testing.T) {
	src := "你好"
	t.Log("原文：", src)
	//16byte密钥
	key := []byte("1234567887654321")
	//调用加密函数
	encryptData, err := aesCtrEncrypt([]byte(src), key)
	if err != nil {
		t.Log("err:", err)
		return
	}
	t.Logf("密文: %x\n", encryptData)

	//调用解密函数
	plainText, err := aesCtrDecrypt(encryptData, key)
	if err != nil {
		t.Log("err:", err)
		return
	}
	t.Logf("解密后明文: %s\n", plainText)
}
