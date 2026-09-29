package postgres

import (
	"context"
	"fmt"

	domain "github.com/daidat02/server/internal/domain/category"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCategoryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryRepository(db *pgxpool.Pool) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{
		db: db,
	}
}

func (r *PostgresCategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	_, err := r.db.Exec(
		ctx,
		createCategoryQuery,
		category.ID,
		category.OrganizationID,
		category.Type,
		category.Name,
		category.Description,
		category.IsActive,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo danh mục: %w", err)
	}

	return nil
}

func scanCategoryRow(scan func(dest ...any) error) (*domain.Category, error) {
	var c domain.Category
	var description *string

	if err := scan(
		&c.ID,
		&c.OrganizationID,
		&c.Type,
		&c.Name,
		&description,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if description != nil {
		c.Description = *description
	}

	return &c, nil
}

func (r *PostgresCategoryRepository) GetCategoriesByOrgID(ctx context.Context, orgID string, limit, offset int) ([]*domain.Category, error) {
	rows, err := r.db.Query(
		ctx,
		listCategoriesByOrgIDQuery,
		orgID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách danh mục: %w", err)
	}
	defer rows.Close()

	categories := make([]*domain.Category, 0)

	for rows.Next() {
		c, err := scanCategoryRow(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng danh mục: %w", err)
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) CountCategoriesByOrgID(ctx context.Context, orgID string) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, countCategoriesByOrgIDQuery, orgID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách danh mục: %w", err)
	}
	return total, nil
}

func (r *PostgresCategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	tag, err := r.db.Exec(
		ctx,
		updateCategoryQuery,
		category.ID,
		category.OrganizationID,
		category.Type,
		category.Name,
		category.Description,
		category.IsActive,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật danh mục: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật danh mục %s: %w", category.ID, domain.ErrCategoryNotFound)
	}

	return nil
}

func (r *PostgresCategoryRepository) DeleteCategory(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(
		ctx,
		deleteCategoryQuery,
		id,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa danh mục: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa danh mục %s: %w", id, domain.ErrCategoryNotFound)
	}

	return nil
}
