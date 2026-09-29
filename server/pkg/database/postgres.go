package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


func NewPostgresDatabase(dsn string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Lỗi parse cấu hình PostgreSQL DSN: %v", err)
	}
	config.MaxConns = 25
	config.MinConns = 5
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Lỗi kết nối PostgreSQL: %v", err)
	}

	if err := dbPool.Ping(ctx); err != nil{
		log.Fatalf(" Ping tới PostgreSQL thất bại: %v", err)
	}
	log.Println("✅ Kết nối PostgreSQL thành công")
	return dbPool
}