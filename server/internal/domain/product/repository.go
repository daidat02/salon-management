package product

import (
	"context"
	"errors"
)

// ErrProductNotFound dùng để usecase map đúng 404 khi UPDATE/DELETE không chạm dòng nào
// (sai id hoặc khác organization).
var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductsByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*Product, error)
	CountProductsByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, orgID, id string) error
}
