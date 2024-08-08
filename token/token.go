package token

import (
	"github.com/dgrijalva/jwt-go"
	"homepage-authorization/user"
)

func CreateTokenByUser(claims jwt.MapClaims, userStruct user.User) (string, error) {
	claims["user_id"] = userStruct.ID
	claims["roles"] = userStruct.Roles
	token, err := SignToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
