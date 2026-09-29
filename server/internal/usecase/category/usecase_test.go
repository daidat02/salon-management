package category

import (
	"context"
	"errors"
	"testing"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/category"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

func validCategory() *domain.Category {
	return &domain.Category{
		Type:        domain.TypeService,
		Name:        "Cắt tóc",
		Description: "Các dịch vụ cắt tóc",
		IsActive:    true,
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

func TestCategoryUsecase_CreateCategory(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		input         *domain.Category
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name: "Thành công - Tạo danh mục dịch vụ", ctx: withOrgCtx("org-123"),
			input: validCategory(), callRepo: true, expectedError: false,
		},
		{
			name: "Thành công - Tạo danh mục sản phẩm", ctx: withOrgCtx("org-123"),
			input: func() *domain.Category {
				c := validCategory()
				c.Type = domain.TypeProduct
				c.Name = "Dầu gội"
				return c
			}(),
			callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			input: validCategory(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (type không hợp lệ)", ctx: withOrgCtx("org-123"),
			input: func() *domain.Category {
				c := validCategory()
				c.Type = "unknown"
				return c
			}(),
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thiếu tên)", ctx: withOrgCtx("org-123"),
			input: func() *domain.Category {
				c := validCategory()
				c.Name = ""
				return c
			}(),
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Trùng tên trong org (unique)", ctx: withOrgCtx("org-123"),
			input: validCategory(), mockReturnErr: errors.New("duplicate key value violates unique constraint"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			if tt.callRepo {
				mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewCategoryUsecase(mockRepo)
			res, err := uc.CreateCategory(tt.ctx, tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, "org-123", res.OrganizationID)
				assert.Equal(t, tt.input.Type, res.Type)
				assert.Equal(t, tt.input.Name, res.Name)
				assert.Equal(t, tt.input.Description, res.Description)
				assert.Equal(t, tt.input.IsActive, res.IsActive)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCategoryUsecase_GetCategoriesByOrgID(t *testing.T) {
	categoryList := []*domain.Category{
		{ID: "cat-1", OrganizationID: "org-123", Type: domain.TypeService, Name: "Cắt tóc", IsActive: true},
		{ID: "cat-2", OrganizationID: "org-123", Type: domain.TypeProduct, Name: "Dầu gội", IsActive: true},
	}

	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		pageSize      int
		mockReturn    []*domain.Category
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
			page: 1, pageSize: 20, mockReturn: categoryList, mockTotal: 2,
			callRepo: true, callCount: true, expectedError: false, expectedLen: 2, expectedTotal: 2,
		},
		{
			name: "Thành công - Danh sách rỗng", ctx: withOrgCtx("org-123"),
			page: 1, pageSize: 20, mockReturn: []*domain.Category{}, mockTotal: 0,
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
			page: 1, pageSize: 20, mockReturn: categoryList, mockCountErr: errors.New("db count failed"),
			callRepo: true, callCount: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			if tt.callRepo {
				mockRepo.On("GetCategoriesByOrgID", mock.Anything, "org-123", mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockListErr)
			}
			if tt.callCount {
				mockRepo.On("CountCategoriesByOrgID", mock.Anything, "org-123").
					Return(tt.mockTotal, tt.mockCountErr)
			}

			uc := NewCategoryUsecase(mockRepo)
			res, err := uc.GetCategoriesByOrgID(tt.ctx, tt.page, tt.pageSize)

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

func TestCategoryUsecase_UpdateCategory(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		id            string
		input         *domain.Category
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name: "Thành công - Cập nhật danh mục", ctx: withOrgCtx("org-123"),
			id: "cat-1", input: validCategory(), callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", input: validCategory(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "cat-1", input: validCategory(), callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Validate (thiếu tên)", ctx: withOrgCtx("org-123"),
			id: "cat-1",
			input: func() *domain.Category {
				c := validCategory()
				c.Name = ""
				return c
			}(),
			callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Không tìm thấy danh mục", ctx: withOrgCtx("org-123"),
			id: "cat-missing", input: validCategory(), mockReturnErr: domain.ErrCategoryNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "cat-1", input: validCategory(), mockReturnErr: errors.New("db error"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			if tt.callRepo {
				mockRepo.On("UpdateCategory", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewCategoryUsecase(mockRepo)
			res, err := uc.UpdateCategory(tt.ctx, tt.id, tt.input)

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
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCategoryUsecase_DeleteCategory(t *testing.T) {
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
			name: "Thành công - Xóa danh mục", ctx: withOrgCtx("org-123"),
			id: "cat-1", callRepo: true, expectedError: false,
		},
		{
			name: "Thất bại - Thiếu id", ctx: withOrgCtx("org-123"),
			id: "", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Thiếu organization_id", ctx: context.Background(),
			id: "cat-1", callRepo: false, expectedError: true, expectedCode: 400,
		},
		{
			name: "Thất bại - Không tìm thấy danh mục", ctx: withOrgCtx("org-123"),
			id: "cat-missing", mockReturnErr: domain.ErrCategoryNotFound,
			callRepo: true, expectedError: true, expectedCode: 404,
		},
		{
			name: "Thất bại - Lỗi DB", ctx: withOrgCtx("org-123"),
			id: "cat-1", mockReturnErr: errors.New("db error"),
			callRepo: true, expectedError: true, expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			if tt.callRepo {
				mockRepo.On("DeleteCategory", mock.Anything, "org-123", tt.id).Return(tt.mockReturnErr)
			}

			uc := NewCategoryUsecase(mockRepo)
			deletedID, err := uc.DeleteCategory(tt.ctx, tt.id)

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
