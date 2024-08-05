package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homepage-authorization/routes/handlers"
	"homepage-authorization/routes/middleware"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Configure CORS
	r.Use(middleware.Cors())

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	i18nGroup := r.Group("/login")
	oauthGroup := i18nGroup.Group("/oauth")
	{
		oauthGroup.POST("/google", handlers.GoogleOauthLogin)
	}

	return r
}
