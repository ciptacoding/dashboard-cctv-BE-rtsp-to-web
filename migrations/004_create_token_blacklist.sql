-- Migration: Create token blacklist table
-- File: migrations/004_create_token_blacklist.sql

CREATE TABLE IF NOT EXISTS token_blacklist (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_hash VARCHAR(64) NOT NULL UNIQUE, -- SHA256 hash of JWT token
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(100) DEFAULT 'LOGOUT',
    blacklisted_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL -- Token expiration time
);

-- Index untuk performance
CREATE INDEX idx_token_blacklist_token_hash ON token_blacklist(token_hash);
CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);
CREATE INDEX idx_token_blacklist_user_id ON token_blacklist(user_id);

-- Auto cleanup expired tokens (optional, bisa pakai cron job)
-- Function untuk cleanup
CREATE OR REPLACE FUNCTION cleanup_expired_blacklist()
RETURNS void AS $$
BEGIN
    DELETE FROM token_blacklist WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

COMMENT ON TABLE token_blacklist IS 'Tabel untuk menyimpan JWT token yang di-blacklist (logout)';
COMMENT ON COLUMN token_blacklist.reason IS 'Reason: LOGOUT, REVOKED, COMPROMISED, etc';