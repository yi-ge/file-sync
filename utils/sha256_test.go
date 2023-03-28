package utils

import "testing"

func TestGetSha256Str(t *testing.T) {
	sha256 := GetSha256Str("123456")
	if sha256 != "7c4a8d09ca3762af61e59520943dc26494f8941b9962fba4e8d64d8" {
		t.Error("sha256 error")
	}
}
