# 🚀 Quick Start Guide - CCTV Monitoring Backend

Panduan lengkap untuk menjalankan project dari awal sampai bisa digunakan.

## 📋 Prerequisites

Pastikan kamu sudah install:
- Docker Desktop (https://www.docker.com/products/docker-desktop)
- Git
- Text editor (VS Code recommended)
- Postman atau curl untuk testing API

## 🏁 Step 1: Setup Project

### 1.1 Clone atau Download Project

```bash
# Buat folder project
mkdir cctv-monitoring-backend
cd cctv-monitoring-backend
```

### 1.2 Buat Struktur Folder

Buat folder sesuai struktur berikut:

```
cctv-monitoring-backend/
├── cmd/
│   └── api/
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   ├── middleware/
│   └── utils/
├── migrations/
```

### 1.3 Copy Semua File

Copy semua file Go yang sudah dibuat ke folder yang sesuai:
- `cmd/api/main.go`
- `internal/config/config.go`
- `internal/database/postgres.go`
- dll... (semua file yang sudah dibuat)

### 1.4 Copy File Konfigurasi

Copy file-file berikut ke root project:
- `.env.example`
- `.gitignore`
- `go.mod`
- `Dockerfile`
- `docker-compose.yml`
- `README.md`

## 🔧 Step 2: Setup Environment

### 2.1 Buat File .env

```bash
cp .env.example .env
```

### 2.2 Edit .env (Optional)

File `.env` sudah siap digunakan dengan default values. Kamu bisa edit jika perlu:

```bash
# Buka dengan text editor
nano .env   # atau
code .env   # jika pakai VS Code
```

**PENTING**: Untuk production, ganti `JWT_SECRET` dengan nilai yang lebih secure!

## 🐳 Step 3: Jalankan dengan Docker

### 3.1 Start Services

```bash
# Pastikan Docker Desktop sudah running
# Kemudian jalankan:
docker-compose up -d
```

Command ini akan:
- ✅ Download images yang diperlukan
- ✅ Build Go application
- ✅ Start PostgreSQL database
- ✅ Start RTSPtoWeb service
- ✅ Start Backend API

### 3.2 Check Status

```bash
# Lihat status containers
docker-compose ps

# Expected output:
# NAME              STATUS        PORTS
# cctv_backend      Up           0.0.0.0:8080->8080/tcp
# cctv_postgres     Up (healthy) 0.0.0.0:5432->5432/tcp
# cctv_rtsptoweb    Up           0.0.0.0:8083->8083/tcp
```

### 3.3 View Logs

```bash
# Lihat logs backend
docker-compose logs -f backend

# Kamu harus lihat:
# ✓ Successfully connected to PostgreSQL database
# ✓ Database migrations completed successfully
# ✓ Server is running on http://localhost:8080
```

## 🧪 Step 4: Test API

### 4.1 Health Check

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "message": "CCTV Monitoring API is running"
}
```

### 4.2 Login dengan Default Admin

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid-here",
      "username": "admin",
      "email": "admin@cctv-monitoring.com",
      "role": "admin",
      "is_active": true
    }
  }
}
```

**PENTING**: Copy token dari response! Kamu akan butuh ini untuk request berikutnya.

### 4.3 Set Token sebagai Variable

```bash
# Linux/Mac
export TOKEN="paste-token-dari-login-response"

# Windows CMD
set TOKEN=paste-token-dari-login-response

# Windows PowerShell
$env:TOKEN="paste-token-dari-login-response"
```

### 4.4 Create Camera Pertama

```bash
# Linux/Mac
curl -X POST http://localhost:8080/api/v1/cameras \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Camera Test 1",
    "description": "Camera untuk testing",
    "rtsp_url": "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mp4",
    "latitude": -6.200000,
    "longitude": 106.816666,
    "building": "Building A",
    "zone": "Lobby",
    "tags": ["test", "lobby"]
  }'

# Windows CMD (gunakan file JSON)
# Buat file camera.json dengan isi di atas, lalu:
curl -X POST http://localhost:8080/api/v1/cameras ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d @camera.json
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Camera created successfully",
  "data": {
    "id": "camera-uuid",
    "name": "Camera Test 1",
    "rtsp_url": "rtsp://...",
    "latitude": -6.200000,
    "longitude": 106.816666,
    "status": "UNKNOWN",
    ...
  }
}
```

Copy `id` dari response!

### 4.5 Get All Cameras

```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/cameras
```

### 4.6 Start Stream

```bash
# Ganti {camera-id} dengan ID dari step 4.4
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/cameras/{camera-id}/stream/start
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Stream started successfully",
  "data": {
    "id": "camera-uuid",
    "status": "ONLINE",
    "hls_url": "http://localhost:8083/stream/camera-uuid/channel/0/hls/live/index.m3u8",
    "webrtc_url": "http://localhost:8083/stream/camera-uuid/channel/0/webrtc",
    "snapshot_url": "http://localhost:8083/stream/camera-uuid/channel/0/jpeg"
  }
}
```

### 4.7 Test Stream di Browser

Buka browser dan akses:
```
http://localhost:8083/stream/{camera-uuid}/channel/0/hls/live/index.m3u8
```

Atau untuk snapshot:
```
http://localhost:8083/stream/{camera-uuid}/channel/0/jpeg
```

## 📱 Step 5: Test dengan Postman

### 5.1 Import Collection

Buat Postman Collection baru dengan endpoints berikut:

1. **Login**
   - Method: POST
   - URL: `http://localhost:8080/api/v1/auth/login`
   - Body (JSON):
   ```json
   {
     "username": "admin",
     "password": "admin123"
   }
   ```

2. **Get All Cameras**
   - Method: GET
   - URL: `http://localhost:8080/api/v1/cameras`
   - Headers:
     - `Authorization`: `Bearer {token}`

3. **Create Camera**
   - Method: POST
   - URL: `http://localhost:8080/api/v1/cameras`
   - Headers:
     - `Authorization`: `Bearer {token}`
     - `Content-Type`: `application/json`
   - Body: (sama seperti step 4.4)

4. **Start Stream**
   - Method: POST
   - URL: `http://localhost:8080/api/v1/cameras/{id}/stream/start`
   - Headers:
     - `Authorization`: `Bearer {token}`

## 🔍 Step 6: Verifikasi Database

### 6.1 Connect ke PostgreSQL

```bash
# Masuk ke container PostgreSQL
docker exec -it cctv_postgres psql -U cctv_user -d cctv_monitoring
```

### 6.2 Query Database

```sql
-- Lihat semua users
SELECT * FROM users;

-- Lihat semua cameras
SELECT id, name, status, latitude, longitude FROM cameras;

-- Exit
\q
```

## 🛑 Step 7: Stop Services

```bash
# Stop semua services
docker-compose down

# Stop dan hapus volumes (HATI-HATI: data akan hilang!)
docker-compose down -v

# Start lagi
docker-compose up -d
```

## 🐛 Troubleshooting

### Problem 1: Port Sudah Digunakan

**Error:**
```
Error: bind: address already in use
```

**Solution:**
```bash
# Check port yang digunakan
# Linux/Mac:
lsof -i :8080
lsof -i :5432
lsof -i :8083

# Windows:
netstat -ano | findstr :8080

# Kill process atau ganti port di docker-compose.yml
```

### Problem 2: Database Connection Failed

**Error:**
```
Failed to connect to database
```

**Solution:**
```bash
# Restart PostgreSQL container
docker-compose restart postgres

# Check logs
docker-compose logs postgres

# Tunggu sampai PostgreSQL ready (lihat "database system is ready to accept connections")
```

### Problem 3: Go Dependencies Error

**Error:**
```
go: module not found
```

**Solution:**
```bash
# Pastikan go.mod dan go.sum ada
# Rebuild container
docker-compose build --no-cache backend
docker-compose up -d
```

### Problem 4: RTSPtoWeb Tidak Bisa Start Stream

**Error:**
```
Failed to start stream
```

**Solution:**
```bash
# Check RTSPtoWeb logs
docker-compose logs rtsptoweb

# Restart RTSPtoWeb
docker-compose restart rtsptoweb

# Pastikan RTSP URL valid
```

### Problem 5: JWT Token Invalid

**Error:**
```
Invalid or expired token
```

**Solution:**
- Token expired (default 24 jam), login ulang untuk dapat token baru
- Pastikan token disertakan dengan format: `Bearer {token}`
- Check JWT_SECRET di .env sama dengan yang digunakan saat generate token

## 📊 Monitoring & Logs

### View All Logs
```bash
docker-compose logs -f
```

### View Specific Service Logs
```bash
# Backend only
docker-compose logs -f backend

# PostgreSQL only
docker-compose logs -f postgres

# RTSPtoWeb only
docker-compose logs -f rtsptoweb
```

### Check Container Resource Usage
```bash
docker stats
```

## 🔄 Update Code

Jika kamu ubah kode Go:

```bash
# Rebuild dan restart backend
docker-compose up -d --build backend

# Atau rebuild semua
docker-compose up -d --build
```

## 💾 Backup & Restore Database

### Backup Database
```bash
# Backup database
docker exec cctv_postgres pg_dump -U cctv_user cctv_monitoring > backup.sql
```

### Restore Database
```bash
# Restore database
docker exec -i cctv_postgres psql -U cctv_user cctv_monitoring < backup.sql
```

## 🎯 Next Steps

Setelah berhasil menjalankan backend:

1. **Buat Frontend** untuk menampilkan peta dan camera streams
2. **Tambah Real RTSP Camera** - ganti RTSP URL dengan camera real
3. **Implement WebSocket** untuk real-time notifications
4. **Add More Features**:
   - Recording management
   - Motion detection alerts
   - Camera health monitoring
   - Video playback

## 📖 Useful Commands Cheat Sheet

```bash
# Docker Compose
docker-compose up -d              # Start services
docker-compose down               # Stop services
docker-compose ps                 # List services
docker-compose logs -f [service]  # View logs
docker-compose restart [service]  # Restart service
docker-compose build --no-cache   # Rebuild images

# Database
docker exec -it cctv_postgres psql -U cctv_user -d cctv_monitoring

# Backend API Testing
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Get Cameras (dengan token)
curl -H "Authorization: Bearer {TOKEN}" \
  http://localhost:8080/api/v1/cameras

# Create Camera
curl -X POST http://localhost:8080/api/v1/cameras \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{...}'

# Start Stream
curl -X POST http://localhost:8080/api/v1/cameras/{ID}/stream/start \
  -H "Authorization: Bearer {TOKEN}"
```

## 🎓 Tips untuk Pemula

1. **Pahami Flow Request**:
   ```
   Client → API Handler → Service Layer → Repository → Database
   ```

2. **Selalu Check Logs**:
   - Jika ada error, check logs backend: `docker-compose logs -f backend`
   - Lihat error message untuk debugging

3. **Test Step by Step**:
   - Jangan langsung test semua
   - Test satu endpoint dulu sampai berhasil
   - Baru lanjut ke endpoint berikutnya

4. **Gunakan Postman**:
   - Lebih mudah daripada curl
   - Bisa save collections
   - Bisa save environment variables

5. **Baca Error Message**:
   - Error message biasanya memberitahu apa yang salah
   - Google error message jika tidak paham

## 🆘 Butuh Bantuan?

Jika masih ada masalah:

1. Check logs: `docker-compose logs -f`
2. Pastikan semua services running: `docker-compose ps`
3. Restart services: `docker-compose restart`
4. Rebuild dari awal: `docker-compose down -v && docker-compose up -d --build`

## ✅ Checklist Setup

- [ ] Docker Desktop installed dan running
- [ ] Project structure sudah dibuat
- [ ] Semua file Go sudah di-copy
- [ ] File .env sudah dibuat
- [ ] `docker-compose up -d` berhasil
- [ ] Health check API berhasil (http://localhost:8080/health)
- [ ] Login berhasil dan dapat token
- [ ] Bisa create camera
- [ ] Bisa start stream
- [ ] Stream bisa diakses di browser

Jika semua checklist sudah ✅, congratulations! Backend kamu sudah running dengan baik! 🎉

## 📚 Learning Resources

Untuk belajar lebih lanjut:
- **Golang**: https://go.dev/tour/
- **Fiber Framework**: https://docs.gofiber.io/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **Docker**: https://docs.docker.com/
- **RTSPtoWeb**: https://github.com/deepch/RTSPtoWeb

---

**Happy Coding! 🚀**