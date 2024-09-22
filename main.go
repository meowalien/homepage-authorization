package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"homepage-authorization/config"
	"homepage-authorization/http"
	"homepage-authorization/log"
	"homepage-authorization/oauth"
	"homepage-authorization/postgresql"
	"homepage-authorization/quit"
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

	r := http.SetupRouter()
	srv := http.StartServer(r)

	quit.WaitForQuitSignal()

	http.ShutdownServer(srv)

}
