package service

import "time"

type ServiceMaterialInput struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"` // Định lượng tiêu hao (vd: 0.75 tuýp, 150 gram)
}

type CreateServiceRequest struct {
	CategoryID      *string `json:"category_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	BufferMinutes   int     `json:"buffer_minutes"`
	Price           float64 `json:"price"`
	IsActive        *bool   `json:"is_active"`
	Materials       []ServiceMaterialInput `json:"materials"`
}

type UpdateServiceRequest struct {
	CategoryID      *string `json:"category_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	BufferMinutes   int     `json:"buffer_minutes"`
	Price           float64 `json:"price"`
	IsActive        *bool   `json:"is_active"`
}

type ServiceResponse struct {
	ID              string     `json:"id"`
	OrganizationID  string     `json:"organization_id"`
	CategoryID      *string    `json:"category_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	DurationMinutes int        `json:"duration_minutes"`
	BufferMinutes   int        `json:"buffer_minutes"`
	Price           float64    `json:"price"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
