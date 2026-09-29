package user

import (
	"errors"
	"strings"
	"time"
)
type CreateUserRequest struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone string `json:"phone_number"`
	Role string `json:"role"`
	Status string `json:"status"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
    PhoneNumber    string `json:"phone"`
    Password string `json:"password"`
}

type UserResponse struct {
    ID          string `json:"id"`
    OrganizationID string `json:"organization_id"`
    Email       string `json:"email,omitempty"`
    FullName    string `json:"full_name,omitempty"`
    PhoneNumber string `json:"phone_number"`
    Role        string `json:"role"`
    Status      string `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`   
    LastLoginAt time.Time `json:"last_login_at"`
}



func (req *LoginRequest) Validate() error {
    if strings.TrimSpace(req.PhoneNumber) == "" {
        return errors.New("lỗi: số điện thoại không được để trống")
    }
    if req.Password == "" {
        return errors.New("lỗi: mật khẩu không được để trống")
    }
    return nil
}