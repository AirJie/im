package auth

import "context"

type authService struct {
}

var AuthService = new(authService)

func (a *authService) SignIn(ctx context.Context, appId, userId, deviceId int64, token string, addr string) error {
	err :=
}

func (a *authService)VerifyToken() error {
	app, err := App
}
