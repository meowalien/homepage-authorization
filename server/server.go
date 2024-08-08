package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func StartServer(handler http.Handler) *http.Server {
	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080
	}
	// Create an http.Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Run the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Infof("listen: %s\n", err)
		}
	}()
	return srv
}

func ShutdownServer(srv *http.Server) {
	logrus.Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}
	logrus.Info("Server exiting")
}
