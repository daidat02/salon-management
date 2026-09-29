package customer

import "time"

type CreateCustomerRequest struct {
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	Gender    string     `json:"gender"`
	BirthDate *string `json:"birth_date"`
	Note      string     `json:"note"`
}

type UpdateCustomerRequest struct {
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	Gender    string     `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
	Note      string     `json:"note"`
}

type CustomerResponse struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	FullName       string     `json:"full_name"`
	Phone          string     `json:"phone"`
	Gender         string     `json:"gender"`
	BirthDate      *time.Time `json:"birth_date"`
	Note           string     `json:"note"`
	TotalSpent     float64    `json:"total_spent"`
	TotalVisits    int        `json:"total_visits"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
