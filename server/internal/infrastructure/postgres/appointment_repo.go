package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domain "github.com/daidat02/server/internal/domain/appointment"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAppointmentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAppointmentRepository(db *pgxpool.Pool) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{
		db: db,
	}
}

// nullableUUID chuyển "" thành nil để qua được cột UUID nullable (staff_id).
func nullableUUID(id string) any {
	if strings.TrimSpace(id) == "" {
		return nil
	}
	return id
}

// CreateAppointment ghi lịch hẹn + items trong 1 transaction.
func (r *PostgresAppointmentRepository) CreateAppointment(ctx context.Context, appointment *domain.Appointment, items []*domain.AppointmentItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		createAppointmentQuery,
		appointment.ID,
		appointment.OrganizationID,
		appointment.Code,
		appointment.CustomerID,
		nullableUUID(appointment.StaffID),
		appointment.StartTime,
		appointment.EndTime,
		appointment.Status,
		appointment.Source,
		appointment.TotalAmount,
		appointment.Note,
		appointment.DeviceID,
		appointment.ClientIP,
		appointment.RequiredDepositAmount,
		appointment.DepositStatus,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo cuộc hẹn: %w", err)
	}

	if err := bulkInsertAppointmentItems(ctx, tx, items); err != nil {
		return fmt.Errorf("lỗi khi thêm chi tiết cuộc hẹn: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

func bulkInsertAppointmentItems(ctx context.Context, tx pgx.Tx, items []*domain.AppointmentItem) error {
	if len(items) == 0 {
		return nil
	}

	query := bulkInsertAppointmentItemsQuery
	var args []any
	var placeholders []string
	for _, item := range items {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5,
			len(args)+6, len(args)+7, len(args)+8, len(args)+9,
		))
		args = append(args,
			item.ID,
			item.OrganizationID,
			item.AppointmentID,
			item.ItemType,
			item.ServiceID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.LineTotal,
		)
	}

	query += strings.Join(placeholders, ", ")
	_, err := tx.Exec(ctx, query, args...)
	return err
}

// scanAppointmentRow map 1 dòng appointments sang entity.
// created_at/updated_at format RFC3339 để giữ nguyên kiểu *string của entity.
func scanAppointmentRow(scan func(dest ...any) error) (*domain.Appointment, error) {
	var a domain.Appointment
	var customerID, staffID, note, cancelReason, deviceID, clientIP *string
	var startTime, endTime time.Time
	var createdAt, updatedAt time.Time

	if err := scan(
		&a.ID,
		&a.OrganizationID,
		&a.Code,
		&customerID,
		&staffID,
		&startTime,
		&endTime,
		&a.Status,
		&a.Source,
		&a.TotalAmount,
		&note,
		&cancelReason,
		&deviceID,
		&clientIP,
		&a.RequiredDepositAmount,
		&a.DepositStatus,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	if customerID != nil {
		a.CustomerID = *customerID
	}
	if staffID != nil {
		a.StaffID = *staffID
	}
	if note != nil {
		a.Note = *note
	}
	if cancelReason != nil {
		a.CancelReason = *cancelReason
	}
	if deviceID != nil {
		a.DeviceID = *deviceID
	}
	if clientIP != nil {
		a.ClientIP = *clientIP
	}
	a.StartTime = startTime.Format(time.RFC3339)
	a.EndTime = endTime.Format(time.RFC3339)
	created := createdAt.Format(time.RFC3339)
	updated := updatedAt.Format(time.RFC3339)
	a.CreatedAt = &created
	a.UpdatedAt = &updated

	return &a, nil
}

func (r *PostgresAppointmentRepository) GetAppointmentByID(ctx context.Context, orgID, id string) (*domain.Appointment, error) {
	var (
		customerID, staffID, note, cancelReason, deviceID, clientIP *string
		startTime, endTime  time.Time
		createdAt, updatedAt time.Time
	)
	a := &domain.Appointment{}

	err := r.db.QueryRow(ctx, getAppointmentByIDQuery, orgID, id).Scan(
		&a.ID,
		&a.OrganizationID,
		&a.Code,
		&customerID,
		&staffID,
		&startTime,
		&endTime,
		&a.Status,
		&a.Source,
		&a.TotalAmount,
		&note,
		&cancelReason,
		&deviceID,
		&clientIP,
		&a.RequiredDepositAmount,
		&a.DepositStatus,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("lịch hẹn %s: %w", id, domain.ErrAppointmentNotFound)
		}
		return nil, fmt.Errorf("lỗi khi lấy lịch hẹn: %w", err)
	}
	if customerID != nil {
		a.CustomerID = *customerID
	}
	if staffID != nil {
		a.StaffID = *staffID
	}
	if note != nil {
		a.Note = *note
	}
	if cancelReason != nil {
		a.CancelReason = *cancelReason
	}
	if deviceID != nil {
		a.DeviceID = *deviceID
	}
	if clientIP != nil {
		a.ClientIP = *clientIP
	}
	a.StartTime = startTime.Format(time.RFC3339)
	a.EndTime = endTime.Format(time.RFC3339)
	created := createdAt.Format(time.RFC3339)
	updated := updatedAt.Format(time.RFC3339)
	a.CreatedAt = &created
	a.UpdatedAt = &updated

	return a, nil
}

// appointmentListFilter dựng mệnh đề WHERE động cho list/count.
func appointmentListFilter(status, from, to string, startIdx int) (string, []any) {
	var conds []string
	var args []any
	idx := startIdx
	if strings.TrimSpace(status) != "" {
		conds = append(conds, fmt.Sprintf("status = $%d", idx))
		args = append(args, strings.TrimSpace(status))
		idx++
	}
	if strings.TrimSpace(from) != "" {
		conds = append(conds, fmt.Sprintf("start_time >= $%d", idx))
		args = append(args, strings.TrimSpace(from))
		idx++
	}
	if strings.TrimSpace(to) != "" {
		conds = append(conds, fmt.Sprintf("start_time <= $%d", idx))
		args = append(args, strings.TrimSpace(to))
		idx++
	}
	if len(conds) == 0 {
		return "", args
	}
	return " AND " + strings.Join(conds, " AND "), args
}

func (r *PostgresAppointmentRepository) GetAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string, limit, offset int) ([]*domain.Appointment, error) {
	filter, args := appointmentListFilter(status, from, to, 2)
	query := fmt.Sprintf(`
		SELECT %s
		FROM appointments
		WHERE organization_id = $1%s
		ORDER BY start_time DESC
		LIMIT $%d OFFSET $%d
	`, appointmentColumns, filter, len(args)+2, len(args)+3)
	args = append([]any{orgID}, append(args, limit, offset)...)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách lịch hẹn: %w", err)
	}
	defer rows.Close()

	appointments := make([]*domain.Appointment, 0)
	for rows.Next() {
		a, err := scanAppointmentRow(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng lịch hẹn: %w", err)
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) CountAppointmentsByOrgID(ctx context.Context, orgID, status, from, to string) (int64, error) {
	filter, args := appointmentListFilter(status, from, to, 2)
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM appointments
		WHERE organization_id = $1%s
	`, filter)
	args = append([]any{orgID}, args...)

	var total int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách lịch hẹn: %w", err)
	}
	return total, nil
}

func (r *PostgresAppointmentRepository) GetAppointmentsByCustomerID(ctx context.Context, orgID, customerID string, limit, offset int) ([]*domain.Appointment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+appointmentColumns+`
		FROM appointments
		WHERE organization_id = $1 AND customer_id = $2
		ORDER BY start_time DESC
		LIMIT $3 OFFSET $4
	`, orgID, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn lịch sử lịch hẹn: %w", err)
	}
	defer rows.Close()

	appointments := make([]*domain.Appointment, 0)
	for rows.Next() {
		a, err := scanAppointmentRow(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng lịch hẹn: %w", err)
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) CountAppointmentsByCustomerID(ctx context.Context, orgID, customerID string) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM appointments
		WHERE organization_id = $1 AND customer_id = $2
	`, orgID, customerID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("lỗi khi đếm lịch sử lịch hẹn: %w", err)
	}
	return total, nil
}

func (r *PostgresAppointmentRepository) GetAppointmentItems(ctx context.Context, orgID, appointmentID string) ([]*domain.AppointmentItem, error) {
	rows, err := r.db.Query(ctx, listAppointmentItemsQuery, orgID, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn chi tiết lịch hẹn: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.AppointmentItem, 0)
	for rows.Next() {
		var item domain.AppointmentItem
		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.AppointmentID,
			&item.ItemType,
			&item.ServiceID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.LineTotal,
		); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng chi tiết: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt chi tiết: %w", err)
	}

	return items, nil
}

func (r *PostgresAppointmentRepository) UpdateAppointmentTimes(ctx context.Context, orgID, id, startTime, endTime string) error {
	tag, err := r.db.Exec(ctx, updateAppointmentTimesQuery, orgID, id, startTime, endTime)
	if err != nil {
		return fmt.Errorf("lỗi khi đổi lịch hẹn: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("đổi lịch hẹn %s: %w", id, domain.ErrAppointmentNotFound)
	}
	return nil
}

func (r *PostgresAppointmentRepository) UpdateAppointmentStatus(ctx context.Context, orgID, id, status, cancelReason string) error {
	tag, err := r.db.Exec(ctx, updateAppointmentStatusQuery, orgID, id, status, cancelReason)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật trạng thái lịch hẹn: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cập nhật lịch hẹn %s: %w", id, domain.ErrAppointmentNotFound)
	}
	return nil
}

func (r *PostgresAppointmentRepository) DeleteAppointment(ctx context.Context, orgID, id string) error {
	tag, err := r.db.Exec(ctx, deleteAppointmentQuery, orgID, id)
	if err != nil {
		return fmt.Errorf("lỗi khi xóa lịch hẹn: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("xóa lịch hẹn %s: %w", id, domain.ErrAppointmentNotFound)
	}
	return nil
}
