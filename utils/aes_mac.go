package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// encryptBytes is a function that takes a plain byte slice and a passphrase and returns an encrypted byte slice
func encryptBytes(bytesIn []byte, passphrase string) []byte {
	passHash, _ := passphraseToHash(passphrase)
	targetPassHash := passHash[0:32]

	// Create an AES Cipher
	block, err := aes.NewCipher([]byte(targetPassHash))
	check(err)

	// Create a new gcm block container
	gcm, err := cipher.NewGCM(block)
	check(err)

	// Never use more than 2^32 random nonces with a given key because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	// Seal will encrypt the file using the GCM mode, appending the nonce and tag (MAC value) to the final data, so we can use it to decrypt it later.
	return gcm.Seal(nonce, nonce, bytesIn, nil)
}

// decryptBytes takes in a byte slice from a file and a passphrase then returns if the encrypted byte slice was decrypted, if so the plaintext contents, and any errors
func decryptBytes(bytesIn []byte, passphrase string) (decrypted bool, plaintextBytes []byte, err error) {
	// bytesIn must be decoded from base 64 first
	// b64.StdEncoding.DecodeString(bytesIn)

	passHash, _ := passphraseToHash(passphrase)
	targetPassHash := passHash[0:32]

	// Create an AES Cipher
	block, err := aes.NewCipher([]byte(targetPassHash))
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	// Create a new gcm block container
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	nonce := bytesIn[:gcm.NonceSize()]
	cipherText := bytesIn[gcm.NonceSize():]
	plaintextBytes, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Panic(err)
		return false, []byte{}, err
	}

	// successfully decrypted
	return true, plaintextBytes, nil
}
