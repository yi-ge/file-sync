package utils

import (
	"fmt"
	"strings"
	"testing"
)

const (
	rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwygxDgXdIjbUlhgCuVay
1TOTQhiJ2828bHd+MlDB6YEgr5SlD8l3690hA0qC0CzN1BFCmDmTJwd9JJ/bPdzp
gymtf0o/20pxqruAiZV/dj8ojqbOTuKa09wU/18tyzhWMFCxHOqrt0W0vI+eAQEV
lDN8mo4rfGgMdZ0aktt8YJjA0+NW8Ccpv/OGRZ8sc6eNKAqJ5kQDEUZCzNIt0iq3
hvquv718XXsd/X8vPnFmYobDWnRTZhEgffd0+w48ZmaXK6E1ZeCgX57ElhCLSIw/
tj9ohm6bioyYHQKWKgkIb9UvFlIDmOdvaVT9mEXk1s1jnSnFP8jco3Qn0imk645L
fmZrN7GVgjzn+GmVajzrShaKkvzTXw0yAN6jBNJTPWciiw437y0GmlQLWdOs0AwS
LtCW9kk/V7VlZ7/28oSw2kBTJMoOcoydyM5nVz365cYW8D9m1ZEtLK41XMaRjN5u
ReB9ykrOXpYX71693b9Q4cnM7aZbcBXawpwMynGWurkhcpePtmMgh39eMB8MIPPX
cR7/K+YWJQQakWBzDZFU9s0z3XHD7MK2f23NxPI7iaKnIGHR8AkR2RxiSb2sDBIj
2yRcbk+EOEjzkgFcCKfjtbS2E4dlyM79idsQNqpj7Hii0mthJ6GHhEUVTb2i3RdH
48f4QHKZmFCuqSgnaVsFdkECAwEAAQ==
-----END PUBLIC KEY-----`
	rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
)

func TestGenerateRSAKeypair(t *testing.T) {
	rsaPrivateKeyPassword := "123"

	privKeyPEM, encryptedPrivKeyPEMBase64, publicKeyPEM, err := GenerateRSAKeypairPEM(4096, rsaPrivateKeyPassword)

	if err != nil {
		t.Error(err)
	}

	t.Logf("privKeyPEM: %s", privKeyPEM)
	t.Logf("encryptedPrivKeyPEMBase64: %s", encryptedPrivKeyPEMBase64)
	t.Logf("publicKeyPEM: %s", publicKeyPEM)
}

func TestEncryptAndDecrypt(t *testing.T) {
	data, err := RsaEncryptWithSha1Base64("123456", rsaPublicKey)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))
	res, _ := RsaDecryptWithSha1Base64(data, rsaPrivateKey)
	if strings.Compare(string("123456"), res) != 0 {
		t.Fatal("TestEncryptAndDecrypt Error.")
	}
}

func TestSignAndVery(t *testing.T) {
	var testSignData = `123456`
	data, _ := RsaSignWithSha1Hex(testSignData, rsaPrivateKey)
	fmt.Println(string(data))

	err := RsaVerySignWithSha1Base64("123456", data, rsaPublicKey)
	if err != nil {
		t.Fatal("TestSignAndVery Error.")
	}
}
