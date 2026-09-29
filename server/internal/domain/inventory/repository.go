package inventory

import (
	"context"
	"errors"
)

// ErrInventoryNotFound dùng để usecase map 404 khi đọc sai id/khác organization.
var ErrInventoryNotFound = errors.New("inventory transaction not found")

// ErrInsufficientStock dùng để usecase map 400 khi phiếu làm âm kho.
var ErrInsufficientStock = errors.New("insufficient stock")

type InventoryRepository interface {
	CreateTransaction(ctx context.Context, tx *InventoryTransaction) error
	CreateDocument(ctx context.Context, doc *InventoryDocument, items []*DocumentItemInput) error
	BulkCreateDocumentItems(ctx context.Context, orgID string, documentID string, items []*DocumentItemInput) error
	GetTransactionsByOrgID(ctx context.Context, orgID, productID string, limit, offset int) ([]*InventoryTransaction, error)
	GetDocumentsByOrgID(ctx context.Context, orgID string, search string, sort string, filter string, limit, offset int) ([]*DocumentResponse, error)
	CountDocumentsByOrgID(ctx context.Context, orgID string) (int64, error)
	CountTransactionsByOrgID(ctx context.Context, orgID, productID string) (int64, error)
	GetLowStockProducts(ctx context.Context, orgID string, limit, offset int) ([]*LowStockItem, error)
	CountLowStockProducts(ctx context.Context, orgID string) (int64, error)
	GetDocumentByID(ctx context.Context, orgID string, documentID string) (*DocumentDetailResponse, error)
	GetInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string,limit int,  offset int , sort string, inventory_type string , date string) ([]*InventoryTransactionResponse, error)
	CountInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string) (int64, error)
}
