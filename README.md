# Backend API

> **REST API untuk aplikasi E-Presensi siswa SMK dengan Go, Fiber Framework dan MongoDB Atlas**

Sebuah REST API untuk sistem presensi siswa SMK yang menggunakan teknologi GPS untuk validasi lokasi. Project ini mengimplementasikan sistem absensi modern dengan fitur keamanan mobile dan validasi lokasi berbasis GPS.

## Deskripsi Project

REST API untuk sistem E-Presensi siswa SMK yang menyediakan:

- **Sistem Autentikasi Siswa** - Registrasi dan login dengan data siswa (NIS, Kelas, Jurusan)
- **Absensi Berbasis GPS** - Check-in/out dengan validasi lokasi
- **Validasi Lokasi Real-time** - Hanya bisa absen dalam radius sekolah
- **Keamanan Mobile** - Security untuk aplikasi mobile
- **Riwayat Kehadiran** - History dan statistik kehadiran siswa

Teknologi yang dipelajari:

- **Go (Golang)** - Bahasa pemrograman utama
- **Fiber Framework** - Web framework yang cepat dan minimalis
- **MongoDB Atlas** - Database NoSQL cloud
- **JWT Authentication** - Sistem autentikasi berbasis token
- **GPS Navigation** - Validasi lokasi dengan Haversine formula
- **Mobile Security** - Implementasi keamanan perangkat mobile
- **Docker** - Containerization dan deployment

## Fitur

### Public Endpoints

- **Health Check** - Status kesehatan API
- **User Registration** - Pendaftaran siswa baru (NIS, Kelas, Jurusan)
- **User Login** - Masuk dengan email dan password
- **API Documentation** - Dokumentasi endpoint yang tersedia

### Protected Endpoints (Memerlukan JWT Token)

- **User Profile** - Melihat dan mengubah profil siswa
- **Change Password** - Mengganti password
- **Account Deactivation** - Menonaktifkan akun
- **Logout & Token Refresh** - Manajemen sesi

### Attendance Endpoints (GPS Required)

- **Check-in** - Absen masuk dengan validasi GPS
- **Check-out** - Absen keluar dengan validasi GPS
- **Today's Attendance** - Lihat absensi hari ini
- **Attendance History** - Riwayat kehadiran dengan pagination
- **Attendance Stats** - Statistik kehadiran (hadir, terlambat, tidak hadir)

### Security Features

- **GPS Location Validation** - Validasi lokasi dalam radius sekolah
- **Mobile Device Security** - Keamanan untuk aplikasi mobile
- **Haversine Distance Calculation** - Perhitungan jarak GPS yang akurat
- **Indonesia Territory Validation** - Validasi lokasi dalam wilayah Indonesia

### Testing Endpoints

- **Get All Users** - Endpoint untuk testing (opsional auth)

## Tech Stack

| Teknologi     | Versi    | Kegunaan         |
| ------------- | -------- | ---------------- |
| **Go**        | 1.23.3   | Backend Language |
| **Fiber**     | v2.52.8  | Web Framework    |
| **MongoDB**   | Atlas    | Database         |
| **JWT**       | v5.2.2   | Authentication   |
| **Docker**    | Latest   | Containerization |
| **Validator** | v10.26.0 | Input Validation |
| **Bcrypt**    | Latest   | Password Hashing |

## 📁 Struktur Project

```
go-fiber-auth-api/
├── cmd/
│   └── main.go                 # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go          # Konfigurasi aplikasi & GPS
│   ├── controllers/
│   │   ├── auth.go            # Controller autentikasi
│   │   ├── attendance.go      # Controller absensi GPS
│   │   ├── health.go          # Controller health check
│   │   └── user.go            # Controller user management
│   ├── middleware/
│   │   ├── auth.go            # Middleware JWT
│   │   └── location.go        # Middleware validasi GPS
│   ├── models/
│   │   ├── attendance.go      # Model absensi dan lokasi
│   │   └── user.go            # Model dan struct user
│   ├── routes/
│   │   └── routes.go          # Definisi routing
│   ├── services/
│   │   ├── attendance.go      # Service absensi GPS
│   │   ├── auth.go            # Service autentikasi
│   │   └── user.go            # Service user management
│   └── utils/
│       ├── hash.go            # Utility hashing
│       ├── jwt.go             # Utility JWT
│       ├── mobile.go          # Utility GPS & mobile
│       ├── response.go        # Utility response
│       ├── sanitize.go        # Utility sanitization
│       └── validation.go      # Utility validation
├── pkg/
│   └── database/
│       └── mongodb.go         # Koneksi MongoDB
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose setup
├── .env                       # Environment variables
├── go.mod                     # Go modules
└── go.sum                     # Go dependencies
```

## Quick Start

### 1. Clone Project

```bash
git clone https://github.com/SideeID/go-fiber-auth-api
cd go-fiber-auth-api
```

### 2. Setup Environment

```bash
# Copy dan edit environment variables
cp .env.example .env
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Development

```bash
go run cmd/main.go
```

API akan berjalan di `http://localhost:8080`

## Docker Deployment

### Quick Deploy

```bash
# Build dan jalankan dengan Docker Compose
docker-compose up -d
```

### Manual Docker

```bash
# Build image
docker build -t ujikom-api .

# Run container
docker run -d --name ujikom-backend -p 3006:8080 ujikom-api
```

## API Endpoints

### Base URL

- **Development**: `http://localhost:8080`

### Public Endpoints

| Method | Endpoint                | Deskripsi            |
| ------ | ----------------------- | -------------------- |
| `GET`  | `/api/v1/health`        | Health check API     |
| `GET`  | `/api/v1/docs`          | Dokumentasi API      |
| `GET`  | `/api/v1/auth/test`     | Test endpoint auth   |
| `POST` | `/api/v1/auth/register` | Registrasi user baru |
| `POST` | `/api/v1/auth/login`    | Login user           |

### Protected Endpoints

| Method | Endpoint                       | Deskripsi          |
| ------ | ------------------------------ | ------------------ |
| `GET`  | `/api/v1/user/profile`         | Lihat profil user  |
| `PUT`  | `/api/v1/user/profile`         | Update profil user |
| `POST` | `/api/v1/user/change-password` | Ganti password     |
| `POST` | `/api/v1/user/deactivate`      | Nonaktifkan akun   |
| `POST` | `/api/v1/user/logout`          | Logout user        |
| `POST` | `/api/v1/user/refresh-token`   | Refresh JWT token  |

### Attendance Endpoints (GPS Required)

| Method | Endpoint                      | Deskripsi              |
| ------ | ----------------------------- | ---------------------- |
| `POST` | `/api/v1/attendance/checkin`  | Check-in dengan GPS    |
| `POST` | `/api/v1/attendance/checkout` | Check-out dengan GPS   |
| `GET`  | `/api/v1/attendance/today`    | Lihat absensi hari ini |
| `GET`  | `/api/v1/attendance/history`  | Riwayat kehadiran      |
| `GET`  | `/api/v1/attendance/stats`    | Statistik kehadiran    |

### Testing Endpoints

| Method | Endpoint                | Deskripsi                        |
| ------ | ----------------------- | -------------------------------- |
| `GET`  | `/api/v1/testing/users` | Lihat semua user (opsional auth) |

## Security Features

- **Password Hashing** - Menggunakan bcrypt
- **JWT Authentication** - Token-based auth
- **GPS Location Validation** - Validasi lokasi dalam radius sekolah
- **Mobile Device Security** - Keamanan untuk aplikasi mobile
- **Haversine Distance Calculation** - Perhitungan jarak GPS yang akurat
- **Indonesia Territory Validation** - Validasi lokasi dalam wilayah Indonesia
- **Input Validation** - Validasi input user
- **CORS Enabled** - Cross-origin resource sharing
- **Rate Limiting** - Pembatasan request
- **SQL Injection Protection** - MongoDB native protection
- **Environment Variables** - Sensitive data protection

## 🤝 Contributing

Ini adalah project pembelajaran, jadi kontribusi dan saran sangat diterima!

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## Pembelajaran

Project ini dibuat untuk mempelajari:

1. **Go Fundamentals** - Sintaks dasar, goroutines, channels
2. **Web Development** - HTTP handling, middleware, routing
3. **Database Integration** - MongoDB operations, ODM
4. **Authentication** - JWT, password hashing, security
5. **GPS & Location Services** - Geolocation, distance calculation, navigation
6. **Mobile API Development** - Mobile-first API design, device security
7. **API Design** - RESTful principles, response format
8. **DevOps** - Docker, deployment, monitoring
9. **Best Practices** - Code structure, error handling, validation

## License

Project ini dibuat untuk keperluan pembelajaran dan bersifat open source.

## Author

**SideeID**

- GitHub: [https://github.com/SideeID]
- Website: [https://side.my.id]

---

<p align="center">
  <i>Dibuat dengan ❤️ untuk belajar Go dan teknologi backend modern</i>
</p>

<p align="center">
  <strong>Happy Coding! 🚀</strong>
</p>


Security Infrastructure-based Wireless Network
Cellular Network Security:
	Validasi koneksi melalui jaringan hp
	Deteksi jenis koneksi (4G/5G)
	Enkripsi data saat transmisi via hp
WLAN Security:
	Validasi WiFi sekolah
	Deteksi SSID sekolah yang authorized
	WPA/WPA2 encryption support
	Blocking akses dari WiFi publik/tidak dikenal
	Virtual Private Networks (VPN):
	Deteksi penggunaan VPN
	Blocking atau warning jika ada VPN aktif
	Secure tunnel untuk data transmission
Mobile IP:
	Tracking IP address perangkat
	Validasi IP range sekolah
	Geofencing berdasarkan IP location

Mobile Sensors
	GPS - untuk location tracking (paling cuma ini yang kepakek)
	Accelerometer - deteksi gerakan/shake
	Gyroscope - orientasi perangkat
	Fingerprint/Face ID - biometric authentication (sama ini kalo mau)
	Proximity sensor - deteksi jarak
	Ambient light sensor - deteksi kondisi pencahayaan

Navigasi GPS
	Real-time location tracking
	Geofencing untuk area sekolah
	Distance calculation dari lokasi sekolah
	Maps integration untuk menampilkan lokasi
	Route tracking (perjalanan ke sekolah)

Keamanan Perangkat
Device Security:
	Device fingerprinting (IMEI, device ID)
	Root/Jailbreak detection
	Screen recording/screenshot prevention
	App tampering detection
	Certificate pinning untuk API calls
Authentication Security:
	Multi-factor authentication
	Biometric authentication
	Token-based authentication (JWT)
	Session management
	Auto-logout setelah idle