package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     viper.GetStringSlice("cors.allowedOrigins"),
		AllowMethods:     viper.GetStringSlice("cors.allowedMethods"),
		AllowHeaders:     viper.GetStringSlice("cors.allowedHeaders"),
		AllowCredentials: viper.GetBool("cors.allowCredentials"),
		MaxAge:           viper.GetDuration("cors.maxAge"),
	})
}
