package server

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/sirupsen/logrus"
)

const (
	secretKeySize = 64
)

func getRandomSecretKey() []byte {
	randomBytes := make([]byte, secretKeySize)
	token := make([]byte, base32.StdEncoding.EncodedLen(len(randomBytes)))
	if _, err := rand.Read(randomBytes); err != nil {
		logrus.Error("fail to get secret key for JWTAuth")
	}

	base32.StdEncoding.Encode(token, randomBytes)

	return token
}
