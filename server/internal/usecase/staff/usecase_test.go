package staff

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/daidat02/server/internal/domain/staff"
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/daidat02/server/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// withOrgCtx giả lập middleware VerifyToken: nhét org_id vào context.
func withOrgCtx(orgID string) context.Context {
	return context.WithValue(context.Background(), middlewares.OrgIdKey, orgID)
}

// validStaff trả về input đầy đủ, pass hết Validate():
// Code, FullName, Phone, Position, CommissionRate 0-100, Status active/inactive.
func validStaff() *domain.Staff {
	return &domain.Staff{
		Code:           "NV001",
		FullName:       "Nguyễn Văn A",
		Phone:          "0901234567",
		Position:       "Thợ cắt tóc",
		CommissionRate: 10,
		Status:         "active",
		HireDate:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}

// assertFlatAppError kiểm tra err là *AppError đúng 1 lớp (không lồng nhau).
func assertFlatAppError(t *testing.T, err error, expectedCode int) {
	t.Helper()
	var appErr *apperror.AppError
	assert.True(t, errors.As(err, &appErr), "Lỗi phải là *apperror.AppError")
	if appErr != nil {
		assert.Equal(t, expectedCode, appErr.StatusCode, "Sai StatusCode")
		// Chống lồng: Err gốc bên trong KHÔNG được là *AppError nữa
		var inner *apperror.AppError
		assert.False(t, errors.As(appErr.Err, &inner), "Lỗi bị lồng 2 lớp AppError")
	}
}

func TestStaffUsecase_CreateStaff(f *testing.T) {
	// Định nghĩa các kịch bản test (Table-Driven Test)
	tests := []struct {
		name          string
		ctx           context.Context
		inputStaff    *domain.Staff
		mockReturnErr error
		callRepo      bool // repo có được gọi không (validate/context fail thì không)
		expectedError bool
		expectedCode  int // StatusCode kỳ vọng của *apperror.AppError
	}{
		{
			name:          "Thành công - Tạo nhân viên mới hợp lệ",
			ctx:           withOrgCtx("org-123"),
			inputStaff:    validStaff(),
			mockReturnErr: nil,
			callRepo:      true,
			expectedError: false,
		},
		{
			name: "Thất bại - Lỗi từ database (trùng mã nhân viên)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.Code = "NV001"
				return s
			}(),
			mockReturnErr: errors.New("database duplicate error"),
			callRepo:      true,
			expectedError: true,
			expectedCode:  500,
		},
		{
			name: "Thất bại - Thiếu organization_id trong context",
			ctx:  context.Background(), // không có org_id
			inputStaff: func() *domain.Staff {
				return validStaff()
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - organization_id rỗng trong context",
			ctx:  withOrgCtx(""),
			inputStaff: func() *domain.Staff {
				return validStaff()
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu mã nhân viên)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.Code = ""
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu tên đầy đủ)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.FullName = ""
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu số điện thoại)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.Phone = ""
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu vị trí)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.Position = ""
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (hoa hồng vượt 100)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.CommissionRate = 150
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (trạng thái không hợp lệ)",
			ctx:  withOrgCtx("org-123"),
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.Status = "unknown"
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
	}

	// Chạy vòng lặp qua từng kịch bản
	for _, tt := range tests {
		f.Run(tt.name, func(t *testing.T) {
			// 1. Khởi tạo Mock Repository
			mockRepo := new(MockStaffRepository)

			// 2. Thiết lập kỳ vọng: dùng mock.Anything cho cả ctx lẫn staff
			// vì usecase mutate staff.ID/staff.OrganizationID trước khi gọi repo.
			if tt.callRepo {
				mockRepo.On("CreateStaff", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			// 3. Khởi tạo Usecase THẬT với mock repo vừa tạo
			uc := NewStaffUsecase(mockRepo)
			// 4. Hứng đủ 2 giá trị trả về (res và err)
			res, err := uc.CreateStaff(tt.ctx, tt.inputStaff)

			// 5. Kiểm tra kết quả bằng testify/assert
			if tt.expectedError {
				assert.Error(t, err, "Kỳ vọng phải trả về lỗi nhưng lại không có")
				assert.Nil(t, res, "Khi có lỗi thì kết quả trả về phải là nil")
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "CreateStaff", mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err, "Không kỳ vọng có lỗi nhưng lại phát sinh lỗi")
				assert.NotNil(t, res, "Khi thành công thì phải trả về kết quả")
				assert.NotEmpty(t, res.ID, "ID phải được sinh ra")
				assert.Equal(t, "org-123", res.OrganizationID, "OrganizationID phải lấy từ context")
				assert.Equal(t, tt.inputStaff.Code, res.Code)
				assert.Equal(t, tt.inputStaff.FullName, res.FullName)
				assert.Equal(t, tt.inputStaff.Phone, res.Phone)
				assert.Equal(t, tt.inputStaff.Position, res.Position)
				assert.Equal(t, tt.inputStaff.CommissionRate, res.CommissionRate)
				assert.Equal(t, tt.inputStaff.Status, res.Status)
				assert.True(t, tt.inputStaff.HireDate.Equal(res.HireDate))
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestStaffUsecase_GetStaffByOrgID(f *testing.T) {
	hireDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	staffList := []*domain.Staff{
		{ID: "staff-1", OrganizationID: "org-123", Code: "NV001", FullName: "Nguyễn Văn A", Phone: "0901234567", Position: "Thợ cắt tóc", CommissionRate: 10, Status: "active", HireDate: hireDate},
		{ID: "staff-2", OrganizationID: "org-123", Code: "NV002", FullName: "Trần Thị B", Phone: "0907654321", Position: "Lễ tân", CommissionRate: 5, Status: "inactive", HireDate: hireDate},
	}

	tests := []struct {
		name          string
		ctx           context.Context
		page          int
		pageSize      int
		mockReturn    []*domain.Staff
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
			name:          "Thành công - Lấy danh sách nhân viên theo org (page 1)",
			ctx:           withOrgCtx("org-123"),
			page:          1,
			pageSize:      20,
			mockReturn:    staffList,
			mockTotal:     2,
			mockListErr:   nil,
			mockCountErr:  nil,
			callRepo:      true,
			callCount:     true,
			expectedError: false,
			expectedLen:   2,
			expectedTotal: 2,
		},
		{
			name:          "Thành công - Org chưa có nhân viên (list rỗng)",
			ctx:           withOrgCtx("org-123"),
			page:          1,
			pageSize:      20,
			mockReturn:    []*domain.Staff{},
			mockTotal:     0,
			mockListErr:   nil,
			mockCountErr:  nil,
			callRepo:      true,
			callCount:     true,
			expectedError: false,
			expectedLen:   0,
			expectedTotal: 0,
		},
		{
			name:          "Thành công - Page vượt quá vẫn normalize và trả envelope",
			ctx:           withOrgCtx("org-123"),
			page:          5,
			pageSize:      1,
			mockReturn:    []*domain.Staff{},
			mockTotal:     2,
			mockListErr:   nil,
			mockCountErr:  nil,
			callRepo:      true,
			callCount:     true,
			expectedError: false,
			expectedLen:   0,
			expectedTotal: 2,
		},
		{
			name:          "Thất bại - Thiếu organization_id trong context",
			ctx:           context.Background(),
			page:          1,
			pageSize:      20,
			mockReturn:    nil,
			mockTotal:     0,
			mockListErr:   nil,
			mockCountErr:  nil,
			callRepo:      false,
			callCount:     false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Lỗi database khi lấy danh sách",
			ctx:           withOrgCtx("org-123"),
			page:          1,
			pageSize:      20,
			mockReturn:    nil,
			mockTotal:     0,
			mockListErr:   errors.New("db connection refused"),
			mockCountErr:  nil,
			callRepo:      true,
			callCount:     false,
			expectedError: true,
			expectedCode:  500,
		},
		{
			name:          "Thất bại - Lỗi database khi đếm danh sách",
			ctx:           withOrgCtx("org-123"),
			page:          1,
			pageSize:      20,
			mockReturn:    staffList,
			mockTotal:     0,
			mockListErr:   nil,
			mockCountErr:  errors.New("db count failed"),
			callRepo:      true,
			callCount:     true,
			expectedError: true,
			expectedCode:  500,
		},
	}

	for _, tt := range tests {
		f.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockStaffRepository)
			if tt.callRepo {
				mockRepo.On("GetStaffByOrgID", mock.Anything, "org-123", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockListErr)
			}
			if tt.callCount {
				mockRepo.On("CountStaffsByOrgID", mock.Anything, "org-123", mock.Anything, mock.Anything).
					Return(tt.mockTotal, tt.mockCountErr)
			}

			uc := NewStaffUsecase(mockRepo)
			res, err := uc.GetStaffByOrgID(tt.ctx, "", "", "created_at_desc", tt.page, tt.pageSize)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "GetStaffByOrgID", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotNil(t, res.Data)
				assert.Len(t, res.Data, tt.expectedLen)
				assert.Equal(t, tt.expectedTotal, res.Total)
				for i := range res.Data {
					assert.Equal(t, "org-123", res.Data[i].OrganizationID)
				}
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestStaffUsecase_UpdateStaff(f *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		id            string
		inputStaff    *domain.Staff
		mockReturnErr error
		callRepo      bool
		expectedError bool
		expectedCode  int
	}{
		{
			name:          "Thành công - Cập nhật nhân viên hợp lệ",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-1",
			inputStaff:    validStaff(),
			mockReturnErr: nil,
			callRepo:      true,
			expectedError: false,
		},
		{
			name:          "Thất bại - Thiếu id nhân viên",
			ctx:           withOrgCtx("org-123"),
			id:            "",
			inputStaff:    validStaff(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Thiếu organization_id trong context",
			ctx:           context.Background(),
			id:            "staff-1",
			inputStaff:    validStaff(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (thiếu tên đầy đủ)",
			ctx:  withOrgCtx("org-123"),
			id:   "staff-1",
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.FullName = ""
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name: "Thất bại - Validate (hoa hồng vượt 100)",
			ctx:  withOrgCtx("org-123"),
			id:   "staff-1",
			inputStaff: func() *domain.Staff {
				s := validStaff()
				s.CommissionRate = 120
				return s
			}(),
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Không tìm thấy nhân viên (sai id / khác org / đã xóa)",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-missing",
			inputStaff:    validStaff(),
			mockReturnErr: domain.ErrStaffNotFound,
			callRepo:      true,
			expectedError: true,
			expectedCode:  404,
		},
		{
			name:          "Thất bại - Lỗi database khi cập nhật",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-1",
			inputStaff:    validStaff(),
			mockReturnErr: errors.New("db connection refused"),
			callRepo:      true,
			expectedError: true,
			expectedCode:  500,
		},
	}

	for _, tt := range tests {
		f.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockStaffRepository)
			if tt.callRepo {
				mockRepo.On("UpdateStaff", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			uc := NewStaffUsecase(mockRepo)
			res, err := uc.UpdateStaff(tt.ctx, tt.id, tt.inputStaff)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "UpdateStaff", mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.id, res.ID, "ID phải giữ nguyên id trên path")
				assert.Equal(t, "org-123", res.OrganizationID, "OrganizationID phải lấy từ context")
				assert.Equal(t, tt.inputStaff.Code, res.Code)
				assert.Equal(t, tt.inputStaff.FullName, res.FullName)
				assert.Equal(t, tt.inputStaff.Phone, res.Phone)
				assert.Equal(t, tt.inputStaff.Position, res.Position)
				assert.Equal(t, tt.inputStaff.CommissionRate, res.CommissionRate)
				assert.Equal(t, tt.inputStaff.Status, res.Status)
				assert.True(t, tt.inputStaff.HireDate.Equal(res.HireDate))
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestStaffUsecase_DeleteStaff(f *testing.T) {
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
			name:          "Thành công - Xóa mềm nhân viên",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-1",
			mockReturnErr: nil,
			callRepo:      true,
			expectedError: false,
		},
		{
			name:          "Thất bại - Thiếu id nhân viên",
			ctx:           withOrgCtx("org-123"),
			id:            "",
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Thiếu organization_id trong context",
			ctx:           context.Background(),
			id:            "staff-1",
			mockReturnErr: nil,
			callRepo:      false,
			expectedError: true,
			expectedCode:  400,
		},
		{
			name:          "Thất bại - Không tìm thấy nhân viên (đã xóa / khác org)",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-missing",
			mockReturnErr: domain.ErrStaffNotFound,
			callRepo:      true,
			expectedError: true,
			expectedCode:  404,
		},
		{
			name:          "Thất bại - Lỗi database khi xóa",
			ctx:           withOrgCtx("org-123"),
			id:            "staff-1",
			mockReturnErr: errors.New("db connection refused"),
			callRepo:      true,
			expectedError: true,
			expectedCode:  500,
		},
	}

	for _, tt := range tests {
		f.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockStaffRepository)
			if tt.callRepo {
				mockRepo.On("DeleteStaff", mock.Anything, "org-123", tt.id).Return(tt.mockReturnErr)
			}

			uc := NewStaffUsecase(mockRepo)
			deletedID, err := uc.DeleteStaff(tt.ctx, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, deletedID)
				assertFlatAppError(t, err, tt.expectedCode)
				if tt.callRepo {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "DeleteStaff", mock.Anything, mock.Anything, mock.Anything)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, deletedID)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
