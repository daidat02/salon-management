package category

import (
	"errors"
	"time"
)

const (
	TypeService = "service"
	TypeProduct = "product"
)

type Category struct {
	ID             string
	OrganizationID string
	Type           string
	Name           string
	Description    string
	IsActive       bool
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func (c *Category) Validate() error {
	if c.Type != TypeService && c.Type != TypeProduct {
		return errors.New("Loại danh mục không hợp lệ (service | product)")
	}
	if c.Name == "" {
		return errors.New("Tên danh mục không được để trống")
	}
	return nil
}
