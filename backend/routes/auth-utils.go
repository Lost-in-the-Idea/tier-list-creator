package routes

import (
	"crypto/rand"
	"encoding/base64"
)

func generateStateCookie() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}