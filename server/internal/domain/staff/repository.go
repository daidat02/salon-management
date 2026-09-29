package staff

import (
	"context"
	"errors"
)

var ErrMissingOrganizationID = errors.New("missing organization_id in context")

// ErrStaffNotFound dùng để usecase map đúng 404 khi UPDATE/DELETE không chạm dòng nào
// (sai id, khác organization, hoặc đã bị xóa mềm).
var ErrStaffNotFound = errors.New("staff not found")

type StaffRepository interface {
	CreateStaff(ctx context.Context, staff *Staff) error
	GetStaffByOrgID(ctx context.Context, orgID string, search string, status string, sort string, limit, offset int) ([]*Staff, error)
	CountStaffsByOrgID(ctx context.Context, orgID string, search string, status string) (int64, error)
	UpdateStaff(ctx context.Context, staff *Staff) error
	DeleteStaff(ctx context.Context, orgID, id string) error
}