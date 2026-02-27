# Tutorial: Project Go dengan Gin, GORM, dan Repository Adapter Pattern

## 📚 Pengenalan

Project ini mengimplementasikan aplikasi CRUD sederhana dengan arsitektur clean code yang baik, menggunakan:
- **Gin**: Web framework untuk HTTP routing
- **GORM**: ORM untuk database operations
- **Repository Adapter Pattern**: Untuk mendukung multiple database

## 🏗️ Arsitektur Layer

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Mendefinisikan entity/model yang merepresentasikan struktur data

```go
type Country struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    Code      string    `gorm:"type:varchar(10);uniqueIndex;not null"`
    Name      string    `gorm:"type:varchar(100);not null"`
    Region    string    `gorm:"type:varchar(50)"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

**Key Points**:
- Struct tags `gorm` untuk definisi kolom database
- Struct tags `json` untuk serialisasi JSON
- Entity ini database-agnostic (tidak terikat pada database tertentu)

### 2. DTO Layer (`internal/dto/`)

**Purpose**: Data Transfer Objects untuk request dan response

#### Request DTO:
```go
type CreateCountryRequest struct {
    Code   string `json:"code" binding:"required,min=2,max=10"`
    Name   string `json:"name" binding:"required,min=3,max=100"`
    Region string `json:"region" binding:"max=50"`
}
```

**Key Points**:
- Validasi dengan tag `binding` dari Gin
- Memisahkan struktur request dari domain entity
- Mencegah user mengirim field yang tidak diinginkan (seperti ID, timestamps)

#### Response DTO:
```go
type CountryResponse struct {
    ID        uint      `json:"id"`
    Code      string    `json:"code"`
    Name      string    `json:"name"`
    Region    string    `json:"region"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Key Points**:
- Kontrol eksplisit terhadap data yang di-expose ke client
- Bisa berbeda dari domain entity jika diperlukan

### 3. Repository Layer (`internal/repository/`)

**Purpose**: Data access layer dengan adapter pattern

#### Interface (`country_repository.go`):
```go
type CountryRepository interface {
    Create(country *domain.Country) error
    FindAll() ([]domain.Country, error)
    FindByID(id uint) (*domain.Country, error)
    FindByCode(code string) (*domain.Country, error)
    Update(country *domain.Country) error
    Delete(id uint) error
}
```

**Key Points**:
- Interface mendefinisikan contract untuk semua implementasi
- Memungkinkan dependency injection dan easy testing
- Database-agnostic

#### PostgreSQL Adapter (`repository/postgres/`):
```go
type CountryPostgresRepository struct {
    db *gorm.DB
}

func NewCountryPostgresRepository(db *gorm.DB) *CountryPostgresRepository {
    return &CountryPostgresRepository{db: db}
}

func (r *CountryPostgresRepository) Create(country *domain.Country) error {
    return r.db.Create(country).Error
}
```

#### Oracle Adapter (`repository/oracle/`):
```go
type CountryOracleRepository struct {
    db *gorm.DB
}

func NewCountryOracleRepository(db *gorm.DB) *CountryOracleRepository {
    return &CountryOracleRepository{db: db}
}

func (r *CountryOracleRepository) Create(country *domain.Country) error {
    return r.db.Create(country).Error
}
```

**Key Points**:
- Setiap adapter mengimplementasikan interface yang sama
- Implementasi bisa berbeda untuk setiap database jika diperlukan
- Misalnya: Oracle mungkin perlu handling khusus untuk sequences

### 4. Service Layer (`internal/service/`)

**Purpose**: Business logic layer

```go
type CountryService struct {
    repo repository.CountryRepository
}

func NewCountryService(repo repository.CountryRepository) *CountryService {
    return &CountryService{repo: repo}
}

func (s *CountryService) CreateCountry(req dto.CreateCountryRequest) (*dto.CountryResponse, error) {
    // Validasi bisnis logic
    existing, err := s.repo.FindByCode(req.Code)
    if existing != nil {
        return nil, errors.New("country with this code already exists")
    }
    
    // Convert DTO to domain
    country := &domain.Country{
        Code:   req.Code,
        Name:   req.Name,
        Region: req.Region,
    }
    
    // Simpan ke database
    if err := s.repo.Create(country); err != nil {
        return nil, err
    }
    
    // Convert domain to response DTO
    return s.toResponse(country), nil
}
```

**Key Points**:
- Menerima interface `CountryRepository`, bukan implementasi konkret
- Tidak peduli database apa yang digunakan
- Berisi validasi bisnis logic (misal: unique code check)
- Convert antara DTO dan Domain

### 5. Controller Layer (`internal/controller/`)

**Purpose**: HTTP request handlers

```go
type CountryController struct {
    service *service.CountryService
}

func (c *CountryController) CreateCountry(ctx *gin.Context) {
    var req dto.CreateCountryRequest
    
    // Parse & validate request
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request", err.Error()))
        return
    }
    
    // Call service
    country, err := c.service.CreateCountry(req)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Failed to create", err.Error()))
        return
    }
    
    // Return response
    ctx.JSON(http.StatusCreated, dto.SuccessResponse("Success", country))
}
```

**Key Points**:
- Menangani HTTP-specific concerns (status codes, headers, etc)
- Parse request body dan validate
- Call service layer
- Format response

## 🔄 Repository Adapter Pattern - Penjelasan Detail

### Mengapa Repository Pattern?

1. **Abstraksi**: Memisahkan business logic dari data access
2. **Testability**: Mudah untuk mock repository dalam unit tests
3. **Flexibility**: Mudah switch database tanpa ubah business logic

### Mengapa Adapter Pattern?

Adapter pattern memungkinkan aplikasi bekerja dengan multiple database:

```
                     ┌─────────────────────┐
                     │  Service Layer      │
                     │  (Business Logic)   │
                     └──────────┬──────────┘
                                │
                                │ depends on
                                ↓
                     ┌─────────────────────┐
                     │ CountryRepository   │ ← Interface
                     │    (Interface)      │
                     └──────────┬──────────┘
                                │
                    ┌───────────┴───────────┐
                    │                       │
         ┌──────────▼─────────┐  ┌─────────▼────────┐
         │ PostgresRepository │  │  OracleRepository │
         │    (Adapter)       │  │    (Adapter)      │
         └────────────────────┘  └──────────────────┘
```

### Cara Kerja

1. **Service** hanya tahu interface `CountryRepository`
2. Saat aplikasi start, kita pilih adapter mana yang digunakan:
   ```go
   var countryRepo repository.CountryRepository
   
   switch cfg.Database.Driver {
   case "postgres":
       countryRepo = postgres.NewCountryPostgresRepository(db)
   case "oracle":
       countryRepo = oracle.NewCountryOracleRepository(db)
   }
   ```
3. Service di-inject dengan adapter yang dipilih
4. Service tidak perlu tahu database apa yang digunakan

## 🔧 Configuration & Setup

### Config Layer (`config/`)

```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

func LoadConfig() (*Config, error) {
    godotenv.Load()
    
    return &Config{
        Database: DatabaseConfig{
            Driver: getEnv("DB_DRIVER", "postgres"),
            Postgres: PostgresConfig{...},
            Oracle: OracleConfig{...},
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
    }, nil
}
```

**Key Points**:
- Load dari environment variables atau .env file
- Provide default values
- Type-safe configuration

### Database Package (`pkg/database/`)

```go
func ConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
    switch cfg.Database.Driver {
    case "postgres":
        return connectPostgres(cfg.Database.Postgres)
    case "oracle":
        return connectOracle(cfg.Database.Oracle)
    }
}
```

**Key Points**:
- Factory pattern untuk create database connection
- Auto migration untuk create/update tables
- Database-specific connection handling

## 📝 Main Application (`main.go`)

Dependency Injection flow:

```go
func main() {
    // 1. Load config
    cfg := config.LoadConfig()
    
    // 2. Connect to database
    db := database.ConnectDatabase(cfg)
    
    // 3. Choose repository adapter
    var countryRepo repository.CountryRepository
    if cfg.Database.Driver == "postgres" {
        countryRepo = postgres.NewCountryPostgresRepository(db)
    } else {
        countryRepo = oracle.NewCountryOracleRepository(db)
    }
    
    // 4. Inject repository to service
    countryService := service.NewCountryService(countryRepo)
    
    // 5. Inject service to controller
    countryController := controller.NewCountryController(countryService)
    
    // 6. Setup routes
    router := setupRouter(countryController)
    
    // 7. Start server
    router.Run(":8080")
}
```

## 🎯 Best Practices yang Diimplementasikan

1. **Dependency Injection**: Semua dependencies di-inject, bukan di-create di dalam
2. **Interface Segregation**: Repository interface hanya berisi methods yang diperlukan
3. **Single Responsibility**: Setiap layer punya tanggung jawab yang jelas
4. **Don't Repeat Yourself**: Shared logic di utils atau helper
5. **Configuration Management**: Centralized configuration
6. **Error Handling**: Consistent error handling across layers

## 🚀 Cara Menjalankan

1. Setup database (PostgreSQL atau Oracle)
2. Copy `.env.example` ke `.env`
3. Sesuaikan konfigurasi di `.env`:
   ```env
   DB_DRIVER=postgres  # atau oracle
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   # ...
   ```
4. Run:
   ```bash
   go mod download
   go run main.go
   ```

## 🧪 Testing API

```bash
# Create country
curl -X POST http://localhost:8080/api/countries \
  -H "Content-Type: application/json" \
  -d '{"code":"ID","name":"Indonesia","region":"Asia"}'

# Get all countries
curl http://localhost:8080/api/countries

# Get by ID
curl http://localhost:8080/api/countries/1

# Update
curl -X PUT http://localhost:8080/api/countries/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Republic of Indonesia"}'

# Delete
curl -X DELETE http://localhost:8080/api/countries/1
```

## 📚 Menambah Entity Baru

Untuk menambah entity baru (misal: City), ikuti langkah berikut:

1. **Domain**: Buat `internal/domain/city.go`
2. **DTO**: Buat `internal/dto/city_request.go` dan `city_response.go`
3. **Repository Interface**: Buat `internal/repository/city_repository.go`
4. **Repository Adapters**: 
   - Buat `internal/repository/postgres/city_postgres_repository.go`
   - Buat `internal/repository/oracle/city_oracle_repository.go`
5. **Service**: Buat `internal/service/city_service.go`
6. **Controller**: Buat `internal/controller/city_controller.go`
7. **Routes**: Tambahkan routes di `main.go`
8. **Migration**: Tambahkan entity di `autoMigrate()` di `pkg/database/database.go`

Selamat belajar! 🎉
