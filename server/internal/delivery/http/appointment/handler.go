package appointment

import (
	"strconv"

	domain "github.com/daidat02/server/internal/domain/appointment"
	uc "github.com/daidat02/server/internal/usecase/appointment"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	uc *uc.AppointmentUseCase
}

func NewAppointmentHandler(uc *uc.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{
		uc: uc,
	}
}

// CreateAppointment godoc
// @Summary      Tạo lịch hẹn (tự tìm/tạo khách theo phone)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.AppointmentCreateRequest  true  "Thông tin lịch hẹn (total do server tính từ items)"
// @Success      200   {object}  response.SuccessResponse         "Tạo lịch hẹn thành công"
// @Failure      400   {object}  response.ErrorResponse           "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse           "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse           "Service/product không thuộc tổ chức"
// @Failure      500   {object}  response.ErrorResponse           "Lỗi server"
// @Router       /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req domain.AppointmentCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.CreateAppointment(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo lịch hẹn thành công", result)
}

// GetAppointments godoc
// @Summary      Lấy danh sách lịch hẹn (phân trang + lọc)
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int     false  "Trang (mặc định 1)"
// @Param        page_size query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Param        status    query     string  false  "Lọc trạng thái (pending|confirmed|...)"
// @Param        from      query     string  false  "Từ thời điểm (ISO8601, lọc start_time)"
// @Param        to        query     string  false  "Đến thời điểm (ISO8601, lọc start_time)"
// @Success      200       {object}  response.SuccessResponse  "Danh sách lịch hẹn kèm phân trang"
// @Failure      400       {object}  response.ErrorResponse    "Bộ lọc không hợp lệ"
// @Failure      401       {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /appointments [get]
func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetAppointments(
		c.Request.Context(), page, pageSize,
		c.Query("status"), c.Query("from"), c.Query("to"),
	)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách lịch hẹn", result)
}

// GetAppointmentDetail godoc
// @Summary      Chi tiết lịch hẹn kèm items
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID lịch hẹn"
// @Success      200  {object}  response.SuccessResponse  "Chi tiết lịch hẹn"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy lịch hẹn"
// @Router       /appointments/{id} [get]
func (h *AppointmentHandler) GetAppointmentDetail(c *gin.Context) {
	result, err := h.uc.GetAppointmentDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Chi tiết lịch hẹn", result)
}

// RescheduleAppointment godoc
// @Summary      Đổi lịch hẹn (chỉ pending/confirmed)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                       true  "ID lịch hẹn"
// @Param        body  body      domain.RescheduleRequest     true  "Giờ mới (ISO8601)"
// @Success      200   {object}  response.SuccessResponse     "Đổi lịch thành công"
// @Failure      400   {object}  response.ErrorResponse       "Giờ invalid / trạng thái không cho đổi"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse       "Không tìm thấy lịch hẹn"
// @Router       /appointments/{id}/reschedule [put]
func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	var req domain.RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.RescheduleAppointment(c.Request.Context(), c.Param("id"), req.StartTime, req.EndTime)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Đổi lịch hẹn thành công", result)
}

// CancelAppointment godoc
// @Summary      Hủy lịch hẹn (chỉ pending/confirmed)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string               true  "ID lịch hẹn"
// @Param        body  body      domain.CancelRequest  true  "Lý do hủy"
// @Success      200   {object}  response.SuccessResponse  "Hủy lịch thành công"
// @Failure      400   {object}  response.ErrorResponse    "Trạng thái không cho hủy"
// @Failure      401   {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse    "Không tìm thấy lịch hẹn"
// @Router       /appointments/{id}/cancel [post]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	var req domain.CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.CancelAppointment(c.Request.Context(), c.Param("id"), req.Reason)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Hủy lịch hẹn thành công", result)
}

// UpdateAppointmentStatus godoc
// @Summary      Chuyển trạng thái lịch hẹn (theo state machine)
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                        true  "ID lịch hẹn"
// @Param        body  body      domain.UpdateStatusRequest    true  "Trạng thái mới"
// @Success      200   {object}  response.SuccessResponse      "Cập nhật trạng thái thành công"
// @Failure      400   {object}  response.ErrorResponse        "Chuyển đổi không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse        "Không tìm thấy lịch hẹn"
// @Router       /appointments/{id}/status [patch]
func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	var req domain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.UpdateAppointmentStatus(c.Request.Context(), c.Param("id"), req.Status)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật trạng thái thành công", result)
}

// DeleteAppointment godoc
// @Summary      Xóa cứng lịch hẹn (chỉ pending/cancelled)
// @Tags         Appointments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID lịch hẹn"
// @Success      200  {object}  response.SuccessResponse  "Xóa lịch hẹn thành công"
// @Failure      400  {object}  response.ErrorResponse    "Trạng thái không cho xóa (hãy hủy lịch)"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy lịch hẹn"
// @Router       /appointments/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	deletedID, err := h.uc.DeleteAppointment(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa lịch hẹn thành công", gin.H{"id": deletedID})
}
