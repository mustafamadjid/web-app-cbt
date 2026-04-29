# 📝 SMAFI CBT — Computer Based Test Web Application

**SMAFI CBT** adalah aplikasi ujian berbasis komputer (_Computer Based Test_) yang dirancang untuk mengelola seluruh proses ujian secara digital — mulai dari pembuatan bank soal, penjadwalan ujian, pelaksanaan ujian oleh siswa, hingga koreksi dan rekap nilai secara otomatis.

Aplikasi ini dibangun menggunakan arsitektur modern dengan **Go (Golang)** sebagai backend dan **React + TypeScript** di sisi frontend, serta menerapkan pola **Hexagonal Architecture** untuk menjaga pemisahan tanggung jawab dan kemudahan pengembangan.

---

## 📑 Daftar Isi

- [Tech Stack](#-tech-stack)
- [Arsitektur Proyek](#-arsitektur-proyek)
- [Fitur-Fitur](#-fitur-fitur)
- [Prasyarat](#-prasyarat)
- [Cara Install & Menjalankan](#-cara-install--menjalankan)
- [Environment Variables](#-environment-variables)
- [Struktur Folder](#-struktur-folder)
- [Peran Pengguna (User Roles)](#-peran-pengguna-user-roles)
- [Lisensi](#-lisensi)

---

## 🛠 Tech Stack

### Backend

| Teknologi             | Keterangan                                     |
| --------------------- | ---------------------------------------------- |
| **Go 1.25**           | Bahasa pemrograman utama backend               |
| **httprouter**        | HTTP router ringan dan cepat                   |
| **pgx/v5**            | Driver PostgreSQL native untuk Go              |
| **golang-jwt/jwt/v5** | Autentikasi berbasis JSON Web Token            |
| **bcrypt**            | Hashing password yang aman                     |
| **rs/cors**           | Middleware untuk Cross-Origin Resource Sharing |
| **golang.org/x/time** | Rate limiting                                  |
| **PostgreSQL**        | Basis data relasional                          |

### Frontend

| Teknologi                | Keterangan                             |
| ------------------------ | -------------------------------------- |
| **React 19**             | Library UI utama                       |
| **TypeScript**           | Superset JavaScript dengan type safety |
| **Vite (Rolldown)**      | Build tool dan dev server super cepat  |
| **TailwindCSS v4**       | Utility-first CSS framework            |
| **MUI (Material UI) v7** | Komponen UI siap pakai                 |
| **Tiptap**               | Rich text editor (WYSIWYG)             |
| **Axios**                | HTTP client untuk komunikasi API       |
| **React Router v7**      | Routing SPA                            |
| **React Hot Toast**      | Notifikasi toast                       |
| **Lucide React**         | Icon library                           |

---

## 🏗 Arsitektur Proyek

Backend menggunakan **Hexagonal Architecture (Ports & Adapters)** yang memisahkan logika bisnis dari detail infrastruktur:

```
backend/internal/
├── adapter/        # Adapters (HTTP handlers, security, dll.)
│   ├── handler/    # HTTP handlers per fitur
│   └── security/   # Implementasi bcrypt, JWT
├── app/            # Application wiring & dependency injection
├── core/           # Business logic (domain, service, port)
│   ├── domain/     # Entity & value objects
│   ├── port/       # Interface (in/out ports)
│   ├── service/    # Use case / service layer
│   └── query/      # Query objects
├── db/             # Database migrations
└── infra/          # Infrastructure (DB connection, logging)
```

Frontend menggunakan struktur modular berbasis fitur:

```
frontend/src/
├── components/     # Reusable UI components
├── contexts/       # React context providers
├── hooks/          # Custom hooks
├── layouts/        # Layout wrappers (Admin, Siswa, dll.)
├── pages/          # Halaman per role (Admin, Siswa, Auth)
├── routes/         # Routing configuration
├── services/       # API service layer (Axios)
├── types/          # TypeScript type definitions
└── helper/         # Utility & helper functions
```

---

## ✨ Fitur-Fitur

### 🔐 Autentikasi & Otorisasi

- Login & logout dengan **JWT** (Access Token + Refresh Token via HTTP-only Cookie)
- Role-based access control: **Administrator**, **Guru**, **Siswa**
- Protected routes berdasarkan role pengguna
- Reset password

### 👥 Manajemen Pengguna

- **Kelola Akun Siswa** — CRUD data akun siswa
- **Kelola Akun Guru** — CRUD data akun guru
- Upload foto profil pengguna
- Lihat & edit profil pengguna

### 📚 Data Master

- **Mata Pelajaran** — Kelola daftar mata pelajaran
- **Kelas** — Kelola data kelas berdasarkan tingkat dan nama kelas
- **Ruang Ujian** — Kelola ruang ujian fisik / virtual
- **Sesi Ujian** — Kelola sesi/waktu ujian

### 📝 Bank Soal

- Buat dan kelola **bank soal** per mata pelajaran
- Tambah soal satu per satu secara manual
- **Import soal dari file DOCX** — parsing otomatis soal dari dokumen Word dengan format khusus
- Mendukung tipe soal **pilihan ganda** dan **essay**
- Mendukung **gambar** dalam soal
- Preview bank soal
- Background worker untuk pemrosesan import soal secara asinkron

### 📋 Penjadwalan Ujian

- Buat jadwal ujian dengan mengaitkan bank soal, kelas, ruang ujian, dan sesi
- Detail konfigurasi ujian (durasi, tanggal mulai, tanggal selesai)
- Lihat daftar jadwal ujian yang akan datang dan selesai

### ✏️ Pelaksanaan Ujian (Siswa)

- Masuk ujian menggunakan **token ujian**
- Mengerjakan soal pilihan ganda dan essay
- Timer ujian berjalan otomatis
- Submit jawaban secara otomatis saat waktu habis

### 📊 Hasil & Koreksi Ujian

- **Koreksi otomatis** untuk soal pilihan ganda
- **Background worker** untuk proses grading secara asinkron
- Rekap hasil ujian per kelas
- Detail jawaban siswa per ujian (attempt detail)
- Statistik nilai ujian
- Cetak hasil ujian

### 📢 Pengumuman

- Buat, edit, dan hapus pengumuman
- Rich text editor (**Tiptap WYSIWYG**) untuk konten pengumuman
- Mendukung upload gambar dalam pengumuman
- Pengumuman terlihat oleh siswa di dashboard

### 🏫 Profil Sekolah

- Kelola informasi profil sekolah (nama, alamat, logo, dll.)
- Konfigurasi pengaturan aplikasi

### 📈 Dashboard

- Dashboard khusus untuk setiap role (Admin, Guru, Siswa)
- Ringkasan statistik dan informasi penting

### 📋 Aktivitas Pengguna

- Pencatatan log aktivitas pengguna secara otomatis
- Monitoring aktivitas login dan operasi penting

---

## 📌 Prasyarat

Pastikan perangkat lunak berikut sudah terinstall di sistem Anda:

| Software       | Versi Minimum | Link Download                                          |
| -------------- | ------------- | ------------------------------------------------------ |
| **Go**         | 1.25+         | [golang.org/dl](https://golang.org/dl/)                |
| **Node.js**    | 18+           | [nodejs.org](https://nodejs.org/)                      |
| **npm**        | 9+            | Terinstall bersama Node.js                             |
| **PostgreSQL** | 14+           | [postgresql.org](https://www.postgresql.org/download/) |
| **Git**        | Terbaru       | [git-scm.com](https://git-scm.com/)                    |

---

## 🚀 Cara Install & Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/mustafamadjid/web-app-cbt.git
cd web-app-cbt
```

### 2. Setup Database PostgreSQL

Buat database baru di PostgreSQL:

```sql
CREATE DATABASE smafi_cbt_db;
```

Jalankan semua file migration yang ada di folder `backend/internal/db/migrations/` secara berurutan ke database yang telah dibuat.

### 3. Setup & Jalankan Backend

```powershell
# Masuk ke folder backend
cd backend

# Salin file environment dan sesuaikan konfigurasi
# Edit file .env sesuai kredensial database Anda
# Pastikan POSTGRES_DBURL mengarah ke database yang benar

# Contoh isi .env:
# ENVIRONMENT=dev
# BASE_URL=http://localhost:8080
# APP_DIR=.
# UPLOAD_DIR=public/uploads
# POSTGRES_DBURL=postgres://username:password@localhost:5432/smafi_cbt_db?sslmode=disable
# ISSUER=web-app-cbt
# ACCESS_TOKEN_SECRET=<your-access-secret>
# REFRESH_TOKEN_SECRET=<your-refresh-secret>
# TRUSTED_ORIGINS=http://localhost:5173

# Export environment variables (PowerShell)
. ./export_env.ps1

# Jalankan server backend
cd cmd/http
go run .
```

> **Alternatif**: Gunakan script `run.ps1` di folder backend:
>
> ```powershell
> cd backend
> ./run.ps1
> ```

Backend akan berjalan di `http://localhost:8080`.

### 4. Setup & Jalankan Frontend

```bash
# Buka terminal baru, masuk ke folder frontend
cd frontend

# Install dependencies
npm install

# Sesuaikan API URL di file .env (opsional, default sudah benar)
# VITE_API_URL=http://localhost:8080

# Jalankan development server
npm run dev
```

Frontend akan berjalan di `http://localhost:5173`.

### 5. Akses Aplikasi

Buka browser dan akses `http://localhost:5173`. Login sesuai peran pengguna yang telah didaftarkan di database.

---

## 🔧 Environment Variables

### Backend (`backend/.env`)

| Variable               | Keterangan                        | Contoh                                             |
| ---------------------- | --------------------------------- | -------------------------------------------------- |
| `ENVIRONMENT`          | Mode environment (`dev` / `prod`) | `dev`                                              |
| `BASE_URL`             | Base URL backend server           | `http://localhost:8080`                            |
| `APP_DIR`              | Direktori root aplikasi           | `.`                                                |
| `UPLOAD_DIR`           | Direktori penyimpanan file upload | `public/uploads`                                   |
| `POSTGRES_DBURL`       | Connection string PostgreSQL      | `postgres://user:pass@localhost:5432/smafi_cbt_db` |
| `ISSUER`               | Issuer JWT token                  | `web-app-cbt`                                      |
| `ACCESS_TOKEN_SECRET`  | Secret key untuk access token     | String acak yang panjang                           |
| `REFRESH_TOKEN_SECRET` | Secret key untuk refresh token    | String acak yang panjang                           |
| `TRUSTED_ORIGINS`      | Daftar origin CORS yang diizinkan | `http://localhost:5173`                            |

### Frontend (`frontend/.env`)

| Variable       | Keterangan      | Contoh                  |
| -------------- | --------------- | ----------------------- |
| `VITE_API_URL` | URL backend API | `http://localhost:8080` |

---

## 📂 Struktur Folder

```
web-app-cbt/
├── backend/                        # Backend (Go)
│   ├── cmd/
│   │   └── http/
│   │       ├── main.go             # Entry point aplikasi
│   │       └── public/             # Static files (uploads)
│   ├── internal/
│   │   ├── adapter/                # Adapters layer
│   │   │   ├── handler/http/       # HTTP handlers per fitur
│   │   │   │   ├── features/       # Handler setiap fitur
│   │   │   │   │   ├── aktivitas_user/
│   │   │   │   │   ├── bank_soal/
│   │   │   │   │   ├── dashboard/
│   │   │   │   │   ├── kelas/
│   │   │   │   │   ├── mata_pelajaran/
│   │   │   │   │   ├── pengumuman/
│   │   │   │   │   ├── profil_sekolah/
│   │   │   │   │   ├── ruang_ujian/
│   │   │   │   │   ├── sesi/
│   │   │   │   │   ├── ujian/
│   │   │   │   │   └── user/
│   │   │   │   └── ...
│   │   │   └── security/           # Bcrypt, JWT adapter
│   │   ├── app/                    # DI & module wiring
│   │   ├── core/                   # Business logic
│   │   │   ├── domain/             # Entities
│   │   │   ├── port/               # Interfaces (in/out)
│   │   │   ├── service/            # Use cases
│   │   │   └── query/              # Query models
│   │   ├── db/migrations/          # SQL migrations
│   │   └── infra/                  # DB pool, logger
│   ├── .env                        # Environment variables
│   ├── go.mod                      # Go module definition
│   └── run.ps1                     # Script untuk menjalankan backend
│
├── frontend/                       # Frontend (React + TypeScript)
│   ├── src/
│   │   ├── components/             # Komponen UI reusable
│   │   ├── contexts/               # React Context providers
│   │   ├── hooks/                  # Custom hooks
│   │   ├── layouts/                # Layout components
│   │   ├── pages/                  # Halaman aplikasi
│   │   │   ├── Admin/              # Dashboard & halaman admin
│   │   │   │   ├── Dashboard/
│   │   │   │   │   ├── BankSoal/
│   │   │   │   │   ├── DataMaster/
│   │   │   │   │   ├── KelolaAkun/
│   │   │   │   │   ├── KelolaSesi/
│   │   │   │   │   ├── Pengaturan/
│   │   │   │   │   ├── Pengumuman/
│   │   │   │   │   └── Ujian/
│   │   │   │   └── Cetak/
│   │   │   ├── Auth/               # Halaman login
│   │   │   ├── Profile/            # Halaman profil
│   │   │   └── Siswa/              # Dashboard & halaman siswa
│   │   │       ├── Ujian/
│   │   │       └── Pengumuman/
│   │   ├── routes/                 # Konfigurasi routing
│   │   ├── services/               # API service layer
│   │   ├── types/                  # TypeScript types
│   │   └── helper/                 # Utility functions
│   ├── .env                        # Frontend environment
│   ├── package.json                # Dependencies & scripts
│   ├── vite.config.ts              # Vite configuration
│   └── index.html                  # HTML entry point
│
├── openapi.yaml                    # API documentation (OpenAPI)
└── README.md                       # Dokumentasi proyek (file ini)
```

---

## 👤 Peran Pengguna (User Roles)

Aplikasi mendukung **3 peran pengguna** dengan akses yang berbeda:

### 🔴 Administrator

- Mengelola seluruh data master (mata pelajaran, kelas, ruang ujian, sesi)
- Mengelola akun guru dan siswa
- Mengelola bank soal dan ujian
- Melihat hasil ujian dan statistik
- Membuat pengumuman
- Mengatur profil sekolah dan konfigurasi aplikasi
- Mencetak laporan
- Melihat log aktivitas pengguna

### 🟡 Guru

- Mengelola bank soal (buat, edit, import soal)
- Membuat dan menjadwalkan ujian
- Melihat dan mengoreksi hasil ujian
- Membuat pengumuman
- Mencetak laporan

### 🟢 Siswa

- Melihat jadwal ujian
- Mengerjakan ujian dengan token akses
- Melihat hasil ujian
- Melihat pengumuman

---

## 📄 Lisensi

Proyek ini dibuat sebagai penelitian Tugas Akhir, Program Studi Teknik Informatika ITERA.

---

<p align="center">
  Dibuat oleh <strong>Mustafa Madjid</strong>
</p>
