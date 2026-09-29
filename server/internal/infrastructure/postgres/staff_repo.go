package postgres

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/daidat02/server/internal/domain/staff"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStaffRepository struct {
	db *pgxpool.Pool
}

func NewPostgresStaffRepository(db *pgxpool.Pool) *PostgresStaffRepository {
	return &PostgresStaffRepository{
		db: db,
	}
}



func (r *PostgresStaffRepository) CreateStaff(ctx context.Context, staff *domain.Staff) error {
	_ ,err := r.db.Exec(
		ctx,
		createStaffQuery,
		staff.ID,
		staff.OrganizationID,
		staff.Code,
		staff.FullName,
		staff.Phone,
		staff.Position,
		staff.CommissionRate,
		staff.Status,
		staff.HireDate,
	)
	if err !=nil{
		return fmt.Errorf("Lỗi khi tạo nhân viên: %w", err)
	}

	return nil
}

func (r *PostgresStaffRepository) GetStaffByOrgID(
    ctx context.Context,
    orgID string,
    search string,
    status string,
    sort string,
    limit, offset int,
) ([]*domain.Staff, error) {
    // 1. Nối toàn bộ các khối hằng số lại với nhau theo đúng thứ tự $1, $2, $3, $4
    query := getStaffsByOrgIDQuery +
        searchStaffsFilter +
        filterStaffsStatus +
        sortStaffsFilter

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
        return nil, fmt.Errorf("lỗi khi truy vấn danh sách nhân viên: %w", err)
    }
    
    // Đảm bảo đóng rows để giải phóng kết nối database
    defer rows.Close()

    // 2. Khởi tạo mảng kết quả (non-nil để JSON ra [] thay vì null)
    staffs := make([]*domain.Staff, 0)

    // 3. Vòng lặp duyệt qua TỪNG DÒNG dữ liệu
    for rows.Next() {
        // Tạo một object trống để hứng dữ liệu cho dòng hiện tại
        var s domain.Staff
        
        // Scan dữ liệu của dòng đó vào object s
        err := rows.Scan(
            &s.ID,
            &s.OrganizationID,
            &s.Code,
            &s.FullName,
            &s.Phone,
            &s.Position,
            &s.CommissionRate,
            &s.Status,
            &s.HireDate,
            &s.CreatedAt,
            &s.UpdatedAt,
            &s.DeletedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng nhân viên: %w", err)
        }
        
        // Push object vừa scan được vào mảng kết quả
        staffs = append(staffs, &s)
    }

    // 4. Kiểm tra xem có lỗi ẩn nào xảy ra trong quá trình chạy vòng lặp không
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
    }

    return staffs, nil
}

func (r *PostgresStaffRepository) CountStaffsByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error) {
    formattedSearch := ""
    if strings.TrimSpace(search) != "" {
        formattedSearch = "%" + strings.TrimSpace(search) + "%"
    }
    args := []any{orgID, formattedSearch, status}
    query := countStaffsByOrgIDQuery + searchStaffsFilter + filterStaffsStatus

    var total int64
    if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
        return 0, fmt.Errorf("lỗi khi đếm danh sách nhân viên: %w", err)
    }
    return total, nil
}

func (r *PostgresStaffRepository) UpdateStaff(ctx context.Context, staff *domain.Staff) error {
	tag, err := r.db.Exec(
		ctx,
		updateStaffQuery,
		staff.ID,
		staff.OrganizationID,
		staff.Code,
		staff.FullName,
		staff.Phone,
		staff.Position,
		staff.CommissionRate,
		staff.Status,
		staff.HireDate,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật nhân viên: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật nhân viên %s: %w", staff.ID, domain.ErrStaffNotFound)
	}

	return nil
}

func (r *PostgresStaffRepository) DeleteStaff(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(
		ctx,
		deleteStaffQuery,
		id,
		orgID,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa nhân viên: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa nhân viên %s: %w", id, domain.ErrStaffNotFound)
	}

	return nil
}