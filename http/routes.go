package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"homepage-authorization/http/handlers"
	"homepage-authorization/http/middleware"
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

	authGroup := r.Group("/auth")
	{
		loginGroup := authGroup.Group("/login")
		{
			loginGroup.POST("/google", handlers.GoogleOauthLogin())
		}

		authGroup.POST("/logout", handlers.CleanUpTokenCookie())
	}

	debugGroup := r.Group("/debug")
	{
		debugGroup.POST("/login", handlers.DebugLogin())
	}

	return r
}
