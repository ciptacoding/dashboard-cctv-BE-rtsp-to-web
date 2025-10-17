package middleware

import (
	"strings"

	"cctv-monitoring-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware adalah middleware untuk validasi JWT token
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing authorization header",
			})
		}

		// Extract token dari "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Simpan user info ke context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware adalah middleware untuk validasi role user
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role dari context
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied",
			})
		}

		userRole := role.(string)

		// Check apakah role user termasuk dalam allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Insufficient permissions",
		})
	}
}
