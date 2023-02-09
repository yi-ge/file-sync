package utils

import (
	"testing"
)

func TestAesCrt(t *testing.T) {
	src := "Hello"
	t.Log("Text: ", src)

	// 16byte key
	key := []byte("1234567887654321")

	// Encrypt
	encryptData, err := AESCTREncryptWithBase64([]byte(src), key)
	if err != nil {
		t.Log("err:", err)
		return
	}
	t.Logf("Crypt data: %x\n", encryptData)

	// Decrypt
	plainText, err := AESCTRDecryptWithBase64(encryptData, key)
	if err != nil {
		t.Fatal("err:", err)
		return
	}
	t.Logf("Text: %s\n", plainText)
}
