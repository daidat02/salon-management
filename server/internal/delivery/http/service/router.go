package service

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterServiceRoutes(r *gin.Engine, h *ServiceHandler, basePath string) {
	group := r.Group(basePath + "/services")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetServicesByOrgID)
		group.GET("/:id/details", h.GetServiceDetails)
		group.POST("", h.CreateService)
		group.PUT("/:id", h.UpdateService)
		group.DELETE("/:id", h.DeleteService)
	}
}
