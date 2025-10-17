package handler

import (
	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler menangani HTTP requests untuk authentication
type AuthHandler struct {
	authService   service.AuthService
	jwtSecret     string
	jwtExpiration string
}

// NewAuthHandler membuat instance baru dari AuthHandler
func NewAuthHandler(authService service.AuthService, jwtSecret, jwtExpiration string) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Login handler untuk login user
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validasi input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "Username and password are required",
		})
	}

	// Proses login
	response, err := h.authService.Login(req.Username, req.Password, h.jwtSecret, h.jwtExpiration)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register handler untuk registrasi user baru
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validasi input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "Username, email, and password are required",
		})
	}

	// Set default role jika tidak diisi
	if req.Role == "" {
		req.Role = "operator"
	}

	// Proses registrasi
	user, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Me handler untuk mendapatkan info user yang sedang login
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Ambil user info dari context (sudah diset oleh auth middleware)
	userID := c.Locals("user_id").(string)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "User info retrieved successfully",
		Data: fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
