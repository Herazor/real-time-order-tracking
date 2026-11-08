# Real-Time Order & Job Tracking System

Aplikasi untuk mengelola dan memantau status pesanan atau pekerjaan secara real-time menggunakan Go, React, PostgreSQL, dan WebSocket.

## 1. Teknologi yang Digunakan

### Backend
- Go
- Gin Framework
- GORM
- PostgreSQL
- Gorilla WebSocket

### Frontend
- React
- Vite

### Infrastruktur
- Docker
- Docker Compose

## 2. Struktur Proyek

```text
study-case-1/
├── backend/
│   ├── main.go
│   ├── orders.go
│   ├── server
│   └── ...
├── frontend/
│   ├── src/
│   ├── package.json
│   └── ...
├── docker-compose.yaml
├── README.md
└── System-Design-Answers.pdf
```

## 3. Persyaratan Sistem

Pastikan perangkat sudah memiliki:
- Docker Desktop
- Git

## 4. Cara Menjalankan Aplikasi

### Clone Repository

```bash
git clone https://github.com/Herazor/real-time-order-tracking.git
cd real-time-order-tracking
```

### Menjalankan Aplikasi

Jalankan seluruh aplikasi menggunakan satu perintah:

```bash
docker compose up --build
```

Docker Compose akan menjalankan layanan berikut:
- PostgreSQL sebagai database.
- Go Backend sebagai layanan API.
- React Frontend sebagai antarmuka aplikasi.

### Mengakses Aplikasi

| Layanan | URL |
|---|---|
| Frontend | http://localhost:5173 |
| Backend | http://localhost:8080 |
| PostgreSQL | Jaringan internal Docker |

### Menghentikan Aplikasi

Untuk menghentikan seluruh layanan:

```bash
docker compose down
```

Untuk menghentikan seluruh layanan sekaligus menghapus data database:

```bash
docker compose down -v
```

**Catatan:** Perintah kedua akan menghapus volume database sehingga seluruh data yang tersimpan di dalamnya ikut terhapus.

## 5. Fitur Utama

- Membuat pesanan baru.
- Menampilkan seluruh pesanan.
- Menampilkan detail pesanan.
- Memperbarui status pesanan.
- Membatalkan pesanan.
- Meninjau permintaan pembatalan pesanan.
- Memperbarui status pesanan secara real-time menggunakan WebSocket.
- Menyimpan riwayat peristiwa pesanan.
- Melacak urutan peristiwa menggunakan nomor urut.
- Menggunakan Optimistic Concurrency Control melalui pemeriksaan versi data.

## 6. Daftar API Endpoint

| Metode | Endpoint | Keterangan |
|---|---|---|
| POST | `/api/orders` | Membuat pesanan baru |
| GET | `/api/orders` | Menampilkan seluruh pesanan |
| GET | `/api/orders/:id` | Menampilkan detail pesanan |
| GET | `/api/orders/:id/events` | Menampilkan riwayat peristiwa pesanan |
| PATCH | `/api/orders/:id/status` | Memperbarui status pesanan |
| POST | `/api/orders/:id/cancel` | Mengajukan permintaan pembatalan pesanan |
| PATCH | `/api/cancellations/:id/review` | Meninjau permintaan pembatalan |
| GET | `/ws/orders/:id` | Menghubungkan klien ke pembaruan pesanan secara real-time |

## 7. Komunikasi Real-Time

Aplikasi menggunakan WebSocket untuk mengirimkan perubahan status pesanan kepada seluruh klien yang sedang terhubung.

Ketika status pesanan berubah, proses yang dilakukan adalah:

1. Backend memproses perubahan status pesanan.
2. Database menyimpan data pesanan yang telah diperbarui.
3. Sistem mencatat peristiwa perubahan pesanan.
4. Klien yang terhubung melalui WebSocket menerima pembaruan tersebut.

Dengan mekanisme ini, frontend dapat menerima perubahan status pesanan tanpa perlu memuat ulang halaman.

## 8. Status Implementasi

| Fitur | Status |
|---|---|
| REST API | Sudah Diimplementasikan |
| Integrasi PostgreSQL | Sudah Diimplementasikan |
| Riwayat Peristiwa Pesanan | Sudah Diimplementasikan |
| Pembaruan WebSocket | Sudah Diimplementasikan |
| Pelacakan Urutan Peristiwa | Sudah Diimplementasikan |
| Optimistic Concurrency Control | Sudah Diimplementasikan |
| Docker Compose | Sudah Diimplementasikan |
| Autentikasi dan Otorisasi | Belum Diimplementasikan |
| Deduplicasi Permintaan dengan Idempotency | Belum Diimplementasikan |
| Penyambungan Ulang WebSocket dan Pemutaran Ulang Peristiwa Otomatis | Belum Diimplementasikan |
| Redis Pub/Sub untuk Penyiaran Multi-Instance | Belum Diimplementasikan |
| Transactional Outbox | Belum Diimplementasikan |

## 9. Keterbatasan Sistem

Beberapa keterbatasan yang masih terdapat dalam implementasi saat ini:

- Autentikasi dan otorisasi berbasis peran belum sepenuhnya diimplementasikan.
- Frontend masih menggunakan ID pengguna demo.
- Idempotency Key belum digunakan untuk mencegah pemrosesan permintaan yang sama secara berulang.
- Penyiaran WebSocket menggunakan registri koneksi berbasis memori sehingga hanya mendukung satu instance backend.
- Penyambungan ulang WebSocket dan pemutaran ulang peristiwa secara otomatis belum diimplementasikan.
- Redis Pub/Sub belum diintegrasikan.
- Transactional Outbox belum diimplementasikan.
- Docker image backend menggunakan binary Linux hasil kompilasi yang telah disertakan dalam repository.

## 10. Pengendalian Konkruensi

Sistem menggunakan metode Optimistic Concurrency Control untuk mencegah konflik ketika beberapa pengguna memperbarui pesanan secara bersamaan.

Setiap pesanan memiliki nomor versi (*version number*).

Ketika melakukan perubahan status, permintaan harus menyertakan versi data yang diharapkan.

Backend kemudian memeriksa apakah versi yang dikirimkan masih sesuai dengan versi terbaru di database sebelum melakukan pembaruan.

Mekanisme ini membantu mencegah perubahan data yang saling bertentangan dan menghindari penimpaan pembaruan dari pengguna lain.

## 11. Riwayat Peristiwa dan Auditabilitas

Setiap perubahan status pesanan dicatat dalam riwayat peristiwa pesanan (*order event history*).

Setiap peristiwa menyimpan informasi seperti:

- ID Pesanan.
- Jenis Peristiwa.
- Status Sebelumnya.
- Status Terbaru.
- Nomor Urut Peristiwa.
- ID Aktor atau Pengguna.
- Waktu Terjadinya Peristiwa.

Dengan mekanisme ini, sistem dapat menyimpan jejak perubahan pesanan sehingga proses pemantauan dan penelusuran aktivitas menjadi lebih mudah.

## 12. Catatan

Proyek ini merupakan implementasi dari studi kasus pengembangan sistem pelacakan pesanan dan pekerjaan secara real-time.

Beberapa fitur lanjutan yang dibutuhkan untuk penggunaan pada lingkungan produksi masih menjadi bagian dari rencana pengembangan berikutnya dan belum diimplementasikan sepenuhnya.
