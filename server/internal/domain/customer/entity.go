package customer

import (
	"errors"
	"time"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

type Customer struct {
	ID             string
	OrganizationID string
	FullName       string
	Phone          string
	Gender         string
	BirthDate      *time.Time
	Note           string
	TotalSpent     float64
	TotalVisits    int
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
}

func (c *Customer) Validate() error {
	if c.FullName == "" {
		return errors.New("Tên khách hàng không được để trống")
	}
	if c.Phone == "" {
		return errors.New("Số điện thoại không được để trống")
	}
	if c.Gender != "" && c.Gender != GenderMale && c.Gender != GenderFemale && c.Gender != GenderOther {
		return errors.New("Giới tính không hợp lệ")
	}
	return nil
}
