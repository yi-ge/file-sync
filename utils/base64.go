package utils

import "encoding/base64"

// Base64EncodeBytesToStr converts a byte slice to a Base64 Encoded String
func Base64EncodeBytesToStr(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// Base64EncodeStrToStr converts a string to a Base64 Encoded String
func Base64EncodeStrToStr(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64DecodeBytesToStr converts a Base64 byte slice to a Base64 Decoded Byte slice
func Base64DecodeBytesToBytes(input []byte) ([]byte, error) {
	return Base64DecodeStrToBytes(string(input))
}

// Base64DecodeStrToBytes converts a Base64 string to a Base64 Decoded Byte slice
func Base64DecodeStrToBytes(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// Base64DecodeStrToStr converts a Base64 string to a Base64 Decoded String
func Base64DecodeStrToStr(input string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func Base64SafetyEncode(text string) string {
	// $text = str_replace(['+', '/'], ['-', '_'], $text);
	return base64.URLEncoding.EncodeToString([]byte(text))
}

func Base64SafetyDecode(data string) ([]byte, error) {
	// $text = str_replace(['-', '_'], ['+', '/'], $text);
	return base64.URLEncoding.DecodeString(data)
}
