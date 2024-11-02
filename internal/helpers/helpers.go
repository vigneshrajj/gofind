package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetB64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func GetB64Decode(data string) string {
	decoded, _ := base64.StdEncoding.DecodeString(data)
	return string(decoded)
}

func Sha256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedBytes := hash.Sum(nil)
	hashedHex := hex.EncodeToString(hashedBytes)
	return hashedHex
}
