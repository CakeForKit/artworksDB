package testobj

import (
	"fmt"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/services/auth"
	"github.com/google/uuid"
)

type AuthUserRequestMother interface {
	RegisterDefault() auth.RegisterUserRequest
	Login(login, password string) auth.LoginUserRequest
	// RegisterDefaultP() *auth.RegisterUserRequest
}

func NewRegisterUserRequestMother() AuthUserRequestMother {
	return &authUserRequestMother{}
}

type authUserRequestMother struct{}

func (um *authUserRequestMother) RegisterDefault() auth.RegisterUserRequest {
	return auth.RegisterUserRequest{
		Username:       "test-username-" + uuid.NewString(),
		Login:          "test-login-" + uuid.NewString(),
		Password:       "test-password-" + uuid.NewString(),
		Email:          fmt.Sprintf("user%s@test.com", uuid.NewString()),
		SubscribeEmail: true,
	}
}

func (um *authUserRequestMother) Login(login, password string) auth.LoginUserRequest {
	return auth.LoginUserRequest{
		Login:    login,
		Password: password,
	}
}

// func (um *authUserRequestMother) RegisterDefaultP() *auth.RegisterUserRequest {
// 	return &auth.RegisterUserRequest{
// 		Username:       "test-username-" + uuid.NewString(),
// 		Login:          "test-login-" + uuid.NewString(),
// 		Password:       "test-password-" + uuid.NewString(),
// 		Email:          fmt.Sprintf("user%s@test.com", uuid.NewString()),
// 		SubscribeEmail: true,
// 	}
// }
