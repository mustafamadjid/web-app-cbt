# Menjalankan Test dengan Coverage

Jalankan perintah Go dari direktori `backend`:

```powershell
Set-Location backend
```

## Service

Jalankan seluruh test service dan ukur hanya statement pada layer service:

```powershell
go test -count=1 -cover `
  -coverpkg "github.com/mustafamadjid/web-app-cbt/internal/core/service/..." `
  -coverprofile "service_coverage.out" `
  "./internal/core/service/..."
```

Lihat total coverage dan buat laporan HTML:

```powershell
go tool cover -func service_coverage.out | Select-Object -Last 1
go tool cover -html service_coverage.out -o service_coverage.html
```

## Ringkasan Skenario Unit Test Service

Gunakan `run-service-test.ps1` untuk menjalankan unit test service dan menghitung persentase kelulusan skenario Basis Path serta Branch Coverage. Skrip ini menghitung **495 skenario** yang didokumentasikan dalam `test.md` dan mengecualikan package `integration_test`.

Jalankan skrip dari direktori root proyek, bukan dari `backend`:

```powershell
Set-Location ..
.\run-service-test.ps1
```

Jika terminal sudah berada di direktori root proyek, cukup jalankan:

```powershell
.\run-service-test.ps1
```

Hasil eksekusi menampilkan jumlah skenario `Passed`, `Failed`, `Skipped`, total yang terdeteksi, target, dan persentase kelulusan. Jika ada skenario gagal, nama package dan test yang gagal juga ditampilkan.

Skrip mengembalikan exit code `1` jika terdapat test gagal atau jumlah skenario tidak sama dengan 495. Jika seluruh pemeriksaan berhasil, skrip mengembalikan exit code `0`.

## HTTP Handler

Jalankan seluruh test HTTP handler dan ukur hanya statement pada layer HTTP handler:

```powershell
go test -count=1 -cover `
  -coverpkg "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/..." `
  -coverprofile "handler_http_coverage.out" `
  "./internal/adapter/handler/http/..."
```

Lihat total coverage dan buat laporan HTML:

```powershell
go tool cover -func handler_http_coverage.out | Select-Object -Last 1
go tool cover -html handler_http_coverage.out -o handler_http_coverage.html
```

## Keseluruhan Backend

Jalankan seluruh test dan ukur statement pada semua package backend:

```powershell
go test -count=1 -cover `
  -coverpkg "github.com/mustafamadjid/web-app-cbt/..." `
  -coverprofile "coverage.out" `
  "./..."
```

Lihat total coverage dan buat laporan HTML:

```powershell
go tool cover -func coverage.out | Select-Object -Last 1
go tool cover -html coverage.out -o coverage.html
```

## Hasil Test Verbose

Jalankan seluruh test dalam mode verbose untuk melihat test yang `PASS` dan `FAIL`:

```powershell
go test -count=1 -v "./..."
```

Simpan sekaligus tampilkan hasil test verbose:

```powershell
go test -count=1 -v "./..." 2>&1 | Tee-Object -FilePath test_verbose.log
```

## Test Verbose untuk Package Tertentu

Jalankan seluruh test service dalam mode verbose:

```powershell
go test -count=1 -v "github.com/mustafamadjid/web-app-cbt/internal/core/service/..."
```

Gunakan path package relatif dari direktori `backend`:

```powershell
go test -count=1 -v "./path/ke/package"
```

Contoh untuk package test service pembuatan user:

```powershell
go test -count=1 -v "./internal/core/service/user/create/test"
```

Contoh untuk package test HTTP handler kelas:

```powershell
go test -count=1 -v "./internal/adapter/handler/http/features/kelas/test"
```

Jalankan satu test tertentu dalam package menggunakan nama test yang sama persis:

```powershell
go test -count=1 -v "./path/ke/package" -run "^TestNamaTest$"
```
