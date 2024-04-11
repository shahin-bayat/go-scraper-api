package services

import "github.com/shahin-bayat/scraper-api/internal/config"

type Services struct {
	AuthService *AuthService
}

func NewServices(appConfig *config.AppConfig) (*Services, error) {
	authService, err := NewAuthService(appConfig)
	if err != nil {
		return nil, err
	}

	return &Services{
		AuthService: authService,
	}, nil
}
