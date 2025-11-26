package main

import (
	"log"

	"go-auth-service/config"
	"go-auth-service/internal/delivery/http"
	"go-auth-service/internal/delivery/http/middleware"
	"go-auth-service/internal/infrastructure"
	"go-auth-service/internal/repository"
	"go-auth-service/internal/service"
	"go-auth-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load config file: %v. Using environment variables.", err)
	}

	db := infrastructure.NewDatabase(cfg)
	redisClient := infrastructure.NewRedisClient(cfg)

	userRepo := repository.NewUserRepository(db)
	tokenService := service.NewTokenService(cfg)
	passwordService := service.NewPasswordService()

	authUsecase := usecase.NewAuthUsecase(userRepo, tokenService, passwordService, redisClient)
	authMiddleware := middleware.NewAuthMiddleware(tokenService, redisClient)

	app := fiber.New()
	app.Use(logger.New())

	http.RegisterUserRoutes(app, authUsecase, authMiddleware)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
