package order

import "time"


type ItemMaterialInput struct {
	ID             string `json:"id"`
	OrderItemID    string `json:"order_item_id"`
	ProductID      string `json:"product_id"`
	Quantity       float64 `json:"quantity"`
	UnitCostAtTime float64 `json:"unit_cost_at_time"`
}

type OrderItemInput struct {
	ItemType  string  `json:"item_type"`
	ServiceID *string `json:"service_id"`
	ProductID *string `json:"product_id"`
	StaffID   *string `json:"staff_id"`
	Quantity  float64 `json:"quantity"`
	Materials []ItemMaterialInput `json:"materials"`
}

type CreateOrderRequest struct {
	CustomerID     *string          `json:"customer_id"`
	AppointmentID  *string          `json:"appointment_id"`
	DiscountAmount float64          `json:"discount_amount"`
	Items          []OrderItemInput `json:"items"`
}

type CreatePaymentRequest struct {
	Method        string  `json:"method"`
	Provider      *string `json:"provider"`
	ProviderTxnID *string `json:"provider_txn_id"`
	Amount        float64 `json:"amount"`
}

type OrderItemResponse struct {
	ID        string   `json:"id"`
	ItemType  string   `json:"item_type"`
	ServiceID *string  `json:"service_id"`
	ProductID *string  `json:"product_id"`
	StaffID   *string  `json:"staff_id"`
	Quantity  float64  `json:"quantity"`
	UnitPrice float64  `json:"unit_price"`
	LineTotal float64  `json:"line_total"`
}

type OrderMaterialResponse struct {
	ID             string  `json:"id"`
	OrderItemID    string  `json:"order_item_id"`
	ProductID      string  `json:"product_id"`
	Quantity       float64 `json:"quantity"`
	UnitCostAtTime float64 `json:"unit_cost_at_time"`
}

type PaymentResponse struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"order_id"`
	Method        string     `json:"method"`
	Provider      *string    `json:"provider"`
	ProviderTxnID *string    `json:"provider_txn_id"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     *time.Time `json:"created_at"`
}

type OrderResponse struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	Code           string     `json:"code"`
	CustomerID     *string    `json:"customer_id"`
	AppointmentID  *string    `json:"appointment_id"`
	SubtotalAmount float64    `json:"subtotal_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	TotalAmount    float64    `json:"total_amount"`
	Status         string     `json:"status"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type OrderDetailResponse struct {
	Order     *OrderResponse          `json:"order"`
	Items     []*OrderItemResponse    `json:"items"`
	Materials []*OrderMaterialResponse `json:"materials"`
	Payments  []*PaymentResponse      `json:"payments"`
}

// OrderServeResponse kết quả bắt đầu phục vụ kèm phiếu xuất vừa sinh.
type OrderServeResponse struct {
	Order        *OrderResponse `json:"order"`
	DocumentID   string         `json:"document_id"`
	DocumentCode string         `json:"document_code"`
}

type PosItemResponse struct{
	ID        string  `json:"id" db:"id"`
	OrganizationID string  `json:"organization_id" db:"organization_id"`
	CategoryID *string `json:"category_id" db:"category_id"`
	Type      string  `json:"type" db:"type"`
	Name	  string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	DurationMinutes *int `json:"duration_minutes" db:"duration_minutes"`
	Unit *string `json:"unit" db:"unit"`
	Price    float64 `json:"price" db:"price"`
	StockQuantity *float64 `json:"stock_quantity" db:"stock_quantity"`
	IsActive  bool    `json:"is_active" db:"is_active"`
}