package middlewares

import (
	"server/src/models"
	"strings"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Extract token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("No token provided")
	}

	// Remove "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token
	claims, err := models.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// Set user information in context
	c.Locals("user", claims.Username)
	return c.Next()
}