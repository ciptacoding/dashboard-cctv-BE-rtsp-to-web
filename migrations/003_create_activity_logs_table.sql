-- Migration: Create activity logs table
-- File: migrations/003_create_activity_logs_table.sql

-- Create activity_logs table
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

-- Create indexes
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_camera_id ON activity_logs(camera_id);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at DESC);
CREATE INDEX idx_activity_logs_details ON activity_logs USING GIN(details);

COMMENT ON TABLE activity_logs IS 'Tabel untuk audit trail semua aktivitas';
COMMENT ON COLUMN activity_logs.action IS 'Action: LOGIN, LOGOUT, CREATE_CAMERA, UPDATE_CAMERA, DELETE_CAMERA, VIEW_STREAM, etc';