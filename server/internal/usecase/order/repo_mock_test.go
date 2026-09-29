package order

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/order"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository giả lập lại interface của OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, items []*domain.OrderItem, materials []*domain.OrderItemMaterial, deductions []domain.StockDeduction) error {
	args := m.Called(ctx, order, items, materials, deductions)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersByOrgID(ctx context.Context, orgID, status string, limit, offset int) ([]*domain.Order, error) {
	args := m.Called(ctx, orgID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) CountOrdersByOrgID(ctx context.Context, orgID, status string) (int64, error) {
	args := m.Called(ctx, orgID, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrderRepository) GetOrderByID(ctx context.Context, orgID, id string) (*domain.Order, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrderItems(ctx context.Context, orgID, orderID string) ([]*domain.OrderItem, error) {
	args := m.Called(ctx, orgID, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OrderItem), args.Error(1)
}

func (m *MockOrderRepository) GetOrderMaterials(ctx context.Context, orderID string) ([]*domain.OrderItemMaterial, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OrderItemMaterial), args.Error(1)
}

func (m *MockOrderRepository) GetOrderPayments(ctx context.Context, orgID, orderID string) ([]*domain.Payment, error) {
	args := m.Called(ctx, orgID, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockOrderRepository) ApplyCancel(ctx context.Context, orgID, orderID, newStatus string, restores []domain.StockDeduction, releases []domain.StockDeduction) error {
	args := m.Called(ctx, orgID, orderID, newStatus, restores, releases)
	return args.Error(0)
}

func (m *MockOrderRepository) ServeOrderWithSlip(ctx context.Context, orgID, orderID string, lines []domain.SaleSlipLine) (domain.SaleSlipResult, error) {
	args := m.Called(ctx, orgID, orderID, lines)
	if args.Get(0) == nil {
		return domain.SaleSlipResult{}, args.Error(1)
	}
	return args.Get(0).(domain.SaleSlipResult), args.Error(1)
}

func (m *MockOrderRepository) CustomerExistsInOrg(ctx context.Context, orgID, customerID string) (bool, error) {
	args := m.Called(ctx, orgID, customerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) AppointmentExistsInOrg(ctx context.Context, orgID, appointmentID string) (bool, error) {
	args := m.Called(ctx, orgID, appointmentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) GetProductPricing(ctx context.Context, orgID string, productIDs []string) (map[string]domain.ProductPricing, error) {
	args := m.Called(ctx, orgID, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]domain.ProductPricing), args.Error(1)
}

func (m *MockOrderRepository) GetServicePrice(ctx context.Context, orgID, serviceID string) (float64, error) {
	args := m.Called(ctx, orgID, serviceID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderRepository) GetServiceMaterials(ctx context.Context, orgID string, serviceIDs []string) (map[string][]domain.ServiceMaterialUsage, error) {
	args := m.Called(ctx, orgID, serviceIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]domain.ServiceMaterialUsage), args.Error(1)
}

func (m *MockOrderRepository) FindPaymentByProviderTxn(ctx context.Context, provider, txnID string) (*domain.Payment, error) {
	args := m.Called(ctx, provider, txnID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockOrderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockOrderRepository) GetTotalPaidByOrder(ctx context.Context, orgID, orderID string) (float64, error) {
	args := m.Called(ctx, orgID, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, orgID, orderID, status string) error {
	args := m.Called(ctx, orgID, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) GetPosItemsByOrgID(ctx context.Context, orgID string, limit int, name, itemType string, id string, search string, categoryFilter string) ([]*domain.PosItemResponse, error) {
	args := m.Called(ctx, orgID, limit, name, itemType, id, search, categoryFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PosItemResponse), args.Error(1)
}

// MockCatalogChecker giả lập CatalogChecker
type MockCatalogChecker struct {
	mock.Mock
}

func (m *MockCatalogChecker) ServicesExistInOrg(ctx context.Context, orgID string, serviceIDs []string) ([]string, error) {
	args := m.Called(ctx, orgID, serviceIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCatalogChecker) ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error) {
	args := m.Called(ctx, orgID, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
