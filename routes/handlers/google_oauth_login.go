package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"homepage-authorization/config"
	"homepage-authorization/oauth"
	"homepage-authorization/token"
	"homepage-authorization/user"
	"net/http"
)

type UserClaims struct {
	jwt.StandardClaims
	User string `json:"user_id"`
}

var cookieConfig []config.CookieConfig

func initCookieConfig() {
	if cookieConfig == nil {
		if err := viper.UnmarshalKey("cookies", &cookieConfig); err != nil {
			logrus.Fatalf("Failed to unmarshal cookies: %v", err)
		}
	}
}

func GoogleOauthLogin() gin.HandlerFunc {
	initCookieConfig()

	return func(c *gin.Context) {
		var googleOauthToken oauth.GoogleOauth2Credential
		if err := c.BindJSON(&googleOauthToken); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		userInfo, err := oauth.GetUserInfo(googleOauthToken)
		if err != nil {
			logrus.Errorf("Failed to get user info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}

		userStruct, err := user.GetUserByUserInfo(userInfo)
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

func CleanUpTokenCookie() gin.HandlerFunc {
	initCookieConfig()

	return func(c *gin.Context) {
		for _, cookie := range cookieConfig {
			c.SetCookie(cookie.Name, "", -1, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		c.Status(http.StatusOK)
	}
}
