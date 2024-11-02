package service

import "encoding/base64"

func GetB64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func GetB64Decode(data string) string {
	decoded, _ := base64.StdEncoding.DecodeString(data)
	return string(decoded)
}
