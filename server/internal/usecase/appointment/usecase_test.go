package appointment

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/appointment"
	customerDomain "github.com/daidat02/server/internal/domain/customer"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

func validRequest() *domain.AppointmentCreateRequest {
	return &domain.AppointmentCreateRequest{
		CustomerName:  "Nguyễn Thị C",
		CustomerPhone: "0912345678",
		StaffID:       "staff-1",
		StartTime:     "2026-10-01T14:00:00+07:00",
		EndTime:       "2026-10-01T15:00:00+07:00",
		Status:        "confirmed",
		Source:        "online",
		TotalAmount:   650000,
		Items: []domain.ApoimentItemInput{
			{ItemType: "service", ServiceID: "svc-1", Quantity: 1, UnitPrice: 500000},
			{ItemType: "product", ProductID: "prod-1", Quantity: 1, UnitPrice: 150000},
		},
	}
}

func existingCustomer() *customerDomain.Customer {
	return &customerDomain.Customer{
		ID: "cus-1", OrganizationID: "org-123",
		FullName: "Nguyễn Thị C", Phone: "0912345678",
	}
}

// catalogOK giả lập items đều thuộc org (svc-1, prod-1).
func catalogOK() *MockCatalogChecker {
	catalog := new(MockCatalogChecker)
	catalog.On("ServicesExistInOrg", mock.Anything, "org-123", mock.Anything).
		Return(nil, nil)
	catalog.On("ProductsExistInOrg", mock.Anything, "org-123", mock.Anything).
		Return(nil, nil)
	return catalog
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

func TestAppointmentUseCase_CreateAppointment(t *testing.T) {
	t.Run("Thành công - Khách cũ, total do server tính, bỏ qua client gửi sai", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(existingCustomer(), nil).Once()
		repo.On("CreateAppointment", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.Code = "APT-CLIENT-1"
		req.TotalAmount = 1 // client gửi sai, server phải tính lại 500000 + 150000
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "cus-1", res.CustomerID, "Phải gắn id khách đã tồn tại")
		assert.Equal(t, "staff-1", res.StaffID, "Phải gán thợ")
		assert.Equal(t, "APT-CLIENT-1", res.Code)
		assert.Equal(t, "org-123", res.OrganizationID)
		assert.Equal(t, 650000.0, res.TotalAmount, "Total phải do server tính từ items")
		finder.AssertNotCalled(t, "CreateCustomer", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
		finder.AssertExpectations(t)
		catalog.AssertExpectations(t)
	})

	t.Run("Thành công - Khách mới, tự tạo + tự sinh mã", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(nil, nil).Once()
		finder.On("CreateCustomer", mock.Anything, mock.Anything).
			Return(nil).Once()
		repo.On("CreateAppointment", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.CustomerID, "Phải có id khách vừa tạo")
		assert.True(t, strings.HasPrefix(res.Code, "APT-"), "Code trống phải tự sinh, got %q", res.Code)
		assert.NotEmpty(t, res.ID)
		assert.Equal(t, 650000.0, res.TotalAmount)
		repo.AssertExpectations(t)
		finder.AssertExpectations(t)
		catalog.AssertExpectations(t)
	})

	t.Run("Thành công - Race trùng phone thì find lại", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(nil, nil).Once()
		finder.On("CreateCustomer", mock.Anything, mock.Anything).
			Return(&pgconn.PgError{Code: "23505"}).Once()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(existingCustomer(), nil).Once()
		repo.On("CreateAppointment", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.NoError(t, err)
		assert.Equal(t, "cus-1", res.CustomerID)
		repo.AssertExpectations(t)
		finder.AssertExpectations(t)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Thiếu organization_id", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(context.Background(), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		finder.AssertNotCalled(t, "FindCustomerByPhone", mock.Anything, mock.Anything, mock.Anything)
		catalog.AssertNotCalled(t, "ServicesExistInOrg", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Status sai thì không được tạo customer", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.Status = "unknown"
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		finder.AssertNotCalled(t, "FindCustomerByPhone", mock.Anything, mock.Anything, mock.Anything)
		finder.AssertNotCalled(t, "CreateCustomer", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "CreateAppointment", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Thiếu phone khách hàng", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.CustomerPhone = "  "
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		finder.AssertNotCalled(t, "FindCustomerByPhone", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Khách mới thiếu tên", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(nil, nil).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.CustomerName = ""
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		finder.AssertNotCalled(t, "CreateCustomer", mock.Anything, mock.Anything)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Giờ kết thúc trước giờ bắt đầu", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.StartTime, req.EndTime = req.EndTime, req.StartTime
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Item sai loại", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		uc := NewAppointmentUseCase(repo, finder, catalog)
		req := validRequest()
		req.Items = []domain.ApoimentItemInput{{ItemType: "combo", Quantity: 1, UnitPrice: 100000}}
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Service không thuộc tổ chức", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		catalog.On("ServicesExistInOrg", mock.Anything, "org-123", mock.Anything).
			Return([]string{"svc-x"}, nil).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
		finder.AssertNotCalled(t, "FindCustomerByPhone", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "CreateAppointment", mock.Anything, mock.Anything, mock.Anything)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Product không thuộc tổ chức", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		catalog.On("ServicesExistInOrg", mock.Anything, "org-123", mock.Anything).
			Return(nil, nil).Once()
		catalog.On("ProductsExistInOrg", mock.Anything, "org-123", mock.Anything).
			Return([]string{"prod-x"}, nil).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Lỗi DB khi kiểm tra catalog", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := new(MockCatalogChecker)
		catalog.On("ServicesExistInOrg", mock.Anything, "org-123", mock.Anything).
			Return(nil, errors.New("db error")).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Lỗi DB khi tìm khách hàng", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(nil, errors.New("db error")).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Lỗi DB khi tạo khách hàng", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(nil, nil).Once()
		finder.On("CreateCustomer", mock.Anything, mock.Anything).
			Return(errors.New("db error")).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
		catalog.AssertExpectations(t)
	})

	t.Run("Thất bại - Lỗi DB khi tạo lịch hẹn", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		finder := new(MockCustomerFinder)
		catalog := catalogOK()
		finder.On("FindCustomerByPhone", mock.Anything, "org-123", "0912345678").
			Return(existingCustomer(), nil).Once()
		repo.On("CreateAppointment", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("db error")).Once()
		uc := NewAppointmentUseCase(repo, finder, catalog)
		res, err := uc.CreateAppointment(withOrgCtx("org-123"), validRequest())

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
		catalog.AssertExpectations(t)
	})
}

func validAppointment(status string) *domain.Appointment {
	return &domain.Appointment{
		ID: "appt-1", OrganizationID: "org-123", Code: "APT-CLIENT-1",
		CustomerID: "cus-1", StaffID: "staff-1",
		StartTime: "2026-10-01T14:00:00+07:00", EndTime: "2026-10-01T15:00:00+07:00",
		Status: status, Source: "online", TotalAmount: 650000,
	}
}

func newUC(repo *MockAppointmentRepository) (*AppointmentUseCase, *MockCustomerFinder, *MockCatalogChecker) {
	finder := new(MockCustomerFinder)
	catalog := new(MockCatalogChecker)
	return NewAppointmentUseCase(repo, finder, catalog), finder, catalog
}

func TestAppointmentUseCase_GetAppointments(t *testing.T) {
	list := []*domain.Appointment{validAppointment("confirmed"), validAppointment("pending")}

	t.Run("Thành công - List phân trang kèm total", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentsByOrgID", mock.Anything, "org-123", "", "", "", 20, 0).
			Return(list, nil).Once()
		repo.On("CountAppointmentsByOrgID", mock.Anything, "org-123", "", "", "").
			Return(int64(2), nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(withOrgCtx("org-123"), 1, 20, "", "", "")

		assert.NoError(t, err)
		assert.NotNil(t, res.Data)
		assert.Len(t, res.Data, 2)
		assert.Equal(t, int64(2), res.Total)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Status lọc không hợp lệ", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(withOrgCtx("org-123"), 1, 20, "unknown", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		repo.AssertNotCalled(t, "GetAppointmentsByOrgID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Khoảng thời gian lọc sai format", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(withOrgCtx("org-123"), 1, 20, "", "not-a-date", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Lỗi DB khi list", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentsByOrgID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(withOrgCtx("org-123"), 1, 20, "", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
	})

	t.Run("Thất bại - Lỗi DB khi count", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentsByOrgID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(list, nil).Once()
		repo.On("CountAppointmentsByOrgID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(int64(0), errors.New("db error")).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(withOrgCtx("org-123"), 1, 20, "", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
	})

	t.Run("Thất bại - Thiếu organization_id", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointments(context.Background(), 1, 20, "", "", "")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})
}

func TestAppointmentUseCase_GetAppointmentDetail(t *testing.T) {
	items := []*domain.AppointmentItem{
		{ID: "item-1", OrganizationID: "org-123", AppointmentID: "appt-1", ItemType: "service", Quantity: 1, UnitPrice: 500000, LineTotal: 500000},
	}

	t.Run("Thành công - Chi tiết kèm items", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("confirmed"), nil).Once()
		repo.On("GetAppointmentItems", mock.Anything, "org-123", "appt-1").
			Return(items, nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointmentDetail(withOrgCtx("org-123"), "appt-1")

		assert.NoError(t, err)
		assert.Equal(t, "appt-1", res.Appointment.ID)
		assert.Len(t, res.Items, 1)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Không tìm thấy lịch hẹn", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "missing").
			Return(nil, domain.ErrAppointmentNotFound).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointmentDetail(withOrgCtx("org-123"), "missing")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
		repo.AssertNotCalled(t, "GetAppointmentItems", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Lỗi DB khi lấy items", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("confirmed"), nil).Once()
		repo.On("GetAppointmentItems", mock.Anything, "org-123", "appt-1").
			Return(nil, errors.New("db error")).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointmentDetail(withOrgCtx("org-123"), "appt-1")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 500)
	})

	t.Run("Thất bại - Thiếu id", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.GetAppointmentDetail(withOrgCtx("org-123"), "  ")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})
}

func TestAppointmentUseCase_RescheduleAppointment(t *testing.T) {
	newStart, newEnd := "2026-10-02T14:00:00+07:00", "2026-10-02T15:00:00+07:00"

	t.Run("Thành công - Đổi lịch pending", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("pending"), nil).Once()
		repo.On("UpdateAppointmentTimes", mock.Anything, "org-123", "appt-1", newStart, newEnd).
			Return(nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.RescheduleAppointment(withOrgCtx("org-123"), "appt-1", newStart, newEnd)

		assert.NoError(t, err)
		assert.Equal(t, newStart, res.StartTime)
		assert.Equal(t, newEnd, res.EndTime)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Lịch completed không được đổi", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("completed"), nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.RescheduleAppointment(withOrgCtx("org-123"), "appt-1", newStart, newEnd)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		repo.AssertNotCalled(t, "UpdateAppointmentTimes", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Giờ mới invalid", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.RescheduleAppointment(withOrgCtx("org-123"), "appt-1", newEnd, newStart)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Không tìm thấy lịch hẹn", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "missing").
			Return(nil, domain.ErrAppointmentNotFound).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.RescheduleAppointment(withOrgCtx("org-123"), "missing", newStart, newEnd)

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
	})
}

func TestAppointmentUseCase_CancelAppointment(t *testing.T) {
	t.Run("Thành công - Hủy lịch confirmed kèm lý do", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("confirmed"), nil).Once()
		repo.On("UpdateAppointmentStatus", mock.Anything, "org-123", "appt-1", "cancelled", "Khách bận").
			Return(nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.CancelAppointment(withOrgCtx("org-123"), "appt-1", "Khách bận")

		assert.NoError(t, err)
		assert.Equal(t, "cancelled", res.Status)
		assert.Equal(t, "Khách bận", res.CancelReason)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Lịch completed không được hủy", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("completed"), nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.CancelAppointment(withOrgCtx("org-123"), "appt-1", "Muộn")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})

	t.Run("Thất bại - Không tìm thấy lịch hẹn", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "missing").
			Return(nil, domain.ErrAppointmentNotFound).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.CancelAppointment(withOrgCtx("org-123"), "missing", "Muộn")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 404)
	})
}

func TestAppointmentUseCase_UpdateAppointmentStatus(t *testing.T) {
	t.Run("Thành công - confirmed sang checked_in", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("confirmed"), nil).Once()
		repo.On("UpdateAppointmentStatus", mock.Anything, "org-123", "appt-1", "checked_in", "").
			Return(nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.UpdateAppointmentStatus(withOrgCtx("org-123"), "appt-1", "checked_in")

		assert.NoError(t, err)
		assert.Equal(t, "checked_in", res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Nhảy cóc pending sang completed", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("pending"), nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.UpdateAppointmentStatus(withOrgCtx("org-123"), "appt-1", "completed")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		repo.AssertNotCalled(t, "UpdateAppointmentStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Trạng thái không tồn tại", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		uc, _, _ := newUC(repo)

		res, err := uc.UpdateAppointmentStatus(withOrgCtx("org-123"), "appt-1", "flying")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
		repo.AssertNotCalled(t, "GetAppointmentByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Lịch cancelled là trạng thái cuối", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("cancelled"), nil).Once()
		uc, _, _ := newUC(repo)

		res, err := uc.UpdateAppointmentStatus(withOrgCtx("org-123"), "appt-1", "confirmed")

		assert.Error(t, err)
		assert.Nil(t, res)
		assertFlatAppError(t, err, 400)
	})
}

func TestAppointmentUseCase_DeleteAppointment(t *testing.T) {
	t.Run("Thành công - Xóa lịch pending", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("pending"), nil).Once()
		repo.On("DeleteAppointment", mock.Anything, "org-123", "appt-1").
			Return(nil).Once()
		uc, _, _ := newUC(repo)

		deletedID, err := uc.DeleteAppointment(withOrgCtx("org-123"), "appt-1")

		assert.NoError(t, err)
		assert.Equal(t, "appt-1", deletedID)
		repo.AssertExpectations(t)
	})

	t.Run("Thất bại - Lịch confirmed còn giá trị, phải hủy thay vì xóa", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "appt-1").
			Return(validAppointment("confirmed"), nil).Once()
		uc, _, _ := newUC(repo)

		deletedID, err := uc.DeleteAppointment(withOrgCtx("org-123"), "appt-1")

		assert.Error(t, err)
		assert.Empty(t, deletedID)
		assertFlatAppError(t, err, 400)
		repo.AssertNotCalled(t, "DeleteAppointment", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Thất bại - Không tìm thấy lịch hẹn", func(t *testing.T) {
		repo := new(MockAppointmentRepository)
		repo.On("GetAppointmentByID", mock.Anything, "org-123", "missing").
			Return(nil, domain.ErrAppointmentNotFound).Once()
		uc, _, _ := newUC(repo)

		deletedID, err := uc.DeleteAppointment(withOrgCtx("org-123"), "missing")

		assert.Error(t, err)
		assert.Empty(t, deletedID)
		assertFlatAppError(t, err, 404)
	})
}
