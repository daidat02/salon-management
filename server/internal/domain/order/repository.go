package order

import (
	"context"
	"errors"
)

// ErrOrderNotFound dùng để usecase map đúng 404 khi sai id/khác organization.
var ErrOrderNotFound = errors.New("order not found")

// ErrOrderNotReady dùng khi đơn chưa ở trạng thái paid nên không thể phục vụ.
var ErrOrderNotReady = errors.New("order not ready to serve")

// ErrInsufficientStock dùng để usecase map 400 khi bán vượt tồn.
var ErrInsufficientStock = errors.New("insufficient stock")

// StockDeduction lượng tồn cần trừ (số dương), repo tự đổi dấu khi ghi sổ bán.
type StockDeduction struct {
	ProductID string
	Quantity  float64
}

// SaleSlipLine dòng hàng cho phiếu xuất của đơn (giá đã chốt, số dương).
type SaleSlipLine struct {
	ProductID string
	Quantity  float64
	UnitPrice float64
}

// SaleSlipResult phiếu xuất vừa sinh cho đơn.
type SaleSlipResult struct {
	DocumentID   string
	DocumentCode string
}

// ProductPricing giá tại thời điểm bán, đọc từ bảng products.
type ProductPricing struct {
	SellPrice float64
	CostPrice float64
}

// ServiceMaterialUsage định mức vật tư của 1 đơn vị dịch vụ kèm giá vốn hiện tại.
type ServiceMaterialUsage struct {
	ProductID         string
	QuantityPerService float64
	UnitCost          float64
}

// ErrServiceNotFound dùng để usecase map 404 khi giá dịch vụ không còn.
var ErrServiceNotFound = errors.New("service not found")

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order, items []*OrderItem, materials []*OrderItemMaterial, deductions []StockDeduction) error
	GetOrdersByOrgID(ctx context.Context, orgID, status string, limit, offset int) ([]*Order, error)
	CountOrdersByOrgID(ctx context.Context, orgID, status string) (int64, error)
	GetOrderByID(ctx context.Context, orgID, id string) (*Order, error)
	GetOrderItems(ctx context.Context, orgID, orderID string) ([]*OrderItem, error)
	GetOrderMaterials(ctx context.Context, orderID string) ([]*OrderItemMaterial, error)
	GetOrderPayments(ctx context.Context, orgID, orderID string) ([]*Payment, error)
	// ApplyCancel chuyển trạng thái + hoàn kho (restores) + nhả giữ chỗ (releases)
	// + vô hiệu phiếu xuất + sổ trả trong 1 transaction.
	// Hủy đơn chưa phục vụ: restores rỗng, chỉ nhả giữ chỗ.
	ApplyCancel(ctx context.Context, orgID, orderID, newStatus string, restores []StockDeduction, releases []StockDeduction) error
	// ServeOrderWithSlip chuyển paid -> serving: trừ kho thật + ghi sổ sale +
	// sinh phiếu xuất, tất cả 1 tx. Trả về mã phiếu vừa sinh.
	ServeOrderWithSlip(ctx context.Context, orgID, orderID string, lines []SaleSlipLine) (SaleSlipResult, error)
	// CustomerExistsInOrg / AppointmentExistsInOrg kiểm tra tham chiếu thuộc đúng salon.
	CustomerExistsInOrg(ctx context.Context, orgID, customerID string) (bool, error)
	AppointmentExistsInOrg(ctx context.Context, orgID, appointmentID string) (bool, error)
	GetProductPricing(ctx context.Context, orgID string, productIDs []string) (map[string]ProductPricing, error)
	GetServicePrice(ctx context.Context, orgID, serviceID string) (float64, error)
	GetServiceMaterials(ctx context.Context, orgID string, serviceIDs []string) (map[string][]ServiceMaterialUsage, error)
	FindPaymentByProviderTxn(ctx context.Context, provider, txnID string) (*Payment, error)
	CreatePayment(ctx context.Context, payment *Payment) error
	GetTotalPaidByOrder(ctx context.Context, orgID, orderID string) (float64, error)
	UpdateOrderStatus(ctx context.Context, orgID, orderID, status string) error
	GetPosItemsByOrgID(ctx context.Context, orgID string, limit int, name, itemType string , id string, search string, categoryFilter string)([]*PosItemResponse, error) 
}
