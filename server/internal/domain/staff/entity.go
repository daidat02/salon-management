package staff

import (
	"errors"
	"time"
)
const (
	statusActive = "active"
	statusInactive = "inactive"
)

type Staff struct {
	ID string
	OrganizationID string
	Code string
	FullName string
	Phone string
	Position string
	CommissionRate float64
	Status string
	HireDate time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (s *Staff) Validate() error {
	if s.Code == "" {
		return errors.New("Mã nhân viên không được để trống")
	}
	if s.FullName == "" {
		return errors.New("Tên đầy đủ không được để trống")
	}
	if s.Phone == "" {
		return errors.New("Số điện thoại không được để trống")
	}
	if s.Position == "" {
		return errors.New("Vị trí không được để trống")
	}
	if s.CommissionRate < 0 || s.CommissionRate > 100 {
		return errors.New("Tỷ lệ hoa hồng phải nằm trong khoảng từ 0 đến 100")
	}
	if s.Status != statusActive && s.Status != statusInactive {
		return errors.New("Trạng thái không hợp lệ")
	}
	return nil
}