package utils

import (
	"testing"
)

func TestAESMAC(t *testing.T) {
	src := "Hello"
	t.Log("Text: ", src)
	// 16byte key
	key := []byte("1234567887654321")
	// Encrypt
	encryptData := AESMACEncryptBytes([]byte(src), string(key))
	t.Logf("Crypt data: %x\n", encryptData)

	// Decrypt
	decrypted, plainText, err := AESMACDecryptBytes(encryptData, string(key))
	if !decrypted || err != nil {
		t.Fatal("err:", err)
	}
	t.Logf("Text: %s\n", plainText)
}
