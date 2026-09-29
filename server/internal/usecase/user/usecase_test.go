package user

import (
	"context"
	"errors"
	"testing"
	"time"

	userdomain "github.com/daidat02/server/internal/domain/user"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// assertFlatAppError kiểm tra err là *AppError đúng 1 lớp (không lồng nhau).
func assertFlatAppError(t *testing.T, err error, expectedCode int) {
	t.Helper()
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr), "Lỗi phải là *apperror.AppError")
	if appErr != nil {
		assert.Equal(t, expectedCode, appErr.StatusCode, "Sai StatusCode")
		var inner *apperror.AppError
		assert.False(t, errors.As(appErr.Err, &inner), "Lỗi bị lồng 2 lớp AppError")
	}
}

// validRegisterUser trả về input đầy đủ, pass hết Validate().
// Lưu ý: usecase tự set Status=active trước khi validate nên không cần set ở đây.
func validRegisterUser() *userdomain.User {
	return &userdomain.User{
		OrganizationID: "org-123",
		Email:          "owner@salon.test",
		Password:       "MatKhau123",
		FullName:       "Chủ Salon",
		Phone:          "0901234567",
		Role:           userdomain.RoleOwner,
	}
}

func TestUserUsecase_RegisterUser(t *testing.T) {
	now := time.Now()
	existing := &userdomain.User{
		ID:             "user-exists",
		OrganizationID: "org-123",
		Email:          "old@salon.test",
		FullName:       "Người Cũ",
		Phone:          "0901234567",
		Role:           userdomain.RoleOwner,
		Status:         userdomain.StatusActive,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}

	tests := []struct {
		name          string
		inputUser     *userdomain.User
		findUser      *userdomain.User
		findErr       error
		createErr     error
		expectFind    bool
		expectCreate  bool
		expectedError bool
		expectedCode  int
	}{
		{
			name:          "Thành công - Đăng ký user mới hợp lệ",
			inputUser:     validRegisterUser(),
			findUser:      nil,
			findErr:       nil,
			createErr:     nil,
			expectFind:    true,
			expectCreate:  true,
			expectedError: false,
		},
		{
			name: "Thất bại - Validate (thiếu mật khẩu)",
			inputUser: func() *userdomain.User {
				u := validRegisterUser()
				u.Password = ""
				return u
			}(),
			findUser:      nil,
			findErr:       nil,
			createErr:     nil,
			expectFind:    false,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu tên đầy đủ)",
			inputUser: func() *userdomain.User {
				u := validRegisterUser()
				u.FullName = ""
				return u
			}(),
			findUser:      nil,
			findErr:       nil,
			createErr:     nil,
			expectFind:    false,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu số điện thoại)",
			inputUser: func() *userdomain.User {
				u := validRegisterUser()
				u.Phone = ""
				return u
			}(),
			findUser:      nil,
			findErr:       nil,
			createErr:     nil,
			expectFind:    false,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (vai trò không hợp lệ)",
			inputUser: func() *userdomain.User {
				u := validRegisterUser()
				u.Role = "superadmin"
				return u
			}(),
			findUser:      nil,
			findErr:       nil,
			createErr:     nil,
			expectFind:    false,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Số điện thoại đã tồn tại",
			inputUser:     validRegisterUser(),
			findUser:      existing,
			findErr:       nil,
			createErr:     nil,
			expectFind:    true,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Lỗi khi kiểm tra số điện thoại",
			inputUser:     validRegisterUser(),
			findUser:      nil,
			findErr:       errors.New("db connection refused"),
			createErr:     nil,
			expectFind:    true,
			expectCreate:  false,
			expectedError: true,
			expectedCode:  500,
		},
		{
			name:          "Thành công - Find timeout được coi như chưa tồn tại",
			inputUser:     validRegisterUser(),
			findUser:      nil,
			findErr:       context.DeadlineExceeded,
			createErr:     nil,
			expectFind:    true,
			expectCreate:  true,
			expectedError: false,
		},
		{
			name:          "Thất bại - Lỗi khi tạo người dùng",
			inputUser:     validRegisterUser(),
			findUser:      nil,
			findErr:       nil,
			createErr:     errors.New("db insert failed"),
			expectFind:    true,
			expectCreate:  true,
			expectedError: true,
			expectedCode:  500,
		},
		{
			name:          "Thất bại - Hết thời gian chờ khi tạo người dùng",
			inputUser:     validRegisterUser(),
			findUser:      nil,
			findErr:       nil,
			createErr:     context.DeadlineExceeded,
			expectFind:    true,
			expectCreate:  true,
			expectedError: true,
			expectedCode:  408,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			if tt.expectFind {
				mockRepo.On("FindUserByPhoneNumber", mock.Anything, tt.inputUser.Phone).
					Return(tt.findUser, tt.findErr)
			}
			if tt.expectCreate {
				mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(tt.createErr)
			}

			uc := NewCreateUserUsecase(mockRepo)
			plainPassword := tt.inputUser.Password
			res, err := uc.RegisterUser(context.Background(), tt.inputUser)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.ID, "ID phải được sinh ra")
				assert.Equal(t, tt.inputUser.OrganizationID, res.OrganizationID)
				assert.Equal(t, tt.inputUser.Email, res.Email)
				assert.Equal(t, tt.inputUser.FullName, res.FullName)
				assert.Equal(t, tt.inputUser.Phone, res.PhoneNumber)
				assert.Equal(t, tt.inputUser.Role, res.Role)
				assert.Equal(t, userdomain.StatusActive, res.Status)
				// Mật khẩu phải được băm, không lưu plaintext
				assert.NotEqual(t, plainPassword, tt.inputUser.Password, "Mật khẩu phải được băm")
				assert.True(t, utils.CheckPasswordHash(plainPassword, tt.inputUser.Password))
			}

			if tt.expectFind {
				mockRepo.AssertExpectations(t)
			} else {
				mockRepo.AssertNotCalled(t, "FindUserByPhoneNumber", mock.Anything, mock.Anything)
				mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
			}
			if !tt.expectCreate {
				mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestUserUsecase_LoginUser(t *testing.T) {
	utils.InitJWT("test-secret-key", time.Hour, 24*time.Hour)

	now := time.Now()
	hash, err := utils.HashPassword("MatKhau123")
	assert.NoError(t, err)
	existing := &userdomain.User{
		ID:             "user-1",
		OrganizationID: "org-123",
		Email:          "owner@salon.test",
		Password:       hash,
		FullName:       "Chủ Salon",
		Phone:          "0901234567",
		Role:           userdomain.RoleOwner,
		Status:         userdomain.StatusActive,
		LastLoginAt:    &now,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}

	tests := []struct {
		name          string
		inputReq      *userdomain.LoginRequest
		findUser      *userdomain.User
		findErr       error
		expectFind    bool
		expectedError bool
		expectedCode  int
	}{
		{
			name: "Thành công - Đăng nhập đúng số điện thoại + mật khẩu",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0901234567",
				Password:    "MatKhau123",
			},
			findUser:      existing,
			findErr:       nil,
			expectFind:    true,
			expectedError: false,
		},
		{
			name: "Thất bại - Validate (thiếu số điện thoại)",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "",
				Password:    "MatKhau123",
			},
			findUser:      nil,
			findErr:       nil,
			expectFind:    false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu mật khẩu)",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0901234567",
				Password:    "",
			},
			findUser:      nil,
			findErr:       nil,
			expectFind:    false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Số điện thoại không tồn tại",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0999999999",
				Password:    "MatKhau123",
			},
			findUser:      nil,
			findErr:       nil,
			expectFind:    true,
			expectedError: true,
			expectedCode:  404,
		},
		{
			name: "Thất bại - Mật khẩu không đúng",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0901234567",
				Password:    "SaiMatKhau",
			},
			findUser:      existing,
			findErr:       nil,
			expectFind:    true,
			expectedError: true,
			expectedCode:  401,
		},
		{
			name: "Thất bại - Lỗi khi tìm kiếm người dùng",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0901234567",
				Password:    "MatKhau123",
			},
			findUser:      nil,
			findErr:       errors.New("db connection refused"),
			expectFind:    true,
			expectedError: true,
			expectedCode:  500,
		},
		{
			name: "Thất bại - Hết thời gian chờ khi tìm kiếm",
			inputReq: &userdomain.LoginRequest{
				PhoneNumber: "0901234567",
				Password:    "MatKhau123",
			},
			findUser:      nil,
			findErr:       context.DeadlineExceeded,
			expectFind:    true,
			expectedError: true,
			expectedCode:  408,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			if tt.expectFind {
				mockRepo.On("FindUserByPhoneNumber", mock.Anything, tt.inputReq.PhoneNumber).
					Return(tt.findUser, tt.findErr)
			}

			uc := NewCreateUserUsecase(mockRepo)
			res, accessToken, refreshToken, err := uc.LoginUser(context.Background(), tt.inputReq)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Nil(t, accessToken)
				assert.Nil(t, refreshToken)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, existing.ID, res.ID)
				assert.Equal(t, existing.OrganizationID, res.OrganizationID)
				assert.Equal(t, existing.Email, res.Email)
				assert.Equal(t, existing.FullName, res.FullName)
				assert.Equal(t, existing.Phone, res.PhoneNumber)
				assert.Equal(t, existing.Role, res.Role)
				assert.Equal(t, existing.Status, res.Status)
				assert.NotNil(t, accessToken)
				assert.NotNil(t, refreshToken)
				assert.NotEmpty(t, *accessToken, "Access token phải được cấp")
				assert.NotEmpty(t, *refreshToken, "Refresh token phải được cấp")
			}

			if tt.expectFind {
				mockRepo.AssertExpectations(t)
			} else {
				mockRepo.AssertNotCalled(t, "FindUserByPhoneNumber", mock.Anything, mock.Anything)
			}
		})
	}
}
