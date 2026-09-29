package response

import (
	"errors"
	"net/http"

	"github.com/daidat02/server/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}


func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := SuccessResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, resp)
}

func Error(c *gin.Context,err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		// Chỉ trả code/message phẳng, Err gốc chỉ log server-side để debug.
		c.JSON(appErr.StatusCode, ErrorResponse{
			Code:    appErr.StatusCode,
			Message: appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Lỗi hệ thống nội bộ, vui lòng thử lại sau",
	})
}
