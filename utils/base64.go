package utils

import "encoding/base64"

// B64EncodeBytesToStr converts a byte slice to a Base64 Encoded String
func B64EncodeBytesToStr(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// B64EncodeStrToStr converts a string to a Base64 Encoded String
func B64EncodeStrToStr(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// B64DecodeBytesToStr converts a Base64 byte slice to a Base64 Decoded Byte slice
func B64DecodeBytesToBytes(input []byte) ([]byte, error) {
	return B64DecodeStrToBytes(string(input))
}

// B64DecodeStrToBytes converts a Base64 string to a Base64 Decoded Byte slice
func B64DecodeStrToBytes(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// B64DecodeStrToStr converts a Base64 string to a Base64 Decoded String
func B64DecodeStrToStr(input string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
