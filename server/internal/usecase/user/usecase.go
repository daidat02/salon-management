package user

import (
	"context"
	"errors"
	"time"

	userdomain "github.com/daidat02/server/internal/domain/user"
	"github.com/daidat02/server/pkg/apperror"
	utils "github.com/daidat02/server/pkg/utils"
)
type UserUsecase struct{
	userRepo userdomain.UserRepository
}


func NewCreateUserUsecase(userRepo userdomain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uc *UserUsecase) RegisterUser(c context.Context, user *userdomain.User) (*userdomain.UserResponse, error ) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	user.Status = userdomain.StatusActive
	if err := user.Validate(); err != nil {
		return nil, err
	}

	existingUser,err := uc.userRepo.FindUserByPhoneNumber(ctx, user.Phone)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return nil, apperror.New(500, "Lỗi khi kiểm tra số điện thoại", err)
	}
	if existingUser != nil {
		return nil, apperror.New(400, "Số điện thoại đã tồn tại", nil)
	}
	
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi băm mật khẩu", err)
	}

	user.Password = hashPassword
	if user.ID == "" {
		user.ID = utils.NewID()
	}

	err = uc.userRepo.CreateUser(ctx,user)
	if err != nil {
		if errors.Is(err,context.DeadlineExceeded){
			return nil, apperror.New(408, "Lỗi: Hết thời gian chờ khi tạo người dùng", err)
		}

		return nil, apperror.New(500, "Lỗi khi tạo người dùng", err)
	}
	
	response := &userdomain.UserResponse{
		ID: user.ID,
		OrganizationID: user.OrganizationID,
		Email: user.Email,
		FullName: user.FullName,
		PhoneNumber: user.Phone,
		Role: user.Role,
		Status: user.Status,
	}
	return response, nil
}


func (uc *UserUsecase)LoginUser(c context.Context, req *userdomain.LoginRequest) (*userdomain.UserResponse,*string,*string, error) {
	ctx, cancel := context.WithTimeout(c, 5 * time.Second)
	defer cancel()
	err := req.Validate()
	if err != nil {
		return nil, nil,nil, apperror.New(400, "Dữ liệu không hợp lệ", err)
	}

	existingUser, err := uc.userRepo.FindUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil,nil, apperror.New(408, "Lỗi: Hết thời gian chờ khi tìm kiếm người dùng", err)
		}
		return nil, nil, nil, apperror.New(500, "Lỗi khi tìm kiếm người dùng", err)
	}
	if existingUser == nil {
		return nil, nil,nil, apperror.New(404, "Số điện thoại không tồn tại", nil)
	}

	if !utils.CheckPasswordHash(req.Password,existingUser.Password){
		return nil, nil, nil, apperror.New(401, "Mật khẩu không đúng", nil)
	}

	acccessToken,err := utils.GenerateToken(existingUser.ID,existingUser.OrganizationID)
	if err != nil {
		return nil, nil,nil, apperror.New(500, "Lỗi khi tạo token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(existingUser.ID,existingUser.OrganizationID)
    if err != nil {
        return nil, nil, nil, apperror.New(500, "Lỗi khi tạo Refresh Token", err)
    }


	var lastLogin time.Time
    if existingUser.LastLoginAt != nil {
        lastLogin = *existingUser.LastLoginAt
    }

	return &userdomain.UserResponse{
		ID: existingUser.ID,
		OrganizationID: existingUser.OrganizationID,
		Email: existingUser.Email,
		FullName: existingUser.FullName,
		PhoneNumber: existingUser.Phone,
		Role: existingUser.Role,
		Status: existingUser.Status,
		LastLoginAt: lastLogin,
		CreatedAt: *existingUser.CreatedAt,
		UpdatedAt: *existingUser.UpdatedAt,
	},&acccessToken, &refreshToken, nil
}

func (uc *UserUsecase) RefreshToken(c context.Context, refreshToken string) (*string,*string, error) {

	refreshClaims, err := utils.ValidateToken(refreshToken)
	if err !=nil{
		return nil,nil, apperror.New(401, "Refresh Token không hợp lệ", err)
	}

	accessToken, err := utils.GenerateToken(refreshClaims.UserID, refreshClaims.OrgID)
	if err != nil {
		return nil,nil, apperror.New(500, "Lỗi khi tạo Access Token", err)
	}
	
	newRefreshToken, err := utils.GenerateRefreshToken(refreshClaims.UserID, refreshClaims.OrgID)
	if err != nil {
		return nil,nil, apperror.New(500, "Lỗi khi tạo Refresh Token", err)
	}

	return &accessToken,&newRefreshToken, nil
}