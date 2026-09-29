package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/daidat02/server/configs"
	_ "github.com/daidat02/server/docs"
	appointmenthttp "github.com/daidat02/server/internal/delivery/http/appointment"
	categoryhttp "github.com/daidat02/server/internal/delivery/http/category"
	customerhttp "github.com/daidat02/server/internal/delivery/http/customer"
	inventoryhttp "github.com/daidat02/server/internal/delivery/http/inventory"
	orderhttp "github.com/daidat02/server/internal/delivery/http/order"
	producthttp "github.com/daidat02/server/internal/delivery/http/product"
	servicehttp "github.com/daidat02/server/internal/delivery/http/service"
	staffhttp "github.com/daidat02/server/internal/delivery/http/staff"
	userhttp "github.com/daidat02/server/internal/delivery/http/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func NewRouter(cfg *config.Config, dbPool *pgxpool.Pool) *gin.Engine {
	// Tạo một router mới sử dụng Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.1.131"})
	// Định nghĩa các nguồn được phép truy cập (CORS)
	allowedOrigins := []string{"http://localhost:58864", "http://localhost:3000", "http://localhost:3002"}
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	basePath := cfg.App.BasePath

	userDeps := userhttp.NewDependencies(dbPool)
	userhttp.RegisterUserRoutes(r, userDeps.Handler, basePath)

	staffDeps := staffhttp.NewDependencies(dbPool)
	staffhttp.RegisterStaffRoutes(r, staffDeps.Handler, basePath)

	customerDeps := customerhttp.NewDependencies(dbPool)
	customerhttp.RegisterCustomerRoutes(r, customerDeps.Handler, basePath)

	categoryDeps := categoryhttp.NewDependencies(dbPool)
	categoryhttp.RegisterCategoryRoutes(r, categoryDeps.Handler, basePath)

	serviceDeps := servicehttp.NewDependencies(dbPool)
	servicehttp.RegisterServiceRoutes(r, serviceDeps.Handler, basePath)

	productDeps := producthttp.NewDependencies(dbPool)
	producthttp.RegisterProductRoutes(r, productDeps.Handler, basePath)

	appointmentDeps := appointmenthttp.NewDependencies(dbPool)
	appointmenthttp.RegisterAppointmentRoutes(r, appointmentDeps.Handler, basePath)

	inventoryDeps := inventoryhttp.NewDependencies(dbPool)
	inventoryhttp.RegisterInventoryRoutes(r, inventoryDeps.Handler, basePath)

	orderDeps := orderhttp.NewDependencies(dbPool)
	orderhttp.RegisterOrderRoutes(r, orderDeps.Handler, basePath)

	return r
}