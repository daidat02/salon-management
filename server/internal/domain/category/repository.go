package category

import (
	"context"
	"errors"
)

// ErrCategoryNotFound dùng để usecase map đúng 404 khi UPDATE/DELETE không chạm dòng nào
// (sai id hoặc khác organization).
var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategoriesByOrgID(ctx context.Context, orgID string, limit, offset int) ([]*Category, error)
	CountCategoriesByOrgID(ctx context.Context, orgID string) (int64, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, orgID, id string) error
}
