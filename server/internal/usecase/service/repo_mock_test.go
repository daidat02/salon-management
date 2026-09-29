package service

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/service"
	"github.com/stretchr/testify/mock"
)

// MockServiceRepository giả lập lại interface của ServiceRepository
type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) CreateService(ctx context.Context, service *domain.Service, materials []*domain.ServiceMaterialInput) error {
	args := m.Called(ctx, service, materials)
	return args.Error(0)
}

func (m *MockServiceRepository) GetServicesByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*domain.Service, error) {
	args := m.Called(ctx, orgID, search, status, sort, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Service), args.Error(1)
}

func (m *MockServiceRepository) CountServicesByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error) {
	args := m.Called(ctx, orgID, search, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockServiceRepository) UpdateService(ctx context.Context, service *domain.Service) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *MockServiceRepository) DeleteService(ctx context.Context, orgID, id string) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

func (m *MockServiceRepository) CategoryExistsInOrg(ctx context.Context, orgID, categoryID string) (bool, error) {
	args := m.Called(ctx, orgID, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockServiceRepository) ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error) {
	args := m.Called(ctx, orgID, productIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockServiceRepository) ServicesExistInOrg(ctx context.Context, orgID string, serviceIDs []string) ([]string, error) {
	args := m.Called(ctx, orgID, serviceIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}


func (m *MockServiceRepository) GetServiceDetails(ctx context.Context, orgID, serviceID string) (*domain.ServiceResponse, []*domain.ServiceMaterial, error) {
	args := m.Called(ctx, orgID, serviceID)
	var svc *domain.ServiceResponse
	if args.Get(0) != nil {
		svc = args.Get(0).(*domain.ServiceResponse)
	}
	var mats []*domain.ServiceMaterial
	if args.Get(1) != nil {
		mats = args.Get(1).([]*domain.ServiceMaterial)
	}
	return svc, mats, args.Error(2)
}