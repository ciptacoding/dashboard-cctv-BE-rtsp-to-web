-- Migration: Create cameras table
-- File: migrations/002_create_cameras_table.sql

-- Enable earthdistance extension for geospatial queries
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

-- Create cameras table
CREATE TABLE IF NOT EXISTS cameras (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rtsp_url TEXT NOT NULL,
    stream_id VARCHAR(255) UNIQUE,
    
    -- Location information
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    building VARCHAR(255),
    zone VARCHAR(255),
    
    -- Camera specifications
    ip_address VARCHAR(50),
    port INTEGER,
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    resolution VARCHAR(50),
    fps INTEGER DEFAULT 25,
    
    -- Metadata
    tags TEXT[] DEFAULT ARRAY[]::text[],
    status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',
    last_seen TIMESTAMPTZ,
    
    -- Audit fields
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_cameras_status ON cameras(status);
CREATE INDEX idx_cameras_building ON cameras(building);
CREATE INDEX idx_cameras_zone ON cameras(zone);
CREATE INDEX idx_cameras_is_active ON cameras(is_active);
CREATE INDEX idx_cameras_tags ON cameras USING GIN(tags);

-- Create geospatial index for location-based queries
CREATE INDEX idx_cameras_location ON cameras USING GIST (
    ll_to_earth(latitude, longitude)
);

COMMENT ON TABLE cameras IS 'Tabel untuk menyimpan data kamera CCTV';
COMMENT ON COLUMN cameras.stream_id IS 'ID stream untuk RTSPtoWeb';
COMMENT ON COLUMN cameras.status IS 'Status: ONLINE, OFFLINE, ERROR, UNKNOWN';