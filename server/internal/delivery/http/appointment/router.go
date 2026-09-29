package appointment

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(r *gin.Engine, h *AppointmentHandler, basePath string) {
	group := r.Group(basePath + "/appointments")
	group.Use(middlewares.VerifyToken())
	{
		group.POST("", h.CreateAppointment)
		group.GET("", h.GetAppointments)
		group.GET("/:id", h.GetAppointmentDetail)
		group.PUT("/:id/reschedule", h.RescheduleAppointment)
		group.POST("/:id/cancel", h.CancelAppointment)
		group.PATCH("/:id/status", h.UpdateAppointmentStatus)
		group.DELETE("/:id", h.DeleteAppointment)
	}
}
