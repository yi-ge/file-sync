package utils

import (
	"testing"
)

func TestPassphraseToHash(t *testing.T) {
	pass := "password"

	t.Run("Test with a simple passphrase", func(t *testing.T) {
		hash, hashByte := PassphraseToHash(pass)
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})

	t.Run("Test with an empty passphrase", func(t *testing.T) {
		hash, hashByte := PassphraseToHash("")
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})

	t.Run("Test with a long passphrase", func(t *testing.T) {
		longPass := "This is a very long passphrase that should be more than 64 bytes to test the hash algorithm."
		hash, hashByte := PassphraseToHash(longPass)
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})

	t.Run("Test with special characters in the passphrase", func(t *testing.T) {
		specialPass := "!@#$%^&*()"
		hash, hashByte := PassphraseToHash(specialPass)
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})

	t.Run("Test with non-ASCII characters in the passphrase", func(t *testing.T) {
		nonASCIIPass := "こんにちは"
		hash, hashByte := PassphraseToHash(nonASCIIPass)
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})

	t.Run("Test with a passphrase that contains NULL characters", func(t *testing.T) {
		nullPass := "pass\x00word"
		hash, hashByte := PassphraseToHash(nullPass)
		if hash == "" {
			t.Errorf("Expected a non-empty string, but got an empty string")
		}

		if len(hashByte) == 0 {
			t.Errorf("Expected a non-empty byte slice, but got an empty byte slice")
		}
	})
}
