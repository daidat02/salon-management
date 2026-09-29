package customer

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(r *gin.Engine, h *CustomerHandler, basePath string) {
	group := r.Group(basePath + "/customers")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetCustomersByOrgID)
		group.GET("/:id", h.GetCustomerDetail)
		group.POST("", h.CreateCustomer)
		group.PUT("/:id", h.UpdateCustomer)
		group.DELETE("/:id", h.DeleteCustomer)
	}
}
