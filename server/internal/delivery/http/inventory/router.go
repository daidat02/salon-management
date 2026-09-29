package inventory

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(r *gin.Engine, h *InventoryHandler, basePath string) {
	group := r.Group(basePath + "/inventory")
	group.Use(middlewares.VerifyToken())
	{
		group.POST("/", h.CreateDocument)
		group.GET("/", h.GetDocumentsByOrgID)
		group.GET("/transactions", h.GetTransactionsByOrgID)
		group.GET("/low-stock", h.GetLowStockProducts)
		group.GET("/:document_id", h.GetDocumentByID)
		group.GET("/transactions/:product_id", h.GetTransactionsByProductID)
	}
}
