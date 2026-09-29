package postgres

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	invdomain "github.com/daidat02/server/internal/domain/inventory"
	domain "github.com/daidat02/server/internal/domain/order"
	"github.com/daidat02/server/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
	qx *sqlx.DB
}

func NewPostgresOrderRepository(db *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
		qx: sqlx.NewDb(stdlib.OpenDBFromPool(db), "pgx"),
	}
}

// nullableStr chuyển "" thành nil để qua được cột nullable (customer/appointment/staff).
func nullableStr(s *string) any {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	return strings.TrimSpace(*s)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// sortDeductionIDs sắp xếp product IDs cố định để các tx trừ cùng tập SP
// không deadlock nhau khi chạy đồng thời.
func sortDeductionIDs(deductions []domain.StockDeduction) []string {
	ids := make([]string, 0, len(deductions))
	seen := map[string]struct{}{}
	for _, d := range deductions {
		if _, ok := seen[d.ProductID]; !ok {
			seen[d.ProductID] = struct{}{}
			ids = append(ids, d.ProductID)
		}
	}
	sort.Strings(ids)
	return ids
}

func deductionQtyByProduct(deductions []domain.StockDeduction) map[string]float64 {
	qty := map[string]float64{}
	for _, d := range deductions {
		qty[d.ProductID] += d.Quantity
	}
	return qty
}

// CreateOrder ghi đơn + dòng đơn + snapshot vật tư + giữ chỗ kho trong 1 transaction.
// Chỉ giữ chỗ (reserved), KHÔNG trừ tồn thật và KHÔNG ghi sổ — tồn trừ lúc phục vụ.
func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, items []*domain.OrderItem, materials []*domain.OrderItemMaterial, deductions []domain.StockDeduction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, createOrderQuery,
		order.ID, order.OrganizationID, order.Code,
		nullableStr(order.CustomerID), nullableStr(order.AppointmentID),
		order.SubtotalAmount, order.DiscountAmount, order.TotalAmount, order.Status, nullableStr(order.CreatedBy),
	); err != nil {
		return fmt.Errorf("lỗi khi tạo đơn hàng: %w", err)
	}

	if err := bulkInsertOrderItems(ctx, tx, items); err != nil {
		return fmt.Errorf("lỗi khi thêm dòng đơn hàng: %w", err)
	}
	if err := bulkInsertOrderItemMaterials(ctx, tx, materials); err != nil {
		return fmt.Errorf("lỗi khi lưu vật tư tiêu hao: %w", err)
	}

	qtyByProduct := deductionQtyByProduct(deductions)
	for _, pid := range sortDeductionIDs(deductions) {
		qty := qtyByProduct[pid]
		tag, err := tx.Exec(ctx, reserveStockQuery, pid, order.OrganizationID, qty)
		if err != nil {
			return fmt.Errorf("lỗi khi giữ chỗ kho: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("sản phẩm %s: %w", pid, domain.ErrInsufficientStock)
		}
	}

	if err := tx.Commit(ctx); err != nil { 
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

func bulkInsertOrderItems(ctx context.Context, tx pgx.Tx, items []*domain.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	var args []any
	var placeholders []string
	for _, it := range items {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5,
			len(args)+6, len(args)+7, len(args)+8, len(args)+9, len(args)+10,
		))
		args = append(args,
			it.ID, it.OrganizationID, it.OrderID, it.ItemType,
			nullableStr(it.ServiceID), nullableStr(it.ProductID), nullableStr(it.StaffID),
			it.Quantity, it.UnitPrice, it.LineTotal,
		)
	}
	_, err := tx.Exec(ctx, bulkInsertOrderItemsQuery+strings.Join(placeholders, ", "), args...)
	return err
}

func bulkInsertOrderItemMaterials(ctx context.Context, tx pgx.Tx, materials []*domain.OrderItemMaterial) error {
	if len(materials) == 0 {
		return nil
	}
	var args []any
	var placeholders []string
	for _, m := range materials {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5,
		))
		args = append(args, m.ID, m.OrderItemID, m.ProductID, m.Quantity, m.UnitCostAtTime)
	}
	_, err := tx.Exec(ctx, bulkInsertOrderItemMaterialsQuery+strings.Join(placeholders, ", "), args...)
	return err
}

// ApplyCancel chuyển trạng thái đơn + hoàn kho (restores) + nhả giữ chỗ (releases)
// + vô hiệu phiếu xuất + sổ trả trong 1 transaction.
// - Hủy đơn chưa phục vụ: restores rỗng, chỉ nhả giữ chỗ.
// - Hoàn đơn đã phục vụ: restores đầy đủ, phiếu xuất sang cancelled.
func (r *PostgresOrderRepository) ApplyCancel(ctx context.Context, orgID, orderID, newStatus string, restores []domain.StockDeduction, releases []domain.StockDeduction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, updateOrderStatusQuery, orgID, orderID, newStatus)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật trạng thái đơn: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("đơn hàng %s: %w", orderID, domain.ErrOrderNotFound)
	}

	for _, rs := range restores {
		if _, err := tx.Exec(ctx, restoreStockQuery, rs.ProductID, orgID, rs.Quantity); err != nil {
			return fmt.Errorf("lỗi khi hoàn kho: %w", err)
		}
		if _, err := tx.Exec(ctx, insertInventoryLedgerQuery,
			utils.NewID(), orgID, rs.ProductID, "return", rs.Quantity, "order", orderID, nil,
		); err != nil {
			return fmt.Errorf("lỗi khi ghi sổ trả hàng: %w", err)
		}
	}

	for _, rel := range releases {
		if _, err := tx.Exec(ctx, releaseReserveQuery, rel.ProductID, orgID, rel.Quantity); err != nil {
			return fmt.Errorf("lỗi khi nhả giữ chỗ kho: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, voidSaleSlipQuery, orgID, orderID); err != nil {
		return fmt.Errorf("lỗi khi vô hiệu phiếu xuất: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return nil
}

// slipCodeForOrder sinh mã phiếu xuất từ mã đơn (XK-<suffix>), retry khi trùng.
func slipCodeForOrder(orderCode string, attempt int) string {
	suffix := strings.TrimPrefix(orderCode, "ORD-")
	if attempt == 0 {
		return "XK-" + suffix
	}
	return fmt.Sprintf("XK-%s-%d", suffix, attempt+1)
}

// ServeOrderWithSlip chuyển paid -> serving: trừ kho thật + ghi sổ sale +
// sinh phiếu xuất, tất cả trong 1 transaction. Trả về mã phiếu vừa sinh.
func (r *PostgresOrderRepository) ServeOrderWithSlip(ctx context.Context, orgID, orderID string, lines []domain.SaleSlipLine) (domain.SaleSlipResult, error) {
	var empty domain.SaleSlipResult
	if len(lines) == 0 {
		return empty, fmt.Errorf("phiếu xuất phải có ít nhất 1 dòng hàng")
	}
	items := make([]*invdomain.DocumentItemInput, 0, len(lines))
	for _, l := range lines {
		if l.Quantity <= 0 {
			return empty, fmt.Errorf("số lượng sản phẩm %s phải lớn hơn 0", l.ProductID)
		}
		if l.UnitPrice < 0 {
			return empty, fmt.Errorf("đơn giá sản phẩm %s không được âm", l.ProductID)
		}
		items = append(items, &invdomain.DocumentItemInput{
			ProductID: l.ProductID,
			Quantity:  l.Quantity,
			UnitPrice: l.UnitPrice,
		})
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return empty, fmt.Errorf("lỗi khi bắt đầu transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Khóa dòng đơn để serialize các lệnh serve đồng thời.
	var status string
	if err := tx.QueryRow(ctx, lockOrderRowQuery, orgID, orderID).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return empty, fmt.Errorf("đơn hàng %s: %w", orderID, domain.ErrOrderNotFound)
		}
		return empty, fmt.Errorf("lỗi khi khóa đơn hàng: %w", err)
	}
	if status != domain.StatusPaid {
		return empty, fmt.Errorf("đơn hàng %s chưa thanh toán: %w", orderID, domain.ErrOrderNotReady)
	}

	tag, err := tx.Exec(ctx, flipToServingQuery, orgID, orderID)
	if err != nil {
		return empty, fmt.Errorf("lỗi khi chuyển trạng thái phục vụ: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return empty, fmt.Errorf("đơn hàng %s chưa thanh toán: %w", orderID, domain.ErrOrderNotReady)
	}

	// Gom số lượng theo sản phẩm, trừ theo thứ tự cố định chống deadlock.
	qtyByProduct := map[string]float64{}
	priceByProduct := map[string]float64{}
	for _, l := range lines {
		if l.Quantity <= 0 {
			return empty, fmt.Errorf("số lượng sản phẩm %s phải lớn hơn 0", l.ProductID)
		}
		qtyByProduct[l.ProductID] += l.Quantity
		priceByProduct[l.ProductID] = l.UnitPrice
	}
	pids := make([]string, 0, len(qtyByProduct))
	for pid := range qtyByProduct {
		pids = append(pids, pid)
	}
	sort.Strings(pids)

	var orderCode string
	if err := tx.QueryRow(ctx, `SELECT code FROM orders WHERE organization_id = $1 AND id = $2`, orgID, orderID).Scan(&orderCode); err != nil {
		return empty, fmt.Errorf("lỗi khi đọc mã đơn hàng: %w", err)
	}

	docID := utils.NewID()
	var docCode string
	var total float64
	for _, pid := range pids {
		qty := qtyByProduct[pid]
		tag, err := tx.Exec(ctx, serveDeductStockQuery, pid, orgID, qty)
		if err != nil {
			return empty, fmt.Errorf("lỗi khi trừ kho sản phẩm %s: %w", pid, err)
		}
		if tag.RowsAffected() == 0 {
			return empty, fmt.Errorf("sản phẩm %s: %w", pid, domain.ErrInsufficientStock)
		}
		if _, err := tx.Exec(ctx, insertInventoryLedgerQuery,
			utils.NewID(), orgID, pid, "sale", -qty, "order", orderID, nil,
		); err != nil {
			return empty, fmt.Errorf("lỗi khi ghi sổ bán hàng: %w", err)
		}
		total += qty * priceByProduct[pid]
	}

	// Sinh header phiếu sale, retry khi trùng mã.
	for attempts := 0; ; attempts++ {
		docCode = slipCodeForOrder(orderCode, attempts)
		_, err = tx.Exec(ctx, createInventoryDocumentQuery,
			docID, orgID, nil, orderID, docCode, "sale", total, nil, nil, nil,
		)
		if err == nil {
			break
		}
		if !isUniqueViolation(err) || attempts >= 2 {
			return empty, fmt.Errorf("lỗi khi tạo phiếu xuất: %w", err)
		}
	}

	products, err := fetchProductSnapshots(ctx, tx, orgID, pids)
	if err != nil {
		return empty, err
	}
	if err := insertDocumentItemsCopy(ctx, tx, docID, items, products); err != nil {
		return empty, err
	}

	if err := tx.Commit(ctx); err != nil {
		return empty, fmt.Errorf("lỗi khi commit transaction: %w", err)
	}
	return domain.SaleSlipResult{DocumentID: docID, DocumentCode: docCode}, nil
}

func scanOrderRow(scan func(dest ...any) error) (*domain.Order, error) {
	var o domain.Order
	if err := scan(
		&o.ID, &o.OrganizationID, &o.Code, &o.CustomerID, &o.AppointmentID,
		&o.SubtotalAmount, &o.DiscountAmount, &o.TotalAmount, &o.Status,
		&o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &o, nil
}

func orderListFilter(status string, startIdx int) (string, []any) {
	if strings.TrimSpace(status) == "" {
		return "", nil
	}
	return fmt.Sprintf(" AND status = $%d", startIdx), []any{strings.TrimSpace(status)}
}

func (r *PostgresOrderRepository) GetOrdersByOrgID(ctx context.Context, orgID, status string, limit, offset int) ([]*domain.Order, error) {
	filter, args := orderListFilter(status, 2)
	query := fmt.Sprintf(`
		SELECT `+orderColumns+`
		FROM orders
		WHERE organization_id = $1%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, filter, len(args)+2, len(args)+3)
	args = append([]any{orgID}, append(args, limit, offset)...)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách đơn hàng: %w", err)
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		o, err := scanOrderRow(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng đơn hàng: %w", err)
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt danh sách: %w", err)
	}
	return orders, nil
}

func (r *PostgresOrderRepository) CountOrdersByOrgID(ctx context.Context, orgID, status string) (int64, error) {
	filter, args := orderListFilter(status, 2)
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM orders
		WHERE organization_id = $1%s
	`, filter)
	args = append([]any{orgID}, args...)

	var total int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi đếm danh sách đơn hàng: %w", err)
	}
	return total, nil
}

func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, orgID, id string) (*domain.Order, error) {
	o, err := scanOrderRow(func(dest ...any) error {
		return r.db.QueryRow(ctx, getOrderByIDQuery, orgID, id).Scan(dest...)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("đơn hàng %s: %w", id, domain.ErrOrderNotFound)
		}
		return nil, fmt.Errorf("lỗi khi lấy đơn hàng: %w", err)
	}
	return o, nil
}

func (r *PostgresOrderRepository) GetOrderItems(ctx context.Context, orgID, orderID string) ([]*domain.OrderItem, error) {
	rows, err := r.db.Query(ctx, listOrderItemsQuery, orgID, orderID)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn dòng đơn hàng: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.OrderItem, 0)
	for rows.Next() {
		var it domain.OrderItem
		if err := rows.Scan(
			&it.ID, &it.OrganizationID, &it.OrderID, &it.ItemType,
			&it.ServiceID, &it.ProductID, &it.StaffID,
			&it.Quantity, &it.UnitPrice, &it.LineTotal,
		); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu dòng đơn: %w", err)
		}
		items = append(items, &it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt dòng đơn: %w", err)
	}
	return items, nil
}

func (r *PostgresOrderRepository) GetOrderMaterials(ctx context.Context, orderID string) ([]*domain.OrderItemMaterial, error) {
	// Lọc salon ở tầng usecase qua đơn cha, ở đây chỉ cần đúng đơn.
	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.order_item_id, m.product_id, m.quantity, m.unit_cost_at_time
		FROM order_item_materials m
		JOIN order_items i ON i.id = m.order_item_id
		WHERE i.order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn vật tư đơn hàng: %w", err)
	}
	defer rows.Close()

	materials := make([]*domain.OrderItemMaterial, 0)
	for rows.Next() {
		var m domain.OrderItemMaterial
		if err := rows.Scan(&m.ID, &m.OrderItemID, &m.ProductID, &m.Quantity, &m.UnitCostAtTime); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu vật tư: %w", err)
		}
		materials = append(materials, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt vật tư: %w", err)
	}
	return materials, nil
}

func (r *PostgresOrderRepository) GetOrderPayments(ctx context.Context, orgID, orderID string) ([]*domain.Payment, error) {
	rows, err := r.db.Query(ctx, listOrderPaymentsQuery, orgID, orderID)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn thanh toán: %w", err)
	}
	defer rows.Close()

	payments := make([]*domain.Payment, 0)
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.OrderID, &p.Method, &p.Provider,
			&p.ProviderTxnID, &p.Amount, &p.Status, &p.PaidAt, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc dữ liệu thanh toán: %w", err)
		}
		payments = append(payments, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình duyệt thanh toán: %w", err)
	}
	return payments, nil
}

func (r *PostgresOrderRepository) CustomerExistsInOrg(ctx context.Context, orgID, customerID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, customerExistsInOrgQuery, orgID, customerID).Scan(&exists); err != nil {
		return false, fmt.Errorf("lỗi khi kiểm tra khách hàng: %w", err)
	}
	return exists, nil
}

func (r *PostgresOrderRepository) AppointmentExistsInOrg(ctx context.Context, orgID, appointmentID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, appointmentExistsInOrgQuery, orgID, appointmentID).Scan(&exists); err != nil {
		return false, fmt.Errorf("lỗi khi kiểm tra lịch hẹn: %w", err)
	}
	return exists, nil
}

func (r *PostgresOrderRepository) GetProductPricing(ctx context.Context, orgID string, productIDs []string) (map[string]domain.ProductPricing, error) {
	rows, err := r.db.Query(ctx, productPricingQuery, orgID, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy giá sản phẩm: %w", err)
	}
	defer rows.Close()

	pricing := make(map[string]domain.ProductPricing, len(productIDs))
	for rows.Next() {
		var id string
		var p domain.ProductPricing
		if err := rows.Scan(&id, &p.SellPrice, &p.CostPrice); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc giá sản phẩm: %w", err)
		}
		pricing[id] = p
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc giá: %w", err)
	}
	return pricing, nil
}

func (r *PostgresOrderRepository) GetServiceMaterials(ctx context.Context, orgID string, serviceIDs []string) (map[string][]domain.ServiceMaterialUsage, error) {
	rows, err := r.db.Query(ctx, serviceMaterialsWithCostQuery, orgID, serviceIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy định mức vật tư: %w", err)
	}
	defer rows.Close()

	usages := make(map[string][]domain.ServiceMaterialUsage)
	for rows.Next() {
		var serviceID string
		var u domain.ServiceMaterialUsage
		if err := rows.Scan(&serviceID, &u.ProductID, &u.QuantityPerService, &u.UnitCost); err != nil {
			return nil, fmt.Errorf("lỗi khi đọc định mức vật tư: %w", err)
		}
		usages[serviceID] = append(usages[serviceID], u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc định mức: %w", err)
	}
	return usages, nil
}

func (r *PostgresOrderRepository) GetServicePrice(ctx context.Context, orgID, serviceID string) (float64, error) {
	var price float64
	if err := r.db.QueryRow(ctx, servicePriceQuery, orgID, serviceID).Scan(&price); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("dịch vụ %s: %w", serviceID, domain.ErrServiceNotFound)
		}
		return 0, fmt.Errorf("lỗi khi lấy giá dịch vụ: %w", err)
	}
	return price, nil
}

func scanPaymentRow(scan func(dest ...any) error) (*domain.Payment, error) {
	var p domain.Payment
	if err := scan(
		&p.ID, &p.OrganizationID, &p.OrderID, &p.Method, &p.Provider,
		&p.ProviderTxnID, &p.Amount, &p.Status, &p.PaidAt, &p.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresOrderRepository) FindPaymentByProviderTxn(ctx context.Context, provider, txnID string) (*domain.Payment, error) {
	p, err := scanPaymentRow(func(dest ...any) error {
		return r.db.QueryRow(ctx, findPaymentByProviderTxnQuery, provider, txnID).Scan(dest...)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("lỗi khi tìm thanh toán: %w", err)
	}
	return p, nil
}

func (r *PostgresOrderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	_, err := r.db.Exec(ctx, createPaymentQuery,
		payment.ID, payment.OrganizationID, payment.OrderID, payment.Method,
		nullableStr(payment.Provider), nullableStr(payment.ProviderTxnID),
		payment.Amount, payment.Status, payment.PaidAt,
	)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo thanh toán: %w", err)
	}
	return nil
}

func (r *PostgresOrderRepository) GetTotalPaidByOrder(ctx context.Context, orgID, orderID string) (float64, error) {
	var total float64
	if err := r.db.QueryRow(ctx, totalPaidByOrderQuery, orgID, orderID).Scan(&total); err != nil {
		return 0, fmt.Errorf("lỗi khi tính tổng đã thu: %w", err)
	}
	return total, nil
}

func (r *PostgresOrderRepository) UpdateOrderStatus(ctx context.Context, orgID, orderID, status string) error {
	tag, err := r.db.Exec(ctx, updateOrderStatusQuery, orgID, orderID, status)
	if err != nil {
		return fmt.Errorf("lỗi khi cập nhật trạng thái đơn: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("đơn hàng %s: %w", orderID, domain.ErrOrderNotFound)
	}
	return nil
}

func (r *PostgresOrderRepository) GetPosItemsByOrgID(ctx context.Context, orgID string, limit int, name, itemType string , id string, search string, categoryFilter string)([]*domain.PosItemResponse, error) {
    formattedSearch := ""
    if strings.TrimSpace(search) != "" {
        formattedSearch = "%" + strings.TrimSpace(search) + "%"
    }
    query := listItemPosQuery + searchItemPosQuery + filterItemPosQuery + paginateItemPosQuery

    args := []any{
        orgID,          // $1
        "retail",       // $2 
        "both",         // $3
        formattedSearch,// $4
        categoryFilter, // $5
        itemType,       // $6 (cursor type)
        name,           // $7 (cursor name)
        id,             // $8 (cursor id)
        limit,          // $9
    }

	var items []*domain.PosItemResponse
	
	if err := r.qx.SelectContext(ctx,&items, query, args...); err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn danh sách sản phẩm POS: %w", err)
	}

	return items, nil
}