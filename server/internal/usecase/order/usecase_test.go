package order

import (
	"context"
	"errors"
	"testing"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/order"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

func strP(s string) *string { return &s }

func catalogOK() *MockCatalogChecker {
	c := new(MockCatalogChecker)
	c.On("ServicesExistInOrg", mock.Anything, "org-123", mock.Anything).
		Return([]string{}, nil)
	c.On("ProductsExistInOrg", mock.Anything, "org-123", mock.Anything).
		Return([]string{}, nil)
	return c
}

func assertFlatAppError(t *testing.T, err error, expectedCode int) {
	t.Helper()
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr), "Lỗi phải là *apperror.AppError")
	if appErr != nil {
		assert.Equal(t, expectedCode, appErr.StatusCode, "Sai StatusCode")
		var inner *apperror.AppError
		assert.False(t, errors.As(appErr.Err, &inner), "Lỗi bị lồng 2 lớp AppError")
	}
}

func validProductOrderReq() *domain.CreateOrderRequest {
	return &domain.CreateOrderRequest{
		DiscountAmount: 0,
		Items: []domain.OrderItemInput{
			{ItemType: domain.ItemTypeProduct, ProductID: strP("prod-1"), Quantity: 2},
		},
	}
}

func TestOrderUsecase_CreateOrder(t *testing.T) {
	t.Run("Thành công - Đơn sản phẩm, giá server tính, trừ kho", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetProductPricing", mock.Anything, "org-123", mock.Anything).
			Return(map[string]domain.ProductPricing{"prod-1": {SellPrice: 150000, CostPrice: 90000}}, nil)
		repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CreateOrder(withOrgCtx("org-123"), validProductOrderReq())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Code)
		assert.Equal(t, float64(300000), res.SubtotalAmount)
		assert.Equal(t, float64(300000), res.TotalAmount)
		assert.Equal(t, domain.StatusPendingPayment, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Thiếu organization_id", func(t *testing.T) {
		uc := NewOrderUsecase(new(MockOrderRepository), catalogOK())
		res, err := uc.CreateOrder(context.Background(), validProductOrderReq())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Đơn rỗng", func(t *testing.T) {
		uc := NewOrderUsecase(new(MockOrderRepository), catalogOK())
		res, err := uc.CreateOrder(withOrgCtx("org-123"), &domain.CreateOrderRequest{})

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Giảm giá vượt tổng", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetProductPricing", mock.Anything, "org-123", mock.Anything).
			Return(map[string]domain.ProductPricing{"prod-1": {SellPrice: 150000, CostPrice: 90000}}, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		req := validProductOrderReq()
		req.DiscountAmount = 999999999
		res, err := uc.CreateOrder(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Tồn kho không đủ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetProductPricing", mock.Anything, "org-123", mock.Anything).
			Return(map[string]domain.ProductPricing{"prod-1": {SellPrice: 150000, CostPrice: 90000}}, nil)
		repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(domain.ErrInsufficientStock)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CreateOrder(withOrgCtx("org-123"), validProductOrderReq())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Sản phẩm khác tổ chức", func(t *testing.T) {
		catalog := new(MockCatalogChecker)
		catalog.On("ProductsExistInOrg", mock.Anything, "org-123", mock.Anything).
			Return([]string{"prod-x"}, nil)

		uc := NewOrderUsecase(new(MockOrderRepository), catalog)
		req := &domain.CreateOrderRequest{
			Items: []domain.OrderItemInput{
				{ItemType: domain.ItemTypeProduct, ProductID: strP("prod-x"), Quantity: 1},
			},
		}
		res, err := uc.CreateOrder(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
	})
}

func TestOrderUsecase_CancelRefund(t *testing.T) {
	pendingOrder := &domain.Order{ID: "ord-1", OrganizationID: "org-123", Code: "ORD-1", Status: domain.StatusPendingPayment, TotalAmount: 300000}
	paidOrder := &domain.Order{ID: "ord-2", OrganizationID: "org-123", Code: "ORD-2", Status: domain.StatusPaid, TotalAmount: 300000}
	servingOrder := &domain.Order{ID: "ord-4", OrganizationID: "org-123", Code: "ORD-4", Status: domain.StatusServing, TotalAmount: 300000}

	t.Run("Thành công - Hủy đơn chưa phục vụ chỉ nhả giữ chỗ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-1").Return(pendingOrder, nil)
		repo.On("GetOrderItems", mock.Anything, "org-123", "ord-1").
			Return([]*domain.OrderItem{
				{ID: "it-1", OrganizationID: "org-123", OrderID: "ord-1", ItemType: domain.ItemTypeProduct, ProductID: strP("prod-1"), Quantity: 2, UnitPrice: 150000, LineTotal: 300000},
			}, nil)
		repo.On("GetOrderMaterials", mock.Anything, "ord-1").
			Return([]*domain.OrderItemMaterial{}, nil)
		repo.On("ApplyCancel", mock.Anything, "org-123", "ord-1", domain.StatusCancelled, mock.Anything, mock.Anything).
			Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CancelOrder(withOrgCtx("org-123"), "ord-1")

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusCancelled, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thành công - Hủy đơn đã thu nhưng chưa phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-2").Return(paidOrder, nil)
		repo.On("GetOrderItems", mock.Anything, "org-123", "ord-2").
			Return([]*domain.OrderItem{}, nil)
		repo.On("GetOrderMaterials", mock.Anything, "ord-2").
			Return([]*domain.OrderItemMaterial{}, nil)
		repo.On("ApplyCancel", mock.Anything, "org-123", "ord-2", domain.StatusCancelled, mock.Anything, mock.Anything).
			Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CancelOrder(withOrgCtx("org-123"), "ord-2")

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusCancelled, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Hủy đơn đang phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-4").Return(servingOrder, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CancelOrder(withOrgCtx("org-123"), "ord-4")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Hoàn đơn mới thu tiền nhưng chưa phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-2").Return(paidOrder, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.RefundOrder(withOrgCtx("org-123"), "ord-2")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thành công - Hoàn đơn đang phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-4").Return(servingOrder, nil)
		repo.On("GetOrderItems", mock.Anything, "org-123", "ord-4").
			Return([]*domain.OrderItem{}, nil)
		repo.On("GetOrderMaterials", mock.Anything, "ord-4").
			Return([]*domain.OrderItemMaterial{}, nil)
		repo.On("ApplyCancel", mock.Anything, "org-123", "ord-4", domain.StatusRefunded, mock.Anything, mock.Anything).
			Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.RefundOrder(withOrgCtx("org-123"), "ord-4")

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusRefunded, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Không tìm thấy đơn", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "missing").
			Return(nil, domain.ErrOrderNotFound)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CancelOrder(withOrgCtx("org-123"), "missing")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
	})
}

func TestOrderUsecase_ServeComplete(t *testing.T) {
	paidOrder := &domain.Order{ID: "ord-2", OrganizationID: "org-123", Code: "ORD-2", Status: domain.StatusPaid, TotalAmount: 300000}
	servingOrder := &domain.Order{ID: "ord-4", OrganizationID: "org-123", Code: "ORD-4", Status: domain.StatusServing, TotalAmount: 300000}

	t.Run("Thành công - Phục vụ đơn đã thu, sinh phiếu xuất", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-2").Return(paidOrder, nil)
		repo.On("GetOrderItems", mock.Anything, "org-123", "ord-2").
			Return([]*domain.OrderItem{
				{ID: "it-1", OrganizationID: "org-123", OrderID: "ord-2", ItemType: domain.ItemTypeProduct, ProductID: strP("prod-1"), Quantity: 2, UnitPrice: 150000, LineTotal: 300000},
			}, nil)
		repo.On("GetOrderMaterials", mock.Anything, "ord-2").
			Return([]*domain.OrderItemMaterial{
				{ID: "m-1", OrderItemID: "it-1", ProductID: "prod-2", Quantity: 1, UnitCostAtTime: 20000},
			}, nil)
		repo.On("ServeOrderWithSlip", mock.Anything, "org-123", "ord-2", mock.Anything).
			Return(domain.SaleSlipResult{DocumentID: "doc-1", DocumentCode: "XK-ORD-2"}, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.ServeOrder(withOrgCtx("org-123"), "ord-2")

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, domain.StatusServing, res.Order.Status)
		assert.Equal(t, "XK-ORD-2", res.DocumentCode)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Phục vụ đơn chưa thu tiền", func(t *testing.T) {
		pending := &domain.Order{ID: "ord-1", OrganizationID: "org-123", Code: "ORD-1", Status: domain.StatusPendingPayment}
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-1").Return(pending, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.ServeOrder(withOrgCtx("org-123"), "ord-1")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Tồn kho không đủ lúc phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-2").Return(paidOrder, nil)
		repo.On("GetOrderItems", mock.Anything, "org-123", "ord-2").
			Return([]*domain.OrderItem{}, nil)
		repo.On("GetOrderMaterials", mock.Anything, "ord-2").
			Return([]*domain.OrderItemMaterial{}, nil)
		repo.On("ServeOrderWithSlip", mock.Anything, "org-123", "ord-2", mock.Anything).
			Return(domain.SaleSlipResult{}, domain.ErrInsufficientStock)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.ServeOrder(withOrgCtx("org-123"), "ord-2")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thành công - Hoàn tất đơn đang phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-4").Return(servingOrder, nil)
		repo.On("UpdateOrderStatus", mock.Anything, "org-123", "ord-4", domain.StatusCompleted).Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CompleteOrder(withOrgCtx("org-123"), "ord-4")

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusCompleted, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Hoàn tất đơn chưa phục vụ", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-2").Return(paidOrder, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CompleteOrder(withOrgCtx("org-123"), "ord-2")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})
}

func TestOrderUsecase_CreatePayment(t *testing.T) {
	pendingOrder := &domain.Order{ID: "ord-1", OrganizationID: "org-123", Code: "ORD-1", Status: domain.StatusPendingPayment, TotalAmount: 300000}

	t.Run("Thành công - Thu đủ tiền mặt, đơn sang paid", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-1").Return(pendingOrder, nil)
		repo.On("CreatePayment", mock.Anything, mock.Anything).Return(nil)
		repo.On("GetTotalPaidByOrder", mock.Anything, "org-123", "ord-1").Return(float64(300000), nil)
		repo.On("UpdateOrderStatus", mock.Anything, "org-123", "ord-1", domain.StatusPaid).Return(nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CreatePayment(withOrgCtx("org-123"), "ord-1", &domain.CreatePaymentRequest{
			Method: domain.PaymentMethodCash, Amount: 300000,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Thành công - Trùng mã giao dịch trả bản ghi cũ", func(t *testing.T) {
		existing := &domain.Payment{ID: "pay-old", OrganizationID: "org-123", OrderID: "ord-1",
			Method: domain.PaymentMethodGateway, Provider: strP("vnpay"), ProviderTxnID: strP("txn-1"),
			Amount: 300000, Status: domain.PaymentStatusSucceeded}
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-1").Return(pendingOrder, nil)
		repo.On("FindPaymentByProviderTxn", mock.Anything, "vnpay", "txn-1").Return(existing, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CreatePayment(withOrgCtx("org-123"), "ord-1", &domain.CreatePaymentRequest{
			Method: domain.PaymentMethodGateway, Provider: strP("vnpay"), ProviderTxnID: strP("txn-1"), Amount: 300000,
		})

		assert.NoError(t, err)
		assert.Equal(t, "pay-old", res.ID)
		repo.AssertNotCalled(t, "CreatePayment", mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Thanh toán đơn đã hủy", func(t *testing.T) {
		cancelled := &domain.Order{ID: "ord-3", OrganizationID: "org-123", Status: domain.StatusCancelled}
		repo := new(MockOrderRepository)
		repo.On("GetOrderByID", mock.Anything, "org-123", "ord-3").Return(cancelled, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.CreatePayment(withOrgCtx("org-123"), "ord-3", &domain.CreatePaymentRequest{
			Method: domain.PaymentMethodCash, Amount: 100000,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})
}

func TestOrderUsecase_GetOrders(t *testing.T) {
	orders := []*domain.Order{
		{ID: "ord-1", OrganizationID: "org-123", Code: "ORD-1", Status: domain.StatusPaid},
	}

	t.Run("Thành công - Danh sách phân trang", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrdersByOrgID", mock.Anything, "org-123", "", mock.Anything, mock.Anything).
			Return(orders, nil)
		repo.On("CountOrdersByOrgID", mock.Anything, "org-123", "").
			Return(int64(1), nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetOrdersByOrgID(withOrgCtx("org-123"), 1, 20, "")

		assert.NoError(t, err)
		assert.Len(t, res.Data, 1)
		assert.Equal(t, int64(1), res.Total)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Lỗi DB khi đếm", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetOrdersByOrgID", mock.Anything, "org-123", "", mock.Anything, mock.Anything).
			Return(orders, nil)
		repo.On("CountOrdersByOrgID", mock.Anything, "org-123", "").
			Return(int64(0), errors.New("db count failed"))

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetOrdersByOrgID(withOrgCtx("org-123"), 1, 20, "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
	})
}

func TestOrderUsecase_GetPosItemsByOrgID(t *testing.T) {
	stock := 50.0
	items := []*domain.PosItemResponse{
		{ID: "prod-1", OrganizationID: "org-123", Type: "product", Name: "Dầu gội", Price: 150000, StockQuantity: &stock, IsActive: true},
		{ID: "svc-1", OrganizationID: "org-123", Type: "service", Name: "Cắt tóc", Price: 200000, StockQuantity: nil, IsActive: true},
	}

	t.Run("Thành công - Trả danh sách gồm product và service", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetPosItemsByOrgID", mock.Anything, "org-123", 20, "", "", "", "", "").
			Return(items, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetPosItemsByOrgID(withOrgCtx("org-123"), 20, "", "", "", "", "")

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "product", res[0].Type)
		assert.NotNil(t, res[0].StockQuantity)
		assert.Nil(t, res[1].StockQuantity)
		repo.AssertExpectations(t)
	})

	t.Run("Thành công - limit ngoài khoảng được clamp", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetPosItemsByOrgID", mock.Anything, "org-123", 20, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(items, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetPosItemsByOrgID(withOrgCtx("org-123"), 0, "", "", "", "", "")

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		repo.AssertExpectations(t)
	})

	t.Run("Thành công - limit vượt trần được clamp", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetPosItemsByOrgID", mock.Anything, "org-123", 100, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(items, nil)

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetPosItemsByOrgID(withOrgCtx("org-123"), 999, "", "", "", "", "")

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Thiếu organization_id", func(t *testing.T) {
		uc := NewOrderUsecase(new(MockOrderRepository), catalogOK())
		res, err := uc.GetPosItemsByOrgID(context.Background(), 20, "", "", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Lỗi DB", func(t *testing.T) {
		repo := new(MockOrderRepository)
		repo.On("GetPosItemsByOrgID", mock.Anything, "org-123", 20, "", "", "", "", "").
			Return(nil, errors.New("db query failed"))

		uc := NewOrderUsecase(repo, catalogOK())
		res, err := uc.GetPosItemsByOrgID(withOrgCtx("org-123"), 20, "", "", "", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
	})
}
