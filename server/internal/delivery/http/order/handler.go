package order

import (
	"strconv"

	domain "github.com/daidat02/server/internal/domain/order"
	uc "github.com/daidat02/server/internal/usecase/order"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	uc *uc.OrderUsecase
}

func NewOrderHandler(uc *uc.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		uc: uc,
	}
}

// CreateOrder godoc
// @Summary      Tạo đơn hàng (trừ kho trong transaction)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateOrderRequest  true  "Thông tin đơn hàng"
// @Success      200   {object}  response.SuccessResponse   "Tạo đơn hàng thành công"
// @Failure      400   {object}  response.ErrorResponse     "Dữ liệu không hợp lệ / tồn không đủ"
// @Failure      401   {object}  response.ErrorResponse     "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse     "Không tìm thấy khách/dịch vụ/sản phẩm"
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo đơn hàng thành công", result)
}

// GetOrdersByOrgID godoc
// @Summary      Lấy danh sách đơn hàng (phân trang)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        status     query     string  false  "Lọc theo trạng thái"
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200        {object}  response.SuccessResponse  "Danh sách đơn hàng kèm phân trang"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /orders [get]
func (h *OrderHandler) GetOrdersByOrgID(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetOrdersByOrgID(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách đơn hàng", result)
}

// GetOrderDetail godoc
// @Summary      Lấy chi tiết đơn hàng (dòng hàng + vật tư + thanh toán)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID đơn hàng"
// @Success      200  {object}  response.SuccessResponse  "Chi tiết đơn hàng"
// @Failure      400  {object}  response.ErrorResponse    "Thiếu id"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy đơn hàng"
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.GetOrderDetail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Chi tiết đơn hàng", result)
}

// CancelOrder godoc
// @Summary      Hủy đơn chưa phục vụ (nhả giữ chỗ kho, không hoàn tồn)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID đơn hàng"
// @Success      200  {object}  response.SuccessResponse  "Hủy đơn thành công"
// @Failure      400  {object}  response.ErrorResponse    "Trạng thái không cho hủy"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy đơn hàng"
// @Router       /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.CancelOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Hủy đơn hàng thành công", result)
}

// RefundOrder godoc
// @Summary      Hoàn đơn đang phục vụ (hoàn kho + vô hiệu phiếu xuất)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID đơn hàng"
// @Success      200  {object}  response.SuccessResponse  "Hoàn đơn thành công"
// @Failure      400  {object}  response.ErrorResponse    "Trạng thái không cho hoàn"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy đơn hàng"
// @Router       /orders/{id}/refund [post]
func (h *OrderHandler) RefundOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.RefundOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Hoàn đơn hàng thành công", result)
}

// ServeOrder godoc
// @Summary      Bắt đầu phục vụ đơn đã thanh toán (trừ kho + sinh phiếu xuất)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID đơn hàng"
// @Success      200  {object}  response.SuccessResponse  "Bắt đầu phục vụ, đã sinh phiếu xuất"
// @Failure      400  {object}  response.ErrorResponse    "Chưa thanh toán / tồn kho không đủ"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy đơn hàng"
// @Router       /orders/{id}/serve [post]
func (h *OrderHandler) ServeOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.ServeOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Bắt đầu phục vụ đơn hàng thành công", result)
}

// CompleteOrder godoc
// @Summary      Hoàn tất phục vụ đơn hàng
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "ID đơn hàng"
// @Success      200  {object}  response.SuccessResponse  "Hoàn tất đơn hàng thành công"
// @Failure      400  {object}  response.ErrorResponse    "Trạng thái không cho hoàn tất"
// @Failure      401  {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404  {object}  response.ErrorResponse    "Không tìm thấy đơn hàng"
// @Router       /orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.CompleteOrder(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Hoàn tất đơn hàng thành công", result)
}

// CreatePayment godoc
// @Summary      Ghi nhận thanh toán cho đơn (idempotent theo mã giao dịch)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                      true  "ID đơn hàng"
// @Param        body  body      domain.CreatePaymentRequest  true  "Thông tin thanh toán"
// @Success      200   {object}  response.SuccessResponse     "Ghi nhận thanh toán thành công"
// @Failure      400   {object}  response.ErrorResponse       "Dữ liệu không hợp lệ / đơn đã hủy"
// @Failure      401   {object}  response.ErrorResponse       "Thiếu token / token không hợp lệ"
// @Failure      404   {object}  response.ErrorResponse       "Không tìm thấy đơn hàng"
// @Router       /orders/{id}/payments [post]
func (h *OrderHandler) CreatePayment(c *gin.Context) {
	id := c.Param("id")

	var req domain.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}

	result, err := h.uc.CreatePayment(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Ghi nhận thanh toán thành công", result)
}



// GetPosItemsByOrgID godoc
// @Summary      Lấy danh sách sản phẩm/dịch vụ cho màn POS (tìm kiếm + lọc danh mục + phân trang cursor)
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        search    query     string  false  "Tìm theo tên hoặc SKU (ILIKE)"
// @Param        category  query     string  false  "Lọc theo ID danh mục"
// @Param        type      query     string  false  "Con trỏ phân trang: type của dòng cuối trang trước"
// @Param        name      query     string  false  "Con trỏ phân trang: name của dòng cuối trang trước"
// @Param        id        query     string  false  "Con trỏ phân trang: id của dòng cuối trang trước"
// @Param        page_size query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200       {object}  response.SuccessResponse  "Danh sách sản phẩm/dịch vụ POS"
// @Failure      400       {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401       {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /orders/pos-items [get]
func (h *OrderHandler) GetPosItemsByOrgID(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	itemType := c.Query("type")
	search := c.Query("search")
	category := c.Query("category")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetPosItemsByOrgID(c.Request.Context(), pageSize, name, itemType, id, search, category)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách sản phẩm POS", result)
}