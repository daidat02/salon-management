package staff

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	domain "github.com/daidat02/server/internal/domain/staff"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

type StaffUsecase struct {
	userRepo domain.StaffRepository
}

func NewStaffUsecase(repo domain.StaffRepository) *StaffUsecase{
	return &StaffUsecase{
		userRepo: repo,
	}
}

func toResponse(s *domain.Staff) *domain.StaffResponse {
	return &domain.StaffResponse{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		Code:           s.Code,
		FullName:       s.FullName,
		Phone:          s.Phone,
		Position:       s.Position,
		CommissionRate: s.CommissionRate,
		Status:         s.Status,
		HireDate:       s.HireDate,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		DeletedAt:      s.DeletedAt,
	}
}


func (u *StaffUsecase) CreateStaff(ctx context.Context, staff *domain.Staff) (*domain.StaffResponse, error) {

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	} 

	if err :=staff.Validate(); err != nil {
		// Validate trả lỗi thường -> tạo AppError 400 duy nhất tại đây.
		return nil, apperror.New(400, err.Error(), err)
	}
	staff.ID = utils.NewID()
	staff.OrganizationID = orgID

	err = u.userRepo.CreateStaff(ctx, staff)
	if err != nil {
		log.Println("Error creating staff:", err)
		// Thiếu org_id -> 400, còn lại -> 500. Mỗi nhánh chỉ New 1 lần.
		if errors.Is(err, domain.ErrMissingOrganizationID) {
			return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
		}
		return nil, apperror.New(500, "Lỗi khi tạo nhân viên", err)
	}
	
	return toResponse(staff), nil
}


func (u *StaffUsecase) GetStaffByOrgID(ctx context.Context, search string, status string, sort string, page, pageSize int) (*pagination.Result[*domain.StaffResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)

	staffs, err := u.userRepo.GetStaffByOrgID(ctx, orgID, search, status, sort, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách nhân viên", err)
	}

	staffsResponse := []*domain.StaffResponse{}

	for _, staff := range staffs {
		staffsResponse = append(staffsResponse, toResponse(staff))
	}

	total, err := u.userRepo.CountStaffsByOrgID(ctx, orgID, search, status)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách nhân viên", err)
	}

	return pagination.NewResult(staffsResponse, params, total), nil
}

func (u *StaffUsecase) UpdateStaff(ctx context.Context, id string, staff *domain.Staff) (*domain.StaffResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id nhân viên", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := staff.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	staff.ID = id
	staff.OrganizationID = orgID

	if err := u.userRepo.UpdateStaff(ctx, staff); err != nil {
		log.Println("Error updating staff:", err)
		if errors.Is(err, domain.ErrStaffNotFound) {
			return nil, apperror.New(404, "Không tìm thấy nhân viên", err)
		}
		return nil, apperror.New(500, "Lỗi khi cập nhật nhân viên", err)
	}

	return toResponse(staff), nil
}

func (u *StaffUsecase) DeleteStaff(ctx context.Context, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", apperror.New(400, "Thiếu id nhân viên", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return "", apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := u.userRepo.DeleteStaff(ctx, orgID, id); err != nil {
		log.Println("Error deleting staff:", err)
		if errors.Is(err, domain.ErrStaffNotFound) {
			return "", apperror.New(404, "Không tìm thấy nhân viên", err)
		}
		return "", apperror.New(500, "Lỗi khi xóa nhân viên", err)
	}

	return id, nil
}

