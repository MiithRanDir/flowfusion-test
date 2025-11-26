package middleware

import (
	"context"
	"strings"

	constant "go-auth-service/internal/constants"
	"go-auth-service/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	tokenManager domain.TokenManager
	redisClient  *redis.Client
}

func NewAuthMiddleware(tokenManager domain.TokenManager, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
		redisClient:  redisClient,
	}
}

func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		tokenString := parts[1]

		// Check blacklist
		if m.redisClient != nil {
			val, _ := m.redisClient.Get(context.Background(), constant.STR_BLACKLIST+tokenString).Result()
			if val != "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is blacklisted"})
			}
		}

		claims, err := m.tokenManager.ValidateToken(tokenString, false)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("accessToken", tokenString) // Store for logout

		return c.Next()
	}
}
