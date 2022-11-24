package utils

import (
	"testing"
)

func TestBase64(t *testing.T) {
	src := "123456789, abc, 中文，中文"
	t.Log("Text: ", src)

	data := Base64SafetyEncode(src)
	t.Logf("Data: %s", data)

	text, err := Base64SafetyDecode(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Text: %s", string(text))
}
