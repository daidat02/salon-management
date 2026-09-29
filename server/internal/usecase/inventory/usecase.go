package inventory

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/inventory"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

// ProductChecker kiểm tra sản phẩm thuộc đúng organization.
// PostgresServiceRepository thỏa mãn ngầm.
type ProductChecker interface {
	ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error)
}

type InventoryUsecase struct {
	repo    domain.InventoryRepository
	checker ProductChecker
}

func NewInventoryUsecase(repo domain.InventoryRepository, checker ProductChecker) *InventoryUsecase {
	return &InventoryUsecase{
		repo:    repo,
		checker: checker,
	}
}

func (u *InventoryUsecase) CreateDocument(ctx context.Context, req *domain.CreateDocumentRequest) (*domain.DocumentResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	if req == nil {
		return nil, apperror.New(400, "Thiếu thông tin phiếu kho", nil)
	}
	if len(req.Items) == 0 {
		return nil, apperror.New(400, "Phiếu kho phải có ít nhất 1 sản phẩm", nil)
	}
	
	doc := &domain.InventoryDocument{
		ID:             utils.NewID(),
		OrganizationID: orgID,
		SupplierID:     req.SupplierID,
		OrderID:        req.OrderID,
		DocumentCode:   strings.TrimSpace(req.DocumentCode),
		Type:           req.Type,
		Reason:         req.Reason,
		Note:           req.Note,
		CreatedBy:      req.CreatedBy,
	}
	items := make([]*domain.DocumentItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, &domain.DocumentItemInput{
			ProductID: strings.TrimSpace(it.ProductID),
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
		})
	}

	if err := u.repo.CreateDocument(ctx, doc, items); err != nil {
		log.Println("Error creating inventory document:", err)
		if errors.Is(err, domain.ErrInsufficientStock) {
			return nil, apperror.New(400, "Tồn kho không đủ cho phiếu này", err)
		}
		return nil, apperror.New(500, "Lỗi khi tạo phiếu kho", err)
	}

	return &domain.DocumentResponse{
		ID:             doc.ID,
		OrganizationID: doc.OrganizationID,
		SupplierID:     doc.SupplierID,
		OrderID:        doc.OrderID,
		DocumentCode:   doc.DocumentCode,
		Type:           doc.Type,
		TotalAmount:    doc.TotalAmount,
		Note:           doc.Note,
		CreatedAt:      doc.CreatedAt,
	}, nil
}

func (u *InventoryUsecase) GetTransactionsByOrgID(ctx context.Context, productID string, page, pageSize int) (*pagination.Result[*domain.InventoryTransaction], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)
	productID = strings.TrimSpace(productID)

	txs, err := u.repo.GetTransactionsByOrgID(ctx, orgID, productID, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy lịch sử kho", err)
	}

	total, err := u.repo.CountTransactionsByOrgID(ctx, orgID, productID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm lịch sử kho", err)
	}

	return pagination.NewResult(txs, params, total), nil
}

func (u *InventoryUsecase) GetDocumentsByOrgID(ctx context.Context, search string, sort string, filter string, page, pageSize int) (*pagination.Result[*domain.DocumentResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	params := pagination.Normalize(page, pageSize)
	inventoryDocs, err := u.repo.GetDocumentsByOrgID(ctx, orgID, search, sort, filter, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách phiếu kho", err)
	}

	total, err := u.repo.CountDocumentsByOrgID(ctx, orgID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách phiếu kho", err)
	}

	return pagination.NewResult(inventoryDocs, params, total), nil
}

func (u *InventoryUsecase) GetLowStockProducts(ctx context.Context, page, pageSize int) (*pagination.Result[*domain.LowStockItem], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)

	items, err := u.repo.GetLowStockProducts(ctx, orgID, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách tồn thấp", err)
	}

	total, err := u.repo.CountLowStockProducts(ctx, orgID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách tồn thấp", err)
	}

	return pagination.NewResult(items, params, total), nil
}

func (u *InventoryUsecase) GetDocumentByID(ctx context.Context, documentID string) (*domain.DocumentDetailResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	
	document, err := u.repo.GetDocumentByID(ctx, orgID, documentID)
	if err != nil {
		if errors.Is(err, domain.ErrInventoryNotFound) {
			return nil, apperror.New(404, "Phiếu kho không tồn tại", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy chi tiết phiếu kho", err)
	}

	return document, nil
}

func (u *InventoryUsecase) GetTransactionsByProductID(ctx context.Context, productID string, page, pageSize int, sort, inventory_type, date string) (*pagination.Result[*domain.InventoryTransactionResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	
	params := pagination.Normalize(page, pageSize)
	productID = strings.TrimSpace(productID)

	txs, err := u.repo.GetInventoryTransactionsByProdID(ctx, orgID, productID, params.PageSize, params.Offset(),sort, inventory_type, date)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách lịch sử kho", err)
	}

	total, err := u.repo.CountInventoryTransactionsByProdID(ctx, orgID, productID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm lịch sử kho", err)
	}

	return pagination.NewResult(txs, params, total), nil
}