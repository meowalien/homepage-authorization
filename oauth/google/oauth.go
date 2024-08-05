package google

import (
	"context"
	"golang.org/x/oauth2"
	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleOauthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	AuthUser    string `json:"authuser"`
	Prompt      string `json:"prompt"`
}

func GetUserInfo(user GoogleOauthToken) (*googleOauth2.Userinfo, error) {
	ctx := context.Background()
	oauth2Service, err := googleOauth2.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: user.AccessToken})))
	if err != nil {
		return nil, err
	}

	userInfoService := googleOauth2.NewUserinfoV2MeService(oauth2Service)
	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
