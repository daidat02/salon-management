package category

import (
	"strconv"

	domain "github.com/daidat02/server/internal/domain/category"
	uc "github.com/daidat02/server/internal/usecase/category"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc *uc.CategoryUsecase
}

func NewCategoryHandler(uc *uc.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		uc: uc,
	}
}

func boolOrDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

// CreateCategory godoc
// @Summary      Tạo danh mục
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateCategoryRequest  true  "Thông tin danh mục (type: service | product)"
// @Success      200   {object}  response.SuccessResponse      "Tạo danh mục thành công"
// @Failure      400   {object}  response.ErrorResponse        "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Failure      500   {object}  response.ErrorResponse        "Lỗi server (vd trùng tên trong org)"
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	category := &domain.Category{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    boolOrDefault(req.IsActive, true),
	}

	result, err := h.uc.CreateCategory(c.Request.Context(), category)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo danh mục thành công", result)
}

// GetCategoriesByOrgID godoc
// @Summary      Lấy danh sách danh mục (phân trang)
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Trang (mặc định 1)"
// @Param        page_size  query     int  false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200        {object}  response.SuccessResponse  "Danh sách danh mục kèm phân trang"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /categories [get]
func (h *CategoryHandler) GetCategoriesByOrgID(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetCategoriesByOrgID(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách danh mục", result)
}

// UpdateCategory godoc
// @Summary      Cập nhật danh mục
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                       true  "ID danh mục"
// @Param        body  body      domain.UpdateCategoryRequest  true  "Thông tin cập nhật"
// @Success      200   {object}  response.SuccessResponse      "Cập nhật danh mục thành công"
// @Failure      400   {object}  response.ErrorResponse        "Thiếu id / dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse        "Không tìm thấy danh mục"
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	category := &domain.Category{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    boolOrDefault(req.IsActive, true),
	}

	result, err := h.uc.UpdateCategory(c.Request.Context(), id, category)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật danh mục thành công", result)
}

// DeleteCategory godoc
// @Summary      Xóa danh mục
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID danh mục"
// @Success      200  {object}  response.SuccessResponse  "Xóa danh mục thành công"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy danh mục"
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	deletedID, err := h.uc.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa danh mục thành công", gin.H{"id": deletedID})
}
