package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	domain "github.com/daidat02/server/internal/domain/service"
	"github.com/daidat02/server/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresServiceRepository struct {
	db *pgxpool.Pool
}

func NewPostgresServiceRepository(db *pgxpool.Pool) *PostgresServiceRepository {
	return &PostgresServiceRepository{
		db: db,
	}
}
 
func (r *PostgresServiceRepository) bulkUpsertServiceMaterials(ctx context.Context,tx pgx.Tx, orgID,serviceID string, materials []*domain.ServiceMaterialInput) error {
	if len(materials) == 0 {
		return nil
	}
	query := bulkInsertServiceMaterialsQuery
	var args []interface{}
		var placeholders []string
	for _, m:= range materials{
		ID := utils.NewID()
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5))
		args = append(args, ID, orgID, serviceID, m.ProductID, m.Quantity)
	}

	query += strings.Join(placeholders, ", ")
	query += ` ON CONFLICT (organization_id, service_id, product_id) DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = now()`
	_, err := tx.Exec(ctx, query, args...)
	return err
}
func (r *PostgresServiceRepository) CreateService(ctx context.Context, service *domain.Service, materials []*domain.ServiceMaterialInput) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		createServiceQuery,
		service.ID,
		service.OrganizationID,
		service.CategoryID,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.BufferMinutes,
		service.Price,
		service.IsActive,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo dịch vụ: %w", err)
	}

	if(len(materials) > 0){
		err = r.bulkUpsertServiceMaterials(ctx, tx, service.OrganizationID, service.ID, materials)
		if err != nil {
			return fmt.Errorf("lỗi khi thêm nguyên vật liệu cho dịch vụ: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

func scanServiceRow(scan func(dest ...any) error) (*domain.Service, error) {
	var s domain.Service
	var description *string

	if err := scan(
		&s.ID,
		&s.OrganizationID,
		&s.CategoryID,
		&s.Name,
		&description,
		&s.DurationMinutes,
		&s.BufferMinutes,
		&s.Price,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if description != nil {
		s.Description = *description
	}

	return &s, nil
}

func (r *PostgresServiceRepository) GetServicesByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*domain.Service, error) {
	query := listServicesByOrgIDQuery + searchServicesFilter + filterServicesStatus + sortServicesFilter

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

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách dịch vụ: %w", err)
	}
	defer rows.Close()

	services := make([]*domain.Service, 0)

	for rows.Next() {
		s, err := scanServiceRow(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng dịch vụ: %w", err)
		}
		services = append(services, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}

	return services, nil
}

func (r *PostgresServiceRepository) CountServicesByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error) {
	formattedSearch := ""
	if strings.TrimSpace(search) != "" {
		formattedSearch = "%" + strings.TrimSpace(search) + "%"
	}
	args := []any{orgID, formattedSearch, status}
	query := countServicesByOrgIDQuery + searchServicesFilter + filterServicesStatus

	var total int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách dịch vụ: %w", err)
	}
	return total, nil
}

func (r *PostgresServiceRepository) CategoryExistsInOrg(ctx context.Context, orgID, categoryID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, categoryExistsInOrgQuery, categoryID, orgID).Scan(&exists); err != nil {
		return false, fmt.Errorf("lỗi khi kiểm tra danh mục: %w", err)
	}
	return exists, nil
}

// ProductsExistInOrg trả về danh sách productIDs không thuộc org (rỗng = tất cả hợp lệ).
func (r *PostgresServiceRepository) ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error) {
	rows, err := r.db.Query(ctx, productsExistInOrgQuery, orgID, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi kiểm tra sản phẩm: %w", err)
	}
	defer rows.Close()

	found := make(map[string]struct{}, len(productIDs))
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu sản phẩm: %w", err)
		}
		found[id] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình kiểm tra sản phẩm: %w", err)
	}

	var missing []string
	for _, id := range productIDs {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing, nil
}

// ServicesExistInOrg trả về danh sách serviceIDs không thuộc org (rỗng = tất cả hợp lệ).
func (r *PostgresServiceRepository) ServicesExistInOrg(ctx context.Context, orgID string, serviceIDs []string) ([]string, error) {
	rows, err := r.db.Query(ctx, servicesExistInOrgQuery, orgID, serviceIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi kiểm tra dịch vụ: %w", err)
	}
	defer rows.Close()

	found := make(map[string]struct{}, len(serviceIDs))
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dịch vụ: %w", err)
		}
		found[id] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình kiểm tra dịch vụ: %w", err)
	}

	var missing []string
	for _, id := range serviceIDs {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing, nil
}

func (r *PostgresServiceRepository) UpdateService(ctx context.Context, service *domain.Service) error {
	tag, err := r.db.Exec(
		ctx,
		updateServiceQuery,
		service.ID,
		service.OrganizationID,
		service.CategoryID,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.BufferMinutes,
		service.Price,
		service.IsActive,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật dịch vụ: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật dịch vụ %s: %w", service.ID, domain.ErrServiceNotFound)
	}

	return nil
}

func (r *PostgresServiceRepository) DeleteService(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(
		ctx,
		deleteServiceQuery,
		id,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa dịch vụ: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa dịch vụ %s: %w", id, domain.ErrServiceNotFound)
	}

	return nil
}



func (r *PostgresServiceRepository) GetServiceDetails(ctx context.Context, orgID, serviceID string) (*domain.ServiceResponse,[]*domain.ServiceMaterial, error) {
	service := &domain.ServiceResponse{}
	var materials []*domain.ServiceMaterial
	var materialsJSON []byte
	err := r.db.QueryRow(
		ctx,
		getServiceDetailsQuery,
		orgID,
		serviceID).Scan(
		&service.ID,
		&service.OrganizationID,
		&service.CategoryID,
		&service.Name,
		&service.Description,
		&service.DurationMinutes,
		&service.BufferMinutes,
		&service.Price,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
		&materialsJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("chi tiết dịch vụ %s: %w", serviceID, domain.ErrServiceNotFound)
		}
		return nil, nil, fmt.Errorf("lỗi query chi tiết dịch vụ: %w", err)
	}
	
	log.Printf("DEBUG: materialsJSON: %s", string(materialsJSON)) // Log the raw JSON for debugging
	if err := json.Unmarshal(materialsJSON, &materials); err != nil {
		return nil, nil, fmt.Errorf("lỗi unmarshal materials json: %w", err)
	}

	return service, materials, nil
}