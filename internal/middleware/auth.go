package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Auth 验证 Bearer Token
func Auth(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if token == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server token not configured",
			})
		}
		auth := c.Get("Authorization")
		if len(auth) < 8 || auth[:7] != "Bearer " || auth[7:] != token {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		return c.Next()
	}
}
