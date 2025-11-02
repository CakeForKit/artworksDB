package testobj

import (
	"fmt"

	authmodels "github.com/CakeForKit/artworksDB.git/internal/services/auth/auth_models"
	"github.com/google/uuid"
)

type AuthUserRequestMother interface {
	RegisterDefault() authmodels.RegisterUserRequest
	Login(login, password string) authmodels.LoginUserRequest
	// RegisterDefaultP() *authmodels.RegisterUserRequest
}

func NewRegisterUserRequestMother() AuthUserRequestMother {
	return &authUserRequestMother{}
}

type authUserRequestMother struct{}

func (um *authUserRequestMother) RegisterDefault() authmodels.RegisterUserRequest {
	return authmodels.RegisterUserRequest{
		Username:       "test-username-" + uuid.NewString(),
		Login:          "test-login-" + uuid.NewString(),
		Password:       "test-password-" + uuid.NewString(),
		Email:          fmt.Sprintf("user%s@test.com", uuid.NewString()),
		SubscribeEmail: true,
	}
}

func (um *authUserRequestMother) Login(login, password string) authmodels.LoginUserRequest {
	return authmodels.LoginUserRequest{
		Login:    login,
		Password: password,
	}
}

// func (um *authUserRequestMother) RegisterDefaultP() *authmodels.RegisterUserRequest {
// 	return &authmodels.RegisterUserRequest{
// 		Username:       "test-username-" + uuid.NewString(),
// 		Login:          "test-login-" + uuid.NewString(),
// 		Password:       "test-password-" + uuid.NewString(),
// 		Email:          fmt.Sprintf("user%s@test.com", uuid.NewString()),
// 		SubscribeEmail: true,
// 	}
// }
