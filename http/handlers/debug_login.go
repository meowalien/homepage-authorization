package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homepage-authorization/token"
	"homepage-authorization/user"
	"net/http"
)

type DebugLoginBody struct {
	UserID string `json:"user_id"`
}

func DebugLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body DebugLoginBody
		if err := c.BindJSON(&body); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		userStruct, err := user.GetUserByUserID(body.UserID)
		if err != nil {
			logrus.Errorf("Failed to get user ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}
		exp := token.DefaultExpireDuration
		claims := token.DefaultClaims(exp)
		signedToken, err := token.CreateTokenByUser(claims, userStruct)
		if err != nil {
			logrus.Errorf("Failed to create token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		_, err = token.VerifyToken(signedToken)
		if err != nil {
			logrus.Errorf("Failed to verify token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			return
		}

		logrus.Debugf("User %s logged in, token: %s", userStruct.Email, signedToken)

		// transform nanoseconds to seconds
		maxAge := int(exp / 1000000000)
		for _, cookie := range cookieConfig {
			c.SetCookie(cookie.Name, signedToken, maxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		c.Status(http.StatusOK)
	}
}
