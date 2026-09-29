package product

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/product"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

type ProductUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

func toResponse(p *domain.Product) *domain.ProductResponse {
	return &domain.ProductResponse{
		ID:             p.ID,
		OrganizationID: p.OrganizationID,
		CategoryID:     p.CategoryID,
		SKU:            p.SKU,
		Name:           p.Name,
		Unit:           p.Unit,
		NetUnit:        p.NetUnit,
		NetAmount:      p.NetAmount,
		ProductType:    p.ProductType,
		ShowOnWeb:      p.ShowOnWeb,
		CostPrice:      p.CostPrice,
		SellPrice:      p.SellPrice,
		StockQuantity:  p.StockQuantity,
		MinStock:       p.MinStock,
		IsActive:       p.IsActive,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// normalizeDefaults gán giá trị mặc định theo schema khi client không gửi.
func normalizeDefaults(p *domain.Product) {
	if p.ProductType == "" {
		p.ProductType = domain.ProductTypeRetail
	}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.ProductResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	normalizeDefaults(product)
	if err := product.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	product.ID = utils.NewID()
	product.OrganizationID = orgID

	if product.NetAmount != nil && *product.NetAmount > 0 {
		product.StockQuantity = product.StockQuantity * *product.NetAmount
		product.MinStock = product.MinStock * *product.NetAmount
	}
	if err := u.repo.CreateProduct(ctx, product); err != nil {
		log.Println("Error creating product:", err)
		return nil, apperror.New(500, "Lỗi khi tạo sản phẩm", err)
	}

	return toResponse(product), nil
}

func (u *ProductUsecase) GetProductsByOrgID(ctx context.Context, search string, status string, sort string, page, pageSize int) (*pagination.Result[*domain.ProductResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)

	products, err := u.repo.GetProductsByOrgID(ctx, orgID, search, status, sort, params.PageSize, params.Offset())
	if err != nil {
		log.Println("Error getting products by org ID:", err)
		return nil, apperror.New(500, "Lỗi khi lấy danh sách sản phẩm", err)
	}

	productsResponse := []*domain.ProductResponse{}

	for _, product := range products {
		productsResponse = append(productsResponse, toResponse(product))
	}

	total, err := u.repo.CountProductsByOrgID(ctx, orgID, search, status)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách sản phẩm", err)
	}
	return pagination.NewResult(productsResponse, params, total), nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id string, product *domain.Product) (*domain.ProductResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id sản phẩm", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	normalizeDefaults(product)
	if err := product.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	product.ID = id
	product.OrganizationID = orgID

	if err := u.repo.UpdateProduct(ctx, product); err != nil {
		log.Println("Error updating product:", err)
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, apperror.New(404, "Không tìm thấy sản phẩm", err)
		}
		return nil, apperror.New(500, "Lỗi khi cập nhật sản phẩm", err)
	}

	return toResponse(product), nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", apperror.New(400, "Thiếu id sản phẩm", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return "", apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := u.repo.DeleteProduct(ctx, orgID, id); err != nil {
		log.Println("Error deleting product:", err)
		if errors.Is(err, domain.ErrProductNotFound) {
			return "", apperror.New(404, "Không tìm thấy sản phẩm", err)
		}
		return "", apperror.New(500, "Lỗi khi xóa sản phẩm", err)
	}

	return id, nil
}
