package http

import (
	"go-auth-service/internal/delivery/http/middleware"
	"go-auth-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, authUsecase domain.AuthUsecase, authMiddleware *middleware.AuthMiddleware) {
	handler := NewAuthHandler(authUsecase)

	auth := app.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", authMiddleware.Protected(), handler.Logout)

	app.Get("/me", authMiddleware.Protected(), handler.GetMe)
}
