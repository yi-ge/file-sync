package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// PassphraseToHash returns a hexadecimal string of an SHA1 checksumed passphrase
func PassphraseToHash(pass string) (string, []byte) {
	// The salt is used as a unique string to defeat rainbow table attacks
	saltHash := md5.New()
	saltHash.Write([]byte(pass))
	saltyBytes := saltHash.Sum(nil)
	salt := hex.EncodeToString(saltyBytes)

	saltyPass := []byte(pass + salt)
	hash := sha1.New()
	hash.Write(saltyPass)

	hashByte := hash.Sum(nil)

	return hex.EncodeToString(hashByte), hashByte
}
