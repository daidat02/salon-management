package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func NewRouter() *gin.Engine {
	// Tạo một router mới sử dụng Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.1.131"})
	// Định nghĩa các nguồn được phép truy cập (CORS)
	allowedOrigins := []string{"http://localhost:58864", "http://localhost:3001", "http://localhost:3002"}
	// Cấu hình CORS để cho phép các yêu cầu từ các nguồn được chỉ định
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	r.GET("/health", func (c *gin.Context)  {
		c.JSON(200, gin.H{"status": "ok",})	
	})

	return r
}