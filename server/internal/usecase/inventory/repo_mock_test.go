package inventory

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/inventory"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepository giả lập lại interface của InventoryRepository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) CreateTransaction(ctx context.Context, tx *domain.InventoryTransaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockInventoryRepository) CreateDocument(ctx context.Context, doc *domain.InventoryDocument, items []*domain.DocumentItemInput) error {
	args := m.Called(ctx, doc, items)
	return args.Error(0)
}

func (m *MockInventoryRepository) BulkCreateDocumentItems(ctx context.Context, orgID string, documentID string, items []*domain.DocumentItemInput) error {
	args := m.Called(ctx, orgID, documentID, items)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetTransactionsByOrgID(ctx context.Context, orgID, productID string, limit, offset int) ([]*domain.InventoryTransaction, error) {
	args := m.Called(ctx, orgID, productID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.InventoryTransaction), args.Error(1)
}

func (m *MockInventoryRepository) CountTransactionsByOrgID(ctx context.Context, orgID, productID string) (int64, error) {
	args := m.Called(ctx, orgID, productID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInventoryRepository) GetLowStockProducts(ctx context.Context, orgID string, limit, offset int) ([]*domain.LowStockItem, error) {
	args := m.Called(ctx, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LowStockItem), args.Error(1)
}

func (m *MockInventoryRepository) CountLowStockProducts(ctx context.Context, orgID string) (int64, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInventoryRepository) GetDocumentsByOrgID(ctx context.Context, orgID, search string, sort string, filter string, limit, offset int) ([]*domain.DocumentResponse, error) {
	args := m.Called(ctx, orgID, search, sort, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DocumentResponse), args.Error(1)
}

func (m *MockInventoryRepository) CountDocumentsByOrgID(ctx context.Context, orgID string) (int64, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).(int64), args.Error(1)
}
// MockProductChecker giả lập ProductChecker
type MockProductChecker struct {
	mock.Mock
}

func (m *MockProductChecker) ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error) {
	args := m.Called(ctx, orgID, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockInventoryRepository) GetDocumentByID(ctx context.Context, orgID string, documentID string) (*domain.DocumentDetailResponse, error) {
	args := m.Called(ctx, orgID, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DocumentDetailResponse), args.Error(1)
}

func (m *MockInventoryRepository) GetInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string,limit int,  offset int , sort string, inventory_type string , date string) ([]*domain.InventoryTransactionResponse, error) {
	args := m.Called(ctx, orgID, productID, limit, offset, sort, inventory_type, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.InventoryTransactionResponse), args.Error(1)
}

func (m *MockInventoryRepository) CountInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string) (int64, error) {
	args := m.Called(ctx, orgID, productID)
	return args.Get(0).(int64), args.Error(1)
}