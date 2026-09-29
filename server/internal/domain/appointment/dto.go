package appointment

import "errors"


type ApoimentItemInput struct {
	ItemType string  `json:"item_type" binding:"required,oneof=product service"`
	ProductID string  `json:"product_id"`
	ServiceID string  `json:"service_id"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type AppointmentCreateRequest struct {
	Code         string `json:"code"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	
	StaffID        string `json:"staff_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	TotalAmount    float64 `json:"total_amount"`
	Note           string `json:"note"`
	DeviceID       string `json:"device_id"`
	ClientIP       string `json:"client_ip"`
	RequiredDepositAmount float64 `json:"required_deposit_amount"`
	DepositStatus string `json:"deposit_status"`
	Items []ApoimentItemInput `json:"items"`
}

type RescheduleRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type CancelRequest struct {
	Reason string `json:"reason"`
}

type AppointmentDetailResponse struct {
	Appointment *Appointment       `json:"appointment"`
	Items       []*AppointmentItem `json:"items"`
}


func (req *AppointmentCreateRequest) Validate() error {
	if req.Status != statusPending && req.Status != statusConfirmed && req.Status != statusCheckedIn && req.Status != statusInProgress && req.Status != statusCompleted && req.Status != statusCancelled && req.Status != statusNoShow {
		return errors.New("Trạng thái cuộc hẹn không hợp lệ")
	}
	return nil
}