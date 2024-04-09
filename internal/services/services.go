package services

type Services struct {
	AuthService *AuthService
}

func NewServices() (*Services, error) {
	authService, err := NewAuthService()
	if err != nil {
		return nil, err
	}

	return &Services{
		AuthService: authService,
	}, nil
}
