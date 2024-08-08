package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	DefaultExpireDuration = time.Hour * 72 // 3 days
)

func DefaultClaims(exp time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"exp": time.Now().Add(exp).Unix(),
	}
}
