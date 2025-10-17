-- Migration: Create users table
-- File: migrations/001_create_users_table.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
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

-- Create index for faster queries
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Insert default admin user (password: admin123)
-- Password hash adalah bcrypt hash dari "admin123"
INSERT INTO users (username, email, password_hash, role) 
VALUES (
    'admin',
    'admin@cctv-monitoring.com',
    '$2a$10$XZaEQzKqL8qKqO5K5YvLDu3.N8N5K9k7LZQH8YvLDu3N8N5K9k7LZ',
    'admin'
) ON CONFLICT (username) DO NOTHING;

COMMENT ON TABLE users IS 'Tabel untuk menyimpan data pengguna sistem';
COMMENT ON COLUMN users.role IS 'Role: admin, operator, viewer';