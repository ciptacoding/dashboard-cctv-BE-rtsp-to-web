package repository

import (
	"database/sql"
	"fmt"

	"cctv-monitoring-backend/internal/models"
)

// UserRepository adalah interface untuk operasi database user
type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance baru dari UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create membuat user baru di database
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByUsername mencari user berdasarkan username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByID mencari user berdasarkan ID
func (r *userRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail mencari user berdasarkan email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
