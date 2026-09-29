package product

import (
	"fmt"
	"strconv"

	domain "github.com/daidat02/server/internal/domain/product"
	uc "github.com/daidat02/server/internal/usecase/product"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc *uc.ProductUsecase
}

func NewProductHandler(uc *uc.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		uc: uc,
	}
}

func boolOrDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

// CreateProduct godoc
// @Summary      Tạo sản phẩm
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateProductRequest  true  "Thông tin sản phẩm"
// @Success      200   {object}  response.SuccessResponse     "Tạo sản phẩm thành công"
// @Failure      400   {object}  response.ErrorResponse       "Dữ liệu không hợp lệ / thiếu organization"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      500   {object}  response.ErrorResponse       "Lỗi server (vd trùng SKU)"
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error binding JSON:", err)
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	product := &domain.Product{
		CategoryID:    req.CategoryID,
		SKU:           req.SKU,
		Name:          req.Name,
		Unit:          req.Unit,
		NetUnit:       req.	NetUnit,
		NetAmount:     req.NetAmount,
		ProductType:   req.ProductType,
		ShowOnWeb:     boolOrDefault(req.ShowOnWeb, true),
		CostPrice:     req.CostPrice,
		SellPrice:     req.SellPrice,
		StockQuantity: req.StockQuantity,
		MinStock:      req.MinStock,
		IsActive:      boolOrDefault(req.IsActive, true),
	}

	result, err := h.uc.CreateProduct(c.Request.Context(), product)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo sản phẩm thành công", result)
}

// GetProductsByOrgID godoc
// @Summary      Lấy danh sách sản phẩm (phân trang + lọc + sắp xếp)
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Param        search     query     string  false  "Tìm theo tên, SKU hoặc loại sản phẩm (ILIKE)"
// @Param        status     query     string  false  "Lọc trạng thái" Enums(active,low_stock,out_of_stock)
// @Param        sort       query     string  false  "Sắp xếp" Enums(created_at_desc,price_asc,price_desc,stock_asc)
// @Success      200        {object}  response.SuccessResponse  "Danh sách sản phẩm kèm phân trang {data, page, page_size, total, total_pages}"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      500        {object}  response.ErrorResponse    "Lỗi server"
// @Router       /products [get]
func (h *ProductHandler) GetProductsByOrgID(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")
	sort := c.DefaultQuery("sort", "created_at_desc")
	result, err := h.uc.GetProductsByOrgID(c.Request.Context(), search, status, sort, page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách sản phẩm", result)
}

// UpdateProduct godoc
// @Summary      Cập nhật sản phẩm
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                      true  "ID sản phẩm"
// @Param        body  body      domain.UpdateProductRequest  true  "Thông tin cập nhật"
// @Success      200   {object}  response.SuccessResponse     "Cập nhật sản phẩm thành công"
// @Failure      400   {object}  response.ErrorResponse       "Thiếu id / dữ liệu không hợp lệ"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse       "Không tìm thấy sản phẩm"
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	product := &domain.Product{
		CategoryID:    req.CategoryID,
		SKU:           req.SKU,
		Name:          req.Name,
		Unit:          req.Unit,
		ProductType:   req.ProductType,
		ShowOnWeb:     boolOrDefault(req.ShowOnWeb, true),
		CostPrice:     req.CostPrice,
		SellPrice:     req.SellPrice,
		StockQuantity: req.StockQuantity,
		MinStock:      req.MinStock,
		IsActive:      boolOrDefault(req.IsActive, true),
	}

	result, err := h.uc.UpdateProduct(c.Request.Context(), id, product)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Cập nhật sản phẩm thành công", result)
}

// DeleteProduct godoc
// @Summary      Xóa sản phẩm
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "ID sản phẩm"
// @Success      200  {object}  response.SuccessResponse  "Xóa sản phẩm thành công"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy sản phẩm"
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	deletedID, err := h.uc.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Xóa sản phẩm thành công", gin.H{"id": deletedID})
}
