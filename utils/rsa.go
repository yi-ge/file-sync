package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

// GenerateRSAKeypair returns a private RSA key pair object
func GenerateRSAKeypair(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if keySize == 0 {
		keySize = 4096
	}
	// create our private and public key
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// PemEncodeRSAPrivateKey creates a PEM from an RSA Private key, and optionally returns an encrypted version
func PemEncodeRSAPrivateKey(privKey *rsa.PrivateKey, rsaPrivateKeyPassword string) (privKeyPEM *bytes.Buffer, encryptedPrivKeyPEMBase64 string) {
	privKeyPEM = new(bytes.Buffer)
	b := new(bytes.Buffer)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	pem.Encode(privKeyPEM, privateKeyBlock) // or EncodeToMemory

	/*
		Legacy encryption, insecure, replaced with AES-GCM encryption
		if rsaPrivateKeyPassword != "" {
			privateKeyBlock, _ = x509.EncryptPEMBlock(rand.Reader, privateKeyBlock.Type, privateKeyBlock.Bytes, []byte(rsaPrivateKeyPassword), x509.PEMCipherAES256)
		}
	*/

	if rsaPrivateKeyPassword != "" {
		encBytes := AESMACEncryptBytes(privKeyPEM.Bytes(), rsaPrivateKeyPassword)
		b.Write(encBytes)

		encryptedPrivKeyPEMBase64 = base64.RawURLEncoding.EncodeToString(b.Bytes())
	}

	return privKeyPEM, encryptedPrivKeyPEMBase64
}

// PemToEncryptedBytes takes a PEM byte buffer and encrypts it
func PemToEncryptedBytes(pem *bytes.Buffer, passphrase string) (b *bytes.Buffer) {
	b = new(bytes.Buffer)

	encBytes := AESMACEncryptBytes(pem.Bytes(), passphrase)
	b.Write(encBytes)

	return b
}

// PemEncodeRSAPublicKey takes a DER formatted RSA Public Key object and converts it to PEM format
func PemEncodeRSAPublicKey(caPubKey *rsa.PublicKey) (*bytes.Buffer, error) {
	derPkix, err := x509.MarshalPKIXPublicKey(caPubKey)
	if err != nil {
		return nil, err
	}
	caPubKeyPEM := new(bytes.Buffer)
	pem.Encode(caPubKeyPEM, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	})
	return caPubKeyPEM, nil
}

func GenerateRSAKeypairPEM(keySize int, rsaPrivateKeyPassword string) (privKeyPEM *bytes.Buffer, encryptedPrivKeyPEMBase64 string, publicKeyPEM *bytes.Buffer, err error) {
	privateKey, publicKey, err := GenerateRSAKeypair(keySize)
	if err != nil {
		return nil, "", nil, err
	}

	privKeyPEM, encryptedPrivKeyPEMBase64 = PemEncodeRSAPrivateKey(privateKey, rsaPrivateKeyPassword)

	publicKeyPEM, err = PemEncodeRSAPublicKey(publicKey)
	if err != nil {
		return nil, "", nil, err
	}

	return privKeyPEM, encryptedPrivKeyPEMBase64, publicKeyPEM, nil
}

// RsaEncryptWithSha1Base64 ???????????????sha1??????????????????base64??????
func RsaEncryptWithSha1Base64(originalData, publicKey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(publicKey)
	pubKey, _ := x509.ParsePKIXPublicKey(key)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

// RsaDecryptWithSha1Base64 ??????????????????sha1??????????????????base64????????????????????????????????????PKCS1?????????
func RsaDecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	key, _ := base64.StdEncoding.DecodeString(privateKey)
	prvKey, _ := x509.ParsePKCS1PrivateKey(key)
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	return string(originalData), err
}

// RsaSignWithSha1Hex ???????????????sha1??????????????????????????????hex???????????????PKCS8?????????
func RsaSignWithSha1Hex(data string, prvKey string) (string, error) {
	keyBytes, err := hex.DecodeString(prvKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}

// RsaSignWithSha1HexPkcs1 ???????????????sha1??????????????????????????????hex???????????????PKCS1?????????
func RsaSignWithSha1HexPkcs1(data string, pkcs1keyStr string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pkcs1keyStr))
	if block == nil {
		return nil, errors.New("no PrivateKey found")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	h := sha1.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return nil, err
	}

	return signature, nil
}

// RsaVerySignWithSha1Base64 ??????????????????sha1????????????????????????base64???????????????????????????
func RsaVerySignWithSha1Base64(originalData, signData, pubKey string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

// RsaSignWithSha256
func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

// RsaVerySignWithSha256
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// ????????????
func RsaEncrypt(data, keyBytes []byte) []byte {
	//??????pem???????????????
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// ????????????
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// ????????????
	pub := pubInterface.(*rsa.PublicKey)
	//??????
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return cipherText
}

// ????????????
func RsaDecrypt(cipherText, keyBytes []byte) []byte {
	//????????????
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	//??????PKCS1???????????????
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// ??????
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
	if err != nil {
		panic(err)
	}
	return data
}
