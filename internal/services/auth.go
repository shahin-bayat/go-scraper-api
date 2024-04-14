package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService interface {
	GetUserInfo(ctx context.Context, token *oauth2.Token) (models.GoogleUserInfo, error)
	GetAuthCodeUrl(state, verifier string) string
	ExchangeToken(ctx context.Context, code, verifier string) (*oauth2.Token, error)
	TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
	TokenValid(token *oauth2.Token) bool
	RevokeToken(token string) error
}

type authService struct {
	oauth2               *oauth2.Config
	userInfoUrl          string
	googleTokenRevokeURL string
}

func NewAuthService(appConfig *config.AppConfig) AuthService {
	// INFO: you can add more providers here

	return &authService{
		oauth2: &oauth2.Config{
			ClientID:     appConfig.GoogleClientID,
			ClientSecret: appConfig.GoogleClientSecret,
			RedirectURL:  appConfig.GoogleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
		},
		userInfoUrl:          appConfig.GoogleUserInfoURL,
		googleTokenRevokeURL: appConfig.GoogleTokenRevokeURL,
	}
}

func (as *authService) GetUserInfo(ctx context.Context, token *oauth2.Token) (models.GoogleUserInfo, error) {
	client := as.oauth2.Client(ctx, token)
	resp, err := client.Get(as.userInfoUrl)
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	userInfo := models.GoogleUserInfo{}
	err = utils.DecodeResponseBody(resp.Body, &userInfo)
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to get user info: %w", err)
	}
	return userInfo, nil
}

func (as *authService) GetAuthCodeUrl(state, verifier string) string {
	return as.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}

func (as *authService) ExchangeToken(ctx context.Context, code, verifier string) (*oauth2.Token, error) {
	token, err := as.oauth2.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", verifier), oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	return token, nil
}

func (as *authService) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return as.oauth2.TokenSource(ctx, token)
}

func (as *authService) TokenValid(token *oauth2.Token) bool {
	if token == nil {
		return false
	}
	return token.Valid() && token.Expiry.After(time.Now())
}

func (as *authService) RevokeToken(token string) error {
	url := fmt.Sprintf("%s?token=%s", as.googleTokenRevokeURL, token)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to revoke token, status code: %d", resp.StatusCode)
	}
	return nil
}
