package category

import (
	"github.com/daidat02/server/internal/delivery/http/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(r *gin.Engine, h *CategoryHandler, basePath string) {
	group := r.Group(basePath + "/categories")
	group.Use(middlewares.VerifyToken())
	{
		group.GET("", h.GetCategoriesByOrgID)
		group.POST("", h.CreateCategory)
		group.PUT("/:id", h.UpdateCategory)
		group.DELETE("/:id", h.DeleteCategory)
	}
}
