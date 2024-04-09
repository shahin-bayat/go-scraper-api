package services

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	Google *oauth2.Config
}

func NewAuthService() (*AuthService, error) {
	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if os.Getenv("GOOGLE_REDIRECT_URL") == "" {
		return nil, fmt.Errorf("GOOGLE_REDIRECT_URL is required")
	}

	// INFO: you can add more providers here

	return &AuthService{
		Google: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		},
	}, nil
}
