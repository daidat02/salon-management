package postgres

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/daidat02/server/internal/domain/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (r *PostgresProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.db.Exec(
		ctx,
		createProductQuery,
		product.ID,
		product.OrganizationID,
		product.CategoryID,
		product.SKU,
		product.Name,
		product.Unit,
		product.NetUnit,
		product.NetAmount,
		product.ProductType,
		product.ShowOnWeb,
		product.CostPrice,
		product.SellPrice,
		product.StockQuantity,
		product.MinStock,
		product.IsActive,
	)	
	if err != nil {
		return fmt.Errorf("lỗi khi tạo sản phẩm: %w", err)
	}

	return nil
}

func scanProductRow(scan func(dest ...any) error) (*domain.Product, error) {
	var p domain.Product
	var unit *string

	if err := scan(
		&p.ID,
		&p.OrganizationID,
		&p.CategoryID,
		&p.SKU,
		&p.Name,
		&unit,
		&p.NetUnit,
		&p.NetAmount,
		&p.ProductType,
		&p.ShowOnWeb,
		&p.CostPrice,
		&p.SellPrice,
		&p.StockQuantity,
		&p.MinStock,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if unit != nil {
		p.Unit = *unit
	}

	return &p, nil
}

func (r *PostgresProductRepository) GetProductsByOrgID(
    ctx context.Context, 
    orgID string, 
    search string, 
    status string, 
    sort string, 
    limit, offset int,
) ([]*domain.Product, error) {

    // 1. Nối toàn bộ các khối hằng số lại với nhau theo đúng thứ tự $1, $2, $3, $4
    query := listProductsByOrgIDQuery + 
             searchProductsFilter + 
             filterProductsStatus + 
             sortProductsFilter

    // 2. Chuẩn bị mảng tham số cố định 4 vị trí đầu tiên ($1, $2, $3, $4)
    // - $1: orgID
    // - $2: search (nếu rỗng thì truyền "", SQL nhận diện qua ($2 = '' OR ...))
    // - $3: status (nếu rỗng thì truyền "", SQL nhận diện qua ($3 = '' OR ...))
    // - $4: sort   (truyền thẳng chuỗi sort vào, CASE WHEN sẽ tự so khớp)
    
    formattedSearch := ""
    if strings.TrimSpace(search) != "" {
        formattedSearch = "%" + strings.TrimSpace(search) + "%"
    }

    args := []any{
        orgID,           // $1
        formattedSearch, // $2
        status,          // $3
        sort,            // $4
    }

    // 3. Tự động cộng thêm phần Phân trang ở cuối dựa trên số lượng args hiện tại ($5, $6)
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
    args = append(args, limit, offset)

    // 4. Thực thi câu lệnh
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("lỗi khi truy vấn danh sách sản phẩm: %w", err)
    }
    defer rows.Close()

    products := make([]*domain.Product, 0)

    for rows.Next() {
        p, err := scanProductRow(rows.Scan)
        if err != nil {
            return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng sản phẩm: %w", err)
        }
        products = append(products, p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
    }

    return products, nil
}

func (r *PostgresProductRepository) CountProductsByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error) {
	formattedSearch := ""
	if strings.TrimSpace(search) != "" {
		formattedSearch = "%" + strings.TrimSpace(search) + "%"
	}
	args := []any{orgID, formattedSearch, status}
	query := countProductsByOrgIDQuery + searchProductsFilter + filterProductsStatus

	var total int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách sản phẩm: %w", err)
	}
	return total, nil
}

func (r *PostgresProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	tag, err := r.db.Exec(
		ctx,
		updateProductQuery,
		product.ID,
		product.OrganizationID,
		product.CategoryID,
		product.SKU,
		product.Name,
		product.Unit,
		product.ProductType,
		product.ShowOnWeb,
		product.CostPrice,
		product.SellPrice,
		product.StockQuantity,
		product.MinStock,
		product.IsActive,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật sản phẩm: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật sản phẩm %s: %w", product.ID, domain.ErrProductNotFound)
	}

	return nil
}

func (r *PostgresProductRepository) DeleteProduct(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(
		ctx,
		deleteProductQuery,
		id,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa sản phẩm: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa sản phẩm %s: %w", id, domain.ErrProductNotFound)
	}

	return nil
}
