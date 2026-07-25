package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"akademik-service/internal/config"
	"akademik-service/internal/siakad"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		slog.Error("gagal membuka koneksi database", "error", err)
		log.Fatalf("gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("gagal ping database", "error", err)
		log.Fatalf("gagal ping database: %v", err)
	}
	slog.Info("koneksi database berhasil")

	mhsRepo := siakad.NewMahasiswaRepository(db)
	nRepo := siakad.NewNilaiRepository(db)
	svc := siakad.NewService(mhsRepo, nRepo)
	handler := siakad.NewHandler(svc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(loggingMiddleware("akademik-service"))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		slog.InfoContext(c.Request.Context(), "health check")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	handler.Register(api)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	slog.Info("server berjalan", "addr", addr)
	if err := router.Run(addr); err != nil {
		slog.Error("gagal menjalankan server", "error", err)
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

func loggingMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		requestID, _ := c.Get("request_id")
		slog.Info("request",
			"service", serviceName,
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
		)
	}
}
