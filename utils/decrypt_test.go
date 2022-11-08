package utils

import (
	"fmt"
	"strings"
	"testing"
)

const (
	rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3CsRKp3+OAq3P2pbEp2/
gEiewPw0TK+eh13hVhwj/xeqFumGj+9U/0v7mHq09ZpGYGgdOyqX5Qslzc/CEPj0
pwVbcQoKfPFfPu4bn09uEikT9LRhfVstV1vmIOTLSsTF3d6j5WRrFsz0iMxmJOh+
M2RT4CZtcShxqt5OA+qPJtFXIcC9wfkZN+s7HEaN+3XRGlrV8grY30CQsJcWF2hb
xQaRVDVYxt3kS/EPQ7/PpmM8ylf6vyEehFomKhr8inIc9WK+fIh85MfmGO3q/WRs
Wxt8nikMbGGQZovtg5wb4PFoG9Q7hikzfaO5jbSGn/q3/sGDRISoQ72h9JCQ5GFh
NQIDAQAB
-----END PUBLIC KEY-----`
	rsaPrivateKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxEUN1DZZ/XU2J6+3EoCX
6ZQExSyGrJlmcq2s4sxAqThJVGAv4BYqCQjnigUGaLF4+2khGHVXrx4LhwnW54iq
V3V3Xq59H0Cj3oGGWgKxSOM62xxfizmc1Og/6uAwZTAX4oCsgx5SMaFQbAU5ensM
VEX9CetXSGhc1bbS23kEHAkjJ0NryRl7DR/ilFKO5pAjTGEzP4aTkF/D3Eu3z15U
wdkf2WisEsANVTEnNHu2qvdiXGzRSLNF4mVFNO3AsgfnbgXzlN0feQ1HbH+J7Ue5
eHleCGhfS/PGFP3lQ4sA0hB4B/5eZ6ROo8YEuQiNTz+UMFteeGymTgFu2sOwLE10
wQIDAQAB
-----END PUBLIC KEY-----`
)

func TestEncryptAndDecrypt(t *testing.T) {
	data, _ := RsaEncryptWithSha1Base64("123456", rsaPublicKey)
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
