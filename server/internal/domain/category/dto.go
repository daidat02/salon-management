package category

import "time"

type CreateCategoryRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CategoryResponse struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
