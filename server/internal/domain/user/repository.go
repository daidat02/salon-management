package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context,user *User) error
	
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
}