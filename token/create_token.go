package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
)

var (
	signKey *rsa.PrivateKey
)

func InitSignKey(signKeyFilePath string) {
	// Read the private key
	privateKeyData, err := os.ReadFile(signKeyFilePath)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	// Parse the private key
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}
}
func SignToken(claim jwt.Claims) (string, error) {
	if signKey == nil {
		return "", fmt.Errorf("sign key is not initialized")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
