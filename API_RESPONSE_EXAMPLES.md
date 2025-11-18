# API Response Examples untuk Frontend Implementation

Dokumen ini berisi contoh response untuk semua endpoint API yang perlu diimplementasikan di Frontend.

## Base Response Structure

### Success Response
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Technical details (optional)"
  }
}
```

## Error Codes

- `NOT_FOUND` - Resource tidak ditemukan
- `VALIDATION_FAILED` - Validasi input gagal
- `MISSING_FIELDS` - Field yang diperlukan tidak ada
- `INTERNAL_ERROR` - Server error
- `SERVICE_UNAVAILABLE` - Service tidak tersedia
- `UNAUTHORIZED` - Tidak terautentikasi
- `INVALID_CREDENTIALS` - Kredensial tidak valid

---

## Camera Endpoints

### 1. Get All Cameras
**GET** `/api/v1/cameras?page=1&page_size=10`

#### Success Response (200)
```json
{
  "success": true,
  "message": "Cameras retrieved successfully",
  "data": [
    {
      "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
      "name": "Lobby",
      "rtsp_url": "rtsp://admin:password@192.168.1.100:554/stream",
      "latitude": -6.2088,
      "longitude": 106.8456,
      "fps": 25,
      "tags": ["lobby", "entrance"],
      "status": "ONLINE",
      "status_message": "Camera is online and streaming",
      "is_active": true,
      "hls_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/hls/live/index.m3u8",
      "snapshot_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/jpeg",
      "stream_id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
      "description": "Main lobby camera",
      "building": "Building A",
      "zone": "Zone 1",
      "ip_address": "192.168.1.100",
      "port": 554,
      "manufacturer": "Hikvision",
      "model": "DS-2CD2342WD-I",
      "resolution": "1920x1080",
      "last_seen": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-10T08:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "049d3cb9-2f1a-403c-af3f-e846b80e761d",
      "name": "Data Center",
      "rtsp_url": "rtsp://admin:password@192.168.1.101:554/stream",
      "latitude": -6.2089,
      "longitude": 106.8457,
      "fps": 30,
      "tags": ["datacenter", "server"],
      "status": "OFFLINE",
      "status_message": "Camera is offline (disconnected 5 minute(s) ago). Attempting to reconnect...",
      "is_active": true,
      "hls_url": "",
      "snapshot_url": "",
      "stream_id": "",
      "description": "Data center monitoring",
      "building": "Building B",
      "zone": "Zone 2",
      "ip_address": "192.168.1.101",
      "port": 554,
      "manufacturer": "Dahua",
      "model": "IPC-HDW2431T-AS-S2",
      "resolution": "1920x1080",
      "last_seen": "2024-01-15T10:25:00Z",
      "created_at": "2024-01-10T08:00:00Z",
      "updated_at": "2024-01-15T10:25:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 25,
    "total_pages": 3
  }
}
```

---

### 2. Get Camera By ID
**GET** `/api/v1/cameras/:id`

#### Success Response - Camera Online (200)
```json
{
  "success": true,
  "message": "Camera retrieved successfully",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "name": "Lobby",
    "rtsp_url": "rtsp://admin:password@192.168.1.100:554/stream",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "fps": 25,
    "tags": ["lobby", "entrance"],
    "status": "ONLINE",
    "status_message": "Camera is online and streaming",
    "is_active": true,
    "hls_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/hls/live/index.m3u8",
    "snapshot_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/jpeg",
    "stream_id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "description": "Main lobby camera",
    "building": "Building A",
    "zone": "Zone 1",
    "ip_address": "192.168.1.100",
    "port": 554,
    "manufacturer": "Hikvision",
    "model": "DS-2CD2342WD-I",
    "resolution": "1920x1080",
    "last_seen": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Success Response - Camera Offline (200)
```json
{
  "success": true,
  "message": "Camera retrieved successfully. Camera is currently offline",
  "data": {
    "id": "049d3cb9-2f1a-403c-af3f-e846b80e761d",
    "name": "Data Center",
    "rtsp_url": "rtsp://admin:password@192.168.1.101:554/stream",
    "latitude": -6.2089,
    "longitude": 106.8457,
    "fps": 30,
    "tags": ["datacenter", "server"],
    "status": "OFFLINE",
    "status_message": "Camera is offline (disconnected 5 minute(s) ago). Attempting to reconnect...",
    "is_active": true,
    "hls_url": "",
    "snapshot_url": "",
    "stream_id": "",
    "description": "Data center monitoring",
    "building": "Building B",
    "zone": "Zone 2",
    "ip_address": "192.168.1.101",
    "port": 554,
    "manufacturer": "Dahua",
    "model": "IPC-HDW2431T-AS-S2",
    "resolution": "1920x1080",
    "last_seen": "2024-01-15T10:25:00Z",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-15T10:25:00Z"
  }
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

---

### 3. Create Camera
**POST** `/api/v1/cameras`

#### Request Body
```json
{
  "name": "Parking Lot",
  "description": "Main parking lot camera",
  "rtsp_url": "rtsp://admin:password@192.168.1.102:554/stream",
  "latitude": -6.2090,
  "longitude": 106.8458,
  "building": "Building C",
  "zone": "Zone 3",
  "ip_address": "192.168.1.102",
  "port": 554,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2342WD-I",
  "resolution": "1920x1080",
  "fps": 25,
  "tags": ["parking", "outdoor"],
  "status": "UNKNOWN"
}
```

#### Success Response (201)
```json
{
  "success": true,
  "message": "Camera created successfully",
  "data": {
    "id": "7f8e9d0a-1b2c-3d4e-5f6a-7b8c9d0e1f2a",
    "name": "Parking Lot",
    "rtsp_url": "rtsp://admin:password@192.168.1.102:554/stream",
    "latitude": -6.2090,
    "longitude": 106.8458,
    "fps": 25,
    "tags": ["parking", "outdoor"],
    "status": "READY",
    "status_message": "Camera is online and streaming",
    "is_active": true,
    "hls_url": "http://localhost:8083/stream/7f8e9d0a-1b2c-3d4e-5f6a-7b8c9d0e1f2a/channel/0/hls/live/index.m3u8",
    "snapshot_url": "http://localhost:8083/stream/7f8e9d0a-1b2c-3d4e-5f6a-7b8c9d0e1f2a/channel/0/jpeg",
    "stream_id": "7f8e9d0a-1b2c-3d4e-5f6a-7b8c9d0e1f2a",
    "description": "Main parking lot camera",
    "building": "Building C",
    "zone": "Zone 3",
    "ip_address": "192.168.1.102",
    "port": 554,
    "manufacturer": "Hikvision",
    "model": "DS-2CD2342WD-I",
    "resolution": "1920x1080",
    "last_seen": "2024-01-15T10:35:00Z",
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

#### Error Response - Validation Failed (400)
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Invalid request body",
    "details": "json: cannot unmarshal string into Go struct field CreateCameraRequest.fps of type int"
  }
}
```

#### Error Response - Missing Fields (400)
```json
{
  "success": false,
  "message": "Camera name and RTSP URL are required",
  "error": {
    "code": "MISSING_FIELDS",
    "message": "Camera name and RTSP URL are required"
  }
}
```

---

### 4. Update Camera
**PUT** `/api/v1/cameras/:id`

#### Request Body
```json
{
  "name": "Lobby Updated",
  "description": "Updated description",
  "fps": 30
}
```

#### Success Response - Camera Online (200)
```json
{
  "success": true,
  "message": "Camera updated successfully",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "name": "Lobby Updated",
    "rtsp_url": "rtsp://admin:password@192.168.1.100:554/stream",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "fps": 30,
    "tags": ["lobby", "entrance"],
    "status": "ONLINE",
    "status_message": "Camera is online and streaming",
    "is_active": true,
    "hls_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/hls/live/index.m3u8",
    "snapshot_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/jpeg",
    "stream_id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "description": "Updated description",
    "building": "Building A",
    "zone": "Zone 1",
    "ip_address": "192.168.1.100",
    "port": 554,
    "manufacturer": "Hikvision",
    "model": "DS-2CD2342WD-I",
    "resolution": "1920x1080",
    "last_seen": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-15T10:40:00Z"
  }
}
```

#### Success Response - Camera Offline (200)
```json
{
  "success": true,
  "message": "Camera updated successfully. Note: Camera is currently offline",
  "data": {
    "id": "049d3cb9-2f1a-403c-af3f-e846b80e761d",
    "name": "Data Center Updated",
    "status": "OFFLINE",
    "status_message": "Camera is offline (disconnected 5 minute(s) ago). Attempting to reconnect...",
    ...
  }
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

---

### 5. Delete Camera
**DELETE** `/api/v1/cameras/:id`

#### Success Response (200)
```json
{
  "success": true,
  "message": "Camera deleted successfully"
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has already been deleted"
  }
}
```

---

### 6. Start Stream
**POST** `/api/v1/cameras/:id/stream/start`

#### Success Response (200)
```json
{
  "success": true,
  "message": "Stream started successfully",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "name": "Lobby",
    "status": "READY",
    "status_message": "Camera is online and streaming",
    "hls_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/hls/live/index.m3u8",
    "snapshot_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/jpeg",
    "stream_id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "last_seen": "2024-01-15T10:45:00Z",
    ...
  }
}
```

#### Success Response - Camera Offline (200)
```json
{
  "success": true,
  "message": "Stream started but camera appears to be offline. Please check the camera connection",
  "data": {
    "id": "049d3cb9-2f1a-403c-af3f-e846b80e761d",
    "name": "Data Center",
    "status": "OFFLINE",
    "status_message": "Camera is offline. Attempting to reconnect...",
    ...
  }
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

#### Error Response - Stream Start Failed (503)
```json
{
  "success": false,
  "message": "Failed to start stream",
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "Failed to start stream",
    "details": "Unable to connect to camera stream. Please check the RTSP URL and camera connection"
  }
}
```

---

### 7. Stop Stream
**POST** `/api/v1/cameras/:id/stream/stop`

#### Success Response (200)
```json
{
  "success": true,
  "message": "Stream stopped successfully. Camera is now offline"
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

---

### 8. Get Camera Preview
**GET** `/api/v1/cameras/:id/preview`

Endpoint ini digunakan untuk mendapatkan informasi preview video kamera saat pin point di klik di map.

#### Success Response (200)
```json
{
  "success": true,
  "message": "Camera preview retrieved successfully",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "name": "Lobby",
    "status": "ONLINE",
    "status_message": "Camera is online and streaming",
    "hls_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/hls/live/index.m3u8",
    "snapshot_url": "http://localhost:8083/stream/01ed06d5-0b38-4d42-b709-17a89b03a7c0/channel/0/jpeg",
    "has_stream": true,
    "last_seen": "2024-01-15T10:30:00Z"
  }
}
```

#### Success Response - Camera Offline (200)
```json
{
  "success": true,
  "message": "Camera preview retrieved successfully",
  "data": {
    "id": "049d3cb9-2f1a-403c-af3f-e846b80e761d",
    "name": "Data Center",
    "status": "OFFLINE",
    "status_message": "Camera is offline (disconnected 5 minute(s) ago). Attempting to reconnect...",
    "hls_url": "",
    "snapshot_url": "",
    "has_stream": false,
    "last_seen": "2024-01-15T10:25:00Z"
  }
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

---

### 9. Report Stream Error
**POST** `/api/v1/cameras/:id/stream/error`

Endpoint ini digunakan untuk melaporkan error stream dari frontend (timeout, HLS error, dll) dan langsung mengubah status kamera menjadi offline.

#### Request Body
```json
{
  "error_type": "timeout",
  "message": "HLS stream timeout after 10 seconds"
}
```

**Error Types:**
- `timeout` - Stream timeout
- `hls_error` - HLS playback error
- `network_error` - Network connection error
- `decode_error` - Video decode error
- `other` - Other errors

#### Success Response (200)
```json
{
  "success": true,
  "message": "Stream error reported. Camera status updated to offline"
}
```

#### Error Response - Camera Not Found (404)
```json
{
  "success": false,
  "message": "Camera not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "Camera not found",
    "details": "The requested camera does not exist or has been deleted"
  }
}
```

#### Error Response - Missing Fields (400)
```json
{
  "success": false,
  "message": "Error type is required",
  "error": {
    "code": "MISSING_FIELDS",
    "message": "Error type is required"
  }
}
```

**Note:** Setelah error dilaporkan, status kamera akan langsung diubah menjadi `OFFLINE` di database dan broadcast ke semua WebSocket clients.

---

## Camera Status Values

### Status Field Values
- `ONLINE` - Camera online dan streaming
- `OFFLINE` - Camera offline
- `READY` - Stream siap
- `ERROR` - Ada error pada kamera
- `FROZEN` - Stream frozen
- `UNKNOWN` - Status tidak diketahui (default untuk kamera baru)

### Status Message Examples

#### ONLINE
```json
{
  "status": "ONLINE",
  "status_message": "Camera is online and streaming"
}
```

#### OFFLINE - Just Disconnected
```json
{
  "status": "OFFLINE",
  "status_message": "Camera is offline (just disconnected)"
}
```

#### OFFLINE - Disconnected Recently
```json
{
  "status": "OFFLINE",
  "status_message": "Camera is offline (disconnected 2 minute(s) ago)"
}
```

#### OFFLINE - Disconnected Long Time
```json
{
  "status": "OFFLINE",
  "status_message": "Camera is offline (disconnected 1 hour(s) ago). Attempting to reconnect..."
}
```

#### ERROR
```json
{
  "status": "ERROR",
  "status_message": "Camera encountered an error. Please check the connection"
}
```

#### FROZEN
```json
{
  "status": "FROZEN",
  "status_message": "Camera stream appears frozen. Refreshing..."
}
```

#### Stream Not Started
```json
{
  "status": "UNKNOWN",
  "status_message": "Camera stream not started"
}
```

---

## WebSocket Messages

### Connection
**URL:** `ws://localhost:8080/ws?token=YOUR_JWT_TOKEN`

### Message Types

#### 1. Connected Message (on connect)
```json
{
  "type": "connected",
  "data": {
    "message": "Connected to CCTV Monitoring WebSocket",
    "clients": 5
  }
}
```

#### 2. Camera Status Update
```json
{
  "type": "camera_status",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "status": "OFFLINE",
    "last_seen": "2024-01-15T10:30:00Z"
  }
}
```

#### 3. Stream Update
```json
{
  "type": "stream_update",
  "data": {
    "id": "01ed06d5-0b38-4d42-b709-17a89b03a7c0",
    "name": "Lobby",
    "status": "offline",
    "message": "Camera stream is offline. Attempting to reconnect..."
  }
}
```

**Stream Update Status Values:**
- `offline` - Stream offline
- `online` - Stream kembali online
- `frozen` - Stream frozen
- `restarted` - Stream berhasil di-restart
- `restart_failed` - Restart gagal

#### 4. Ping/Pong
**Client sends:**
```json
{
  "type": "ping",
  "data": "2024-01-15T10:30:00Z"
}
```

**Server responds:**
```json
{
  "type": "pong",
  "data": {
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

---

## Frontend Implementation Guide

### 1. Handling Camera Status

```typescript
interface Camera {
  id: string;
  name: string;
  status: 'ONLINE' | 'OFFLINE' | 'READY' | 'ERROR' | 'FROZEN' | 'UNKNOWN';
  status_message: string;
  hls_url?: string;
  snapshot_url?: string;
  stream_id?: string;
  last_seen?: string;
  // ... other fields
}

// Get status color
function getStatusColor(status: string): string {
  switch (status) {
    case 'ONLINE':
    case 'READY':
      return 'green';
    case 'OFFLINE':
      return 'red';
    case 'ERROR':
      return 'orange';
    case 'FROZEN':
      return 'yellow';
    default:
      return 'gray';
  }
}

// Get status icon
function getStatusIcon(status: string): string {
  switch (status) {
    case 'ONLINE':
    case 'READY':
      return '✓';
    case 'OFFLINE':
      return '✗';
    case 'ERROR':
      return '⚠';
    case 'FROZEN':
      return '🧊';
    default:
      return '?';
  }
}
```

### 2. Displaying Camera Status

```typescript
// Display status badge
function CameraStatusBadge({ camera }: { camera: Camera }) {
  return (
    <div className={`status-badge status-${camera.status.toLowerCase()}`}>
      <span className="status-icon">{getStatusIcon(camera.status)}</span>
      <span className="status-text">{camera.status}</span>
      {camera.status_message && (
        <span className="status-message">{camera.status_message}</span>
      )}
    </div>
  );
}
```

### 3. Handling WebSocket Updates

```typescript
// WebSocket connection
const ws = new WebSocket('ws://localhost:8080/ws?token=' + token);

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch (message.type) {
    case 'camera_status':
      // Update camera status in your state
      updateCameraStatus(message.data.id, {
        status: message.data.status,
        last_seen: message.data.last_seen
      });
      break;
      
    case 'stream_update':
      // Show notification
      showNotification({
        type: message.data.status,
        title: message.data.name,
        message: message.data.message
      });
      break;
      
    case 'connected':
      console.log('Connected to WebSocket');
      break;
  }
};
```

### 4. Error Handling

```typescript
async function fetchCamera(id: string) {
  try {
    const response = await fetch(`/api/v1/cameras/${id}`);
    const data = await response.json();
    
    if (!data.success) {
      // Handle error
      switch (data.error.code) {
        case 'NOT_FOUND':
          showError('Camera tidak ditemukan');
          break;
        case 'SERVICE_UNAVAILABLE':
          showError('Service tidak tersedia. Silakan coba lagi.');
          break;
        default:
          showError(data.error.message);
      }
      return null;
    }
    
    return data.data;
  } catch (error) {
    showError('Terjadi kesalahan saat mengambil data kamera');
    return null;
  }
}
```

### 5. Displaying Camera List with Status

```typescript
function CameraList({ cameras }: { cameras: Camera[] }) {
  return (
    <div className="camera-list">
      {cameras.map(camera => (
        <div key={camera.id} className="camera-card">
          <div className="camera-header">
            <h3>{camera.name}</h3>
            <CameraStatusBadge camera={camera} />
          </div>
          
          {camera.status === 'ONLINE' && camera.hls_url && (
            <video src={camera.hls_url} controls />
          )}
          
          {camera.status === 'OFFLINE' && camera.snapshot_url && (
            <img src={camera.snapshot_url} alt={camera.name} />
          )}
          
          {camera.status_message && (
            <p className="status-message">{camera.status_message}</p>
          )}
          
          {camera.last_seen && (
            <p className="last-seen">
              Last seen: {new Date(camera.last_seen).toLocaleString()}
            </p>
          )}
        </div>
      ))}
    </div>
  );
}
```

### 6. Preview Video on Map Pin Click

```typescript
// Fetch preview when pin is clicked
async function handlePinClick(cameraId: string) {
  try {
    const response = await fetch(`/api/v1/cameras/${cameraId}/preview`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const data = await response.json();
    
    if (data.success) {
      const preview = data.data;
      showPreviewModal(preview);
    }
  } catch (error) {
    console.error('Failed to load preview:', error);
  }
}

// Display preview in modal
function PreviewModal({ preview }: { preview: CameraPreview }) {
  return (
    <div className="preview-modal">
      <h3>{preview.name}</h3>
      <CameraStatusBadge camera={preview} />
      
      {preview.status === 'ONLINE' && preview.hls_url ? (
        <VideoPlayer 
          src={preview.hls_url}
          onError={() => handleStreamError(preview.id, 'hls_error')}
          onTimeout={() => handleStreamError(preview.id, 'timeout')}
        />
      ) : (
        <div className="offline-placeholder">
          <p>{preview.status_message}</p>
          {preview.snapshot_url && (
            <img src={preview.snapshot_url} alt={preview.name} />
          )}
        </div>
      )}
    </div>
  );
}
```

### 7. Handle Stream Error and Report to Backend

```typescript
// Report stream error to backend
async function handleStreamError(cameraId: string, errorType: string) {
  try {
    const response = await fetch(`/api/v1/cameras/${cameraId}/stream/error`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        error_type: errorType,
        message: `Stream error: ${errorType}`
      })
    });
    
    const data = await response.json();
    
    if (data.success) {
      // Update local state
      updateCameraStatus(cameraId, { status: 'OFFLINE' });
      
      // Show notification
      showNotification({
        type: 'error',
        message: `Camera stream error reported. Status updated to offline.`
      });
    }
  } catch (error) {
    console.error('Failed to report stream error:', error);
  }
}

// Video player component with error handling
function VideoPlayer({ 
  src, 
  onError, 
  onTimeout 
}: { 
  src: string; 
  onError: () => void;
  onTimeout: () => void;
}) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const timeoutRef = useRef<NodeJS.Timeout>();
  
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;
    
    // Set timeout (10 seconds)
    timeoutRef.current = setTimeout(() => {
      if (video.readyState < 2) { // HAVE_CURRENT_DATA
        onTimeout();
      }
    }, 10000);
    
    // Handle video errors
    const handleError = () => {
      clearTimeout(timeoutRef.current);
      onError();
    };
    
    // Handle video loaded
    const handleLoadedData = () => {
      clearTimeout(timeoutRef.current);
    };
    
    video.addEventListener('error', handleError);
    video.addEventListener('loadeddata', handleLoadedData);
    
    return () => {
      clearTimeout(timeoutRef.current);
      video.removeEventListener('error', handleError);
      video.removeEventListener('loadeddata', handleLoadedData);
    };
  }, [src, onError, onTimeout]);
  
  return (
    <video
      ref={videoRef}
      src={src}
      controls
      autoPlay
      playsInline
      onError={onError}
    />
  );
}
```

---

## Testing Examples

### Test Camera Offline
1. Stop camera stream: `POST /api/v1/cameras/:id/stream/stop`
2. Check status: `GET /api/v1/cameras/:id`
3. Should return `status: "OFFLINE"` with appropriate `status_message`

### Test Camera Not Found
1. Request non-existent camera: `GET /api/v1/cameras/invalid-id`
2. Should return 404 with error code `NOT_FOUND`

### Test WebSocket Updates
1. Connect to WebSocket
2. Stop a camera stream
3. Should receive `camera_status` and `stream_update` messages

---

## Notes

1. **Status Message**: Selalu cek `status_message` untuk informasi detail tentang status kamera
2. **HLS URL**: Hanya tersedia jika `status` adalah `ONLINE` atau `READY` dan `stream_id` ada
3. **Snapshot URL**: Sama seperti HLS URL, hanya tersedia jika stream aktif
4. **Last Seen**: Format ISO 8601 (RFC3339), bisa digunakan untuk menampilkan waktu terakhir kamera online
5. **WebSocket**: Gunakan untuk real-time updates tanpa polling
6. **Error Handling**: Selalu cek `success` field sebelum menggunakan `data`

