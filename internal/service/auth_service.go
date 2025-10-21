package service

import (
	"errors"
	"fmt"
	"time"

	"cctv-monitoring-backend/internal/models"
	"cctv-monitoring-backend/internal/repository"
	"cctv-monitoring-backend/internal/utils"
)

// Custom errors untuk auth service
var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrPasswordMismatch   = errors.New("password does not match")
	ErrTokenBlacklisted   = errors.New("token has been revoked")
)

// AuthService adalah interface untuk business logic authentication
type AuthService interface {
	Login(username, password string, jwtSecret string, jwtExpiration string) (*models.LoginResponse, error)
	Register(req *models.CreateUserRequest) (*models.User, error)
	Logout(token, userID string, jwtExpiration string) error
	VerifyToken(token, jwtSecret string) (*utils.JWTClaims, error)
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Login melakukan authentication user
func (s *authService) Login(username, password, jwtSecret, jwtExpiration string) (*models.LoginResponse, error) {
	// Cari user berdasarkan username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Cek apakah user aktif
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Validasi password
	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
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

// Logout melakukan logout user dengan blacklist token
func (s *authService) Logout(token, userID string, jwtExpiration string) error {
	// Hash token untuk disimpan di blacklist
	tokenHash := utils.HashToken(token)

	// Calculate token expiration time
	expiresAt := time.Now().Add(utils.ParseDuration(jwtExpiration))

	// Blacklist token
	if err := s.tokenRepo.BlacklistToken(tokenHash, userID, "LOGOUT", expiresAt); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// VerifyToken memverifikasi token dan check blacklist
func (s *authService) VerifyToken(token, jwtSecret string) (*utils.JWTClaims, error) {
	// Validate JWT token
	claims, err := utils.ValidateToken(token, jwtSecret)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	tokenHash := utils.HashToken(token)
	isBlacklisted, err := s.tokenRepo.IsTokenBlacklisted(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	return claims, nil
}
