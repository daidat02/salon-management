package staff

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/staff"
	"github.com/stretchr/testify/mock"
)

// MockStaffRepository giả lập lại interface của StaffRepository
type MockStaffRepository struct {
	mock.Mock
}

func (m *MockStaffRepository) CreateStaff(ctx context.Context, staff *domain.Staff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}

func (m *MockStaffRepository) GetStaffByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*domain.Staff, error) {
	args := m.Called(ctx, orgID, search, status, sort, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Staff), args.Error(1)
}

func (m *MockStaffRepository) CountStaffsByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error) {
	args := m.Called(ctx, orgID, search, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStaffRepository) UpdateStaff(ctx context.Context, staff *domain.Staff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}

func (m *MockStaffRepository) DeleteStaff(ctx context.Context, orgID, id string) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}