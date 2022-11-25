package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256Str(str string) string {
	h := sha256.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}
