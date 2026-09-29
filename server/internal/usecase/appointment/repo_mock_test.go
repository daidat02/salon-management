package appointment

import (
	"context"

	domain "github.com/daidat02/server/internal/domain/appointment"
	customerDomain "github.com/daidat02/server/internal/domain/customer"
	"github.com/stretchr/testify/mock"
)

// MockAppointmentRepository giả lập AppoimentRepository.
type MockAppointmentRepository struct {
	mock.Mock
}

func (m *MockAppointmentRepository) CreateAppointment(ctx context.Context, appointment *domain.Appointment, items []*domain.AppointmentItem) error {
	args := m.Called(ctx, appointment, items)
	return args.Error(0)
}

func (m *MockAppointmentRepository) GetAppointmentByID(ctx context.Context, orgID, id string) (*domain.Appointment, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string, limit, offset int) ([]*domain.Appointment, error) {
	args := m.Called(ctx, orgID, status, from, to, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) CountAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string) (int64, error) {
	args := m.Called(ctx, orgID, status, from, to)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentItems(ctx context.Context, orgID, appointmentID string) ([]*domain.AppointmentItem, error) {
	args := m.Called(ctx, orgID, appointmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.AppointmentItem), args.Error(1)
}

func (m *MockAppointmentRepository) UpdateAppointmentTimes(ctx context.Context, orgID, id, startTime, endTime string) error {
	args := m.Called(ctx, orgID, id, startTime, endTime)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateAppointmentStatus(ctx context.Context, orgID, id, status, cancelReason string) error {
	args := m.Called(ctx, orgID, id, status, cancelReason)
	return args.Error(0)
}

func (m *MockAppointmentRepository) DeleteAppointment(ctx context.Context, orgID, id string) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

// MockCustomerFinder giả lập CustomerFinder (find-or-create customer theo phone).
type MockCustomerFinder struct {
	mock.Mock
}

func (m *MockCustomerFinder) FindCustomerByPhone(ctx context.Context, orgID, phone string) (*customerDomain.Customer, error) {
	args := m.Called(ctx, orgID, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customerDomain.Customer), args.Error(1)
}

func (m *MockCustomerFinder) CreateCustomer(ctx context.Context, customer *customerDomain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

// MockCatalogChecker giả lập CatalogChecker (check service/product thuộc org).
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
