package staff

import (
	"strconv"

	domain "github.com/daidat02/server/internal/domain/staff"
	uc "github.com/daidat02/server/internal/usecase/staff"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct{
	uc *uc.StaffUsecase
}

func NewStaffHandler(uc *uc.StaffUsecase) *StaffHandler{
	return &StaffHandler{
		uc: uc,
	}
}

// CreateStaff godoc
// @Summary      Tạo nhân viên
// @Tags         Staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateStaffRequest  true  "Thông tin nhân viên"
// @Success      200   {object}  response.SuccessResponse   "Tạo nhân viên thành công"
// @Failure      400   {object}  response.ErrorResponse     "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse     "Thiếu token / token không hợp lệ"
// @Failure      500   {object}  response.ErrorResponse     "Lỗi server"
// @Router       /staff [post]
func (h *StaffHandler) CreateStaff(c *gin.Context){
	var req domain.CreateStaffRequest

	if err:= c.ShouldBindJSON(&req); err!=nil{
		c.JSON(400, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	staff := &domain.Staff{
		Code: req.Code,
		FullName: req.FullName,
		Phone: req.Phone,
		Position: req.Position,
		CommissionRate: req.CommissionRate,
		Status: req.Status,
		HireDate: req.HireDate,
	}

	result,err:= h.uc.CreateStaff(c.Request.Context(), staff)
	if err !=nil{
		response.Error(c,err)
		c.Abort()
		return 
	}
	response.Success(c,200,"Tạo nhân viên thành công",result)
}


// GetStaffByOrgID godoc
// @Summary      Lấy danh sách nhân viên (phân trang + lọc + sắp xếp)
// @Tags         Staff
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Param        search     query     string  false  "Tìm theo tên, mã, SĐT hoặc vị trí (ILIKE)"
// @Param        status     query     string  false  "Lọc trạng thái" Enums(active,inactive)
// @Param        sort       query     string  false  "Sắp xếp" Enums(created_at_desc,name_asc,name_desc,hire_date_asc,hire_date_desc)
// @Success      200        {object}  response.SuccessResponse  "Danh sách nhân viên kèm phân trang {data, page, page_size, total, total_pages}"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      500        {object}  response.ErrorResponse    "Lỗi server"
// @Router       /staff [get]
func (h *StaffHandler) GetStaffByOrgID(c *gin.Context){
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")
	sort := c.DefaultQuery("sort", "created_at_desc")

	result, err := h.uc.GetStaffByOrgID(c.Request.Context(), search, status, sort, page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách nhân viên", result)
}

// UpdateStaff godoc
// @Summary      Cập nhật nhân viên
// @Tags         Staff
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                    true  "ID nhân viên"
// @Param        body  body      domain.UpdateStaffRequest  true  "Thông tin cập nhật"
// @Success      200   {object}  response.SuccessResponse   "Cập nhật nhân viên thành công"
// @Failure      400   {object}  response.ErrorResponse     "Thiếu id / dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse     "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse     "Không tìm thấy nhân viên"
// @Router       /staff/{id} [put]
func (h *StaffHandler) UpdateStaff(c *gin.Context){
	id := c.Param("id")

	var req domain.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	staff := &domain.Staff{
		Code:           req.Code,
		FullName:       req.FullName,
		Phone:          req.Phone,
		Position:       req.Position,
		CommissionRate: req.CommissionRate,
		Status:         req.Status,
		HireDate:       req.HireDate,
	}

	result, err := h.uc.UpdateStaff(c.Request.Context(), id, staff)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật nhân viên thành công", result)
}

// DeleteStaff godoc
// @Summary      Xóa mềm nhân viên
// @Tags         Staff
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID nhân viên"
// @Success      200  {object}  response.SuccessResponse  "Xóa nhân viên thành công"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy nhân viên"
// @Router       /staff/{id} [delete]
func (h *StaffHandler) DeleteStaff(c *gin.Context){
	id := c.Param("id")

	deletedID, err := h.uc.DeleteStaff(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa nhân viên thành công", gin.H{"id": deletedID})
}