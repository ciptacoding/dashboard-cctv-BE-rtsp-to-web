package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect membuat koneksi ke PostgreSQL database
func Connect(dsn string) (*sql.DB, error) {
	// Buka koneksi database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Maksimal 25 koneksi terbuka
	db.SetMaxIdleConns(5)                  // Maksimal 5 koneksi idle
	db.SetConnMaxLifetime(5 * time.Minute) // Maksimal lifetime koneksi

	// Test koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✓ Successfully connected to PostgreSQL database")

	return db, nil
}

// RunMigrations menjalankan SQL migrations
func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Migration 1: Create users table
	migration1 := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(100) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'operator',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
	`

	if _, err := db.Exec(migration1); err != nil {
		return fmt.Errorf("migration 1 failed: %w", err)
	}

	// Migration 2: Create cameras table
	migration2 := `
		CREATE EXTENSION IF NOT EXISTS cube;
		CREATE EXTENSION IF NOT EXISTS earthdistance;

		CREATE TABLE IF NOT EXISTS cameras (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			rtsp_url TEXT NOT NULL,
			stream_id VARCHAR(255) UNIQUE,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			building VARCHAR(255),
			zone VARCHAR(255),
			ip_address VARCHAR(50),
			port INTEGER,
			manufacturer VARCHAR(100),
			model VARCHAR(100),
			resolution VARCHAR(50),
			fps INTEGER DEFAULT 25,
			tags TEXT[] DEFAULT ARRAY[]::text[],
			status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',
			last_seen TIMESTAMPTZ,
			is_active BOOLEAN DEFAULT true,
			created_by UUID REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_cameras_status ON cameras(status);
		CREATE INDEX IF NOT EXISTS idx_cameras_building ON cameras(building);
		CREATE INDEX IF NOT EXISTS idx_cameras_zone ON cameras(zone);
		CREATE INDEX IF NOT EXISTS idx_cameras_is_active ON cameras(is_active);
		CREATE INDEX IF NOT EXISTS idx_cameras_tags ON cameras USING GIN(tags);
		CREATE INDEX IF NOT EXISTS idx_cameras_location ON cameras USING GIST (ll_to_earth(latitude, longitude));
	`

	if _, err := db.Exec(migration2); err != nil {
		return fmt.Errorf("migration 2 failed: %w", err)
	}

	// Migration 3: Create activity logs table
	migration3 := `
		CREATE TABLE IF NOT EXISTS activity_logs (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE SET NULL,
			camera_id UUID REFERENCES cameras(id) ON DELETE CASCADE,
			action VARCHAR(100) NOT NULL,
			details JSONB,
			ip_address VARCHAR(50),
			user_agent TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_camera_id ON activity_logs(camera_id);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at DESC);
	`

	if _, err := db.Exec(migration3); err != nil {
		return fmt.Errorf("migration 3 failed: %w", err)
	}

	// Migration 4: Create token blacklist table
	migration4 := `
		CREATE TABLE IF NOT EXISTS token_blacklist (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			token_hash VARCHAR(64) NOT NULL UNIQUE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			reason VARCHAR(100) DEFAULT 'LOGOUT',
			blacklisted_at TIMESTAMPTZ DEFAULT NOW(),
			expires_at TIMESTAMPTZ NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_token_blacklist_token_hash ON token_blacklist(token_hash);
		CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);
		CREATE INDEX IF NOT EXISTS idx_token_blacklist_user_id ON token_blacklist(user_id);
	`

	if _, err := db.Exec(migration4); err != nil {
		return fmt.Errorf("migration 4 failed: %w", err)
	}

	log.Println("✓ Database migrations completed successfully")
	return nil
}
