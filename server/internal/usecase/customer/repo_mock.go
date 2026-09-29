package customer

import (
	"context"

	appointmentDomain "github.com/daidat02/server/internal/domain/appointment"
	domain "github.com/daidat02/server/internal/domain/customer"
	"github.com/stretchr/testify/mock"
)

// MockCustomerRepository giả lập lại interface của CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetCustomersByOrgID(ctx context.Context, orgID, search string, limit, offset int) ([]*domain.Customer, error) {
	args := m.Called(ctx, orgID, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) CountCustomersByOrgID(ctx context.Context, orgID, search string) (int64, error) {
	args := m.Called(ctx, orgID, search)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCustomerRepository) GetCustomerByID(ctx context.Context, orgID, id string) (*domain.Customer, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) DeleteCustomer(ctx context.Context, orgID, id string) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

// MockHistoryStore giả lập HistoryStore (lịch sử lịch hẹn của khách).
type MockHistoryStore struct {
	mock.Mock
}

func (m *MockHistoryStore) GetAppointmentsByCustomerID(ctx context.Context, orgID, customerID string, limit, offset int) ([]*appointmentDomain.Appointment, error) {
	args := m.Called(ctx, orgID, customerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*appointmentDomain.Appointment), args.Error(1)
}

func (m *MockHistoryStore) CountAppointmentsByCustomerID(ctx context.Context, orgID, customerID string) (int64, error) {
	args := m.Called(ctx, orgID, customerID)
	return args.Get(0).(int64), args.Error(1)
}
