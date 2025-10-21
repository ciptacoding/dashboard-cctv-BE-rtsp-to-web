package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// TokenRepository adalah interface untuk operasi database token blacklist
type TokenRepository interface {
	BlacklistToken(tokenHash, userID, reason string, expiresAt time.Time) error
	IsTokenBlacklisted(tokenHash string) (bool, error)
	CleanupExpiredTokens() error
	RevokeAllUserTokens(userID string) error
}

type tokenRepository struct {
	db *sql.DB
}

// NewTokenRepository membuat instance baru dari TokenRepository
func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

// BlacklistToken menambahkan token ke blacklist
func (r *tokenRepository) BlacklistToken(tokenHash, userID, reason string, expiresAt time.Time) error {
	query := `
		INSERT INTO token_blacklist (token_hash, user_id, reason, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (token_hash) DO NOTHING
	`

	_, err := r.db.Exec(query, tokenHash, userID, reason, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// IsTokenBlacklisted mengecek apakah token ada di blacklist
func (r *tokenRepository) IsTokenBlacklisted(tokenHash string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM token_blacklist 
			WHERE token_hash = $1 AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.QueryRow(query, tokenHash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return exists, nil
}

// CleanupExpiredTokens menghapus token yang sudah expired dari blacklist
func (r *tokenRepository) CleanupExpiredTokens() error {
	query := "DELETE FROM token_blacklist WHERE expires_at < NOW()"

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("✓ Cleaned up %d expired tokens from blacklist\n", rowsAffected)
	}

	return nil
}

// RevokeAllUserTokens me-revoke semua token user (untuk security breach)
func (r *tokenRepository) RevokeAllUserTokens(userID string) error {
	// Untuk revoke semua token user, kita bisa:
	// 1. Update user's "token_version" di database
	// 2. Atau blacklist semua token yang belum expired

	// Implementasi sederhana: set reason = 'REVOKED_ALL'
	query := `
		UPDATE token_blacklist 
		SET reason = 'REVOKED_ALL'
		WHERE user_id = $1 AND expires_at > NOW()
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	return nil
}
