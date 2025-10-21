package handler

import (
	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/service"
	"errors"
	"strings"

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
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid request body",
				err.Error(),
			),
		)
	}

	// Validasi input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeMissingFields,
				"Username and password are required",
			),
		)
	}

	// Proses login
	response, err := h.authService.Login(req.Username, req.Password, h.jwtSecret, h.jwtExpiration)
	if err != nil {
		// Handle different error types
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.NewErrorResponse(
					models.ErrCodeInvalidCredentials,
					"Invalid username or password",
				),
			)
		case errors.Is(err, service.ErrUserInactive):
			return c.Status(fiber.StatusForbidden).JSON(
				models.NewErrorResponse(
					models.ErrCodeUserInactive,
					"Your account is inactive. Please contact administrator",
				),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.NewErrorResponse(
					models.ErrCodeInternalError,
					"An error occurred during login",
					err.Error(),
				),
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Logout handler untuk logout user
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeUnauthorized,
				"Missing authorization header",
			),
		)
	}

	// Extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeUnauthorized,
				"Invalid authorization header format",
			),
		)
	}

	token := parts[1]

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("user_id").(string)

	// Logout (blacklist token)
	if err := h.authService.Logout(token, userID, h.jwtExpiration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to logout",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// Register handler untuk registrasi user baru
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid request body",
				err.Error(),
			),
		)
	}

	// Validasi input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeMissingFields,
				"Username, email, and password are required",
			),
		)
	}

	// Set default role jika tidak diisi
	if req.Role == "" {
		req.Role = "operator"
	}

	// Proses registrasi
	user, err := h.authService.Register(&req)
	if err != nil {
		// Cek apakah username/email sudah ada
		errMsg := err.Error()
		if errMsg == "username already exists" || errMsg == "email already exists" {
			return c.Status(fiber.StatusConflict).JSON(
				models.NewErrorResponse(
					models.ErrCodeAlreadyExists,
					errMsg,
				),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Registration failed",
				err.Error(),
			),
		)
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
