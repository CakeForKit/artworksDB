package auth

import (
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth/token"
)

type AuthService interface {
}

func NewAuthService() (AuthService, error) {
	tokenMaker, err := token.NewPasetoMaker("")
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &authService{
		tokenMaker: tokenMaker,
	}

	return server, nil
}

type authService struct {
	tokenMaker token.Maker
}
