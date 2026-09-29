package customer

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/daidat02/server/internal/delivery/http/middlewares"
	appointmentDomain "github.com/daidat02/server/internal/domain/appointment"
	domain "github.com/daidat02/server/internal/domain/customer"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/pagination"
	"github.com/daidat02/server/pkg/utils"
)

type CustomerUsecase struct {
	repo    domain.CustomerRepository
	history HistoryStore
}

// HistoryStore lấy lịch sử lịch hẹn của khách.
// PostgresAppointmentRepository thỏa mãn ngầm.
type HistoryStore interface {
	GetAppointmentsByCustomerID(ctx context.Context, orgID, customerID string, limit, offset int) ([]*appointmentDomain.Appointment, error)
	CountAppointmentsByCustomerID(ctx context.Context, orgID, customerID string) (int64, error)
}

func NewCustomerUsecase(repo domain.CustomerRepository, history HistoryStore) *CustomerUsecase {
	return &CustomerUsecase{
		repo:    repo,
		history: history,
	}
}

// CustomerDetail hồ sơ khách + lịch sử lịch hẹn (phân trang riêng).
// Lịch sử mua hàng (orders) sẽ bổ sung khi có module Order.
type CustomerDetail struct {
	Customer     *domain.CustomerResponse                           `json:"customer"`
	Appointments *pagination.Result[*appointmentDomain.Appointment] `json:"appointments"`
}

func toResponse(c *domain.Customer) *domain.CustomerResponse {
	return &domain.CustomerResponse{
		ID:             c.ID,
		OrganizationID: c.OrganizationID,
		FullName:       c.FullName,
		Phone:          c.Phone,
		Gender:         c.Gender,
		BirthDate:      c.BirthDate,
		Note:           c.Note,
		TotalSpent:     c.TotalSpent,
		TotalVisits:    c.TotalVisits,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		DeletedAt:      c.DeletedAt,
	}
}

func (u *CustomerUsecase) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.CustomerResponse, error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := customer.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	customer.ID = utils.NewID()
	customer.OrganizationID = orgID

	if err := u.repo.CreateCustomer(ctx, customer); err != nil {
		log.Println("Error creating customer:", err)
		return nil, apperror.New(500, "Lỗi khi tạo khách hàng", err)
	}
	
	return toResponse(customer), nil
}

func (u *CustomerUsecase) GetCustomersByOrgID(ctx context.Context, page, pageSize int, search string) (*pagination.Result[*domain.CustomerResponse], error) {
	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	params := pagination.Normalize(page, pageSize)
	search = strings.TrimSpace(search)

	customers, err := u.repo.GetCustomersByOrgID(ctx, orgID, search, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy danh sách khách hàng", err)
	}

	var result []*domain.CustomerResponse
	for _,c :=range customers {
		result = append(result, toResponse(c))
	}

	total, err := u.repo.CountCustomersByOrgID(ctx, orgID, search)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm danh sách khách hàng", err)
	}

	return pagination.NewResult(result, params, total), nil
}

func (u *CustomerUsecase) GetCustomerDetail(ctx context.Context, id string, histPage, histPageSize int) (*CustomerDetail, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id khách hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	customer, err := u.repo.GetCustomerByID(ctx, orgID, id)
	if err != nil {
		if errors.Is(err, domain.ErrCustomerNotFound) {
			return nil, apperror.New(404, "Không tìm thấy khách hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi lấy khách hàng", err)
	}

	params := pagination.Normalize(histPage, histPageSize)

	appointments, err := u.history.GetAppointmentsByCustomerID(ctx, orgID, id, params.PageSize, params.Offset())
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi lấy lịch sử lịch hẹn", err)
	}

	total, err := u.history.CountAppointmentsByCustomerID(ctx, orgID, id)
	if err != nil {
		return nil, apperror.New(500, "Lỗi khi đếm lịch sử lịch hẹn", err)
	}

	return &CustomerDetail{
		Customer:     toResponse(customer),
		Appointments: pagination.NewResult(appointments, params, total),
	}, nil
}

func (u *CustomerUsecase) UpdateCustomer(ctx context.Context, id string, customer *domain.Customer) (*domain.CustomerResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperror.New(400, "Thiếu id khách hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return nil, apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := customer.Validate(); err != nil {
		return nil, apperror.New(400, err.Error(), err)
	}
	customer.ID = id
	customer.OrganizationID = orgID

	if err := u.repo.UpdateCustomer(ctx, customer); err != nil {
		log.Println("Error updating customer:", err)
		if errors.Is(err, domain.ErrCustomerNotFound) {
			return nil, apperror.New(404, "Không tìm thấy khách hàng", err)
		}
		return nil, apperror.New(500, "Lỗi khi cập nhật khách hàng", err)
	}

	return toResponse(customer), nil
}

func (u *CustomerUsecase) DeleteCustomer(ctx context.Context, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", apperror.New(400, "Thiếu id khách hàng", nil)
	}

	orgID, err := middlewares.GetOrgID(ctx)
	if err != nil {
		return "", apperror.New(400, "Không tìm thấy organization_id trong context", err)
	}

	if err := u.repo.DeleteCustomer(ctx, orgID, id); err != nil {
		log.Println("Error deleting customer:", err)
		if errors.Is(err, domain.ErrCustomerNotFound) {
			return "", apperror.New(404, "Không tìm thấy khách hàng", err)
		}
		return "", apperror.New(500, "Lỗi khi xóa khách hàng", err)
	}

	return id, nil
}

