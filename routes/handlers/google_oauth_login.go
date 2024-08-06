package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homepage-authorization/oauth"
	"homepage-authorization/token"
	"homepage-authorization/user"
	"net/http"
)

type UserClaims struct {
	jwt.StandardClaims
	User string `json:"user_id"`
}

func GoogleOauthLogin(c *gin.Context) {

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
	//fmt.Println("userInfo: ", userInfo)

	userID, err := user.GetUserIDByUserInfo(userInfo)
	if err != nil {
		logrus.Errorf("Failed to get user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	signedToken, err := token.CreateTokenByUserID(userID)
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

	//c.SetSameSite(http.SameSiteNoneMode)
	//c.SetSameSite(http.SameSiteLaxMode)
	// Set the cookie
	//c.SetCookie("token", signedToken, 3600, "/", ".meowalien.com", false, true)
	c.SetCookie("token", signedToken, 3600, "/", ".meowalien.com", true, true)
	//c.SetCookie("token", signedToken, 3600, "/", "api.meowalien.com", true, true)
	//c.SetCookie("token", signedToken, 3600, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}
