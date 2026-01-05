package idp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GetSecretHash(username, clientId, clientSecret string) string {
	message := username + clientId

	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
