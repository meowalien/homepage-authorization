package token

func CreateTokenByUserID(userID string) (string, error) {
	claims := DefaultClaims()
	claims["user_id"] = userID
	token, err := SignToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
