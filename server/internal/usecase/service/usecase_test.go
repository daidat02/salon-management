package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/service"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

func validService() *domain.Service {
	catID := "cat-1"
	return &domain.Service{
		CategoryID:      &catID,
		Name:            "Cắt tóc nam",
		Description:     "Cắt + gội + sấy",
		DurationMinutes: 45,
		BufferMinutes:   10,
		Price:           150000,
		IsActive:        true,
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

func TestServiceUsecase_CreateService(t *testing.T) {
	tests := []struct {
		name                string
		ctx                 context.Context
		input               *domain.Service
		inputMaterials      []*domain.ServiceMaterialInput
		mockCategoryExists  bool
		mockCategoryErr     error
		expectCategoryCheck bool
		mockMissingProducts []string
		mockProductsErr     error
		expectProductsCheck bool
		mockReturnErr       error
		callRepo            bool
		expectedError       bool
		expectedCode        int
	}{
		{
			name: "Thành công - Tạo dịch vụ có danh mục", ctx: withOrgCtx("org-123"),
			input: validService(), inputMaterials: nil,
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: true, expectedError: false,
		},
		{
			name: "Thành công - Tạo dịch vụ không danh mục", ctx: withOrgCtx("org-123"),
			input: func() *domain.Service {
				s := validService()
				s.CategoryID = nil
				return s
			}(),
			inputMaterials: nil, expectCategoryCheck: false,
			callRepo: true, expectedError: false,
		},
		{
			name: "Thành công - Tạo dịch vụ kèm nguyên vật liệu", ctx: withOrgCtx("org-123"),
			input: validService(),
			inputMaterials: []*domain.ServiceMaterialInput{
				{ProductID: "prod-1", Quantity: 1},
				{ProductID: "prod-2", Quantity: 0.5},
			},
			mockCategoryExists: true, expectCategoryCheck: true,
			mockMissingProducts: nil, expectProductsCheck: true,
			callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			input: validService(), inputMaterials: nil, callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thiếu tên)", ctx: withOrgCtx("org-123"),
			input: func() *domain.Service {
				s := validService()
				s.Name = ""
				return s
			}(),
			inputMaterials: nil, callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thời lượng <= 0)", ctx: withOrgCtx("org-123"),
			input: func() *domain.Service {
				s := validService()
				s.DurationMinutes = 0
				return s
			}(),
			inputMaterials: nil, callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (giá âm)", ctx: withOrgCtx("org-123"),
			input: func() *domain.Service {
				s := validService()
				s.Price = -1000
				return s
			}(),
			inputMaterials: nil, callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Material thiếu product_id", ctx: withOrgCtx("org-123"),
			input: validService(),
			inputMaterials: []*domain.ServiceMaterialInput{
				{ProductID: "", Quantity: 1},
			},
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Material quantity <= 0", ctx: withOrgCtx("org-123"),
			input: validService(),
			inputMaterials: []*domain.ServiceMaterialInput{
				{ProductID: "prod-1", Quantity: 0},
			},
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Danh mục không thuộc tổ chức", ctx: withOrgCtx("org-123"),
			input: validService(), inputMaterials: nil,
			mockCategoryExists: false, expectCategoryCheck: true,
			callRepo: false, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB khi kiểm tra danh mục", ctx: withOrgCtx("org-123"),
			input: validService(), inputMaterials: nil,
			mockCategoryErr: errors.New("db error"), expectCategoryCheck: true,
			callRepo: false, expectedError: true, expectedCode: 500,
		},
		{
			name: "Thất bại - Sản phẩm không thuộc tổ chức", ctx: withOrgCtx("org-123"),
			input: validService(),
			inputMaterials: []*domain.ServiceMaterialInput{
				{ProductID: "prod-missing", Quantity: 1},
			},
			mockCategoryExists: true, expectCategoryCheck: true,
			mockMissingProducts: []string{"prod-missing"}, expectProductsCheck: true,
			callRepo: false, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB khi kiểm tra sản phẩm", ctx: withOrgCtx("org-123"),
			input: validService(),
			inputMaterials: []*domain.ServiceMaterialInput{
				{ProductID: "prod-1", Quantity: 1},
			},
			mockCategoryExists: true, expectCategoryCheck: true,
			mockProductsErr: errors.New("db error"), expectProductsCheck: true,
			callRepo: false, expectedError: true, expectedCode: 500,
		},
		{
			name: "Thất bại - Lỗi DB khi tạo dịch vụ", ctx: withOrgCtx("org-123"),
			input: validService(), inputMaterials: nil,
			mockCategoryExists: true, expectCategoryCheck: true,
			mockReturnErr: errors.New("db insert failed"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockServiceRepository)
			if tt.expectCategoryCheck {
				mockRepo.On("CategoryExistsInOrg", mock.Anything, "org-123", mock.Anything).
					Return(tt.mockCategoryExists, tt.mockCategoryErr)
			}
			if tt.expectProductsCheck {
				mockRepo.On("ProductsExistInOrg", mock.Anything, "org-123", mock.Anything).
					Return(tt.mockMissingProducts, tt.mockProductsErr)
			}
			if tt.callRepo {
				mockRepo.On("CreateService", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewServiceUsecase(mockRepo)
			res, err := uc.CreateService(tt.ctx, tt.input, tt.inputMaterials)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, "org-123", res.OrganizationID)
				assert.Equal(t, tt.input.Name, res.Name)
				assert.Equal(t, tt.input.DurationMinutes, res.DurationMinutes)
				assert.Equal(t, tt.input.Price, res.Price)
				if tt.input.CategoryID == nil {
					assert.Nil(t, res.CategoryID)
				} else {
					assert.Equal(t, *tt.input.CategoryID, *res.CategoryID)
				}
			}
			if !tt.expectCategoryCheck {
				mockRepo.AssertNotCalled(t, "CategoryExistsInOrg", mock.Anything, mock.Anything, mock.Anything)
			}
			if !tt.expectProductsCheck {
				mockRepo.AssertNotCalled(t, "ProductsExistInOrg", mock.Anything, mock.Anything, mock.Anything)
			}
			if !tt.callRepo {
				mockRepo.AssertNotCalled(t, "CreateService", mock.Anything, mock.Anything, mock.Anything)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceUsecase_GetServicesByOrgID(t *testing.T) {
	serviceList := []*domain.Service{
		{ID: "svc-1", OrganizationID: "org-123", Name: "Cắt tóc nam", DurationMinutes: 45, Price: 150000, IsActive: true},
		{ID: "svc-2", OrganizationID: "org-123", Name: "Nhuộm tóc", DurationMinutes: 120, Price: 500000, IsActive: true},
	}

	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		pageSize      int
		mockReturn    []*domain.Service
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
			name: "Thành công - Lấy danh sách phân trang", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, mockReturn: serviceList, mockTotal: 2,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 2, expectedTotal: 2,
		},
		{
			name: "Thành công - Danh sách rỗng", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, mockReturn: []*domain.Service{}, mockTotal: 0,
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
			page: 1, pageSize: 20, mockReturn: serviceList, mockCountErr: errors.New("db count failed"),
			callRepo: true, callCount: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockServiceRepository)
			if tt.callRepo {
				mockRepo.On("GetServicesByOrgID", mock.Anything, "org-123", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockListErr)
			}
			if tt.callCount {
				mockRepo.On("CountServicesByOrgID", mock.Anything, "org-123", mock.Anything, mock.Anything).
					Return(tt.mockTotal, tt.mockCountErr)
			}

			uc := NewServiceUsecase(mockRepo)
			res, err := uc.GetServicesByOrgID(tt.ctx, "", "", "created_at_desc", tt.page, tt.pageSize)

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
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceUsecase_UpdateService(t *testing.T) {
	tests := []struct {
		name                string
		ctx                 context.Context
		id                  string
		input               *domain.Service
		mockCategoryExists  bool
		mockCategoryErr     error
		expectCategoryCheck bool
		mockReturnErr       error
		callRepo            bool
		expectedError       bool
		expectedCode        int
	}{
		{
			name: "Thành công - Cập nhật dịch vụ", ctx: withOrgCtx("org-123"),
			id: "svc-1", input: validService(),
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", input: validService(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "svc-1", input: validService(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thiếu tên)", ctx: withOrgCtx("org-123"),
			id: "svc-1",
			input: func() *domain.Service {
				s := validService()
				s.Name = ""
				return s
			}(),
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Danh mục không thuộc tổ chức", ctx: withOrgCtx("org-123"),
			id: "svc-1", input: validService(),
			mockCategoryExists: false, expectCategoryCheck: true,
			callRepo: false, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Không tìm thấy dịch vụ", ctx: withOrgCtx("org-123"),
			id: "svc-missing", input: validService(), mockReturnErr: domain.ErrServiceNotFound,
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "svc-1", input: validService(), mockReturnErr: errors.New("db error"),
			mockCategoryExists: true, expectCategoryCheck: true,
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockServiceRepository)
			if tt.expectCategoryCheck {
				mockRepo.On("CategoryExistsInOrg", mock.Anything, "org-123", mock.Anything).
					Return(tt.mockCategoryExists, tt.mockCategoryErr)
			}
			if tt.callRepo {
				mockRepo.On("UpdateService", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewServiceUsecase(mockRepo)
			res, err := uc.UpdateService(tt.ctx, tt.id, tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.id, res.ID)
				assert.Equal(t, "org-123", res.OrganizationID)
				assert.Equal(t, tt.input.Name, res.Name)
				assert.Equal(t, tt.input.Price, res.Price)
			}
			if !tt.expectCategoryCheck {
				mockRepo.AssertNotCalled(t, "CategoryExistsInOrg", mock.Anything, mock.Anything, mock.Anything)
			}
			if !tt.callRepo {
				mockRepo.AssertNotCalled(t, "UpdateService", mock.Anything, mock.Anything)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestServiceUsecase_DeleteService(t *testing.T) {
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
			name: "Thành công - Xóa dịch vụ", ctx: withOrgCtx("org-123"),
			id: "svc-1", callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "svc-1", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Không tìm thấy dịch vụ", ctx: withOrgCtx("org-123"),
			id: "svc-missing", mockReturnErr: domain.ErrServiceNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "svc-1", mockReturnErr: errors.New("db error"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockServiceRepository)
			if tt.callRepo {
				mockRepo.On("DeleteService", mock.Anything, "org-123", tt.id).Return(tt.mockReturnErr)
			}

			uc := NewServiceUsecase(mockRepo)
			deletedID, err := uc.DeleteService(tt.ctx, tt.id)

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

func TestServiceUsecase_GetServiceDetails(t *testing.T) {
	now := time.Now()
	catID := "cat-1"
	mockService := &domain.ServiceResponse{
		ID: "svc-1", OrganizationID: "org-123", CategoryID: &catID,
		Name: "Nhuộm phủ bạc", Description: "Nhuộm + hấp dầu",
		DurationMinutes: 120, BufferMinutes: 15, Price: 500000,
		IsActive: true, CreatedAt: &now, UpdatedAt: &now,
	}
	mockMaterials := []*domain.ServiceMaterial{
		{ID: "sm-1", OrganizationID: "org-123", ServiceID: "svc-1", ProductID: "prod-2", Quantity: 1},
		{ID: "sm-2", OrganizationID: "org-123", ServiceID: "svc-1", ProductID: "prod-3", Quantity: 0.5},
	}

	tests := []struct {
		name           string
		ctx            context.Context
		serviceID      string
		mockService    *domain.ServiceResponse
		mockMaterials  []*domain.ServiceMaterial
		mockErr        error
		callRepo       bool
		expectedError  bool
		expectedCode   int
		expectedMatLen int
	}{
		{
			name: "Thành công - Dịch vụ kèm 2 nguyên vật liệu", ctx: withOrgCtx("org-123"),
			serviceID: "svc-1", mockService: mockService, mockMaterials: mockMaterials,
			callRepo: true, expectedError: false, expectedMatLen: 2,
		},
		{
			name: "Thành công - Dịch vụ không có nguyên vật liệu", ctx: withOrgCtx("org-123"),
			serviceID: "svc-2", mockService: mockService, mockMaterials: []*domain.ServiceMaterial{},
			callRepo: true, expectedError: false, expectedMatLen: 0,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			serviceID: "svc-1", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			serviceID: "svc-1", mockErr: errors.New("db connection refused"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
		{
			name: "Thất bại - Không tìm thấy dịch vụ", ctx: withOrgCtx("org-123"),
			serviceID: "svc-missing", mockErr: domain.ErrServiceNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockServiceRepository)
			if tt.callRepo {
				mockRepo.On("GetServiceDetails", mock.Anything, "org-123", tt.serviceID).
					Return(tt.mockService, tt.mockMaterials, tt.mockErr)
			}

			uc := NewServiceUsecase(mockRepo)
			svc, mats, err := uc.GetServiceDetails(tt.ctx, tt.serviceID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, svc)
				assert.Nil(t, mats)
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "GetServiceDetails", mock.Anything, mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, svc)
				assert.Equal(t, "svc-1", svc.ID)
				assert.Equal(t, "org-123", svc.OrganizationID)
				assert.Equal(t, "Nhuộm phủ bạc", svc.Name)
				assert.Equal(t, 120, svc.DurationMinutes)
				assert.Equal(t, 500000.0, svc.Price)
				assert.NotNil(t, mats)
				assert.Len(t, mats, tt.expectedMatLen)
				for _, m := range mats {
					assert.Equal(t, "svc-1", m.ServiceID)
					assert.Equal(t, "org-123", m.OrganizationID)
					assert.Greater(t, m.Quantity, 0.0)
				}
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
