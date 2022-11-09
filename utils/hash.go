package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// passphraseToHash returns a hexadecimal string of an SHA1 checksumed passphrase
func passphraseToHash(pass string) (string, []byte) {
	// The salt is used as a unique string to defeat rainbow table attacks
	saltHash := md5.New()
	saltHash.Write([]byte(pass))
	saltyBytes := saltHash.Sum(nil)
	salt := hex.EncodeToString(saltyBytes)

	saltyPass := []byte(pass + salt)
	hasher := sha1.New()
	hasher.Write(saltyPass)

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), hash
}
