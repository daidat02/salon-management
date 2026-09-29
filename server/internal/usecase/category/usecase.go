package category

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/category"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

type CategoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func toResponse(c *domain.Category) *domain.CategoryResponse {
	return &domain.CategoryResponse{
		ID:             c.ID,
		OrganizationID: c.OrganizationID,
		Type:           c.Type,
		Name:           c.Name,
		Description:    c.Description,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func (u *CategoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) (*domain.CategoryResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := category.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	category.ID = utils.NewID()
	category.OrganizationID = orgID

	if err := u.repo.CreateCategory(ctx, category); err != nil {
		log.Println("Error creating category:", err)
		return nil, apperror.New(500, "Lỗi khi tạo danh mục", err)
	}

	return toResponse(category), nil
}

func (u *CategoryUsecase) GetCategoriesByOrgID(ctx context.Context, page, pageSize int) (*pagination.Result[*domain.CategoryResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)

	categories, err := u.repo.GetCategoriesByOrgID(ctx, orgID, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách danh mục", err)
	}

	catResponses := []*domain.CategoryResponse{}
	for _, cat := range categories {
		catResponses = append(catResponses, toResponse(cat))
	}
	total, err := u.repo.CountCategoriesByOrgID(ctx, orgID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách danh mục", err)
	}

	return pagination.NewResult(catResponses, params, total), nil
}

func (u *CategoryUsecase) UpdateCategory(ctx context.Context, id string, category *domain.Category) (*domain.CategoryResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id danh mục", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := category.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	category.ID = id
	category.OrganizationID = orgID

	if err := u.repo.UpdateCategory(ctx, category); err != nil {
		log.Println("Error updating category:", err)
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return nil, apperror.New(404, "Không tìm thấy danh mục", err)
		}
		return nil, apperror.New(500, "Lỗi khi cập nhật danh mục", err)
	}

	return toResponse(category), nil
}

func (u *CategoryUsecase) DeleteCategory(ctx context.Context, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", apperror.New(400, "Thiếu id danh mục", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return "", apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := u.repo.DeleteCategory(ctx, orgID, id); err != nil {
		log.Println("Error deleting category:", err)
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return "", apperror.New(404, "Không tìm thấy danh mục", err)
		}
		return "", apperror.New(500, "Lỗi khi xóa danh mục", err)
	}

	return id, nil
}
