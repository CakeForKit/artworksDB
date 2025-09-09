package hasher

import "github.com/stretchr/testify/mock"

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CheckPassword(password string, hashedPassword string) error {
	args := m.Called(password, hashedPassword)
	return args.Error(0)
}
