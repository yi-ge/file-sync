package utils

import "testing"

func TestGetSha1Str(t *testing.T) {
	sha1 := GetSha1Str("123456")
	if sha1 != "7c4a8d09ca3762af61e59520943dc26494f8941b" {
		t.Error("sha1 is not correct")
	}
}
