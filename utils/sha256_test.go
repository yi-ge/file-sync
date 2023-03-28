package utils

import "testing"

func TestGetSha256Str(t *testing.T) {
	sha256 := GetSha256Str("123456")
	if sha256 != "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92" {
		t.Error("sha256 error")
	}
}
