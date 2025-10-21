package main

import (
	"fmt"
	"log"
	"time"

	"cctv-monitoring-backend/internal/config"
	"cctv-monitoring-backend/internal/database"
	"cctv-monitoring-backend/internal/handler"
	"cctv-monitoring-backend/internal/middleware"
	"cctv-monitoring-backend/internal/repository"
	"cctv-monitoring-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load konfigurasi
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// Koneksi ke database
	db, err := database.Connect(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Jalankan migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	cameraRepo := repository.NewCameraRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo)
	rtspService := service.NewRTSPService(cfg.RTSP.APIURL, cfg.RTSP.PublicBaseURL, cfg.RTSP.Username, cfg.RTSP.Password)
	cameraService := service.NewCameraService(cameraRepo, rtspService)

	// Start cleanup job for expired tokens (run every 1 hour)
	cleanupService := service.NewCleanupService(tokenRepo)
	cleanupService.StartCleanupJob(1 * time.Hour)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg.JWT.Secret, cfg.JWT.Expiration.String())
	cameraHandler := handler.NewCameraHandler(cameraService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Middleware untuk inject dependencies ke context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("jwt_secret", cfg.JWT.Secret)
		c.Locals("auth_service", authService)
		return c.Next()
	})

	// Routes
	setupRoutes(app, authHandler, cameraHandler, authService)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("✓ Server is running on http://localhost%s", addr)
	log.Fatal(app.Listen(addr))
}

// setupRoutes mengatur semua routing aplikasi
func setupRoutes(app *fiber.App, authHandler *handler.AuthHandler, cameraHandler *handler.CameraHandler, authService service.AuthService) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "CCTV Monitoring API is running",
		})
	})

	// API v1
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	// Protected routes dengan auth middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Auth routes (protected)
	auth.Get("/me", authMiddleware, authHandler.Me)
	auth.Post("/logout", authMiddleware, authHandler.Logout)

	// Camera routes
	cameras := api.Group("/cameras", authMiddleware)
	cameras.Get("/", cameraHandler.GetAll)
	cameras.Get("/:id", cameraHandler.GetByID)
	cameras.Post("/", cameraHandler.Create)
	cameras.Put("/:id", cameraHandler.Update)
	cameras.Delete("/:id", cameraHandler.Delete)

	// Camera filter routes
	cameras.Get("/zone/filter", cameraHandler.GetByZone)
	cameras.Get("/nearby", cameraHandler.GetNearby)

	// Stream routes
	cameras.Post("/:id/stream/start", cameraHandler.StartStream)
	cameras.Post("/:id/stream/stop", cameraHandler.StopStream)
}

// customErrorHandler adalah custom error handler untuk Fiber
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": "An error occurred",
		"error":   err.Error(),
	})
}
