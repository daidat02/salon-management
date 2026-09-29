package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AppError chứa mã HTTP và thông báo lỗi.
// Quy ước chống lồng: mỗi lỗi HTTP chỉ được New 1 lần (ở usecase).
// Các tầng khác trả lỗi thường (fmt.Errorf/sentinel), không New tiếp.
type AppError struct {
	StatusCode int
	Message    string
	Err        error // (Tùy chọn) Lưu lỗi gốc để debug, không trả ra client
}

// Bắt buộc phải có hàm Error() để struct này được Go công nhận là kiểu 'error'
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap để errors.As/Is xuyên qua được lớp bọc ngoài.
func (e *AppError) Unwrap() error {
	return e.Err
}

// MarshalJSON chỉ trả code/message phẳng, không lộ cây Err lồng nhau.
func (e *AppError) MarshalJSON() ([]byte, error) {
	type flat struct {
		StatusCode int    `json:"StatusCode"`
		Message    string `json:"Message"`
	}
	return json.Marshal(flat{
		StatusCode: e.StatusCode,
		Message:    e.Message,
	})
}

// Hàm tiện ích để tạo lỗi nhanh
func New(statusCode int, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// From trả về *AppError sẵn có nếu err đã là AppError (propagate, không bọc lại).
// Ngược lại tạo mới với fallbackStatus/fallbackMsg.
func From(err error, fallbackStatus int, fallbackMsg string) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return New(fallbackStatus, fallbackMsg, err)
}