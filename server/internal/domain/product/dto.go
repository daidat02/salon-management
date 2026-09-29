package product

import "time"

type CreateProductRequest struct {
	CategoryID    *string `json:"category_id"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Unit          string  `json:"unit"`
	NetUnit       *string  `json:"net_unit"`
	NetAmount     *float64 `json:"net_amount"`
	ProductType   string  `json:"product_type"`
	ShowOnWeb     *bool   `json:"show_on_web"`
	CostPrice     float64 `json:"cost_price"`
	SellPrice     float64 `json:"sell_price"`
	StockQuantity float64 `json:"stock_quantity"`
	MinStock      float64 `json:"min_stock"`
	IsActive      *bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	CategoryID    *string `json:"category_id"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Unit          string  `json:"unit"`
	NetUnit       *string  `json:"net_unit"`
	NetAmount     *float64 `json:"net_amount"`
	ProductType   string  `json:"product_type"`
	ShowOnWeb     *bool   `json:"show_on_web"`
	CostPrice     float64 `json:"cost_price"`
	SellPrice     float64 `json:"sell_price"`
	StockQuantity float64 `json:"stock_quantity"`
	MinStock      float64 `json:"min_stock"`
	IsActive      *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID            string     `json:"id"`
	OrganizationID string    `json:"organization_id"`
	CategoryID    *string    `json:"category_id"`
	SKU           string     `json:"sku"`
	Name          string     `json:"name"`
	Unit          string     `json:"unit"`
	NetUnit       *string     `json:"net_unit"`
	NetAmount     *float64    `json:"net_amount"`
	ProductType   string     `json:"product_type"`
	ShowOnWeb     bool       `json:"show_on_web"`
	CostPrice     float64    `json:"cost_price"`
	SellPrice     float64    `json:"sell_price"`
	StockQuantity float64    `json:"stock_quantity"`
	MinStock      float64    `json:"min_stock"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
