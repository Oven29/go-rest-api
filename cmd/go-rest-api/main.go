package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-rest-api/internal/api"
	"go-rest-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		cfg.DB.Host, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	l := getLogLevel(cfg)
	logger := slog.New(slog.NewTextHandler(log.Default().Writer(), &slog.HandlerOptions{Level: l}))

	router := api.NewRouter(db, logger)
	registerCustomError(router)
	if cfg.EnableSwagger {
		registerSwagger(router)
	}

	address := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	logger.Info("Starting server on http://" + address)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getLogLevel(cfg *config.Config) slog.Level {
	switch cfg.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func registerSwagger(router *gin.Engine) {
	router.Static("/docs", "/home/ivan/projects/go-rest-api/docs")
	url := ginSwagger.URL("/docs/openapi.yml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, url))
}

func registerCustomError(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found on this server.",
			"path":    c.Request.URL.Path,
		})
	})
}
