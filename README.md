# Backend Dashboard - Data Coloris

Backend API untuk dashboard manajemen data Coloris menggunakan Golang dan MongoDB.

## Fitur

- CRUD (Create, Read, Update, Delete) data Coloris
- Import data dari file Excel (.xlsx)
- Export data ke file Excel
- Filter data berdasarkan region, cabang, dan bulan
- Pagination untuk list data
- RESTful API

## Teknologi

- **Language**: Go 1.25+
- **Database**: MongoDB
- **Framework**: Gin (HTTP Web Framework)
- **Excel Processing**: Excelize
- **MongoDB Driver**: Official MongoDB Go Driver

## Struktur Project

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
├── config/
│   ├── config.go                # Load konfigurasi dari environment
│   └── database.go              # MongoDB connection setup
├── internal/
│   ├── handlers/
│   │   ├── coloris_handler.go  # HTTP handlers
│   │   └── router.go            # Route setup
│   ├── models/
│   │   └── coloris.go           # Data models dan DTOs
│   ├── repository/
│   │   └── coloris_repository.go # Database operations
│   └── service/
│       └── coloris_service.go   # Business logic
├── pkg/
│   └── utils/
│       ├── excel.go             # Excel processing utilities
│       └── time.go              # Time parsing utilities
├── uploads/                      # Directory untuk temporary file uploads
├── .env.example                 # Example environment variables
├── .gitignore                   # Git ignore file
├── go.mod                       # Go module dependencies
└── README.md                    # Dokumentasi

```

## Setup dan Instalasi

### Prerequisites

- Go 1.25 atau lebih tinggi
- MongoDB 4.0 atau lebih tinggi

### Langkah-langkah

1. Clone repository

2. Install dependencies
```bash
go mod download
```

3. Copy `.env.example` ke `.env` dan sesuaikan konfigurasi
```bash
cp .env.example .env
```

4. Edit file `.env` sesuai environment Anda
```env
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=dashboard_db
SERVER_PORT=8080
ALLOWED_ORIGINS=*
```

5. Jalankan aplikasi
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /api/v1/health
```

### Data Coloris

#### 1. Create Data (Input Manual)
```
POST /api/v1/coloris
Content-Type: application/json

{
  "timestamp": "1/17/2025 14:47:50",
  "bulan": "17 Januari 2025",
  "region": "Jakarta",
  "cabang": "Cabang A",
  "materi": "Training",
  "nama_atasan_langsung": "John Doe",
  "nama_toko": "Toko ABC",
  "nama_lengkap_sesuai_ktp": "Jane Smith",
  "nilai_pg": 85.5,
  "nilai_akhir": 90.0,
  "total": 175.5
}
```

#### 2. Get All Data (dengan pagination dan filter)
```
GET /api/v1/coloris?page=1&per_page=10
GET /api/v1/coloris?region=Jakarta&cabang=Cabang A&bulan=17 Januari 2025
```

#### 3. Get Data by ID
```
GET /api/v1/coloris/:id
```

#### 4. Update Data
```
PUT /api/v1/coloris/:id
Content-Type: application/json

{
  "timestamp": "1/17/2025 14:47:50",
  "bulan": "17 Januari 2025",
  "region": "Jakarta",
  "cabang": "Cabang A",
  "materi": "Training Updated",
  "nama_atasan_langsung": "John Doe",
  "nama_toko": "Toko ABC",
  "nama_lengkap_sesuai_ktp": "Jane Smith",
  "nilai_pg": 85.5,
  "nilai_akhir": 90.0,
  "total": 175.5
}
```

#### 5. Delete Data
```
DELETE /api/v1/coloris/:id
```

#### 6. Import dari Excel
```
POST /api/v1/coloris/import
Content-Type: multipart/form-data

Form data:
- file: [Excel file .xlsx]
```

Format Excel harus memiliki kolom:
- Timestamp (1/17/2025 14:47:50)
- Bulan (17 Januari 2025)
- Region
- Cabang
- Materi
- Nama Atasan Langsung
- Nama Toko
- Nama Lengkap Sesuai KTP
- Nilai PG
- Nilai Akhir
- Total

#### 7. Export ke Excel
```
GET /api/v1/coloris/export
GET /api/v1/coloris/export?filename=data_2025
```

## Response Format

### Success Response
```json
{
  "message": "Success message"
}
```

### Data List Response
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "per_page": 10,
  "total_pages": 10
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

## Development

### Build
```bash
go build -o app cmd/api/main.go
```

### Run
```bash
./app
```

### Testing
```bash
go test ./...
```

## Catatan

- Pastikan MongoDB sudah berjalan sebelum start aplikasi
- Default port adalah 8080, dapat diubah di file .env
- File Excel yang diimport harus memiliki header di baris pertama
- Format timestamp yang didukung: `1/2/2006 15:04:05`, `2006-01-02`, dll

## License

Private Project
