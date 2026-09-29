package category

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/category"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository giả lập lại interface của CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetCategoriesByOrgID(ctx context.Context, orgID string, limit, offset int) ([]*domain.Category, error) {
	args := m.Called(ctx, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) CountCategoriesByOrgID(ctx context.Context, orgID string) (int64, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, orgID, id string) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
