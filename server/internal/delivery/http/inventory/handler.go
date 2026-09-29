package inventory

import (
	"strconv"

	"github.com/bytedance/gopkg/util/logger"
	domain "github.com/daidat02/server/internal/domain/inventory"
	uc "github.com/daidat02/server/internal/usecase/inventory"
	"github.com/daidat02/server/pkg/apperror"
	"github.com/daidat02/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	uc *uc.InventoryUsecase
}

func NewInventoryHandler(uc *uc.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{
		uc: uc,
	}
}

// CreateDocument godoc
// @Summary      Tạo phiếu kho (header + chi tiết + ghi sổ tồn)
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      domain.CreateDocumentRequest  true  "Thông tin phiếu kho kèm danh sách sản phẩm"
// @Success      200   {object}  response.SuccessResponse      "Tạo phiếu kho thành công"
// @Failure      400   {object}  response.ErrorResponse        "Dữ liệu không hợp lệ / tồn không đủ"
// @Failure      401   {object}  response.ErrorResponse        "Thiếu token / token không hợp lệ"
// @Router       /inventory/documents [post]
func (h *InventoryHandler) CreateDocument(c *gin.Context) {
	var req domain.CreateDocumentRequest
	var userID = c.GetString("userId")
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.New(400, "Dữ liệu không hợp lệ", err))
		c.Abort()
		return
	}
	req.CreatedBy = &userID
	result, err := h.uc.CreateDocument(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Tạo phiếu kho thành công", result)
}

// GetTransactionsByOrgID godoc
// @Summary      Lấy lịch sử kho (phân trang, lọc theo sản phẩm)
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        product_id  query     string  false  "Lọc theo ID sản phẩm"
// @Param        page        query     int     false  "Trang (mặc định 1)"
// @Param        page_size   query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200         {object}  response.SuccessResponse  "Lịch sử kho kèm phân trang"
// @Failure      400         {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401         {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /inventory/transactions [get]
func (h *InventoryHandler) GetTransactionsByOrgID(c *gin.Context) {
	productID := c.Query("product_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetTransactionsByOrgID(c.Request.Context(), productID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Lịch sử kho", result)
}

// GetLowStockProducts godoc
// @Summary      Lấy danh sách sản phẩm tồn thấp
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Trang (mặc định 1)"
// @Param        page_size  query     int  false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200        {object}  response.SuccessResponse  "Danh sách tồn thấp kèm phân trang"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /inventory/low-stock [get]
func (h *InventoryHandler) GetLowStockProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.uc.GetLowStockProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách tồn thấp", result)
}

// GetDocumentsByOrgID godoc
// @Summary      Lấy danh sách phiếu kho (tìm kiếm + lọc loại + sắp xếp + phân trang)
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        search     query     string  false  "Tìm theo mã phiếu, tên NCC, mã đơn, tên/SĐT người tạo"
// @Param        sort       query     string  false  "Sắp xếp" Enums(total_amount_asc,total_amount_desc,created_at_asc,created_at_desc)
// @Param        filter     query     string  false  "Lọc theo loại phiếu" Enums(import,export,adjust,sale,return)
// @Param        page       query     int     false  "Trang (mặc định 1)"
// @Param        page_size  query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200        {object}  response.SuccessResponse  "Danh sách phiếu kho kèm phân trang"
// @Failure      400        {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401        {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /inventory/ [get]
func (h *InventoryHandler) GetDocumentsByOrgID(c *gin.Context){
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "")
	filter := c.DefaultQuery("filter", "")

	result, err := h.uc.GetDocumentsByOrgID(c.Request.Context(), search, sort, filter, page, pageSize)
	if err != nil {
		response.Error(c, err)
		logger.Error("Lỗi khi lấy danh sách phiếu kho: ", err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Danh sách phiếu kho", result)
}

// GetDocumentByID godoc
// @Summary      Lấy chi tiết phiếu kho theo ID (kèm danh sách dòng sản phẩm)
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        document_id  path      string  true  "ID phiếu kho"
// @Success      200          {object}  response.SuccessResponse  "Chi tiết phiếu kho"
// @Failure      400          {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401          {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Failure      404          {object}  response.ErrorResponse    "Phiếu kho không tồn tại"
// @Router       /inventory/{document_id} [get]
func (h *InventoryHandler) GetDocumentByID(c *gin.Context) {
	documentID := c.Param("document_id")
	result, err := h.uc.GetDocumentByID(c.Request.Context(), documentID)
	if err != nil {
		response.Error(c, err)
		logger.Error("Lỗi khi lấy chi tiết phiếu kho: ", err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Chi tiết phiếu kho", result)
}

// GetTransactionsByProductID godoc
// @Summary      Lấy lịch sử kho theo sản phẩm (lọc loại + ngày + sắp xếp + phân trang)
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        product_id      path      string  true   "ID sản phẩm"
// @Param        sort            query     string  false  "Sắp xếp" Enums(created_desc,created_asc)
// @Param        inventory_type  query     string  false  "Lọc theo loại giao dịch" Enums(import,export,adjust,sale,return)
// @Param        date            query     string  false  "Lọc theo ngày (YYYY-MM-DD)"
// @Param        page            query     int     false  "Trang (mặc định 1)"
// @Param        page_size       query     int     false  "Số dòng/trang (mặc định 20, tối đa 100)"
// @Success      200             {object}  response.SuccessResponse  "Lịch sử kho theo sản phẩm kèm phân trang"
// @Failure      400             {object}  response.ErrorResponse    "Thiếu organization"
// @Failure      401             {object}  response.ErrorResponse    "Thiếu token / token không hợp lệ"
// @Router       /inventory/transactions/{product_id} [get]
func (h *InventoryHandler) GetTransactionsByProductID(c *gin.Context) {
	productID := c.Param("product_id")
	page,_ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize,_ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sort := c.DefaultQuery("sort", "")
	inventoryType := c.DefaultQuery("inventory_type", "")
	date := c.DefaultQuery("date", "")

	result, err := h.uc.GetTransactionsByProductID(c.Request.Context(), productID, page, pageSize, sort, inventoryType, date)
	if err != nil {
		response.Error(c, err)
		logger.Error("Lỗi khi lấy lịch sử kho theo sản phẩm: ", err)
		c.Abort()
		return
	}
	response.Success(c, 200, "Lịch sử kho theo sản phẩm", result)
}