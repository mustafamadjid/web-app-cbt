# Backend Go Project Structure

Dokumen ini mendeskripsikan struktur direktori yang diusulkan untuk backend Go. Struktur ini bisa dijadikan landasan saat mengembangkan layanan.

```
├── cmd
│   └── http
│       └── main.go
├── docs
├── internal
│   ├── adapter
│   │   ├── handler
│   │   │   └── http
│   │   │       ├── auth
│   │   │       │   ├── auth_handler.go
│   │   │       │   ├── auth_request.go
│   │   │       │   └── auth_response.go
│   │   │       ├── cookie
│   │   │       │   └── cookie.go
│   │   │       ├── features
│   │   │       │   ├── profil_sekolah
│   │   │       │   │   ├── get
│   │   │       │   │   │   └── get.go
│   │   │       │   │   └── update
│   │   │       │   │       └── update.go
│   │   │       │   └── user
│   │   │       │       ├── create
│   │   │       │       │   └── create.go
│   │   │       │       ├── delete
│   │   │       │       │   └── delete.go
│   │   │       │       ├── get
│   │   │       │       │   ├── get_guru.go
│   │   │       │       │   └── get_siswa.go
│   │   │       │       ├── update
│   │   │       │       │   └── update.go
│   │   │       │       ├── user_request.go
│   │   │       │       └── user_response.go
│   │   │       ├── helper
│   │   │       │   ├── response_envelope
│   │   │       │   │   ├── envelope.go
│   │   │       │   │   └── response_helper.go
│   │   │       │   └── image_upload_saver.go
│   │   │       ├── middleware
│   │   │       │   ├── auth_middleware.go
│   │   │       │   └── token_middleware.go
│   │   │       └── validation
│   │   │           └── input_validation.go
│   │   ├── repository
│   │   │   └── postgres
│   │   │       ├── auth_user_repo.go
│   │   │       ├── executor.go
│   │   │       ├── profil_guru_repo.go
│   │   │       ├── profil_sekolah_repo.go
│   │   │       ├── profil_siswa_repo.go
│   │   │       ├── session_repo.go
│   │   │       ├── tx.go
│   │   │       ├── tx_manager.go
│   │   │       └── user_repo.go
│   │   ├── securtity
│   │   │   └── bcrypt
│   │   │       └── hasher.go
│   │   └── token
│   │       ├── jwt_access.go
│   │       └── jwt_refresh.go
│   ├── app
│   │   ├── auth_module.go
│   │   ├── build.go
│   │   ├── config.go
│   │   ├── http_module.go
│   │   ├── infra_module.go
│   │   ├── profil_sekolah_module.go
│   │   ├── token_module.go
│   │   └── user_module.go
│   ├── core
│   │   ├── core_error
│   │   │   └── core_error.go
│   │   ├── domain
│   │   │   ├── auth
│   │   │   │   └── session
│   │   │   │       └── session.go
│   │   │   ├── profil_sekolah
│   │   │   │   └── profil_sekolah.go
│   │   │   └── user
│   │   │       ├── cmd
│   │   │       ├── profil_guru.go
│   │   │       ├── profil_siswa.go
│   │   │       └── user.go
│   │   ├── port
│   │   │   └── in
│   │   │       ├── auth_port_in
│   │   │       │   └── auth_use_case.go
│   │   │       └── user
│   │   │           ├── profil_guru_use_case.go
│   │   │           └── profil_siswa_use_case.go
│   │   ├── query
│   │   │   └── user
│   │   │       ├── list_guru.go
│   │   │       └── list_siswa.go
│   │   ├── service
│   │   │   ├── auth_service
│   │   │   │   ├── auth_service.go
│   │   │   │   ├── auth_service_test.go
│   │   │   │   └── login_cmd.go
│   │   │   ├── profil_sekolah
│   │   │   │   ├── get
│   │   │   │   │   ├── service.go
│   │   │   │   │   └── service_test.go
│   │   │   │   └── update
│   │   │   │       ├── service.go
│   │   │   │       └── service_test.go
│   │   │   └── user
│   │   │       ├── create
│   │   │       │   ├── coverage.html
│   │   │       │   ├── coverage.out
│   │   │       │   ├── create.go
│   │   │       │   ├── create_cmd.go
│   │   │       │   ├── create_res.go
│   │   │       │   └── create_test.go
│   │   │       ├── delete
│   │   │       │   ├── delete.go
│   │   │       │   └── delete_test.go
│   │   │       ├── get
│   │   │       │   ├── get_guru.go
│   │   │       │   ├── get_guru_test.go
│   │   │       │   ├── get_siswa.go
│   │   │       │   └── get_siswa_test.go
│   │   │       └── update
│   │   │           ├── update.go
│   │   │           ├── update_cmd.go
│   │   │           └── update_test.go
│   │   └── util
│   ├── db
│   │   └── migrations│   
│   └── infra
│       └── db
│           └── postgres.go
├── .gitignore
├── export_env.ps1
├── go.mod
└── go.sum
```

## Penjelasan Singkat
- `cmd/`: Entry point untuk menjalankan aplikasi (misalnya API dan worker).
- `internal/domain/`: Definisi domain model dan policy bisnis.
- `internal/usecase/`: Use case dengan kontrak port untuk repository atau gateway.
- `internal/adapter/`: Implementasi adapter seperti handler HTTP, repository PostgreSQL, generator ID, dan utilitas waktu.
- `migrations/`: Skrip migrasi database.
- `configs/`: Berkas konfigurasi aplikasi.
- `tests/`: Direktori untuk pengujian integrasi dan end-to-end.
