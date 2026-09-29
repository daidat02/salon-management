package appointment

import (
	"context"
	"errors"

	"github.com/daidat02/server/internal/domain/customer"
)

// ErrAppointmentNotFound dùng để usecase map đúng 404.
var ErrAppointmentNotFound = errors.New("appointment not found")

type AppoimentRepository interface {
	CreateAppointment(ctx context.Context, appointment *Appointment, items []*AppointmentItem) error
	GetAppointmentByID(ctx context.Context, orgID, id string) (*Appointment, error)
	GetAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string, limit, offset int) ([]*Appointment, error)
	CountAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string) (int64, error)
	GetAppointmentItems(ctx context.Context, orgID, appointmentID string) ([]*AppointmentItem, error)
	UpdateAppointmentTimes(ctx context.Context, orgID, id, startTime, endTime string) error
	UpdateAppointmentStatus(ctx context.Context, orgID, id, status, cancelReason string) error
	DeleteAppointment(ctx context.Context, orgID, id string) error
}

type CustomerFinder interface {
	CreateCustomer(ctx context.Context, customer *customer.Customer) error
    FindCustomerByPhone(ctx context.Context, orgID, phone string) (*customer.Customer, error)
}

type ServiceFinder interface {
	
}