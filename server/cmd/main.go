package main

import (
	"log"

	config "github.com/daidat02/server/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi khi tải cấu hình: %v", err)
	}
	log.Println("🚀 Starting Salon Management server...")
	r := NewRouter()
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)

	}
	log.Printf("Salon Management server Đang Chạy trên cổng :%s", cfg.App.Port)
}