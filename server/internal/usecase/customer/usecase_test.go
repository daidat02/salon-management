package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/customer"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

func validCustomer() *domain.Customer {
	birth := time.Date(1995, 5, 20, 0, 0, 0, 0, time.UTC)
	return &domain.Customer{
		FullName:  "Nguyễn Thị C",
		Phone:     "0912345678",
		Gender:    domain.GenderFemale,
		BirthDate: &birth,
		Note:      "Khách VIP",
	}
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

func TestCustomerUsecase_CreateCustomer(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		input         *domain.Customer
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name:          "Thành công - Tạo khách hàng mới hợp lệ",
			ctx:           withOrgCtx("org-123"),
			input:         validCustomer(),
			mockReturnErr: nil,
			callRepo:      true,
			expectedError: false,
		},
		{
			name:          "Thất bại - Thiếu organization_id trong context",
			ctx:           context.Background(),
			input:         validCustomer(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu tên)",
			ctx:  withOrgCtx("org-123"),
			input: func() *domain.Customer {
				c := validCustomer()
				c.FullName = ""
				return c
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu phone)",
			ctx:  withOrgCtx("org-123"),
			input: func() *domain.Customer {
				c := validCustomer()
				c.Phone = ""
				return c
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (giới tính không hợp lệ)",
			ctx:  withOrgCtx("org-123"),
			input: func() *domain.Customer {
				c := validCustomer()
				c.Gender = "unknown"
				return c
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Trùng phone trong org (unique)",
			ctx:           withOrgCtx("org-123"),
			input:         validCustomer(),
			mockReturnErr: errors.New("duplicate key value violates unique constraint"),
			callRepo:      true,
			expectedError: true,
			expectedCode:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			if tt.callRepo {
				mockRepo.On("CreateCustomer", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewCustomerUsecase(mockRepo, new(MockHistoryStore))
			res, err := uc.CreateCustomer(tt.ctx, tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "CreateCustomer", mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, "org-123", res.OrganizationID)
				assert.Equal(t, tt.input.FullName, res.FullName)
				assert.Equal(t, tt.input.Phone, res.Phone)
				assert.Equal(t, tt.input.Gender, res.Gender)
				assert.Equal(t, tt.input.Note, res.Note)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestCustomerUsecase_GetCustomersByOrgID(t *testing.T) {
	now := time.Now()
	customerList := []*domain.Customer{
		{ID: "cus-1", OrganizationID: "org-123", FullName: "Nguyễn Thị C", Phone: "0912345678", Gender: domain.GenderFemale, TotalSpent: 1500000, TotalVisits: 5, CreatedAt: &now, UpdatedAt: &now},
		{ID: "cus-2", OrganizationID: "org-123", FullName: "Lê Văn D", Phone: "0987654321", Gender: domain.GenderMale, TotalSpent: 500000, TotalVisits: 2, CreatedAt: &now, UpdatedAt: &now},
	}

	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		pageSize      int
		search        string
		mockReturn    []*domain.Customer
		mockTotal     int64
		mockListErr   error
		mockCountErr  error
		callRepo      bool
		callCount     bool
		expectedError bool
		expectedCode  int
		expectedLen   int
		expectedTotal int64
	}{
		{
			name: "Thành công - Lấy danh sách khách hàng phân trang", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, search: "", mockReturn: customerList, mockTotal: 2,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 2, expectedTotal: 2,
		},
		{
			name: "Thành công - Tìm kiếm theo tên", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, search: "Nguyễn", mockReturn: customerList[:1], mockTotal: 1,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 1, expectedTotal: 1,
		},
		{
			name: "Thành công - Tìm kiếm theo phone", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, search: "0912", mockReturn: customerList[:1], mockTotal: 1,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 1, expectedTotal: 1,
		},
		{
			name: "Thành công - Danh sách rỗng", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, search: "", mockReturn: []*domain.Customer{}, mockTotal: 0,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 0, expectedTotal: 0,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			page: 1, pageSize: 20, callRepo: false, callCount: false,
			expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Lỗi DB khi lấy danh sách", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, mockListErr: errors.New("db connection refused"),
			callRepo: true, callCount: false, expectedError: true, expectedCode: 500,
		},
		{
			name: "Thất bại - Lỗi DB khi đếm", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, mockReturn: customerList, mockCountErr: errors.New("db count failed"),
			callRepo: true, callCount: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			if tt.callRepo {
				mockRepo.On("GetCustomersByOrgID", mock.Anything, "org-123", tt.search, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockListErr)
			}
			if tt.callCount {
				mockRepo.On("CountCustomersByOrgID", mock.Anything, "org-123", tt.search).
					Return(tt.mockTotal, tt.mockCountErr)
			}

			uc := NewCustomerUsecase(mockRepo, new(MockHistoryStore))
			res, err := uc.GetCustomersByOrgID(tt.ctx, tt.page, tt.pageSize, tt.search)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotNil(t, res.Data)
				assert.Len(t, res.Data, tt.expectedLen)
				assert.Equal(t, tt.expectedTotal, res.Total)
				for i := range res.Data {
					assert.Equal(t, "org-123", res.Data[i].OrganizationID)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerUsecase_UpdateCustomer(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		id            string
		input         *domain.Customer
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name: "Thành công - Cập nhật khách hàng", ctx: withOrgCtx("org-123"),
			id: "cus-1", input: validCustomer(), callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", input: validCustomer(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "cus-1", input: validCustomer(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thiếu phone)", ctx: withOrgCtx("org-123"),
			id: "cus-1",
			input: func() *domain.Customer {
				c := validCustomer()
				c.Phone = ""
				return c
			}(),
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Không tìm thấy khách hàng", ctx: withOrgCtx("org-123"),
			id: "cus-missing", input: validCustomer(), mockReturnErr: domain.ErrCustomerNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "cus-1", input: validCustomer(), mockReturnErr: errors.New("db error"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			if tt.callRepo {
				mockRepo.On("UpdateCustomer", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewCustomerUsecase(mockRepo, new(MockHistoryStore))
			res, err := uc.UpdateCustomer(tt.ctx, tt.id, tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.id, res.ID)
				assert.Equal(t, "org-123", res.OrganizationID)
				assert.Equal(t, tt.input.FullName, res.FullName)
				assert.Equal(t, tt.input.Phone, res.Phone)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerUsecase_DeleteCustomer(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		id            string
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name: "Thành công - Xóa mềm khách hàng", ctx: withOrgCtx("org-123"),
			id: "cus-1", callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "cus-1", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Không tìm thấy khách hàng", ctx: withOrgCtx("org-123"),
			id: "cus-missing", mockReturnErr: domain.ErrCustomerNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "cus-1", mockReturnErr: errors.New("db error"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			if tt.callRepo {
				mockRepo.On("DeleteCustomer", mock.Anything, "org-123", tt.id).Return(tt.mockReturnErr)
			}

			uc := NewCustomerUsecase(mockRepo, new(MockHistoryStore))
			deletedID, err := uc.DeleteCustomer(tt.ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, deletedID)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, deletedID)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
