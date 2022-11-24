package utils

import (
	"testing"
)

func TestAesCBC(t *testing.T) {
	key := []byte("abcdabcdabcdabcdabcdabcdabcdabcd")

	src := "1234567891111111"
	t.Logf("Text: %s", src)

	enc, err := AESCBCEncrypt(key, src)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Encrypt data: %s", enc)

	decryptData, err := AESCBCDecrypt(key, enc)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("Decrypt text: %s", decryptData)

	// PHP AES-256-CBC encryptText
	encryptText := "DF8xTlJ03oyOohmsoEB33UsybnM30vopCIiRxr4-DVzAajKBaxOrOypRHPkCltfM"

	rawText, err := AESCBCDecryptSafety(key, encryptText)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("Raw text is %s \n", rawText)
}
