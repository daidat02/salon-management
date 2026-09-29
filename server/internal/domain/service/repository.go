package service

import (
	"context"
	"errors"
)

// ErrServiceNotFound dùng để usecase map đúng 404 khi UPDATE/DELETE không chạm dòng nào
// (sai id hoặc khác organization).
var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository interface {
	CreateService(ctx context.Context, service *Service, materials []*ServiceMaterialInput) error
	GetServicesByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*Service, error)
	CountServicesByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error)
	UpdateService(ctx context.Context, service *Service) error
	DeleteService(ctx context.Context, orgID, id string) error
	CategoryExistsInOrg(ctx context.Context, orgID, categoryID string) (bool, error)
	ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error)
	ServicesExistInOrg(ctx context.Context, orgID string, serviceIDs []string) ([]string, error)
	GetServiceDetails(ctx context.Context, orgID, serviceID string) (*ServiceResponse,[]*ServiceMaterial, error)
}
