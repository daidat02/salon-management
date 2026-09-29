package customer

import (
	"log"
	"strconv"
	"time"

	domain "github.com/daidat02/server/internal/domain/customer"
	uc "github.com/daidat02/server/internal/usecase/customer"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	uc *uc.CustomerUsecase
}

func NewCustomerHandler(uc *uc.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		uc: uc,
	}
}

// CreateCustomer godoc
// @Summary      Tạo khách hàng
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateCustomerRequest  true  "Thông tin khách hàng"
// @Success      200   {object}  response.SuccessResponse      "Tạo khách hàng thành công"
// @Failure      400   {object}  response.ErrorResponse        "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Failure      500   {object}  response.ErrorResponse        "Lỗi server (vd trùng số điện thoại)"
// @Router       /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req domain.CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		log.Println("Error binding JSON:", err)
		c.Abort()
		return
	}
	var birthDateParsed *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
    // 1. Phải hứng cả 2 giá trị trả về: parsedTime và err
    parsedTime, err := time.Parse("2006-01-02", *req.BirthDate)
    if err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
    }
    
    // 3. Gán địa chỉ nếu parse thành công
    birthDateParsed = &parsedTime
}
	customer := &domain.Customer{
		FullName:  req.FullName,
		Phone:     req.Phone,
		Gender:    req.Gender,
		BirthDate: birthDateParsed,
		Note:      req.Note,
	}

	result, err := h.uc.CreateCustomer(c.Request.Context(), customer)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo khách hàng thành công", result)
}

// GetCustomersByOrgID godoc
// @Summary      Lấy danh sách khách hàng (phân trang + tìm kiếm)
// @Tags         Customers
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Param        search     query     string  false  "Tìm theo tên hoặc số điện thoại"
// @Success      200        {object}  response.SuccessResponse  "Danh sách khách hàng kèm phân trang"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /customers [get]
func (h *CustomerHandler) GetCustomersByOrgID(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetCustomersByOrgID(c.Request.Context(), page, pageSize, c.Query("search"))
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách khách hàng", result)
}

// GetCustomerDetail godoc
// @Summary      Chi tiết khách hàng kèm lịch sử lịch hẹn
// @Tags         Customers
// @Produce      json
// @Security     BearerAuth
// @Param        id               path   string  true   "ID khách hàng"
// @Param        hist_page        query  int     false  "Trang lịch sử (mặc định 1)"
// @Param        hist_page_size   query  int     false  "Số dòng lịch sử/trang (mặc định 10, tối đa 100)"
// @Success      200              {object}  response.SuccessResponse  "Hồ sơ + lịch sử lịch hẹn"
// @Failure      400              {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401              {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404              {object}  response.ErrorResponse    "Không tìm thấy khách hàng"
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetCustomerDetail(c *gin.Context) {
	histPage, _ := strconv.Atoi(c.DefaultQuery("hist_page", "1"))
	histPageSize, _ := strconv.Atoi(c.DefaultQuery("hist_page_size", "10"))

	result, err := h.uc.GetCustomerDetail(c.Request.Context(), c.Param("id"), histPage, histPageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Chi tiết khách hàng", result)
}

// UpdateCustomer godoc
// @Summary      Cập nhật khách hàng
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                       true  "ID khách hàng"
// @Param        body  body      domain.UpdateCustomerRequest  true  "Thông tin cập nhật"
// @Success      200   {object}  response.SuccessResponse      "Cập nhật khách hàng thành công"
// @Failure      400   {object}  response.ErrorResponse        "Thiếu id / dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse        "Không tìm thấy khách hàng"
// @Router       /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	customer := &domain.Customer{
		FullName:  req.FullName,
		Phone:     req.Phone,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		Note:      req.Note,
	}

	result, err := h.uc.UpdateCustomer(c.Request.Context(), id, customer)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật khách hàng thành công", result)
}

// DeleteCustomer godoc
// @Summary      Xóa mềm khách hàng
// @Tags         Customers
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID khách hàng"
// @Success      200  {object}  response.SuccessResponse  "Xóa khách hàng thành công"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy khách hàng"
// @Router       /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	deletedID, err := h.uc.DeleteCustomer(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa khách hàng thành công", gin.H{"id": deletedID})
}
