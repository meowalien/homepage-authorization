package oauth

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/people/v1"
	"log"
	"os"
	"time"
)

type GoogleOauth2Credential struct {
	Credential string `json:"credential"`
	ClientID   string `json:"clientId"`
	SelectBy   string `json:"select_by"`
}

type GoogleUserInfo struct {
	Aud           string  `json:"aud"`
	Azp           string  `json:"azp"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	Exp           float64 `json:"exp"`
	GivenName     string  `json:"given_name"`
	Iat           float64 `json:"iat"`
	Iss           string  `json:"iss"`
	Jti           string  `json:"jti"`
	Name          string  `json:"name"`
	Nbf           float64 `json:"nbf"`
	Picture       string  `json:"picture"`
	Sub           string  `json:"sub"`
}

var clientCredentialsBytes []byte

func Init() {
	var err error
	clientCredentialsBytes, err = os.ReadFile(viper.GetString("oauth.oauthCredentials"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
}

func GetUserInfo(user GoogleOauth2Credential) (userInfo GoogleUserInfo, err error) {

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(clientCredentialsBytes, people.UserinfoEmailScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	credentialJWT := user.Credential
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	validator, err := idtoken.NewValidator(ctx, idtoken.WithCredentialsJSON(clientCredentialsBytes))
	if err != nil {
		return userInfo, err
	}
	payload, err := validator.Validate(ctx, credentialJWT, config.ClientID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return userInfo, fmt.Errorf("timeout during token validation")
		}
		return userInfo, err
	}

	if payload.Claims["email_verified"] != true {
		return userInfo, fmt.Errorf("email is not verified")
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &userInfo,
	})
	if err != nil {
		return userInfo, err
	}
	err = decoder.Decode(payload.Claims)
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}
