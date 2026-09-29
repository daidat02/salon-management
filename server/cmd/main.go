package main

import (
	"log"

	config "github.com/daidat02/server/configs"
	"github.com/daidat02/server/pkg/database"
	"github.com/daidat02/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

// @title           Salon Management API
// @version         1.0
// @description     API quản lý salon: user, nhân viên, khách hàng.
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     "Bearer {access_token}"

func main() {
	// Set Gin mode to Debug for development purposes
	log.Println("Setting Gin mode to Debug...")
	gin.SetMode(gin.DebugMode)

	log.Println("🚀 Starting Salon Management server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi khi tải cấu hình: %v", err)
	}

	// Kết nối cơ sở dữ liệu PostgreSQL
	dsn := cfg.Database.DB_url
	dbPool := database.NewPostgresDatabase(dsn)

	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.AccessTokenExpire, cfg.JWT.RefreshTokenExpire)

	// Khởi tạo router và đăng ký các route
	r := NewRouter(cfg, dbPool)

	// Khởi động server trên cổng được chỉ định trong cấu hình
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)

	}

	log.Printf("Salon Management server Đang Chạy trên cổng :%s", cfg.App.Port)
}