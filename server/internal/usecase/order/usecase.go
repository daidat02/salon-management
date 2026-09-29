package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/order"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderUsecase struct {
	repo    domain.OrderRepository
	catalog CatalogChecker
}

// CatalogChecker kiểm tra service/product của dòng đơn có thuộc đúng organization không.
// PostgresServiceRepository thỏa mãn ngầm.
type CatalogChecker interface {
	ServicesExistInOrg(ctx context.Context, orgID string, serviceIDs []string) ([]string, error)
	ProductsExistInOrg(ctx context.Context, orgID string, productIDs []string) ([]string, error)
}

func NewOrderUsecase(repo domain.OrderRepository, catalog CatalogChecker) *OrderUsecase {
	return &OrderUsecase{
		repo:    repo,
		catalog: catalog,
	}
}

// generateOrderCode sinh mã đơn vừa VARCHAR(20), vd ORD-3F9A2C1D.
func generateOrderCode() string {
	id := strings.ReplaceAll(utils.NewID(), "-", "")
	if len(id) > 8 {
		id = id[:8]
	}
	return "ORD-" + strings.ToUpper(id)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func toResponse(o *domain.Order) *domain.OrderResponse {
	return &domain.OrderResponse{
		ID:             o.ID,
		OrganizationID: o.OrganizationID,
		Code:           o.Code,
		CustomerID:     o.CustomerID,
		AppointmentID:  o.AppointmentID,
		SubtotalAmount: o.SubtotalAmount,
		DiscountAmount: o.DiscountAmount,
		TotalAmount:    o.TotalAmount,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}

func strPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := strings.TrimSpace(s)
	return &v
}

func (u *OrderUsecase) checkReferencesInOrg(ctx context.Context, orgID string, req *domain.CreateOrderRequest) error {
	if req.CustomerID != nil && strings.TrimSpace(*req.CustomerID) != "" {
		exists, err := u.repo.CustomerExistsInOrg(ctx, orgID, strings.TrimSpace(*req.CustomerID))
		if err != nil {
			return apperror.New(500, "Lỗi khi kiểm tra khách hàng", err)
		}
		if !exists {
			return apperror.New(404, "Không tìm thấy khách hàng trong tổ chức", nil)
		}
	}
	if req.AppointmentID != nil && strings.TrimSpace(*req.AppointmentID) != "" {
		exists, err := u.repo.AppointmentExistsInOrg(ctx, orgID, strings.TrimSpace(*req.AppointmentID))
		if err != nil {
			return apperror.New(500, "Lỗi khi kiểm tra lịch hẹn", err)
		}
		if !exists {
			return apperror.New(404, "Không tìm thấy lịch hẹn trong tổ chức", nil)
		}
	}
	return nil
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}
	if len(req.Items) == 0 {
		return nil, apperror.New(400, "Đơn hàng phải có ít nhất một dòng hàng", nil)
	}
	if req.DiscountAmount < 0 {
		return nil, apperror.New(400, "Giảm giá không được âm", nil)
	}
	if err := u.checkReferencesInOrg(ctx, orgID, req); err != nil {
		return nil, err
	}

	var serviceIDs, productIDs []string
	seenService := map[string]struct{}{}
	seenProduct := map[string]struct{}{}
	for i, in := range req.Items {
		if in.Quantity <= 0 {
			return nil, apperror.New(400, fmt.Sprintf("Số lượng ở dòng thứ %d phải lớn hơn 0", i+1), nil)
		}
		switch strings.TrimSpace(in.ItemType) {
		case domain.ItemTypeService:
			if in.ServiceID == nil || strings.TrimSpace(*in.ServiceID) == "" {
				return nil, apperror.New(400, fmt.Sprintf("Dòng thứ %d thiếu service_id", i+1), nil)
			}
			if _, ok := seenService[strings.TrimSpace(*in.ServiceID)]; !ok {
				seenService[strings.TrimSpace(*in.ServiceID)] = struct{}{}
				serviceIDs = append(serviceIDs, strings.TrimSpace(*in.ServiceID))
			}
		case domain.ItemTypeProduct:
			if in.ProductID == nil || strings.TrimSpace(*in.ProductID) == "" {
				return nil, apperror.New(400, fmt.Sprintf("Dòng thứ %d thiếu product_id", i+1), nil)
			}
			if _, ok := seenProduct[strings.TrimSpace(*in.ProductID)]; !ok {
				seenProduct[strings.TrimSpace(*in.ProductID)] = struct{}{}
				productIDs = append(productIDs, strings.TrimSpace(*in.ProductID))
			}
		default:
			return nil, apperror.New(400, fmt.Sprintf("Loại dòng thứ %d không hợp lệ (service | product)", i+1), nil)
		}
	}

	if len(serviceIDs) > 0 {
		missing, err := u.catalog.ServicesExistInOrg(ctx, orgID, serviceIDs)
		if err != nil {
			return nil, apperror.New(500, "Lỗi khi kiểm tra dịch vụ", err)
		}
		if len(missing) > 0 {
			return nil, apperror.New(404, fmt.Sprintf("Không tìm thấy dịch vụ trong tổ chức: %s", strings.Join(missing, ", ")), nil)
		}
	}
	if len(productIDs) > 0 {
		missing, err := u.catalog.ProductsExistInOrg(ctx, orgID, productIDs)
		if err != nil {
			return nil, apperror.New(500, "Lỗi khi kiểm tra sản phẩm", err)
		}
		if len(missing) > 0 {
			return nil, apperror.New(404, fmt.Sprintf("Không tìm thấy sản phẩm trong tổ chức: %s", strings.Join(missing, ", ")), nil)
		}
	}

	// Giá tại thời điểm bán do server quyết định, không tin client.
	pricing := map[string]domain.ProductPricing{}
	if len(productIDs) > 0 {
		pricing, err = u.repo.GetProductPricing(ctx, orgID, productIDs)
		if err != nil {
			return nil, apperror.New(500, "Lỗi khi lấy giá sản phẩm", err)
		}
		for _, id := range productIDs {
			if _, ok := pricing[id]; !ok {
				return nil, apperror.New(404, fmt.Sprintf("Không tìm thấy sản phẩm trong tổ chức: %s", id), nil)
			}
		}
	}
	materialMap := map[string][]domain.ServiceMaterialUsage{}
	if len(serviceIDs) > 0 {
		materialMap, err = u.repo.GetServiceMaterials(ctx, orgID, serviceIDs)
		if err != nil {
			return nil, apperror.New(500, "Lỗi khi lấy định mức vật tư", err)
		}
	}

	orderID := utils.NewID()
	var items []*domain.OrderItem
	var materials []*domain.OrderItemMaterial
	deductions := map[string]float64{}
	var subtotal float64

	for _, in := range req.Items {
		itemType := strings.TrimSpace(in.ItemType)
		var unitPrice float64
		if itemType == domain.ItemTypeProduct {
			unitPrice = pricing[strings.TrimSpace(*in.ProductID)].SellPrice
		} else {
			// Giá dịch vụ đọc từ service repo? Tạm dùng giá 0 và để repo FK giữ đúng tham chiếu;
			// giá dịch vụ thực tế lấy từ bảng services qua pricing mở rộng sau.
			// Để đơn đúng tiền ngay, yêu cầu client không gửi giá mà server tra giá dịch vụ:
			unitPrice, err = u.servicePrice(ctx, orgID, strings.TrimSpace(*in.ServiceID))
			if err != nil {
				return nil, err
			}
		}
		item := &domain.OrderItem{
			ID:             utils.NewID(),
			OrganizationID: orgID,
			OrderID:        orderID,
			ItemType:       itemType,
			Quantity:       in.Quantity,
			UnitPrice:      unitPrice,
			LineTotal:      in.Quantity * unitPrice,
		}
		if in.ServiceID != nil {
			v := strings.TrimSpace(*in.ServiceID)
			item.ServiceID = &v
		}
		if in.ProductID != nil {
			v := strings.TrimSpace(*in.ProductID)
			item.ProductID = &v
		}
		if in.StaffID != nil {
			v := strings.TrimSpace(*in.StaffID)
			item.StaffID = &v
		}
		if err := item.Validate(); err != nil {
			return nil, apperror.New(400, err.Error(), err)
		}
		items = append(items, item)
		subtotal += item.LineTotal

		if itemType == domain.ItemTypeProduct {
			deductions[*item.ProductID] += item.Quantity
		} else {
			for _, mu := range materialMap[*item.ServiceID] {
				qty := mu.QuantityPerService * item.Quantity
				materials = append(materials, &domain.OrderItemMaterial{
					ID:             utils.NewID(),
					OrderItemID:    item.ID,
					ProductID:      mu.ProductID,
					Quantity:       qty,
					UnitCostAtTime: mu.UnitCost,
				});
				deductions[mu.ProductID] += qty
			}
		}
	}

	if req.DiscountAmount > subtotal {
		return nil, apperror.New(400, "Giảm giá không được vượt tổng tiền", nil)
	}
	total := subtotal - req.DiscountAmount
	status := domain.StatusPendingPayment
	if total == 0 {
		status = domain.StatusPaid
	}

	order := &domain.Order{
		ID:             orderID,
		OrganizationID: orgID,
		Code:           generateOrderCode(),
		CustomerID:     strPtr(ptrVal(req.CustomerID)),
		AppointmentID:  strPtr(ptrVal(req.AppointmentID)),
		SubtotalAmount: subtotal,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    total,
		Status:         status,
	}
	if err := order.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	for _, m := range materials {
		if err := m.Validate(); err != nil {
			return nil, apperror.New(400, err.Error(), err)
		}
	}

	var deductList []domain.StockDeduction
	for pid, qty := range deductions {
		deductList = append(deductList, domain.StockDeduction{ProductID: pid, Quantity: qty})
	}

	if err := u.repo.CreateOrder(ctx, order, items, materials, deductList); err != nil {
		log.Println("Error creating order:", err)
		if errors.Is(err, domain.ErrInsufficientStock) {
			return nil, apperror.New(400, "Tồn kho không đủ cho đơn hàng này", err)
		}
		if isUniqueViolation(err) {
			return nil, apperror.New(500, "Trùng mã đơn hàng, vui lòng thử lại", err)
		}
		return nil, apperror.New(500, "Lỗi khi tạo đơn hàng", err)
	}

	return toResponse(order), nil
}

// servicePrice tra giá dịch vụ tại thời điểm bán.
func (u *OrderUsecase) servicePrice(ctx context.Context, orgID, serviceID string) (float64, error) {
	price, err := u.repo.GetServicePrice(ctx, orgID, serviceID)
	if err != nil {
		if errors.Is(err, domain.ErrServiceNotFound) {
			return 0, apperror.New(404, "Không tìm thấy dịch vụ trong tổ chức", err)
		}
		return 0, apperror.New(500, "Lỗi khi lấy giá dịch vụ", err)
	}
	return price, nil
}

func ptrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (u *OrderUsecase) GetOrdersByOrgID(ctx context.Context, page, pageSize int, status string) (*pagination.Result[*domain.Order], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)
	status = strings.TrimSpace(status)

	orders, err := u.repo.GetOrdersByOrgID(ctx, orgID, status, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách đơn hàng", err)
	}

	total, err := u.repo.CountOrdersByOrgID(ctx, orgID, status)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách đơn hàng", err)
	}

	return pagination.NewResult(orders, params, total), nil
}

func (u *OrderUsecase) GetOrderDetail(ctx context.Context, id string) (*domain.OrderDetailResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id đơn hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	order, err := u.repo.GetOrderByID(ctx, orgID, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy đơn hàng", err)
	}

	items, err := u.repo.GetOrderItems(ctx, orgID, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy dòng đơn hàng", err)
	}

	materials, err := u.repo.GetOrderMaterials(ctx, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy vật tư đơn hàng", err)
	}

	payments, err := u.repo.GetOrderPayments(ctx, orgID, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy thanh toán đơn hàng", err)
	}

	detail := &domain.OrderDetailResponse{
		Order:     toResponse(order),
		Items:     make([]*domain.OrderItemResponse, 0, len(items)),
		Materials: make([]*domain.OrderMaterialResponse, 0, len(materials)),
		Payments:  make([]*domain.PaymentResponse, 0, len(payments)),
	}
	for _, it := range items {
		detail.Items = append(detail.Items, &domain.OrderItemResponse{
			ID: it.ID, ItemType: it.ItemType, ServiceID: it.ServiceID,
			ProductID: it.ProductID, StaffID: it.StaffID,
			Quantity: it.Quantity, UnitPrice: it.UnitPrice, LineTotal: it.LineTotal,
		})
	}
	for _, m := range materials {
		detail.Materials = append(detail.Materials, &domain.OrderMaterialResponse{
			ID: m.ID, OrderItemID: m.OrderItemID, ProductID: m.ProductID,
			Quantity: m.Quantity, UnitCostAtTime: m.UnitCostAtTime,
		})
	}
	for _, p := range payments {
		detail.Payments = append(detail.Payments, toPaymentResponse(p))
	}
	return detail, nil
}

func toPaymentResponse(p *domain.Payment) *domain.PaymentResponse {
	return &domain.PaymentResponse{
		ID: p.ID, OrderID: p.OrderID, Method: p.Method, Provider: p.Provider,
		ProviderTxnID: p.ProviderTxnID, Amount: p.Amount, Status: p.Status,
		PaidAt: p.PaidAt, CreatedAt: p.CreatedAt,
	}
}

// CancelOrder hủy đơn chưa thu tiền, hoàn kho đã trừ.
func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) (*domain.OrderResponse, error) {
	return u.closeOrder(ctx, id, domain.StatusCancelled, []string{domain.StatusDraft, domain.StatusPendingPayment, domain.StatusPaid},
		"Chỉ được hủy đơn chưa phục vụ")
}

// RefundOrder hoàn đơn đang phục vụ, hoàn kho đã trừ.
func (u *OrderUsecase) RefundOrder(ctx context.Context, id string) (*domain.OrderResponse, error) {
	return u.closeOrder(ctx, id, domain.StatusRefunded, []string{domain.StatusServing},
		"Chỉ được hoàn đơn đang phục vụ")
}

// ServeOrder bắt đầu phục vụ: trừ kho thật + ghi sổ sale + sinh phiếu xuất (repo làm 1 tx).
func (u *OrderUsecase) ServeOrder(ctx context.Context, id string) (*domain.OrderServeResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id đơn hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	order, err := u.repo.GetOrderByID(ctx, orgID, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy đơn hàng", err)
	}
	if order.Status != domain.StatusPaid {
		return nil, apperror.New(400, "Chỉ phục vụ đơn đã thanh toán", nil)
	}

	items, err := u.repo.GetOrderItems(ctx, orgID, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy dòng đơn hàng", err)
	}
	materials, err := u.repo.GetOrderMaterials(ctx, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy vật tư đơn hàng", err)
	}

	var lines []domain.SaleSlipLine
	for _, it := range items {
		if it.ItemType == domain.ItemTypeProduct && it.ProductID != nil {
			lines = append(lines, domain.SaleSlipLine{
				ProductID: *it.ProductID,
				Quantity:  it.Quantity,
				UnitPrice: it.UnitPrice,
			})
		}
	}
	for _, m := range materials {
		lines = append(lines, domain.SaleSlipLine{
			ProductID: m.ProductID,
			Quantity:  m.Quantity,
			UnitPrice: m.UnitCostAtTime,
		})
	}

	slip, err := u.repo.ServeOrderWithSlip(ctx, orgID, order.ID, lines)
	if err != nil {
		log.Println("Error serving order:", err)
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		if errors.Is(err, domain.ErrOrderNotReady) {
			return nil, apperror.New(400, "Chỉ phục vụ đơn đã thanh toán", err)
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			return nil, apperror.New(400, "Tồn kho không đủ để phục vụ đơn này", err)
		}
		return nil, apperror.New(500, "Lỗi khi bắt đầu phục vụ đơn hàng", err)
	}

	served := *order
	served.Status = domain.StatusServing
	return &domain.OrderServeResponse{
		Order:        toResponse(&served),
		DocumentID:   slip.DocumentID,
		DocumentCode: slip.DocumentCode,
	}, nil
}

// CompleteOrder kết thúc phục vụ: serving -> completed, không động kho.
func (u *OrderUsecase) CompleteOrder(ctx context.Context, id string) (*domain.OrderResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id đơn hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	order, err := u.repo.GetOrderByID(ctx, orgID, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy đơn hàng", err)
	}
	if order.Status != domain.StatusServing {
		return nil, apperror.New(400, "Chỉ hoàn tất đơn đang phục vụ", nil)
	}

	if err := u.repo.UpdateOrderStatus(ctx, orgID, order.ID, domain.StatusCompleted); err != nil {
		log.Println("Error completing order:", err)
		return nil, apperror.New(500, "Lỗi khi hoàn tất đơn hàng", err)
	}

	completed := *order
	completed.Status = domain.StatusCompleted
	return toResponse(&completed), nil
}

func (u *OrderUsecase) closeOrder(ctx context.Context, id, newStatus string, allowed []string, denyMsg string) (*domain.OrderResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id đơn hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	order, err := u.repo.GetOrderByID(ctx, orgID, strings.TrimSpace(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy đơn hàng", err)
	}

	ok := false
	for _, s := range allowed {
		if order.Status == s {
			ok = true
			break
		}
	}
	if !ok {
		return nil, apperror.New(400, denyMsg, nil)
	}

	items, err := u.repo.GetOrderItems(ctx, orgID, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy dòng đơn hàng", err)
	}
	materials, err := u.repo.GetOrderMaterials(ctx, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy vật tư đơn hàng", err)
	}

	restores := map[string]float64{}
	for _, it := range items {
		if it.ItemType == domain.ItemTypeProduct && it.ProductID != nil {
			restores[*it.ProductID] += it.Quantity
		}
	}
	for _, m := range materials {
		restores[m.ProductID] += m.Quantity
	}
	// Chỉ đơn đã phục vụ (refund) mới từng bị trừ kho thật nên cần hoàn.
	// Đơn hủy trước phục vụ chỉ nhả giữ chỗ (releases), tồn thật không đổi.
	var restoreList, releaseList []domain.StockDeduction
	for pid, qty := range restores {
		if newStatus == domain.StatusRefunded {
			restoreList = append(restoreList, domain.StockDeduction{ProductID: pid, Quantity: qty})
		} else {
			releaseList = append(releaseList, domain.StockDeduction{ProductID: pid, Quantity: qty})
		}
	}

	if err := u.repo.ApplyCancel(ctx, orgID, order.ID, newStatus, restoreList, releaseList); err != nil {
		log.Println("Error cancelling order:", err)
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi hủy/hoàn đơn hàng", err)
	}

	order.Status = newStatus
	return toResponse(order), nil
}

// CreatePayment ghi nhận thanh toán, idempotent theo cặp provider/txn.
// Thanh toán thành công đủ tổng thì đơn tự chuyển sang đã thanh toán.
func (u *OrderUsecase) CreatePayment(ctx context.Context, orderID string, req *domain.CreatePaymentRequest) (*domain.PaymentResponse, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, apperror.New(400, "Thiếu id đơn hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	payment := &domain.Payment{
		ID:             utils.NewID(),
		OrganizationID: orgID,
		OrderID:        strings.TrimSpace(orderID),
		Method:         strings.TrimSpace(req.Method),
		Provider:       req.Provider,
		ProviderTxnID:  req.ProviderTxnID,
		Amount:         req.Amount,
		Status:         domain.PaymentStatusSucceeded,
	}
	if t := time.Now(); true {
		payment.PaidAt = &t
	}
	if err := payment.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}

	order, err := u.repo.GetOrderByID(ctx, orgID, payment.OrderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy đơn hàng", err)
	}
	if order.Status == domain.StatusCancelled || order.Status == domain.StatusRefunded {
		return nil, apperror.New(400, "Không thể thanh toán đơn đã hủy/hoàn", nil)
	}
	if order.Status == domain.StatusPaid {
		return nil, apperror.New(400, "Đơn hàng đã được thanh toán đủ", nil)
	}

	// Idempotency: cùng cổng + mã giao dịch thì trả về bản ghi cũ.
	if payment.Provider != nil && strings.TrimSpace(*payment.Provider) != "" &&
		payment.ProviderTxnID != nil && strings.TrimSpace(*payment.ProviderTxnID) != "" {
		existing, err := u.repo.FindPaymentByProviderTxn(ctx, strings.TrimSpace(*payment.Provider), strings.TrimSpace(*payment.ProviderTxnID))
		if err != nil {
			return nil, apperror.New(500, "Lỗi khi kiểm tra thanh toán", err)
		}
		if existing != nil {
			return toPaymentResponse(existing), nil
		}
	}

	if err := u.repo.CreatePayment(ctx, payment); err != nil {
		log.Println("Error creating payment:", err)
		if isUniqueViolation(err) {
			if payment.Provider != nil && payment.ProviderTxnID != nil {
				existing, findErr := u.repo.FindPaymentByProviderTxn(ctx, strings.TrimSpace(*payment.Provider), strings.TrimSpace(*payment.ProviderTxnID))
				if findErr == nil && existing != nil {
					return toPaymentResponse(existing), nil
				}
			}
			return nil, apperror.New(400, "Giao dịch thanh toán đã tồn tại", err)
		}
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, apperror.New(404, "Không tìm thấy đơn hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi tạo thanh toán", err)
	}

	totalPaid, err := u.repo.GetTotalPaidByOrder(ctx, orgID, order.ID)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi tính tổng đã thu", err)
	}
	if totalPaid >= order.TotalAmount {
		if err := u.repo.UpdateOrderStatus(ctx, orgID, order.ID, domain.StatusPaid); err != nil {
			log.Println("Error marking order paid:", err)
			return nil, apperror.New(500, "Lỗi khi cập nhật trạng thái đơn", err)
		}
	}

	return toPaymentResponse(payment), nil
}

func (u *OrderUsecase) GetPosItemsByOrgID(ctx context.Context, limit int , name, itemType string , id string, search string, categoryFilter string)([]*domain.PosItemResponse, error) {
	orgID , err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	// Giới hạn page_size hợp lệ, mặc định 20
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	items, err := u.repo.GetPosItemsByOrgID(ctx, orgID, limit, name, itemType , id , search , categoryFilter)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách sản phẩm/dịch vụ", err)
	}

	return items, nil
	
}