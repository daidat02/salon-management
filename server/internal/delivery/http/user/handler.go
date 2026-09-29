package user

import (
	"net/http"

	domain "github.com/daidat02/server/internal/domain/user"
	uc "github.com/daidat02/server/internal/usecase/user"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCase *uc.UserUsecase
}

func NewCreateUserHandler(useCase *uc.UserUsecase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}


// CreateUser godoc
// @Summary      Đăng ký người dùng
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateUserRequest  true  "Thông tin đăng ký"
// @Success      200   {object}  response.SuccessResponse  "Tạo người dùng thành công"
// @Failure      400   {object}  response.ErrorResponse    "Dữ liệu không hợp lệ / Số điện thoại đã tồn tại"
// @Failure      500   {object}  response.ErrorResponse    "Lỗi server"
// @Router       /register [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email: req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone: req.Phone,
		Role: req.Role,
		Status: req.Status,
	}

	createdUser, err := h.useCase.RegisterUser(c, user);
	if err !=nil{
		response.Error(c,err)
		return
	}
	response.Success(c,http.StatusOK,"Tạo người dùng thành công",createdUser)
}


func setRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*7, // 7 days
		"/",
		"",
		true,  // secure
		true,  // httpOnly
	)
}		
	

// Login godoc
// @Summary      Đăng nhập
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.LoginRequest       true  "Số điện thoại + mật khẩu"
// @Success      200   {object}  response.SuccessResponse  "Đăng nhập thành công (kèm access_token)"
// @Failure      400   {object}  response.ErrorResponse    "Dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse    "Mật khẩu không đúng"
// @Failure      404   {object}  response.ErrorResponse    "Số điện thoại không tồn tại"
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context){
	var req domain.LoginRequest
	if err:= c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user,accessToken, refreshToken, err := h.useCase.LoginUser(c, &req)
	if err != nil {
		response.Error(c,err)
		return
	}

	setRefreshTokenCookie(c, *refreshToken)

	response.Success(c,http.StatusOK,"Đăng nhập thành công",gin.H{
		"user": user,
		"access_token": accessToken,
	});
}


func (h *UserHandler) RefreshToken(c *gin.Context){
	refreshToken, err := c.Cookie("refresh_token")	
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accessToken, newrefreshToken, err := h.useCase.RefreshToken(c,refreshToken)
	if err != nil {
		response.Error(c,err)
		return
	}

	setRefreshTokenCookie(c, *newrefreshToken)
	response.Success(c,http.StatusOK,"Làm mới token thành công",gin.H{
		"access_token": *accessToken,
	});
}

