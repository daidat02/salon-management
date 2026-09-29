package staff

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterStaffRoutes(r *gin.Engine, h *StaffHandler, basePath string) {
	group := r.Group(basePath+ "/staff")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetStaffByOrgID)
		group.POST("", h.CreateStaff)
		group.PUT("/:id", h.UpdateStaff)
		group.DELETE("/:id", h.DeleteStaff)
	}
	
}