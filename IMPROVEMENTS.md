# CCTV Freeze Prevention Improvements

Dokumen ini menjelaskan perbaikan yang telah dilakukan untuk mengatasi masalah freeze CCTV yang sering terjadi.

## Masalah yang Ditemukan

1. **HTTP Client tanpa timeout** - Request bisa hang indefinitely
2. **Tidak ada connection pooling** - Setiap request membuat koneksi baru
3. **Health checks tanpa rate limiting** - Semua camera dicek bersamaan, bisa overwhelm RTSP service
4. **Deteksi frozen stream kurang akurat** - Hanya mengandalkan LastSeen timestamp
5. **Tidak ada retry mechanism dengan backoff** - Restart stream langsung tanpa delay
6. **Error handling kurang baik** - Sulit untuk debugging

## Perbaikan yang Dilakukan

### 1. HTTP Client dengan Timeout dan Connection Pooling ✅

**File:** `internal/service/rtsp_service.go`

- Menambahkan HTTP client dengan timeout 10 detik untuk semua request
- Menggunakan connection pooling dengan:
  - MaxIdleConns: 100
  - MaxIdleConnsPerHost: 10
  - IdleConnTimeout: 90 detik
- Menggunakan context dengan timeout untuk setiap request
- Menambahkan error handling yang lebih baik dengan membaca response body

**Manfaat:**
- Mencegah request yang hang
- Meningkatkan performa dengan connection reuse
- Lebih cepat mendeteksi masalah koneksi

### 2. Rate Limiting untuk Health Checks ✅

**File:** `internal/service/camera_health_monitor.go`

- Menambahkan semaphore untuk membatasi concurrent health checks (max 5 concurrent)
- Menggunakan WaitGroup untuk menunggu semua checks selesai
- Menambahkan timeout untuk health check cycle

**Manfaat:**
- Mencegah overwhelm RTSP service dengan terlalu banyak request bersamaan
- Health checks lebih stabil dan predictable
- Mengurangi beban pada sistem

### 3. Improved Frozen Stream Detection ✅

**File:** `internal/service/camera_health_monitor.go` dan `internal/service/rtsp_service.go`

- Menambahkan method `GetSnapshotHash()` untuk mendapatkan MD5 hash dari snapshot
- Membandingkan snapshot hash untuk mendeteksi frozen stream
- Jika hash sama selama lebih dari 45 detik, stream dianggap frozen
- Fallback ke time-based detection jika snapshot tidak bisa diambil

**Manfaat:**
- Deteksi frozen stream lebih akurat
- Bisa mendeteksi stream yang stuck pada frame yang sama
- Lebih cepat merespons masalah frozen stream

### 4. Retry Mechanism dengan Exponential Backoff ✅

**File:** `internal/service/camera_health_monitor.go`

- Menambahkan tracking restart attempts per camera
- Menggunakan exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, ... (max 5 menit)
- Mencegah restart yang terlalu sering
- Reset counter setelah restart berhasil

**Manfaat:**
- Mengurangi beban pada RTSP service
- Memberi waktu untuk stream recover sendiri
- Mencegah restart loop yang tidak perlu

### 5. Better Error Handling dan Logging ✅

**File:** `internal/service/camera_health_monitor.go` dan `internal/service/rtsp_service.go`

- Menambahkan logging yang lebih detail untuk debugging
- Log ketika snapshot hash berubah (stream updating)
- Log error dengan context yang jelas
- Log restart attempts dengan informasi backoff

**Manfaat:**
- Lebih mudah debugging masalah
- Bisa track status stream secara real-time
- Monitoring lebih baik

## Konfigurasi

### Health Check Interval
Default: 30 detik (bisa diubah di `cmd/api/main.go`)

```go
healthMonitor := service.NewCameraHealthMonitor(
    cameraRepo,
    rtspService,
    wsHub,
    30*time.Second, // Check every 30 seconds
)
```

### Concurrent Health Checks
Default: 5 concurrent checks (bisa diubah di `camera_health_monitor.go`)

### Frozen Detection Threshold
- Snapshot hash comparison: 45 detik
- Time-based fallback: 60 detik

### Exponential Backoff
- Start: 1 detik
- Max: 300 detik (5 menit)
- Formula: `2^(attempt-1)` detik

## Monitoring

Sistem sekarang akan log:
- ✅ Stream yang berhasil update (snapshot hash berubah)
- ⚠️ Error saat check status atau snapshot
- 🧊 Stream yang terdeteksi frozen
- 🔄 Attempt restart dengan attempt number dan backoff time
- ✓ Stream yang berhasil di-restart

## Testing

Untuk test perbaikan:

1. **Test HTTP timeout:**
   - Matikan RTSP service sementara
   - Pastikan request timeout dalam 10 detik

2. **Test rate limiting:**
   - Monitor log saat health check berjalan
   - Pastikan tidak lebih dari 5 concurrent checks

3. **Test frozen detection:**
   - Simulasi frozen stream (stop RTSP source)
   - Pastikan terdeteksi dalam 45-60 detik

4. **Test exponential backoff:**
   - Simulasi camera yang terus gagal
   - Pastikan restart attempts menggunakan backoff yang benar

## Rekomendasi Tambahan

1. **Monitoring Dashboard:**
   - Tambahkan metrics untuk tracking:
     - Jumlah frozen streams
     - Restart attempts per camera
     - Average health check time

2. **Alerting:**
   - Alert jika camera frozen lebih dari X menit
   - Alert jika restart attempts terlalu banyak

3. **Configuration:**
   - Buat threshold bisa dikonfigurasi via environment variable
   - Allow tuning untuk production environment

4. **Performance:**
   - Consider caching snapshot hashes untuk mengurangi load
   - Batch health checks jika jumlah camera sangat banyak

## Catatan Penting

- Perbaikan ini backward compatible, tidak perlu perubahan database
- Semua perbaikan sudah di-test compile
- Pastikan RTSP service memiliki resource yang cukup untuk handle concurrent requests
- Monitor log untuk melihat efektivitas perbaikan

