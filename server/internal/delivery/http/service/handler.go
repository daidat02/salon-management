package service

import (
	"strconv"

	domain "github.com/daidat02/server/internal/domain/service"
	uc "github.com/daidat02/server/internal/usecase/service"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	uc *uc.ServiceUsecase
}

func NewServiceHandler(uc *uc.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		uc: uc,
	}
}

func boolOrDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

// CreateService godoc
// @Summary      Tạo dịch vụ
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateServiceRequest  true  "Thông tin dịch vụ"
// @Success      200   {object}  response.SuccessResponse     "Tạo dịch vụ thành công"
// @Failure      400   {object}  response.ErrorResponse       "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      500   {object}  response.ErrorResponse       "Lỗi server"
// @Router       /services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req domain.CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	service := &domain.Service{
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		BufferMinutes:   req.BufferMinutes,
		Price:           req.Price,
		IsActive:        boolOrDefault(req.IsActive, true),
	}


	materials := make([]*domain.ServiceMaterialInput, len(req.Materials))
	for i, m := range req.Materials {
		materials[i] = &domain.ServiceMaterialInput{
			ProductID: m.ProductID,
			Quantity:  m.Quantity,
		}
	}
	result, err := h.uc.CreateService(c.Request.Context(), service,materials)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo dịch vụ thành công", result)
}

// GetServicesByOrgID godoc
// @Summary      Lấy danh sách dịch vụ (phân trang + lọc + sắp xếp)
// @Tags         Services
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Param        search     query     string  false  "Tìm theo tên hoặc mô tả (ILIKE)"
// @Param        status     query     string  false  "Lọc trạng thái" Enums(active,inactive)
// @Param        sort       query     string  false  "Sắp xếp" Enums(created_at_desc,price_asc,price_desc,duration_asc,name_asc)
// @Success      200        {object}  response.SuccessResponse  "Danh sách dịch vụ kèm phân trang {data, page, page_size, total, total_pages}"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      500        {object}  response.ErrorResponse    "Lỗi server"
// @Router       /services [get]
func (h *ServiceHandler) GetServicesByOrgID(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")
	sort := c.DefaultQuery("sort", "created_at_desc")

	result, err := h.uc.GetServicesByOrgID(c.Request.Context(), search, status, sort, page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách dịch vụ", result)
}

// UpdateService godoc
// @Summary      Cập nhật dịch vụ
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                      true  "ID dịch vụ"
// @Param        body  body      domain.UpdateServiceRequest  true  "Thông tin cập nhật"
// @Success      200   {object}  response.SuccessResponse     "Cập nhật dịch vụ thành công"
// @Failure      400   {object}  response.ErrorResponse       "Thiếu id / dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse       "Không tìm thấy dịch vụ"
// @Router       /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	service := &domain.Service{
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		BufferMinutes:   req.BufferMinutes,
		Price:           req.Price,
		IsActive:        boolOrDefault(req.IsActive, true),
	}

	result, err := h.uc.UpdateService(c.Request.Context(), id, service)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật dịch vụ thành công", result)
}

// DeleteService godoc
// @Summary      Xóa dịch vụ
// @Tags         Services
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID dịch vụ"
// @Success      200  {object}  response.SuccessResponse  "Xóa dịch vụ thành công"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy dịch vụ"
// @Router       /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	deletedID, err := h.uc.DeleteService(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa dịch vụ thành công", gin.H{"id": deletedID})
}


// GetServiceDetails godoc
// @Summary      Lấy chi tiết dịch vụ và nguyên vật liệu
// @Tags         Services
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID dịch vụ"
// @Success      200  {object}  response.SuccessResponse  "Chi tiết dịch vụ và danh sách nguyên vật liệu"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy dịch vụ"
// @Router       /services/{id}/details [get]
func (h *ServiceHandler) GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("id")
	service, materials, err := h.uc.GetServiceDetails(c.Request.Context(), serviceID)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}

	response.Success(c, 200, "Chi tiết dịch vụ", gin.H{
		"service":   service,
		"materials": materials,
	})
}