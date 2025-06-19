# Backend API

> **REST API pembelajaran Go dengan Fiber Framework dan MongoDB Atlas**

Sebuah project pembelajaran untuk memahami pengembangan REST API menggunakan Go (Golang) dengan Fiber web framework dan MongoDB sebagai database. Project ini dibuat sebagai bagian dari proses belajar dan eksplorasi teknologi backend modern.

## Deskripsi Project

REST API sederhana yang menyediakan sistem autentikasi dan manajemen user. Project ini dibuat untuk mempelajari:

- **Go (Golang)** - Bahasa pemrograman utama
- **Fiber Framework** - Web framework yang cepat dan minimalis
- **MongoDB Atlas** - Database NoSQL cloud
- **JWT Authentication** - Sistem autentikasi berbasis token
- **Docker** - Containerization dan deployment
- **Security Best Practices** - Implementasi keamanan dasar

## Fitur

### Public Endpoints

- **Health Check** - Status kesehatan API
- **User Registration** - Pendaftaran user baru
- **User Login** - Masuk dengan email dan password
- **API Documentation** - Dokumentasi endpoint yang tersedia

### Protected Endpoints (Memerlukan JWT Token)

- **User Profile** - Melihat dan mengubah profil user
- **Change Password** - Mengganti password
- **Account Deactivation** - Menonaktifkan akun
- **Logout & Token Refresh** - Manajemen sesi

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
│   │   └── config.go          # Konfigurasi aplikasi
│   ├── controllers/
│   │   ├── auth.go            # Controller autentikasi
│   │   ├── health.go          # Controller health check
│   │   └── user.go            # Controller user management
│   ├── middleware/
│   │   └── auth.go            # Middleware JWT
│   ├── models/
│   │   └── user.go            # Model dan struct user
│   ├── routes/
│   │   └── routes.go          # Definisi routing
│   ├── services/
│   │   ├── auth.go            # Service autentikasi
│   │   └── user.go            # Service user management
│   └── utils/
│       ├── hash.go            # Utility hashing
│       ├── jwt.go             # Utility JWT
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

### Testing Endpoints

| Method | Endpoint                | Deskripsi                        |
| ------ | ----------------------- | -------------------------------- |
| `GET`  | `/api/v1/testing/users` | Lihat semua user (opsional auth) |


## Security Features

- **Password Hashing** - Menggunakan bcrypt
- **JWT Authentication** - Token-based auth
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
5. **API Design** - RESTful principles, response format
6. **DevOps** - Docker, deployment, monitoring
7. **Best Practices** - Code structure, error handling, validation

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
