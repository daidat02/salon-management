package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/service"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

type ServiceUsecase struct {
	repo domain.ServiceRepository
}

func NewServiceUsecase(repo domain.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{
		repo: repo,
	}
}

func toResponse(s *domain.Service) *domain.ServiceResponse {
	return &domain.ServiceResponse{
		ID:              s.ID,
		OrganizationID:  s.OrganizationID,
		CategoryID:      s.CategoryID,
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		BufferMinutes:   s.BufferMinutes,
		Price:           s.Price,
		IsActive:        s.IsActive,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func (u *ServiceUsecase) CreateService(ctx context.Context, service *domain.Service, materials []*domain.ServiceMaterialInput) (*domain.ServiceResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := service.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	if err := u.checkCategoryInOrg(ctx, orgID, service.CategoryID); err != nil {
		return nil, err
	}
	if err := u.validateMaterials(ctx, orgID, materials); err != nil {
		return nil, err
	}
	service.ID = utils.NewID()
	service.OrganizationID = orgID

	if err := u.repo.CreateService(ctx, service, materials); err != nil {
		log.Println("Error creating service:", err)
		return nil, apperror.New(500, "Lỗi khi tạo dịch vụ", err)
	}

	return toResponse(service), nil
}

// validateMaterials kiểm tra từng material ở tầng usecase (không phụ thuộc
// mỗi binding validator của handler): đủ product_id, quantity > 0,
// và mọi product đều thuộc đúng organization.
func (u *ServiceUsecase) validateMaterials(ctx context.Context, orgID string, materials []*domain.ServiceMaterialInput) error {
	seen := make(map[string]struct{}, len(materials))
	var productIDs []string
	for i, m := range materials {
		if m == nil || strings.TrimSpace(m.ProductID) == "" {
			return apperror.New(400, fmt.Sprintf("Thiếu product_id ở nguyên vật liệu thứ %d", i+1), nil)
		}
		if m.Quantity <= 0 {
			return apperror.New(400, "Số lượng nguyên vật liệu phải lớn hơn 0", nil)
		}
		if _, ok := seen[m.ProductID]; !ok {
			seen[m.ProductID] = struct{}{}
			productIDs = append(productIDs, m.ProductID)
		}
	}
	if len(productIDs) == 0 {
		return nil
	}

	missing, err := u.repo.ProductsExistInOrg(ctx, orgID, productIDs)
	if err != nil {
		return apperror.New(500, "Lỗi khi kiểm tra sản phẩm", err)
	}
	if len(missing) > 0 {
		return apperror.New(404, fmt.Sprintf("Không tìm thấy sản phẩm trong tổ chức: %s", strings.Join(missing, ", ")), nil)
	}
	return nil
}

// checkCategoryInOrg đảm bảo category (nếu có) thuộc đúng organization.
func (u *ServiceUsecase) checkCategoryInOrg(ctx context.Context, orgID string, categoryID *string) error {
	if categoryID == nil || strings.TrimSpace(*categoryID) == "" {
		return nil
	}

	exists, err := u.repo.CategoryExistsInOrg(ctx, orgID, *categoryID)
	if err != nil {
		return apperror.New(500, "Lỗi khi kiểm tra danh mục", err)
	}
	if !exists {
		return apperror.New(404, "Không tìm thấy danh mục trong tổ chức", nil)
	}
	return nil
}

func (u *ServiceUsecase) GetServicesByOrgID(ctx context.Context, search string, status string, sort string, page, pageSize int) (*pagination.Result[*domain.ServiceResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)

	services, err := u.repo.GetServicesByOrgID(ctx, orgID, search, status, sort, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách dịch vụ", err)
	}

	servicesResponse := []*domain.ServiceResponse{}
	for _, service := range services {
		servicesResponse = append(servicesResponse, toResponse(service))
	}

	total, err := u.repo.CountServicesByOrgID(ctx, orgID, search, status)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách dịch vụ", err)
	}

	return pagination.NewResult(servicesResponse, params, total), nil
}

func (u *ServiceUsecase) UpdateService(ctx context.Context, id string, service *domain.Service) (*domain.ServiceResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id dịch vụ", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := service.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	if err := u.checkCategoryInOrg(ctx, orgID, service.CategoryID); err != nil {
		return nil, err
	}
	service.ID = id
	service.OrganizationID = orgID

	if err := u.repo.UpdateService(ctx, service); err != nil {
		log.Println("Error updating service:", err)
		if errors.Is(err, domain.ErrServiceNotFound) {
			return nil, apperror.New(404, "Không tìm thấy dịch vụ", err)
		}
		return nil, apperror.New(500, "Lỗi khi cập nhật dịch vụ", err)
	}

	return toResponse(service), nil
}

func (u *ServiceUsecase) DeleteService(ctx context.Context, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", apperror.New(400, "Thiếu id dịch vụ", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return "", apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := u.repo.DeleteService(ctx, orgID, id); err != nil {
		log.Println("Error deleting service:", err)
		if errors.Is(err, domain.ErrServiceNotFound) {
			return "", apperror.New(404, "Không tìm thấy dịch vụ", err)
		}
		return "", apperror.New(500, "Lỗi khi xóa dịch vụ", err)
	}

	return id, nil
}



func (u *ServiceUsecase) GetServiceDetails(ctx context.Context, serviceID string) (*domain.ServiceResponse, []*domain.ServiceMaterial, error) {
	orgID,err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	service, materials, err := u.repo.GetServiceDetails(ctx,orgID,serviceID)
	if err != nil {
		if errors.Is(err, domain.ErrServiceNotFound) {
			return nil, nil, apperror.New(404, "Không tìm thấy dịch vụ", err)
		}
		return nil, nil, apperror.New(500, "Lỗi khi lấy thông tin dịch vụ và nguyên vật liệu", err)
	}

	return service, materials, nil
}