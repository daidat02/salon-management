package appointment

import "errors"

const(
	statusPending = "pending"
	statusConfirmed = "confirmed"
	statusCheckedIn = "checked_in"
	statusInProgress = "in_progress"
	statusCompleted = "completed"
	statusCancelled = "cancelled"
	statusNoShow = "no_show"
)

type Appointment struct {
	ID			 string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Code		   string `json:"code"`
	CustomerID     string `json:"customer_id"`
	StaffID        string `json:"staff_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	TotalAmount    float64 `json:"total_amount"`
	Note           string `json:"note"`
	CancelReason   string `json:"cancel_reason"`
	DeviceID       string `json:"device_id"`
	ClientIP       string `json:"client_ip"`
	RequiredDepositAmount float64 `json:"required_deposit_amount"`
	DepositStatus string `json:"deposit_status"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
}

type AppointmentItem struct {
	ID 		  string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	AppointmentID  string  `json:"appointment_id"`
	ItemType      string  `json:"item_type"`
	ServiceID     *string `json:"service_id,omitempty"`
	ProductID     *string `json:"product_id,omitempty"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	TotalPrice    float64 `json:"total_price"`
	LineTotal     float64 `json:"line_total"`
}


func (a *Appointment) Validate() error {
	if a.Status != statusPending && a.Status != statusConfirmed && a.Status != statusCheckedIn && a.Status != statusInProgress && a.Status != statusCompleted && a.Status != statusCancelled && a.Status != statusNoShow {
		return errors.New("Trạng thái cuộc hẹn không hợp lệ")
	}
	return nil
}

func (ai *AppointmentItem) Validate() error {
	
	if ai.ItemType != "service" && ai.ItemType != "product" {
		return errors.New("Loại mục cuộc hẹn không hợp lệ")
	}
	if ai.ItemType == "service" && ai.ServiceID == nil {
		return errors.New("Dịch vụ ID không được để trống khi loại mục là dịch vụ")
	}
	if ai.ItemType == "product" && ai.ProductID == nil {
		return errors.New("Sản phẩm ID không được để trống khi loại mục là sản phẩm")
	}
	return nil
}

// AllowedTransitions định nghĩa state machine vòng đời lịch hẹn.
var AllowedTransitions = map[string][]string{
	statusPending:    {statusConfirmed, statusCancelled, statusNoShow},
	statusConfirmed:  {statusCheckedIn, statusCancelled, statusNoShow},
	statusCheckedIn:  {statusInProgress, statusNoShow},
	statusInProgress: {statusCompleted},
	statusCompleted:  {},
	statusCancelled:  {},
	statusNoShow:     {},
}

// CanTransitionTo kiểm tra chuyển trạng thái có hợp lệ không.
func (a *Appointment) CanTransitionTo(next string) bool {
	for _, s := range AllowedTransitions[a.Status] {
		if s == next {
			return true
		}
	}
	return false
}

// ReschedulableStatuses các trạng thái được phép đổi lịch.
var ReschedulableStatuses = map[string]bool{
	statusPending:   true,
	statusConfirmed: true,
}

// CancellableStatuses các trạng thái được phép hủy lịch.
var CancellableStatuses = map[string]bool{
	statusPending:   true,
	statusConfirmed: true,
}

// DeletableStatuses các trạng thái được phép xóa cứng (không còn giá trị lịch sử).
var DeletableStatuses = map[string]bool{
	statusPending:   true,
	statusCancelled: true,
}