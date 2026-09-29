package order

import (
	"errors"
	"time"
)

const (
	StatusDraft          = "draft"
	StatusPendingPayment = "pending_payment"
	StatusPaid           = "paid"
	StatusServing        = "serving"
	StatusCompleted      = "completed"
	StatusCancelled      = "cancelled"
	StatusRefunded       = "refunded"
)

const (
	ItemTypeService = "service"
	ItemTypeProduct = "product"
)

const (
	PaymentMethodCash         = "cash"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodGateway      = "gateway"
	PaymentMethodCard         = "card"
)

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

type Order struct {
	ID             string
	OrganizationID string
	Code           string
	CustomerID     *string
	AppointmentID  *string
	SubtotalAmount float64
	DiscountAmount float64
	TotalAmount    float64
	Status         string
	CreatedBy      *string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func (o *Order) Validate() error {
	switch o.Status {
	case StatusDraft, StatusPendingPayment, StatusPaid, StatusServing, StatusCompleted, StatusCancelled, StatusRefunded:
	default:
		return errors.New("Trạng thái đơn hàng không hợp lệ")
	}
	if o.SubtotalAmount < 0 {
		return errors.New("Tổng tiền trước giảm giá không được âm")
	}
	if o.DiscountAmount < 0 {
		return errors.New("Giảm giá không được âm")
	}
	if o.DiscountAmount > o.SubtotalAmount {
		return errors.New("Giảm giá không được vượt tổng tiền")
	}
	if o.TotalAmount < 0 {
		return errors.New("Tổng tiền thanh toán không được âm")
	}
	return nil
}

type OrderItem struct {
	ID             string
	OrganizationID string
	OrderID        string
	ItemType       string
	ServiceID      *string
	ProductID      *string
	StaffID        *string
	Quantity       float64
	UnitPrice      float64
	LineTotal      float64
}

func (it *OrderItem) Validate() error {
	if it.Quantity <= 0 {
		return errors.New("Số lượng dòng hàng phải lớn hơn 0")
	}
	if it.UnitPrice < 0 {
		return errors.New("Đơn giá dòng hàng không được âm")
	}
	switch it.ItemType {
	case ItemTypeService:
		if it.ServiceID == nil || *it.ServiceID == "" {
			return errors.New("Dòng dịch vụ thiếu service_id")
		}
		if it.ProductID != nil {
			return errors.New("Dòng dịch vụ không được kèm product_id")
		}
	case ItemTypeProduct:
		if it.ProductID == nil || *it.ProductID == "" {
			return errors.New("Dòng sản phẩm thiếu product_id")
		}
		if it.ServiceID != nil {
			return errors.New("Dòng sản phẩm không được kèm service_id")
		}
	default:
		return errors.New("Loại dòng hàng không hợp lệ (service | product)")
	}
	return nil
}

// OrderItemMaterial snapshot vật tư tiêu hao + giá vốn tại thời điểm bán.
type OrderItemMaterial struct {
	ID             string
	OrderItemID    string
	ProductID      string
	Quantity       float64
	UnitCostAtTime float64
}

func (m *OrderItemMaterial) Validate() error {
	if m.ProductID == "" {
		return errors.New("Thiếu product_id của vật tư tiêu hao")
	}
	if m.Quantity <= 0 {
		return errors.New("Số lượng vật tư tiêu hao phải lớn hơn 0")
	}
	if m.UnitCostAtTime < 0 {
		return errors.New("Giá vốn tại thời điểm bán không được âm")
	}
	return nil
}

type Payment struct {
	ID             string
	OrganizationID string
	OrderID        string
	Method         string
	Provider       *string
	ProviderTxnID  *string
	Amount         float64
	Status         string
	PaidAt         *time.Time
	CreatedAt      *time.Time
}

func (p *Payment) Validate() error {
	switch p.Method {
	case PaymentMethodCash, PaymentMethodBankTransfer, PaymentMethodGateway, PaymentMethodCard:
	default:
		return errors.New("Hình thức thanh toán không hợp lệ (cash | bank_transfer | gateway | card)")
	}
	switch p.Status {
	case PaymentStatusPending, PaymentStatusSucceeded, PaymentStatusFailed, PaymentStatusRefunded:
	default:
		return errors.New("Trạng thái thanh toán không hợp lệ")
	}
	if p.Amount <= 0 {
		return errors.New("Số tiền thanh toán phải lớn hơn 0")
	}
	return nil
}
