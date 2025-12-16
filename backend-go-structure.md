# Backend Go Project Structure

Dokumen ini mendeskripsikan struktur direktori yang diusulkan untuk backend Go. Struktur ini bisa dijadikan landasan saat mengembangkan layanan.

```
myproject
├── cmd
│   ├── api
│   │   └── main.go
│   └── worker
│       └── main.go
│
├── internal
│   ├── domain
│   │   ├── user
│   │   │   ├── entity.go
│   │   │   └── policy.go
│   │   ├── exam
│   │   │   ├── entity.go
│   │   │   └── policy.go
│   │   ├── question
│   │   │   └── entity.go
│   │   └── attempt
│   │       └── entity.go
│   │
│   ├── usecase
│   │   └── exam
│   │       ├── start_attempt.go
│   │       ├── submit_attempt.go
│   │       ├── grade_attempt.go
│   │       ├── ports.go        // interface repository/gateway yang dibutuhkan usecase
│   │       ├── start_attempt_test.go
│   │       └── submit_attempt_test.go
│   │
│   └── adapter
│       ├── http
│       │   └── handler
│       │       └── exam_handler.go
│       ├── persistence
│       │   └── postgres
│       │       └── exam_repository.go
│       ├── clock
│       └── idgen
│
├── migrations
├── configs
└── tests
    ├── integration
    └── e2e
```

## Penjelasan Singkat
- `cmd/`: Entry point untuk menjalankan aplikasi (misalnya API dan worker).
- `internal/domain/`: Definisi domain model dan policy bisnis.
- `internal/usecase/`: Use case dengan kontrak port untuk repository atau gateway.
- `internal/adapter/`: Implementasi adapter seperti handler HTTP, repository PostgreSQL, generator ID, dan utilitas waktu.
- `migrations/`: Skrip migrasi database.
- `configs/`: Berkas konfigurasi aplikasi.
- `tests/`: Direktori untuk pengujian integrasi dan end-to-end.
