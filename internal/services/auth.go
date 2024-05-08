package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrorInvalidToken               = errors.New("invalid token")
	ErrorDecodeUserInfo             = errors.New("failed to decode user info")
	ErrorMissingAuthorizationHeader = errors.New("required authorization header is missing")
	ErrorRevokeToken                = errors.New("failed to revoke token")
	ErrorExchangeToken              = errors.New("failed to exchange token")
	ErrorMissingToken               = errors.New("token is missing")
	ErrorGenerateAuthState          = errors.New("failed to generate auth state")
	ErrorAuthStateMissmatch         = errors.New("auth state missmatch")
	ErrorUnauthorized               = errors.New("user is not authorized")
)

type AuthService interface {
	GetAuthCodeUrl(state, verifier string) string
	ExchangeToken(ctx context.Context, code, verifier string) (*oauth2.Token, error)
	Token(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)
	ValidateToken(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error)
	RevokeToken(token string) error

	ErrorInvalidToken() error
	ErrorDecodeUserInfo() error
	ErrorMissingAuthorizationHeader() error
	ErrorGenerateAuthState() error
	ErrorAuthStateMissmatch() error
	ErrorUnauthorized() error
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

func (as *authService) GetAuthCodeUrl(state, verifier string) string {
	return as.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}

func (as *authService) ExchangeToken(ctx context.Context, code, verifier string) (*oauth2.Token, error) {
	token, err := as.oauth2.Exchange(
		ctx, code, oauth2.SetAuthURLParam("code_verifier", verifier), oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	if err != nil {
		return nil, ErrorExchangeToken
	}
	return token, nil
}

func (as *authService) Token(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSrc := as.oauth2.TokenSource(ctx, token)
	token, err := oauth2.ReuseTokenSource(token, tokenSrc).Token()
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (as *authService) ValidateToken(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error) {
	if token == nil {
		return &models.GoogleUserInfo{}, ErrorMissingToken
	}
	client := as.oauth2.Client(ctx, token)
	resp, err := client.Get(as.userInfoUrl)
	if err != nil {
		return &models.GoogleUserInfo{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &models.GoogleUserInfo{}, err
	}
	defer resp.Body.Close()

	userInfo := &models.GoogleUserInfo{}
	err = utils.DecodeResponseBody(resp.Body, &userInfo)
	if err != nil {
		return &models.GoogleUserInfo{}, ErrorDecodeUserInfo
	}
	return userInfo, nil

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
		return ErrorRevokeToken
	}
	return nil
}

func (as *authService) ErrorInvalidToken() error {
	return ErrorInvalidToken
}

func (as *authService) ErrorDecodeUserInfo() error {
	return ErrorDecodeUserInfo
}

func (as *authService) ErrorMissingAuthorizationHeader() error {
	return ErrorMissingAuthorizationHeader
}

func (as *authService) ErrorGenerateAuthState() error {
	return ErrorGenerateAuthState
}

func (as *authService) ErrorAuthStateMissmatch() error {
	return ErrorAuthStateMissmatch
}

func (as *authService) ErrorUnauthorized() error {
	return ErrorUnauthorized
}
