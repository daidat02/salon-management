package staff

import "time"

type CreateStaffRequest struct {
	ID             string  `json:"id"`
	OrganizationId string  `json:"organization_id"`
	Code           string  `json:"code"`
	FullName       string  `json:"full_name"`
	Phone          string  `json:"phone"`
	Position       string  `json:"position"`
	CommissionRate float64 `json:"commission_rate"`
	Status         string  `json:"status"`
	HireDate       time.Time `json:"hire_date"`
}

type UpdateStaffRequest struct {
	Code           string    `json:"code"`
	FullName       string    `json:"full_name"`
	Phone          string    `json:"phone"`
	Position       string    `json:"position"`
	CommissionRate float64   `json:"commission_rate"`
	Status         string    `json:"status"`
	HireDate       time.Time `json:"hire_date"`
}

type StaffResponse struct {
	ID            string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	Code           string  `json:"code"`
	FullName       string  `json:"full_name"`
	Phone          string  `json:"phone"`
	Position       string  `json:"position"`
	CommissionRate float64 `json:"commission_rate"`
	Status         string  `json:"status"`
	HireDate       time.Time `json:"hire_date"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
