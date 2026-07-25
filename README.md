# SIAKAD Microservices

## Arsitektur

```
                    ┌──────────────────┐
                    │   docker-compose │
                    └────────┬─────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────┐ ┌───────────────┐ ┌─────────────────┐
│  akademik-db    │ │  akademik-    │ │   rekap-        │
│  (PostgreSQL)   │ │  service      │ │   service       │
│                 │ │  :8081 (ext)  │ │   :8082 (ext)   │
│                 │ │  :8080 (int)  │ │   :8080 (int)   │
└────────┬────────┘ └───────┬───────┘ └────────┬────────┘
         │                  │                   │
         └──────────────────┘                   │
              SQL (pgx)                         │
                                                │
                              HTTP REST ────────┘
                              (AkademikClient)
```

## Pembagian Domain

- **Akademik Service** (Pemilik Data): Mengelola entitas mahasiswa dan nilai di database PostgreSQL. Menyediakan endpoint CRUD dan endpoint internal ringkasan untuk Rekap Service.
- **Rekap Service** (Konsumen Data): Tidak memiliki database. Mengambil data dari Akademik Service melalui HTTP REST untuk menghasilkan laporan rekap per jurusan dan top IPK.

## Cara Menjalankan

### Prasyarat
- Docker & Docker Compose
- Go 1.22+ (untuk development lokal)

### Development Mode (Docker)

```bash
docker compose up --build
```

Akses endpoint:
- Akademik Service: http://localhost:8081
- Rekap Service: http://localhost:8082

### Development Mode (Lokal)

```bash
# Terminal 1 - Akademik Service (butuh PostgreSQL berjalan)
cd akademik-service
go run . 

# Terminal 2 - Rekap Service
cd rekap-service
AKADEMIK_BASE_URL=http://localhost:8080 go run .
```

### Testing

```bash
cd akademik-service && go test ./...
cd rekap-service && go test ./...
```

### Postman Collection

Import file `postman/siakad-microservices.postman_collection.json` ke Postman.

## Timeout Budget

| Lapisan | Timeout |
|---------|---------|
| Postman / Client | 10 detik |
| Rekap Service server | 5 detik |
| HTTP call ke Akademik | 2 detik |
| Database query | 1 detik |

## Error Mapping

| Status | Error |
|--------|-------|
| 400 | `ErrNIMTidakValid`, `ErrNilaiTidakValid`, `ErrSKSTidakValid` |
| 404 | `ErrMahasiswaTidakAda` |
| 409 | `ErrMahasiswaSudahAda` |
| 503 | `ErrAkademikTidakTersedia` |
| 504 | `ErrAkademikTimeout` |

## Checklist

- [x] Kode Minggu 2 di akademik-service/ berfungsi penuh
- [x] GET /health di kedua service
- [x] Endpoint ringkasan mengembalikan JSON sesuai spesifikasi
- [x] Route rekap dihapus dari Akademik Service
- [x] Rekap Service tidak punya database/DSN/SQL
- [x] Interface AkademikClient sesuai spesifikasi
- [x] http.Client dengan Timeout, NewRequestWithContext
- [x] context diteruskan sampai NewRequestWithContext
- [x] defer resp.Body.Close()
- [x] 7 error terpetakan ke status yang benar
- [x] Dockerfile multi-stage, non-root user
- [x] .dockerignore ada di kedua service
- [x] Satu docker-compose.yml di root
- [x] Base URL dari environment variable
- [x] Log JSON dengan service name dan request ID
- [x] go fmt, go vet, go test hijau
- [x] Postman collection tersedia
