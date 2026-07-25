# Product Requirements Document (PRD) - SIAKAD Microservices

## 1. Ringkasan & Tujuan
Memecah aplikasi monolith SIAKAD API (dari tugas Minggu 2) menjadi dua *service* Go (Akademik Service dan Rekap Service) yang berdiri sendiri. Kedua *service* akan saling berkomunikasi secara sinkron melalui REST API dan diorkestrasi menggunakan satu perintah Docker Compose. Aturan domain dan perhitungan IPK tetap sama, perubahan difokuskan pada kepemilikan data dan metode komunikasi antar *service*.

## 2. Arsitektur & Kepemilikan Data
Batas domain dibagi berdasarkan kepemilikan data.
*   **Akademik Service (Pemilik Data):** Memiliki seluruh data akademik (Tabel mahasiswa & nilai). Berfungsi untuk mencatat mahasiswa dan nilainya (mendukung operasi POST/PUT/DELETE/GET). Lapisan arsitekturnya terdiri dari handler -> service -> repository -> DB. Terhubung langsung ke `akademik-db` (PostgreSQL).
*   **Rekap Service (Konsumen Data):** Tidak memiliki data (database) sama sekali. Berfungsi untuk meringkas keadaan akademik menjadi laporan dengan meminjam data lewat HTTP ke Akademik Service. Lapisan arsitekturnya terdiri dari handler -> service -> client -> HTTP. **Dilarang keras memiliki DSN database di env.**

## 3. Spesifikasi API & Endpoint

### 3.1. Akademik Service (Port Publikasi: 8081, Port Internal: 8080)
*   `GET /health`: Liveness check (Status 200).
*   `POST /api/v1/mahasiswa`: Tambah mahasiswa (Status 201 / 400, 409).
*   `GET /api/v1/mahasiswa`: Daftar mahasiswa (Status 200).
*   `GET /api/v1/mahasiswa/:nim`: Detail + IPK (Status 200 / 404).
*   `PUT /api/v1/mahasiswa/:nim`: Ubah mahasiswa (Status 200 / 400, 404).
*   `DELETE /api/v1/mahasiswa/:nim`: Hapus mahasiswa (Status 200 / 404).
*   `POST /api/v1/mahasiswa/:nim/nilai`: Input nilai (Status 201 / 400, 404).
*   `GET /api/v1/mahasiswa/:nim/transkrip`: Daftar nilai + IPK (Status 200 / 404).
*   `GET /api/v1/mahasiswa/:nim/ringkasan`: Kontrak internal (Status 200 / 404).
*   *Catatan: Semua endpoint `GET /api/v1/rekap/*` harus dihapus dari service ini (pindah ke Rekap Service).*

### 3.2. Rekap Service (Port Publikasi: 8082, Port Internal: 8080)
*   `GET /health`: Liveness check (Status 200).
*   `GET /api/v1/rekap/jurusan`: Mengelompokkan mahasiswa per jurusan beserta rata-rata IPK (Status 200 / 503, 504).
*   `GET /api/v1/rekap/top-ipk?n=3`: Mendapatkan n mahasiswa dengan IPK tertinggi, urut menurun (Status 200 / 400, 503, 504). Parameter n bersifat opsional (default 3). Jika n <= 0, kembalikan 400.
*   `GET /api/v1/rekap/mahasiswa/:nim`: Ringkasan satu mahasiswa untuk uji jalur 404 (Status 200 / 404, 503, 504).

### 3.3. Kontrak Internal (Ringkasan)
Digunakan oleh Rekap Service untuk memanggil data dari Akademik Service. Format respons JSON wajib konsisten menggunakan standar amplop respons:
```json
// Sukses (200 OK)
{
  "sukses": true,
  "data": {
    "nim": "23010001",
    "nama": "Bunga",
    "jurusan": "Teknik Informatika",
    "status": "Aktif",
    "total_sks": 5,
    "ipk": 3.60
  }
}

// Gagal (404 Not Found)
{
  "sukses": false,
  "error": "mahasiswa tidak ditemukan"
}
```

## 4. Persyaratan Teknis & Komponen

### 4.1. Interface HTTP Client (Di Rekap Service)
Wajib menggunakan interface agar dapat diuji (mocking) tanpa menyalakan Akademik Service.
```go
type AkademikClient interface {
    DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error)
    Ringkasan(ctx context.Context, nim string) (Ringkasan, error)
}
```
*   Implementasi HTTP menggunakan `http.Client` dengan timeout eksplisit (tidak menggunakan `http.Get` atau `DefaultClient`).
*   Menggunakan `http.NewRequestWithContext` dan meneruskan `context.Context`.
*   Base URL didapat dari environment variable `AKADEMIK_BASE_URL`.

### 4.2. Pemetaan Error & HTTP Status
*   **400 Bad Request:** `ErrNIMTidakValid` (NIM bukan 8 digit), `ErrNilaiTidakValid` (Mutu di luar 0.0-4.0), `ErrSKSTidakValid` (SKS <= 0).
*   **409 Conflict:** `ErrMahasiswaSudahAda` (NIM sudah terdaftar).
*   **404 Not Found:** `ErrMahasiswaTidakAda` (NIM tidak ditemukan).
*   **504 Gateway Timeout (BARU):** `ErrAkademikTimeout` (Panggilan > 2 detik, context.DeadlineExceeded).
*   **503 Service Unavailable (BARU):** `ErrAkademikTidakTersedia` (Connection refused, DNS gagal, atau status 5xx dari Akademik Service).
*   *Catatan: Kegagalan tetangga tidak boleh memunculkan status 500 di end-user.*

### 4.3. Anggaran Timeout Berlapis
*   Batas kesabaran klien / Postman: 10s
*   Server timeout Rekap Service: 5s
*   Panggilan HTTP ke Akademik Service: 2s
*   Query database di Akademik Service: 1s

### 4.4. Infrastruktur & Docker
*   **Dockerfile:** Harus menggunakan multi-stage build (`golang:1.24-alpine` -> `alpine:3.21`). Menggunakan `CGO_ENABLED=0 GOOS=linux`, memiliki sertifikat `ca-certificates`, berjalan sebagai user non-root (`appuser`), dan memiliki `.dockerignore`.
*   **Docker Compose (`docker-compose.yml`):**
    *   Tiga container: `akademik-db`, `akademik-service`, `rekap-service`.
    *   `healthcheck` pada database dan `depends_on: condition: service_healthy` pada service.
    *   Volume bernama untuk persistensi data Postgres.
    *   Komunikasi internal menggunakan DNS internal Compose (contoh: `http://akademik-service:8080`).

## 5. Struktur Direktori
```text
siakad-microservices/
├── docker-compose.yml
├── README.md
├── postman/
├── akademik-service/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── go.mod
│   ├── main.go
│   ├── migrations/
│   └── internal/
└── rekap-service/
    ├── Dockerfile
    ├── .dockerignore
    ├── go.mod
    ├── main.go
    └── internal/
        ├── client.go
        ├── client_test.go
        ├── service.go
        ├── service_test.go
        └── handler.go
```
*Catatan: Tidak boleh ada import silang antar modul.*

## 6. Persyaratan Pengujian (Testing)
1.  **akademik-service/service_test.go:** Table-driven test untuk fitur Hitung IPK (minimal 4 skenario).
2.  **rekap-service/client_test.go:** Pengujian penerjemahan respons dengan `httptest.NewServer`. Kasus uji meliputi respons 200 sukses, 404 (ErrMahasiswaTidakAda), 500 (ErrAkademikTidakTersedia), delay 3 detik (ErrAkademikTimeout), dan 200 JSON rusak.
3.  **rekap-service/service_test.go:** Pengujian fitur TopIPK & PerJurusan menggunakan client tiruan (mock) yang memenuhi `AkademikClient`.

## 7. Kriteria Penerimaan (Skenario Demo)
Sistem harus berhasil menjalankan skenario berikut dari *cold start*:
1.  `docker compose up --build` berjalan lancar dengan semua service sehat.
2.  Endpoint `/health` di port 8081 dan 8082 merespons 200.
3.  Operasi CRUD mahasiswa dan input nilai berjalan dan terakumulasi di kalkulasi IPK secara benar.
4.  Fitur rekap per jurusan dan top-ipk mengembalikan data yang valid dan terkalkulasi akurat.
5.  Pencarian rekap mahasiswa yang tidak ada mengembalikan respons 404, bukan 500/503.
6.  **Resiliency Test 1:** Mematikan `akademik-service` (`docker compose stop akademik-service`) lalu hit endpoint rekap harus memunculkan respons 503 secara cepat tanpa menggantung.
7.  **Resiliency Test 2:** Memperlambat respon `akademik-service` (menyisipkan delay 5 detik) harus memunculkan respons 504 pada pemanggilan rekap setelah batas maksimal 2 detik.
8.  Log JSON termonitor dengan baik (memuat nama service, *request id*, dan durasi).

## 8. Pastikan ini Dibaca Terlebih dahulu (NOTE)
Berikut adalah teks dari gambar image_1d1adf.png yang telah disalin ulang dengan mempertahankan format kode dan penekanan kata (tebal):

Checklist Pengumpulan
[ ] Kode Minggu 2 sudah berada di akademik-service/ dan masih berfungsi penuh.

[ ] GET /health ada di kedua service dan menjawab 200.

[ ] Endpoint ringkasan (3.3) mengembalikan bentuk JSON persis seperti spesifikasi.

[ ] Route rekap sudah dihapus dari Akademik Service — tidak ada dua implementasi.

[ ] Rekap Service tidak punya driver database, DSN, maupun SQL satu baris pun.

[ ] Interface AkademikClient persis Bagian 3.4, dipenuhi implisit oleh implementasi HTTP.

[ ] http.Client dibuat sekali dengan Timeout; bukan http.Get atau DefaultClient.

[ ] context diteruskan dari c.Request.Context() sampai ke NewRequestWithContext; defer cancel() ada.

[ ] defer resp.Body.Close() di setiap panggilan keluar.

[ ] Tujuh error (lima lama + dua baru) terpetakan ke status yang benar: 400 / 404 / 409 / 503 / 504.

[ ] Dockerfile multi-stage di kedua service; image runtime bukan golang:; proses tidak berjalan sebagai root.

[ ] .dockerignore ada di kedua service.

[ ] Satu docker-compose.yml di root menyalakan seluruh sistem dengan satu perintah.

[ ] Base URL antar service dibaca dari environment variable, bukan hardcode.

[ ] Log JSON ke stdout, memuat nama service dan request id yang diteruskan lewat X-Request-ID.

[ ] Lima skenario Bagian 6.1 dan seluruh skenario demo Bagian 7 berjalan.

[ ] go fmt ./..., go vet ./... bersih dan go test ./... HIJAU di kedua modul.

[ ] README memuat diagram singkat, alasan pembagian domain, dan cara menjalankan.

[ ] Postman collection tersimpan di postman/.

[ ] Kode di-push ke repo GitHub Minggu 3 (sertakan .env.example, bukan .env); commit bertahap, bukan satu commit raksasa.