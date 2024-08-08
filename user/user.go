package user

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"homepage-authorization/oauth"
	"homepage-authorization/postgresql"
	"homepage-authorization/role"
)

type User struct {
	ID    string
	Email string
	Name  string
	Roles []string
}

func GetUserByUserInfo(userInfo oauth.GoogleUserInfo) (User, error) {
	db := postgresql.GetClient()
	var user User

	tx, err := db.Begin()
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	user, err = findUserByEmail(tx, userInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Debugf("User does not exist, creating a new user: %s", userInfo.Email)
			user, err = createUser(tx, userInfo)
			if err != nil {
				return User{}, fmt.Errorf("failed to create user: %v", err)
			}

			err = assignRoleToUser(tx, user.ID, role.Reader)
			if err != nil {
				return User{}, fmt.Errorf("failed to assign role to user: %v", err)
			}
			user.Roles = append(user.Roles, "read")
		} else {
			return User{}, fmt.Errorf("failed to find user by email: %v", err)
		}
	} else {
		user.Roles, err = getUserRoles(tx, user.ID)
		if err != nil {
			return User{}, fmt.Errorf("failed to get user roles: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return User{}, err
	}

	return user, nil
}

func findUserByEmail(tx *sql.Tx, email string) (User, error) {
	var user User
	err := tx.QueryRow("SELECT id, email, name FROM users WHERE email = $1;", email).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func createUser(tx *sql.Tx, userInfo oauth.GoogleUserInfo) (User, error) {
	var user User
	err := tx.QueryRow(
		`INSERT INTO users (email, name)
		VALUES ($1, $2) RETURNING id, email, name;`,
		userInfo.Email, userInfo.Name,
	).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func assignRoleToUser(tx *sql.Tx, userID, roleName string) error {
	_, err := tx.Exec(
		`INSERT INTO user_roles (user_id, role_name)
		VALUES ($1, $2);`,
		userID, roleName,
	)
	return err
}

func getUserRoles(tx *sql.Tx, userID string) ([]string, error) {
	rows, err := tx.Query(
		`SELECT role_name FROM user_roles WHERE user_id = $1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if errScanRow := rows.Scan(&role); errScanRow != nil {
			return nil, errScanRow
		}
		roles = append(roles, role)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
