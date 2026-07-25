package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"rekap-service/internal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	akademikBaseURL := env("AKADEMIK_BASE_URL", "http://localhost:8080")
	appPort := env("APP_PORT", "8080")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client := internal.NewAkademikClient(akademikBaseURL)
	svc := internal.NewService(client)
	handler := internal.NewHandler(svc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(loggingMiddleware("rekap-service"))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
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

	addr := fmt.Sprintf(":%s", appPort)
	slog.Info("server berjalan", "addr", addr)
	if err := router.Run(addr); err != nil {
		slog.Error("gagal menjalankan server", "error", err)
		os.Exit(1)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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
