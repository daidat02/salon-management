package order

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine, h *OrderHandler, basePath string) {
	group := r.Group(basePath + "/orders")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetOrdersByOrgID)
		group.POST("", h.CreateOrder)
		group.GET("/:id", h.GetOrderDetail)
		
		group.POST("/:id/cancel", h.CancelOrder)
		group.POST("/:id/refund", h.RefundOrder)
		group.POST("/:id/serve", h.ServeOrder)
		group.POST("/:id/complete", h.CompleteOrder)
		group.POST("/:id/payments", h.CreatePayment)

		group.GET("/pos-items", h.GetPosItemsByOrgID)
	}
}
