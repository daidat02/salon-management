package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	domain "github.com/daidat02/server/internal/domain/inventory"
	prodomain "github.com/daidat02/server/internal/domain/product"
	"github.com/daidat02/server/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresInventoryRepository struct {
	db *pgxpool.Pool
	qx *sqlx.DB
}

func NewPostgresInventoryRepository(db *pgxpool.Pool) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
		db: db,
		qx: sqlx.NewDb(stdlib.OpenDBFromPool(db), "pgx"),
	}
}

func validateDocumentItems(items []*domain.DocumentItemInput, docType string) ([]string, error) {
	if len(items) == 0 {
		return nil, errors.New("danh sách sản phẩm trong phiếu không được để trống")
	}
	productIDs := make([]string, 0, len(items))
	for _, it := range items {
		if it == nil || it.ProductID == "" {
			return nil, errors.New("thiếu product_id trong phiếu kho")
		}
		if it.Quantity == 0 {
			return nil, fmt.Errorf("số lượng sản phẩm %s phải khác 0", it.ProductID)
		}
		// Dấu của tồn kho suy ra từ loại phiếu, nên input luôn nhập số dương
		// (riêng adjust cho phép âm để trừ kho trực tiếp).
		if docType != domain.TypeAdjust && it.Quantity < 0 {
			return nil, fmt.Errorf("số lượng sản phẩm %s phải lớn hơn 0", it.ProductID)
		}
		if it.UnitPrice < 0 {
			return nil, fmt.Errorf("đơn giá sản phẩm %s không được âm", it.ProductID)
		}
		productIDs = append(productIDs, it.ProductID)
	}
	return productIDs, nil
}

// signedQuantity đổi dấu số lượng theo loại phiếu:
// nhập/trả -> +, xuất/bán -> -, điều chỉnh giữ nguyên dấu người dùng nhập.
func signedQuantity(docType string, qty float64) float64 {
	switch docType {
	case domain.TypeExport, domain.TypeSale:
		if qty > 0 {
			return -qty
		}
		return qty
	case domain.TypeImport, domain.TypeReturn:
		if qty < 0 {
			return -qty
		}
		return qty
	default:
		return qty
	}
}

// fetchProductSnapshots lấy snapshot sản phẩm 1 lần duy nhất thay vì query từng dòng (N+1).
func fetchProductSnapshots(ctx context.Context, tx pgx.Tx, orgID string, productIDs []string) (map[string]prodomain.Product, error) {
	rows, err := tx.Query(ctx, getProductsByIDsQuery, orgID, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn sản phẩm: %w", err)
	}
	products := make(map[string]prodomain.Product, len(productIDs))
	for rows.Next() {
		var p prodomain.Product
		if err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Unit,
			&p.NetUnit,
			&p.NetAmount,
		); err != nil {
			rows.Close()
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu sản phẩm: %w", err)
		}
		products[p.ID] = p
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt sản phẩm: %w", err)
	}
	for _, id := range productIDs {
		if _, ok := products[id]; !ok {
			return nil, fmt.Errorf("sản phẩm %s không tồn tại trong tổ chức %s", id, orgID)
		}
	}
	return products, nil
}

// insertDocumentItemsCopy ghi chi tiết phiếu hàng loạt bằng CopyFrom.
func insertDocumentItemsCopy(ctx context.Context, tx pgx.Tx, documentID string, items []*domain.DocumentItemInput, products map[string]prodomain.Product) error {
	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"inventory_document_items"},
		[]string{"id", "document_id", "product_id", "product_name", "sku", "unit", "net_unit", "net_amount", "quantity", "unit_price", "total_price"},
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			it := items[i]
			p := products[it.ProductID]
			lineTotal := it.Quantity * it.UnitPrice
			return []any{
				utils.NewID(),
				documentID,
				p.ID,
				p.Name,
				p.SKU,
				p.Unit,
				p.NetUnit,
				p.NetAmount,
				it.Quantity,
				it.UnitPrice,
				lineTotal,
			}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("lỗi khi ghi chi tiết phiếu kho: %w", err)
	}
	return nil
}

// applyDocumentLedger cập nhật tồn kho + ghi sổ inventory_transactions cho từng dòng phiếu.
func applyDocumentLedger(ctx context.Context, tx pgx.Tx, orgID string, doc *domain.InventoryDocument, items []*domain.DocumentItemInput) error {
	const referenceType = "inventory_document"
	for _, it := range items {
		signedQty := signedQuantity(doc.Type, it.Quantity)
		var newStock float64
		if err := tx.QueryRow(ctx, applyStockChangeQuery, it.ProductID, orgID, signedQty).Scan(&newStock); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("sản phẩm %s: %w", it.ProductID, domain.ErrInsufficientStock)
			}
			return fmt.Errorf("lỗi khi cập nhật tồn kho sản phẩm %s: %w", it.ProductID, err)
		}
		if _, err := tx.Exec(ctx, createInventoryTransactionQuery,
			utils.NewID(), orgID, it.ProductID, doc.Type, signedQty, referenceType, doc.ID, doc.CreatedBy,
		); err != nil {
			return fmt.Errorf("lỗi khi ghi sổ kho sản phẩm %s: %w", it.ProductID, err)
		}
	}
	return nil
}

func (r *PostgresInventoryRepository) BulkCreateDocumentItems(ctx context.Context, orgID string, documentID string, items []*domain.DocumentItemInput) error {
	if documentID == "" {
		return errors.New("thiếu document_id của phiếu kho")
	}
	// Bulk lẻ dòng (không đổi tồn/sổ): coi như phiếu nhập để validate số dương.
	productIDs, err := validateDocumentItems(items, domain.TypeImport)
	if err != nil {
		return err
	}

	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	products, err := fetchProductSnapshots(ctx, dbTx, orgID, productIDs)
	if err != nil {
		return err
	}

	if err := insertDocumentItemsCopy(ctx, dbTx, documentID, items, products); err != nil {
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

// CreateDocument tạo phiếu kho hoàn chỉnh trong 1 transaction:
// header -> chi tiết (snapshot) -> cập nhật tồn + ghi sổ kho từng sản phẩm.
func (r *PostgresInventoryRepository) CreateDocument(ctx context.Context, doc *domain.InventoryDocument, items []*domain.DocumentItemInput) error {
	if doc == nil {
		return errors.New("thiếu thông tin phiếu kho")
	}
	if doc.OrganizationID == "" {
		return errors.New("thiếu organization_id của phiếu kho")
	}
	productIDs, err := validateDocumentItems(items, doc.Type)
	if err != nil {
		return err
	}
	if err := doc.Validate(); err != nil {
		return err
	}

	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	products, err := fetchProductSnapshots(ctx, dbTx, doc.OrganizationID, productIDs)
	if err != nil {
		return err
	}

	// Tính tổng server-side trước để insert header (items FK vào header nên phải tạo header trước).
	var total float64
	for _, it := range items {
		total += it.Quantity * it.UnitPrice
	}
	doc.TotalAmount = total

	if _, err := dbTx.Exec(ctx, createInventoryDocumentQuery,
		doc.ID, doc.OrganizationID, doc.SupplierID, doc.OrderID, doc.DocumentCode, doc.Type, doc.TotalAmount,doc.Reason, doc.Note, doc.CreatedBy,
	); err != nil {
		return fmt.Errorf("lỗi khi tạo phiếu kho: %w", err)
	}

	if err := insertDocumentItemsCopy(ctx, dbTx, doc.ID, items, products); err != nil {
		return err
	}

	if err := applyDocumentLedger(ctx, dbTx, doc.OrganizationID, doc, items); err != nil {
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

func (r *PostgresInventoryRepository) GetDocumentsByOrgID(ctx context.Context, orgID string, search string, sort string, filter string, limit, offset int) ([]*domain.DocumentResponse, error) {
	query := listInventoryDocumentsQuery +searchInventoryFilter + filterInventoryType +sortInventoryFilter
	var formatSearch = ""
	if search != "" {
		formatSearch = `%` + strings.TrimSpace(search) + `%`
	}
	args := []any{orgID, formatSearch, filter, sort}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	inventoryDocuments := make([]*domain.DocumentResponse, 0)
	if err := r.qx.SelectContext(ctx, &inventoryDocuments, query, args...); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách phiếu kho: %w", err)
	}

	return inventoryDocuments, nil
}

// CreateTransaction ghi sổ kho + cập nhật tồn trong 1 transaction.
// Cập nhật tồn dùng điều kiện chống âm nên an toàn khi bán đồng thời.
func (r *PostgresInventoryRepository) CreateTransaction(ctx context.Context, tx *domain.InventoryTransaction) error {
	dbTx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	var newStock float64
	if err := dbTx.QueryRow(ctx, applyStockChangeQuery, tx.ProductID, tx.OrganizationID, tx.Quantity).Scan(&newStock); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("sản phẩm %s: %w", tx.ProductID, domain.ErrInsufficientStock)
		}
		return fmt.Errorf("lỗi khi cập nhật tồn kho: %w", err)
	}

	if _, err := dbTx.Exec(ctx, createInventoryTransactionQuery,
		tx.ID, tx.OrganizationID, tx.ProductID, tx.Type, tx.Quantity, tx.ReferenceType, tx.ReferenceID, tx.CreatedBy,
	); err != nil {
		return fmt.Errorf("lỗi khi ghi sổ kho: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

func (r *PostgresInventoryRepository) GetTransactionsByOrgID(ctx context.Context, orgID, productID string, limit, offset int) ([]*domain.InventoryTransaction, error) {
	txs := make([]*domain.InventoryTransaction, 0)
	if err := r.qx.SelectContext(ctx, &txs, listInventoryTransactionsQuery, orgID, productID, limit, offset); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn lịch sử kho: %w", err)
	}
	return txs, nil
}
func (r *PostgresInventoryRepository) CountDocumentsByOrgID(ctx context.Context, orgID string) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, countInventoryDocumentsQuery, orgID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm lịch sử kho: %w", err)
	}
	return total, nil
}
func (r *PostgresInventoryRepository) CountTransactionsByOrgID(ctx context.Context, orgID, productID string) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, countInventoryTransactionsQuery, orgID, productID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm lịch sử kho: %w", err)
	}
	return total, nil
}

func (r *PostgresInventoryRepository) GetLowStockProducts(ctx context.Context, orgID string, limit, offset int) ([]*domain.LowStockItem, error) {
	items := make([]*domain.LowStockItem, 0)
	if err := r.qx.SelectContext(ctx, &items, listLowStockProductsQuery, orgID, limit, offset); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn sản phẩm tồn thấp: %w", err)
	}
	return items, nil
}

func (r *PostgresInventoryRepository) CountLowStockProducts(ctx context.Context, orgID string) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, countLowStockProductsQuery, orgID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm sản phẩm tồn thấp: %w", err)
	}
	return total, nil
}


func (r *PostgresInventoryRepository) GetDocumentByID(ctx context.Context, orgID string, documentID string) (*domain.DocumentDetailResponse, error) {
	var rawQueryResult domain.DocumentQueryResult

	if err := r.qx.GetContext(ctx, &rawQueryResult, InventoryDocumentDetailsQuery, orgID, documentID); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn chi tiết tài liệu: %w", err)
	}
	var items []*domain.DocumentItemResponse
    if len(rawQueryResult.ItemsRaw) > 0 {
        if err := json.Unmarshal(rawQueryResult.ItemsRaw, &items); err != nil {
            return nil, fmt.Errorf("lỗi khi parse items JSON: %w", err)
        }
    }

	response := &domain.DocumentDetailResponse{
        Document: &rawQueryResult.DocumentResponse,
        Items:    items,
    }
	return response, nil
}


func (r *PostgresInventoryRepository) GetInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string,limit int,  offset int , sort string, inventory_type string , date string) ([]*domain.InventoryTransactionResponse, error) {
	query := listInventoryTransactionsByProductIDQuery + FilterInventoryTransactionsByType + sortInventoryTransactionsFilter
	agrs := []any{orgID,productID, inventory_type, date, sort}

	query += fmt.Sprintf("Limit $%d OFFSET $%d", len(agrs) +1 , len(agrs)+2)
	agrs = append(agrs,limit, offset)

	inventoryTransactions := make([]*domain.InventoryTransactionResponse, 0)
	if err := r.qx.SelectContext(ctx, &inventoryTransactions, query, agrs...); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn lịch sử kho: %w", err)
	}
	return inventoryTransactions, nil
}

func (r *PostgresInventoryRepository) CountInventoryTransactionsByProdID(ctx context.Context, orgID string, productID string) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, countInventoryTransactionsByProductIDQuery, orgID, productID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm lịch sử kho: %w", err)
	}
	return total, nil
}