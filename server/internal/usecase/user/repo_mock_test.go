package user

import (
	"context"

	userdomain "github.com/daidat02/server/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository giả lập lại interface của UserRepository (chỉ dùng trong test).
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *userdomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*userdomain.User, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdomain.User), args.Error(1)
}
