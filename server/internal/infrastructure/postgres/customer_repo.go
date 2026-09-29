package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domain "github.com/daidat02/server/internal/domain/customer"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCustomerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCustomerRepository(db *pgxpool.Pool) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{
		db: db,
	}
}

// nullableGender chuyển "" thành nil để qua được CHECK gender của DB.
func nullableGender(gender string) any {
	if gender == "" {
		return nil
	}
	return gender
}

func (r *PostgresCustomerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	_, err := r.db.Exec(
		ctx,
		createCustomerQuery,
		customer.ID,
		customer.OrganizationID,
		customer.FullName,
		customer.Phone,
		nullableGender(customer.Gender),
		customer.BirthDate,
		customer.Note,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo khách hàng: %w", err)
	}

	return nil
}

func (r *PostgresCustomerRepository) GetCustomersByOrgID(ctx context.Context, orgID, search string, limit, offset int) ([]*domain.Customer, error) {
	query := listCustomersByOrgIDQuery
	args := []any{orgID}
	if strings.TrimSpace(search) != "" {
		query += searchCustomersFilter
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách khách hàng: %w", err)
	}
	defer rows.Close()

	customers := make([]*domain.Customer, 0)

	for rows.Next() {
		var c domain.Customer
		var gender, note *string

		err := rows.Scan(
			&c.ID,
			&c.OrganizationID,
			&c.FullName,
			&c.Phone,
			&gender,
			&c.BirthDate,
			&note,
			&c.TotalSpent,
			&c.TotalVisits,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng khách hàng: %w", err)
		}
		if gender != nil {
			c.Gender = *gender
		}
		if note != nil {
			c.Note = *note
		}

		customers = append(customers, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}

	return customers, nil
}

func (r *PostgresCustomerRepository) CountCustomersByOrgID(ctx context.Context, orgID, search string) (int64, error) {
	query := countCustomersByOrgIDQuery
	args := []any{orgID}
	if strings.TrimSpace(search) != "" {
		query += searchCustomersFilter
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}

	var total int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách khách hàng: %w", err)
	}
	return total, nil
}

func (r *PostgresCustomerRepository) GetCustomerByID(ctx context.Context, orgID, id string) (*domain.Customer, error) {
	var c domain.Customer
	var gender, note *string

	err := r.db.QueryRow(ctx, getCustomerByIDQuery, orgID, id).Scan(
		&c.ID,
		&c.OrganizationID,
		&c.FullName,
		&c.Phone,
		&gender,
		&c.BirthDate,
		&note,
		&c.TotalSpent,
		&c.TotalVisits,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("khách hàng %s: %w", id, domain.ErrCustomerNotFound)
		}
		return nil, fmt.Errorf("lỗi khi lấy khách hàng: %w", err)
	}
	if gender != nil {
		c.Gender = *gender
	}
	if note != nil {
		c.Note = *note
	}

	return &c, nil
}

func (r *PostgresCustomerRepository) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	tag, err := r.db.Exec(
		ctx,
		updateCustomerQuery,
		customer.ID,
		customer.OrganizationID,
		customer.FullName,
		customer.Phone,
		nullableGender(customer.Gender),
		customer.BirthDate,
		customer.Note,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật khách hàng: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật khách hàng %s: %w", customer.ID, domain.ErrCustomerNotFound)
	}

	return nil
}

func (r *PostgresCustomerRepository) DeleteCustomer(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(
		ctx,
		deleteCustomerQuery,
		id,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa khách hàng: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa khách hàng %s: %w", id, domain.ErrCustomerNotFound)
	}

	return nil
}

func (r *PostgresCustomerRepository) FindCustomerByPhone(ctx context.Context, orgID, phone string) (*domain.Customer, error) {
 	c := &domain.Customer{}
	var gender *string
	var note *string

	err := r.db.QueryRow(
		ctx,
		findCustomerByPhoneQuery,
		orgID,
		phone,
	).Scan(
		&c.ID,
		&c.OrganizationID,
		&c.FullName,
		&c.Phone,
		&gender,
		&c.BirthDate,
		&note,
		&c.TotalSpent,
		&c.TotalVisits,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	if err != nil {
		if(errors.Is(err, pgx.ErrNoRows)){
			return nil, nil
		}
		return nil, fmt.Errorf("lỗi khi tìm khách hàng theo số điện thoại: %w", err)
	}
	if gender != nil {
		c.Gender = *gender
	}
	if note != nil {
		c.Note = *note
	}

	return c, nil
}