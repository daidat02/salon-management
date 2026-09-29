package service

import (
	"errors"
	"time"
)

type Service struct {
	ID             string
	OrganizationID string
	CategoryID     *string
	Name           string
	Description    string
	DurationMinutes int
	BufferMinutes   int
	Price           float64
	IsActive        bool
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type ServiceMaterial struct {
    ID             string     `json:"id"`
    OrganizationID string     `json:"organization_id"`
    ServiceID      string     `json:"service_id"`
    ProductID      string     `json:"product_id"`
    Quantity       float64    `json:"quantity"`
    CreatedAt      *time.Time `json:"created_at"`
    UpdatedAt      *time.Time `json:"updated_at"`

    // Thêm các trường lấy từ bảng products (phải khớp 100% tên key trong json_build_object)
    ProductName    string     `json:"product_name"`
    ProductUnit    string     `json:"product_unit"`
    CostPrice      float64    `json:"cost_price"`
}

func (s *Service) Validate() error {
	if s.Name == "" {
		return errors.New("Tên dịch vụ không được để trống")
	}
	if s.DurationMinutes <= 0 {
		return errors.New("Thời lượng dịch vụ phải lớn hơn 0 phút")
	}
	if s.BufferMinutes < 0 {
		return errors.New("Thời gian đệm không được âm")
	}
	if s.Price < 0 {
		return errors.New("Giá dịch vụ không được âm")
	}
	return nil
}


func (sm *ServiceMaterial) Validate() error {
	if sm.Quantity <= 0 {
		return errors.New("Số lượng nguyên vật liệu phải lớn hơn 0")
	}
	return nil
}