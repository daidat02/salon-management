package user

import (
	"net/http"
	"time"

	apperror "github.com/daidat02/server/pkg/apperror"
)
const (
	StatusActive = "active"
	StatusInactive = "inactive"
	// Roles
	RoleOwner = "owner"
	RoleManager = "manager"
	RoleReceptionist = "receptionist"
)

type User struct{
	ID string
	OrganizationID string
	Email string
	Password string
	FullName string
	Phone string
	Role string
	Status string
	LastLoginAt *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
}


func (user *User) Validate( ) error {
	if user.Password == "" {
		return apperror.New(http.StatusBadRequest, "Mật khẩu không được để trống", nil)
	}
	if user.FullName == "" {
		return apperror.New(http.StatusBadRequest, "Tên đầy đủ không được để trống", nil)
	}
	if user.Phone == "" {
		return apperror.New(http.StatusBadRequest, "Số điện thoại không được để trống", nil)
	}
	if role := user.Role; role != RoleOwner && role != RoleManager && role != RoleReceptionist {
		return apperror.New(http.StatusBadRequest, "Lỗi :Vai trò không hợp lệ", nil)
	}
	if status := user.Status; status != StatusActive && status != StatusInactive {
		return apperror.New(http.StatusBadRequest, "Lỗi: Trạng thái không hợp lệ", nil)
	}
	return nil
}