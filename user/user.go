package user

import (
	"database/sql"
	"fmt"
	"homepage-authorization/oauth"
	"homepage-authorization/postgresql"
)

func GetUserIDByUserInfo(userInfo oauth.GoogleUserInfo) (string, error) {
	db := postgresql.GetCollection()
	var userID string

	err := db.QueryRow("SELECT id FROM users WHERE email = $1;", userInfo.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User does not exist, creating a new user")
			// User does not exist, create a new user
			err = db.QueryRow(
				`INSERT INTO users ( email, name) 
				VALUES ($1, $2) RETURNING id;`,
				userInfo.Email, userInfo.Name,
			).Scan(&userID)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return userID, nil
}
