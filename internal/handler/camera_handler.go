package handler

import (
	"strconv"

	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// CameraHandler menangani HTTP requests untuk camera
type CameraHandler struct {
	cameraService service.CameraService
}

// NewCameraHandler membuat instance baru dari CameraHandler
func NewCameraHandler(cameraService service.CameraService) *CameraHandler {
	return &CameraHandler{
		cameraService: cameraService,
	}
}

// Create handler untuk membuat camera baru
func (h *CameraHandler) Create(c *fiber.Ctx) error {
	// Parse request body
	var req models.CreateCameraRequest
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
	if req.Name == "" || req.RTSPUrl == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeMissingFields,
				"Camera name and RTSP URL are required",
			),
		)
	}

	// Ambil user ID dari context
	userID := c.Locals("user_id").(string)

	// Proses create camera
	camera, err := h.cameraService.Create(&req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to create camera",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Camera created successfully",
		Data:    camera,
	})
}

// GetByID handler untuk mengambil camera berdasarkan ID
func (h *CameraHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	camera, err := h.cameraService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.NewErrorResponse(
				models.ErrCodeNotFound,
				"Camera not found",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Camera retrieved successfully",
		Data:    camera,
	})
}

// GetAll handler untuk mengambil semua camera dengan pagination
func (h *CameraHandler) GetAll(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	cameras, meta, err := h.cameraService.GetAll(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to retrieve cameras",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.PaginatedResponse{
		Success:    true,
		Message:    "Cameras retrieved successfully",
		Data:       cameras,
		Pagination: *meta,
	})
}

// Update handler untuk mengupdate camera
func (h *CameraHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse request body
	var req models.UpdateCameraRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid request body",
				err.Error(),
			),
		)
	}

	// Proses update
	camera, err := h.cameraService.Update(id, &req)
	if err != nil {
		// Check if camera not found
		if err.Error() == "camera not found" {
			return c.Status(fiber.StatusNotFound).JSON(
				models.NewErrorResponse(
					models.ErrCodeNotFound,
					"Camera not found",
					err.Error(),
				),
			)
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to update camera",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Camera updated successfully",
		Data:    camera,
	})
}

// Delete handler untuk menghapus camera
func (h *CameraHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.cameraService.Delete(id); err != nil {
		// Check if camera not found
		if err.Error() == "camera not found" {
			return c.Status(fiber.StatusNotFound).JSON(
				models.NewErrorResponse(
					models.ErrCodeNotFound,
					"Camera not found",
					err.Error(),
				),
			)
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to delete camera",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Camera deleted successfully",
	})
}

// GetByZone handler untuk mengambil camera berdasarkan zone
func (h *CameraHandler) GetByZone(c *fiber.Ctx) error {
	zone := c.Query("zone")
	if zone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeMissingFields,
				"Zone parameter is required",
			),
		)
	}

	cameras, err := h.cameraService.GetByZone(zone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to retrieve cameras",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Cameras retrieved successfully",
		Data:    cameras,
	})
}

// GetNearby handler untuk mengambil camera dalam radius tertentu
func (h *CameraHandler) GetNearby(c *fiber.Ctx) error {
	// Parse query parameters
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid latitude parameter",
				err.Error(),
			),
		)
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid longitude parameter",
				err.Error(),
			),
		)
	}

	radius, err := strconv.ParseFloat(c.Query("radius", "5"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeValidationFailed,
				"Invalid radius parameter",
				err.Error(),
			),
		)
	}

	cameras, err := h.cameraService.GetNearby(lat, lng, radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.NewErrorResponse(
				models.ErrCodeInternalError,
				"Failed to retrieve nearby cameras",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Nearby cameras retrieved successfully",
		Data:    cameras,
	})
}

// StartStream handler untuk memulai streaming camera
func (h *CameraHandler) StartStream(c *fiber.Ctx) error {
	id := c.Params("id")

	camera, err := h.cameraService.StartStream(id)
	if err != nil {
		// Check if camera not found
		if err.Error() == "camera not found" {
			return c.Status(fiber.StatusNotFound).JSON(
				models.NewErrorResponse(
					models.ErrCodeNotFound,
					"Camera not found",
					err.Error(),
				),
			)
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeServiceUnavailable,
				"Failed to start stream",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Stream started successfully",
		Data:    camera,
	})
}

// StopStream handler untuk menghentikan streaming camera
func (h *CameraHandler) StopStream(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.cameraService.StopStream(id); err != nil {
		// Check if camera not found
		if err.Error() == "camera not found" {
			return c.Status(fiber.StatusNotFound).JSON(
				models.NewErrorResponse(
					models.ErrCodeNotFound,
					"Camera not found",
					err.Error(),
				),
			)
		}

		return c.Status(fiber.StatusBadRequest).JSON(
			models.NewErrorResponse(
				models.ErrCodeServiceUnavailable,
				"Failed to stop stream",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Message: "Stream stopped successfully",
	})
}
