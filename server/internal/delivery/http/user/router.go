package user

import "github.com/gin-gonic/gin"

func RegisterUserRoutes (r *gin.Engine, h *UserHandler, basePath string) {
	group := r.Group(basePath)
	
	group.POST("/register", h.CreateUser)
	group.POST("/login", h.Login)
	group.POST("/refresh-token",h.RefreshToken)
}