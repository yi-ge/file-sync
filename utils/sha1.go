package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetSha1Str(str string) string {
	h := sha1.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}
