package customer

import (
	"context"
	"errors"
)

// ErrCustomerNotFound dùng để usecase map đúng 404 khi UPDATE/DELETE không chạm dòng nào
// (sai id, khác organization, hoặc đã bị xóa mềm).
var ErrCustomerNotFound = errors.New("customer not found")

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *Customer) error
	GetCustomersByOrgID(ctx context.Context, orgID, search string, limit, offset int) ([]*Customer, error)
	CountCustomersByOrgID(ctx context.Context, orgID, search string) (int64, error)
	GetCustomerByID(ctx context.Context, orgID, id string) (*Customer, error)
	UpdateCustomer(ctx context.Context, customer *Customer) error
	DeleteCustomer(ctx context.Context, orgID, id string) error
}
