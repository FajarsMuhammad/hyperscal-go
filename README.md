# HyperScal Go - CRUD Application with Multiple Database Support

Aplikasi CRUD sederhana menggunakan Gin dan GORM dengan dukungan multiple database (PostgreSQL dan Oracle).

## Arsitektur

Project ini mengikuti clean architecture dengan layer:
- **Controller**: Menangani HTTP requests dan responses
- **Service**: Business logic
- **Repository**: Data access layer dengan adapter pattern
- **Domain**: Entity/Model
- **DTO**: Data Transfer Objects (Request/Response)

## Repository Adapter Pattern

Repository menggunakan adapter pattern dimana:
- `CountryRepository` adalah interface yang mendefinisikan contract
- `CountryPostgresRepository` adalah implementasi untuk PostgreSQL
- `CountryOracleRepository` adalah implementasi untuk Oracle

Database bisa di-switch melalui konfigurasi di file `.env` dengan mengubah `DB_DRIVER`.

## Struktur Folder

```
.
├── config/              # Konfigurasi aplikasi dan database
├── internal/
│   ├── controller/      # HTTP handlers
│   ├── service/         # Business logic layer
│   ├── repository/      # Data access layer
│   │   ├── interface/   # Repository interfaces
│   │   ├── postgres/    # PostgreSQL implementations
│   │   └── oracle/      # Oracle implementations
│   ├── domain/          # Entity models
│   └── dto/             # Request & Response DTOs
├── pkg/                 # Shared packages
│   └── database/        # Database connection utilities
├── .env                 # Environment variables
├── go.mod
└── main.go
```

## Setup

1. Copy `.env.example` ke `.env`:
```bash
cp .env.example .env
```

2. Sesuaikan konfigurasi database di `.env`

3. Install dependencies:
```bash
go mod download
```

4. Jalankan aplikasi:
```bash
go run main.go
```

## API Endpoints

### Countries

- `GET /api/countries` - Get all countries
- `GET /api/countries/:id` - Get country by ID
- `POST /api/countries` - Create new country
- `PUT /api/countries/:id` - Update country
- `DELETE /api/countries/:id` - Delete country

### Example Request

**Create Country:**
```bash
curl -X POST http://localhost:8080/api/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ID",
    "name": "Indonesia",
    "region": "Asia"
  }'
```

**Get All Countries:**
```bash
curl http://localhost:8080/api/countries
```

## Database Support

Aplikasi ini mendukung 2 database:
- **PostgreSQL**: Set `DB_DRIVER=postgres` di `.env`
- **Oracle**: Set `DB_DRIVER=oracle` di `.env`

Implementasi repository akan otomatis menggunakan adapter yang sesuai berdasarkan konfigurasi.
