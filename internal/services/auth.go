package services

import (
	"github.com/shahin-bayat/scraper-api/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	Google *oauth2.Config
}

func NewAuthService(appConfig *config.AppConfig) (*AuthService, error) {
	// INFO: you can add more providers here

	return &AuthService{
		Google: &oauth2.Config{
			ClientID:     appConfig.GoogleClientID,
			ClientSecret: appConfig.GoogleClientSecret,
			RedirectURL:  appConfig.GoogleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
		},
	}, nil
}
