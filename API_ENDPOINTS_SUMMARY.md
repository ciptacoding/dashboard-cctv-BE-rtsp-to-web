# API Endpoints Summary

Dokumen ini berisi ringkasan semua endpoint API yang tersedia.

## Base URL
```
http://localhost:8080
```

## WebSocket URL
```
ws://localhost:8080/ws?token=<jwt_token>
```

---

## Public Endpoints (Tidak Perlu Auth)

### Health Check
- **GET** `/health` - Check server status

### Auth
- **POST** `/api/v1/auth/login` - Login user
- **POST** `/api/v1/auth/register` - Register new user

---

## Protected Endpoints (Perlu Auth Token)

### Auth (Protected)
- **GET** `/api/v1/auth/me` - Get current user info
- **POST** `/api/v1/auth/logout` - Logout user

### Camera Management
- **GET** `/api/v1/cameras` - Get all cameras (with pagination)
- **GET** `/api/v1/cameras/:id` - Get camera by ID
- **POST** `/api/v1/cameras` - Create new camera
- **PUT** `/api/v1/cameras/:id` - Update camera
- **DELETE** `/api/v1/cameras/:id` - Delete camera

### Camera Filtering
- **GET** `/api/v1/cameras/zone/filter?zone=<zone>` - Get cameras by zone
- **GET** `/api/v1/cameras/nearby?lat=<lat>&lng=<lng>&radius=<radius>` - Get nearby cameras

### Stream Management
- **POST** `/api/v1/cameras/:id/stream/start` - Start camera stream
- **POST** `/api/v1/cameras/:id/stream/stop` - Stop camera stream
- **POST** `/api/v1/cameras/:id/stream/error` - Report stream error

### Preview
- **GET** `/api/v1/cameras/:id/preview` - Get camera preview for video display

---

## Request/Response Format

### Request Headers (Protected Endpoints)
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Success Response Format
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... }
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Error message",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": "Technical details (optional)"
  }
}
```

---

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error, missing fields)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (user inactive)
- `404` - Not Found
- `409` - Conflict (resource already exists)
- `500` - Internal Server Error
- `503` - Service Unavailable

---

## Quick Reference

### Get All Cameras
```bash
curl -X GET "http://localhost:8080/api/v1/cameras?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

### Get Camera Preview
```bash
curl -X GET "http://localhost:8080/api/v1/cameras/{id}/preview" \
  -H "Authorization: Bearer <token>"
```

### Report Stream Error
```bash
curl -X POST "http://localhost:8080/api/v1/cameras/{id}/stream/error" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"error_type": "timeout", "message": "HLS timeout"}'
```

### Start Stream
```bash
curl -X POST "http://localhost:8080/api/v1/cameras/{id}/stream/start" \
  -H "Authorization: Bearer <token>"
```

### Stop Stream
```bash
curl -X POST "http://localhost:8080/api/v1/cameras/{id}/stream/stop" \
  -H "Authorization: Bearer <token>"
```

---

## WebSocket Events

### Client → Server
- `ping` - Keep-alive ping

### Server → Client
- `connected` - Connection established
- `pong` - Response to ping
- `camera_status` - Camera status update
- `stream_update` - Stream status update (offline, online, frozen, restarted, restart_failed)

---

## Field Reference

### Camera Object Fields
- `id` - Camera UUID
- `name` - Camera name
- `rtsp_url` - RTSP stream URL
- `latitude` - Latitude coordinate
- `longitude` - Longitude coordinate
- `status` - Status: ONLINE, OFFLINE, READY, ERROR, FROZEN, UNKNOWN
- `status_message` - Human readable status message
- `hls_url` - HLS stream URL (if stream active)
- `snapshot_url` - Snapshot image URL (if stream active)
- `stream_id` - Stream ID from RTSPtoWeb
- `has_stream` - Boolean indicating if stream is active
- `last_seen` - Last seen timestamp (ISO 8601)
- `is_active` - Camera is active
- `fps` - Frames per second
- `tags` - Array of tags
- `description` - Camera description
- `building` - Building name
- `zone` - Zone name
- `ip_address` - Camera IP address
- `port` - Camera port
- `manufacturer` - Camera manufacturer
- `model` - Camera model
- `resolution` - Camera resolution
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

---

## Implementation Checklist

- [ ] Authentication flow (login, register, logout)
- [ ] Get all cameras with pagination
- [ ] Get camera by ID
- [ ] Create/Update/Delete camera
- [ ] Filter cameras by zone
- [ ] Get nearby cameras
- [ ] Start/Stop stream
- [ ] Get camera preview
- [ ] Report stream error
- [ ] WebSocket connection
- [ ] Handle WebSocket messages
- [ ] Display camera status
- [ ] Handle video player errors
- [ ] Auto-report timeout/HLS errors


