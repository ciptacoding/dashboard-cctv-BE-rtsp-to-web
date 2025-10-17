package service

import (
	"fmt"

	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
	"cctv-monitoring-backend/internal/utils"
)

// AuthService adalah interface untuk business logic authentication
type AuthService interface {
	Login(username, password string, jwtSecret string, jwtExpiration string) (*models.LoginResponse, error)
	Register(req *models.CreateUserRequest) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login melakukan authentication user
func (s *authService) Login(username, password, jwtSecret, jwtExpiration string) (*models.LoginResponse, error) {
	// Cari user berdasarkan username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Cek apakah user aktif
	if !user.IsActive {
		return nil, fmt.Errorf("user is inactive")
	}

	// Validasi password
	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		jwtSecret,
		utils.ParseDuration(jwtExpiration),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response
	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Register mendaftarkan user baru
func (s *authService) Register(req *models.CreateUserRequest) (*models.User, error) {
	// Cek apakah username sudah ada
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Cek apakah email sudah ada
	existingEmail, _ := s.userRepo.GetByEmail(req.Email)
	if existingEmail != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Buat user baru
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		IsActive:     true,
	}

	// Simpan ke database
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
