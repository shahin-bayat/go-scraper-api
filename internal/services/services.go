package services

import "github.com/shahin-bayat/scraper-api/internal/config"

type Services struct {
	AuthService AuthService
}

func NewServices(appConfig *config.AppConfig) (*Services, error) {
	authService := NewAuthService(appConfig)

	return &Services{
		AuthService: authService,
	}, nil
}
