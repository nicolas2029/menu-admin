package authorization

import (
	"crypto/rsa"
	"os"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)

func LoadCertificates() error {
	var err error
	once.Do(func() {
		err = loadCertificates()
	})

	return err
}

func loadCertificates() error {
	privateBytes, _ := os.LookupEnv("MENU-RSA")

	publicBytes, _ := os.LookupEnv("MENU-RSA-PUB")

	return parseRSA([]byte(privateBytes), []byte(publicBytes))
}

func parseRSA(privateBytes, publicBytes []byte) error {
	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
