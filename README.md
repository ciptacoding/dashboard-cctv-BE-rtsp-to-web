# CCTV Monitoring Backend API

Backend API untuk sistem monitoring CCTV berbasis Web GIS menggunakan Golang, PostgreSQL, dan RTSPtoWeb.

## 📋 Fitur

- ✅ Authentication & Authorization (JWT)
- ✅ CRUD Camera CCTV
- ✅ Integrasi RTSPtoWeb untuk streaming RTSP
- ✅ Geospatial queries (nearby cameras, zone filtering)
- ✅ Stream management (start/stop)
- ✅ RESTful API dengan response konsisten
- ✅ Database migrations
- ✅ Docker deployment ready

## 🏗️ Arsitektur

```
Frontend (Web GIS)
       ↓
Backend API (Golang + Fiber)
       ↓
   PostgreSQL + RTSPtoWeb
       ↓
  CCTV Cameras (RTSP)
```

## 🛠️ Tech Stack

- **Backend**: Golang 1.21 + Fiber Framework
- **Database**: PostgreSQL 15 (dengan PostGIS extensions)
- **Streaming**: RTSPtoWeb
- **Authentication**: JWT
- **Deployment**: Docker & Docker Compose

## 📦 Struktur Project

```
cctv-monitoring-backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── config/                  # Konfigurasi
│   ├── database/                # Database connection
│   ├── models/                  # Data models
│   ├── repository/              # Database operations
│   ├── service/                 # Business logic
│   ├── handler/                 # HTTP handlers
│   ├── middleware/              # Middlewares
│   └── utils/                   # Utilities
├── migrations/                  # SQL migrations
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (untuk development)

### 1. Clone Repository

```bash
git clone <repository-url>
cd cctv-monitoring-backend
```

### 2. Setup Environment Variables

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

### 3. Jalankan dengan Docker Compose

```bash
docker-compose up -d
```

Services yang akan berjalan:
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **RTSPtoWeb**: http://localhost:8083

### 4. Test API

```bash
# Health check
curl http://localhost:8080/health

# Login (default admin)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 📚 API Documentation

### Authentication

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "operator1",
  "email": "operator1@example.com",
  "password": "password123",
  "role": "operator"
}
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

### Camera Management

#### Create Camera
```http
POST /api/v1/cameras
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Camera Lobby 1",
  "description": "Camera di lobby utama",
  "rtsp_url": "rtsp://username:password@192.168.1.100:554/stream",
  "latitude": -6.200000,
  "longitude": 106.816666,
  "building": "Tower A",
  "zone": "Lobby",
  "ip_address": "192.168.1.100",
  "port": 554,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2143G0-I",
  "resolution": "1920x1080",
  "fps": 25,
  "tags": ["lobby", "entrance", "main"]
}
```

#### Get All Cameras (with pagination)
```http
GET /api/v1/cameras?page=1&page_size=10
Authorization: Bearer <token>
```

#### Get Camera by ID
```http
GET /api/v1/cameras/{id}
Authorization: Bearer <token>
```

#### Update Camera
```http
PUT /api/v1/cameras/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Camera Lobby 1 Updated",
  "status": "ONLINE"
}
```

#### Delete Camera
```http
DELETE /api/v1/cameras/{id}
Authorization: Bearer <token>
```

### Camera Filtering

#### Get Cameras by Zone
```http
GET /api/v1/cameras/zone/filter?zone=Lobby
Authorization: Bearer <token>
```

#### Get Nearby Cameras
```http
GET /api/v1/cameras/nearby?lat=-6.200000&lng=106.816666&radius=5
Authorization: Bearer <token>

# radius dalam kilometer
```

### Stream Management

#### Start Stream
```http
POST /api/v1/cameras/{id}/stream/start
Authorization: Bearer <token>

Response:
{
  "success": true,
  "message": "Stream started successfully",
  "data": {
    "id": "uuid",
    "name": "Camera Lobby 1",
    "status": "ONLINE",
    "hls_url": "http://localhost:8083/stream/{id}/channel/0/hls/live/index.m3u8",
    "webrtc_url": "http://localhost:8083/stream/{id}/channel/0/webrtc",
    "snapshot_url": "http://localhost:8083/stream/{id}/channel/0/jpeg"
  }
}
```

#### Stop Stream
```http
POST /api/v1/cameras/{id}/stream/stop
Authorization: Bearer <token>
```

## 🔧 Development

### Setup Local Development

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run PostgreSQL dengan Docker
docker run -d \
  --name postgres \
  -e POSTGRES_USER=cctv_user \
  -e POSTGRES_PASSWORD=cctv_password_123 \
  -e POSTGRES_DB=cctv_monitoring \
  -p 5432:5432 \
  postgres:15-alpine

# Run RTSPtoWeb
docker run -d \
  --name rtsptoweb \
  -p 8083:8083 \
  ghcr.io/deepch/rtsptoweb:latest

# Run aplikasi
go run cmd/api/main.go
```

### Run Tests

```bash
go test ./...
```

### Build Binary

```bash
go build -o main cmd/api/main.go
```

## 🐳 Docker Commands

```bash
# Build dan jalankan semua services
docker-compose up -d

# Stop semua services
docker-compose down

# View logs
docker-compose logs -f backend

# Rebuild backend setelah code changes
docker-compose up -d --build backend

# Stop dan hapus volumes (HATI-HATI: data akan hilang)
docker-compose down -v
```

## 📊 Database Schema

### Users Table
```sql
- id (UUID, PK)
- username (VARCHAR, UNIQUE)
- email (VARCHAR, UNIQUE)
- password_hash (TEXT)
- role (VARCHAR): admin, operator, viewer
- is_active (BOOLEAN)
- created_at (TIMESTAMPTZ)
- updated_at (TIMESTAMPTZ)
```

### Cameras Table
```sql
- id (UUID, PK)
- name (VARCHAR)
- description (TEXT)
- rtsp_url (TEXT)
- stream_id (VARCHAR, UNIQUE)
- latitude (DOUBLE PRECISION)
- longitude (DOUBLE PRECISION)
- building (VARCHAR)
- zone (VARCHAR)
- ip_address (VARCHAR)
- port (INTEGER)
- manufacturer (VARCHAR)
- model (VARCHAR)
- resolution (VARCHAR)
- fps (INTEGER)
- tags (TEXT[])
- status (VARCHAR): ONLINE, OFFLINE, ERROR, UNKNOWN
- last_seen (TIMESTAMPTZ)
- is_active (BOOLEAN)
- created_by (UUID, FK -> users.id)
- created_at (TIMESTAMPTZ)
- updated_at (TIMESTAMPTZ)
```

### Activity Logs Table
```sql
- id (UUID, PK)
- user_id (UUID, FK -> users.id)
- camera_id (UUID, FK -> cameras.id)
- action (VARCHAR)
- details (JSONB)
- ip_address (VARCHAR)
- user_agent (TEXT)
- created_at (TIMESTAMPTZ)
```

## 🔐 Security

- Password di-hash menggunakan bcrypt
- JWT untuk authentication
- Role-based access control
- CORS protection
- Rate limiting (bisa ditambahkan)
- Input validation

## 🌍 Environment Variables

Lihat `.env.example` untuk daftar lengkap environment variables yang tersedia.

## 📝 TODO / Future Improvements

- [ ] Rate limiting middleware
- [ ] Refresh token mechanism
- [ ] WebSocket untuk real-time notifications
- [ ] Camera health monitoring
- [ ] Video recording management
- [ ] Motion detection alerts
- [ ] Multiple user roles dengan permissions detail
- [ ] API documentation dengan Swagger
- [ ] Unit tests & integration tests
- [ ] Metrics & monitoring (Prometheus)

## 🤝 Contributing

Silakan buat Pull Request atau Issue untuk improvement.

## 📄 License

MIT License

## 👨‍💻 Author

Your Name - CCTV Monitoring System