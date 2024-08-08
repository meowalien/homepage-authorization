package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"homepage-authorization/config"
	"homepage-authorization/log"
	"homepage-authorization/oauth"
	"homepage-authorization/postgresql"
	"homepage-authorization/quit"
	"homepage-authorization/routes"
	"homepage-authorization/server"
	"homepage-authorization/token"
)

func main() {
	defer logrus.Info("Main exiting")
	config.InitConfig()
	log.InitLogger()
	postgresql.ConnectDB()
	defer postgresql.DisconnectDB()
	token.InitSignKey(viper.GetString("token.privateKeyPath"))
	token.InitVerifyKey(viper.GetString("token.publicKeyPath"))
	oauth.Init()

	r := routes.SetupRouter()
	srv := server.StartServer(r)

	quit.WaitForQuitSignal()

	server.ShutdownServer(srv)

}
