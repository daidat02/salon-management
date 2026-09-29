package product

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, h *ProductHandler, basePath string) {
	group := r.Group(basePath + "/products")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetProductsByOrgID)
		group.POST("", h.CreateProduct)
		group.PUT("/:id", h.UpdateProduct)
		group.DELETE("/:id", h.DeleteProduct)
	}
}
