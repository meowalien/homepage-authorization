package user

import (
	"database/sql"
	"google.golang.org/api/oauth2/v2"
	"homepage-authorization/postgresql"
)

func GetUserIDByUserInfo(userInfo *oauth2.Userinfo) (string, error) {
	db := postgresql.GetCollection()
	var userID string

	err := db.QueryRow("SELECT id FROM users WHERE google_id = $1;", userInfo.Id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, create a new user
			err = db.QueryRow(
				`INSERT INTO users (google_id, email, name, picture_url) 
				VALUES ($1, $2, $3, $4) RETURNING id;`,
				userInfo.Id, userInfo.Email, userInfo.Name, userInfo.Picture,
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
