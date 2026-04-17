# Integrasi Dashboard Statistik Frontend untuk Admin dan Guru

## Summary
- Integrasikan statistik dashboard backend ke halaman home shared untuk `ADMIN` dan `GURU`, menggantikan angka hardcoded pada card statistik.
- Pertahankan struktur page yang sekarang: `ADMIN` menampilkan card pengguna + 3 card statistik lain, `GURU` hanya menampilkan 3 card shared.
- Buat `type` khusus dan `service` baru; `hook` wrapper digabung di file service sesuai keputusan yang dikunci.
- Area non-statistik tetap di luar scope tahap ini: `Pengumuman`, `Jadwal Ujian`, `Ujian Berlangsung`, dan `Log Aktivitas` tidak diubah selain tetap hidup berdampingan dengan hook dashboard baru.

## Key Changes
- Tambah type frontend khusus untuk payload dashboard dengan shape mengikuti response backend:
  - `total_siswa`
  - `total_guru`
  - `total_ujian_terlaksana`
  - `total_bank_soal`
  - `total_mapel_aktif`
- Tambah service dashboard baru yang:
  - Menyediakan `getDashboardStatistik(role)` untuk fetch data statistik.
  - Menyediakan `useGetDashboardStatistik(role, enabled?)` berbasis `useFetch`.
  - Memakai resolver endpoint berbasis role agar home admin dan guru bisa berbagi service yang sama.
- Mapping endpoint yang dipakai di plan ini:
  - `ADMIN` -> `/admin/dashboard`
  - `GURU` -> endpoint guru dengan payload identik, direkomendasikan `/guru/dashboard`
- Ubah home shared admin/guru agar:
  - Memanggil `useGetDashboardStatistik(role, role === "ADMIN" || role === "GURU")`
  - Mengganti angka hardcoded di widget statistik dengan data API
  - Menjaga `TotalSiswaGuruWidget` tetap `ADMIN`-only
  - Menjaga 3 card statistik lain tampil di `ADMIN` dan `GURU`
- Ubah card `Total Ujian Terlaksana` dari widget donut menjadi simple stat karena backend hanya memberi total, bukan progress/denominator. Ini menghindari UI yang menampilkan persentase palsu.
- Tambahkan handling state fetch di home:
  - `loading`: tampilkan indikator ringan non-blocking, sambil widget memakai fallback aman
  - `error`: tampilkan pesan ringkas dan pertahankan page tetap render tanpa crash
- Nilai fallback saat data belum tersedia/error:
  - semua angka statistik default ke `0`
  - widget tetap render agar layout stabil

## Public APIs / Interfaces
- Tambahan type:
  - `DashboardStatistik` dengan field snake_case sesuai payload backend
- Tambahan service API:
  - `getDashboardStatistik(role: "ADMIN" | "GURU"): Promise<DashboardStatistik>`
- Tambahan hook:
  - `useGetDashboardStatistik(role: "ADMIN" | "GURU", enabled?: boolean)`
- Kontrak backend yang dibutuhkan untuk guru:
  - payload guru harus identik dengan admin agar frontend tidak butuh type atau mapper terpisah

## Test Plan
- Validasi service memilih endpoint yang benar berdasarkan role.
- Validasi hook mengembalikan `data`, `loading`, `error`, dan `refetch` sesuai pola `useFetch`.
- Validasi home `ADMIN`:
  - card `Total Siswa/Guru` muncul
  - nilai semua card berasal dari API, bukan hardcoded
- Validasi home `GURU`:
  - card `Total Siswa/Guru` tidak muncul
  - card `Total Ujian Terlaksana`, `Total Bank Soal`, dan `Mata Pelajaran Aktif` tampil dari API
- Validasi loading state:
  - halaman tetap render
  - tidak ada crash/null access
- Validasi error state:
  - banner/pesan error tampil
  - widget tetap aman dengan fallback `0`
- Manual acceptance:
  - login sebagai admin dan guru
  - buka home masing-masing
  - refresh page dan pastikan statistik konsisten dengan response backend
  - pastikan role lain (`SISWA`) tidak terdampak

## Assumptions
- Hook wrapper digabung di file service, bukan file hook terpisah.
- Type diletakkan di folder `types`, service di folder `services`.
- Layout home admin/guru tetap memakai page shared yang sama; diferensiasi cukup dengan role check yang sudah ada.
- `Total Ujian Terlaksana` tidak lagi memakai donut/progress sampai backend benar-benar menyediakan field pembagi atau percent.
- Integrasi guru bergantung pada backend membuka endpoint statistik untuk guru dengan payload yang sama seperti admin; tanpa ini, integrasi final guru tidak bisa aktif penuh.
