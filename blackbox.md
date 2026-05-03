# Skenario Blackbox Testing Fitur Kelola Akun Siswa dan Guru

Dokumen ini berisi skenario blackbox testing untuk fitur kelola akun `siswa` dan `guru` pada area admin, mencakup list, detail, tambah akun, ubah akun, hapus akun, hapus massal, reset password, validasi input, otorisasi, dan efek terhadap session login.

## 1. Ruang Lingkup

Fitur yang diuji:

- List akun siswa `GET /admin/siswa`
- Detail akun siswa `GET /admin/siswa/:id`
- Tambah akun siswa `POST /admin/siswa`
- Ubah akun siswa `PATCH /admin/siswa/:id`
- List akun guru `GET /admin/guru`
- Detail akun guru `GET /admin/guru/:id`
- Tambah akun guru `POST /admin/guru`
- Ubah akun guru `PATCH /admin/guru/:id`
- Hapus satu akun pengguna `DELETE /admin/pengguna/:id`
- Hapus banyak akun pengguna `DELETE /admin/pengguna`
- Reset password akun `PUT /admin/pengguna/:idPengguna/reset-password`
- Otorisasi admin terhadap seluruh endpoint kelola akun
- Validasi multipart upload foto profil
- Revokasi session user target setelah update data

Di luar ruang lingkup:

- Pengujian detail modul kelas, dashboard, auth, ujian, pengumuman, dan bank soal
- Pengujian whitebox internal service/repository
- Pengujian performa beban tinggi

## 2. Ringkasan Perilaku Sistem

Berdasarkan implementasi saat ini:

- Semua endpoint kelola akun hanya boleh diakses `ADMIN`.
- Endpoint list dan detail menggunakan method `GET`.
- Endpoint create dan update menggunakan `multipart/form-data`.
- Endpoint reset password menggunakan `application/json`.
- Username wajib 5 sampai 20 karakter.
- Username hanya boleh berisi huruf, angka, titik, dan underscore.
- Password minimal 8 karakter dan maksimal 72 bytes.
- Email opsional, tetapi jika diisi harus valid dan akan dinormalisasi lowercase.
- Nomor HP opsional, tetapi jika diisi harus sesuai pola nomor telepon.
- Foto profil opsional.
- Upload multipart dibatasi 10 MB.
- File foto yang terlalu besar akan ditolak dengan `FILE_TOO_LARGE`.
- Akun guru yang dibuat otomatis memiliki role `GURU` dan status awal `AKTIF`.
- Akun siswa yang dibuat otomatis memiliki role `SISWA` dan status awal `AKTIF`.
- NIP guru harus `18 digit` atau `-`.
- NISN siswa harus `10 digit` atau `-`.
- `angkatan` siswa harus berada pada rentang valid domain.
- `no_absen` siswa harus lebih dari 0.
- `tanggal_lahir` siswa harus format `YYYY-MM-DD`.
- Update data guru/siswa akan merevoke semua session aktif user target setelah berhasil.
- Hapus banyak akun menerima body JSON berisi array `ids`.

## 3. Endpoint yang Diuji

| Fitur | Endpoint | Method | Content-Type |
|---|---|---|---|
| List siswa | `/admin/siswa` | `GET` | - |
| Detail siswa | `/admin/siswa/:id` | `GET` | - |
| Tambah siswa | `/admin/siswa` | `POST` | `multipart/form-data` |
| Ubah siswa | `/admin/siswa/:id` | `PATCH` | `multipart/form-data` |
| List guru | `/admin/guru` | `GET` | - |
| Detail guru | `/admin/guru/:id` | `GET` | - |
| Tambah guru | `/admin/guru` | `POST` | `multipart/form-data` |
| Ubah guru | `/admin/guru/:id` | `PATCH` | `multipart/form-data` |
| Hapus satu pengguna | `/admin/pengguna/:id` | `DELETE` | - |
| Hapus banyak pengguna | `/admin/pengguna` | `DELETE` | `application/json` |
| Reset password | `/admin/pengguna/:idPengguna/reset-password` | `PUT` | `application/json` |

## 4. Data Uji

| Kode | Tipe | Kondisi |
|---|---|---|
| DU-ADM-01 | Admin | Admin valid dan sedang login |
| DU-GRU-01 | Guru | Guru aktif valid tanpa session |
| DU-GRU-02 | Guru | Guru aktif valid dengan session aktif |
| DU-GRU-03 | Guru | Guru dengan NIP unik |
| DU-GRU-04 | Guru | Guru dengan username/email/no_hp/NIP yang sudah terpakai |
| DU-SIS-01 | Siswa | Siswa aktif valid tanpa session |
| DU-SIS-02 | Siswa | Siswa aktif valid dengan session aktif |
| DU-SIS-03 | Siswa | Siswa dengan NISN unik |
| DU-SIS-04 | Siswa | Siswa dengan username/email/no_hp/NISN yang sudah terpakai |
| DU-REF-01 | Referensi kelas | `id_nama_kelas` valid tersedia |
| DU-REF-02 | Referensi kelas | `id_nama_kelas` tidak valid/tidak tersedia |

## 5. Format Response Umum

Response sukses:

```json
{
  "data": {},
  "message": "Success",
  "error": null
}
```

Response gagal:

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "error message"
  }
}
```

## 6. Skenario List dan Detail Guru

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| GRU-GET-001 | Admin melihat list guru | Login admin, akses `GET /admin/guru` | Status `200`, daftar guru tampil |
| GRU-GET-002 | Admin melihat detail guru valid | Akses `GET /admin/guru/:id` dengan id valid | Status `200`, data guru sesuai id |
| GRU-GET-003 | Filter search guru dengan `q` | Akses `GET /admin/guru?q=guru01` | Status `200`, hasil terfilter |
| GRU-GET-004 | Filter search guru dengan `search` | Akses `GET /admin/guru?search=matematika` | Status `200`, hasil terfilter |
| GRU-GET-005 | Filter status guru | Akses `GET /admin/guru?status=AKTIF` | Status `200`, hanya data status aktif tampil |
| GRU-GET-006 | Filter bidang studi | Akses `GET /admin/guru?bidang_studi=Matematika` | Status `200`, hasil sesuai filter |
| GRU-GET-007 | Pagination guru | Akses `GET /admin/guru?limit=10&offset=0` | Status `200`, jumlah hasil mengikuti parameter |
| GRU-GET-008 | Sorting guru | Akses `GET /admin/guru?sort_by=username&sort_desc=true` | Status `200`, urutan data sesuai sort |
| GRU-GET-009 | Detail guru dengan id tidak ada | Akses `GET /admin/guru/999999` | Status `404`, code `NOT_FOUND` |
| GRU-GET-010 | Detail guru dengan id bukan angka | Akses `GET /admin/guru/abc` | Status `400`, code `BAD_REQUEST` |
| GRU-GET-011 | Filter guru limit bukan angka | Akses `GET /admin/guru?limit=abc` | Status `400`, code `INVALID_INPUT` |
| GRU-GET-012 | Filter guru offset bukan angka | Akses `GET /admin/guru?offset=abc` | Status `400`, code `INVALID_INPUT` |
| GRU-GET-013 | Filter guru sort_desc bukan boolean | Akses `GET /admin/guru?sort_desc=abc` | Status `400`, code `INVALID_INPUT` |
| GRU-GET-014 | Guest mengakses list guru | Tanpa login akses endpoint | Status `401`, code `UNAUTHORIZED` |
| GRU-GET-015 | Non-admin mengakses detail guru | Login guru/siswa lalu akses endpoint admin | Status `403`, code `FORBIDDEN` |

## 7. Skenario List dan Detail Siswa

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SIS-GET-001 | Admin melihat list siswa | Login admin, akses `GET /admin/siswa` | Status `200`, daftar siswa tampil |
| SIS-GET-002 | Admin melihat detail siswa valid | Akses `GET /admin/siswa/:id` dengan id valid | Status `200`, data siswa sesuai id |
| SIS-GET-003 | Filter search siswa dengan `q` | Akses `GET /admin/siswa?q=siswa01` | Status `200`, hasil terfilter |
| SIS-GET-004 | Filter status siswa | Akses `GET /admin/siswa?status=AKTIF` | Status `200`, hasil sesuai filter |
| SIS-GET-005 | Filter angkatan | Akses `GET /admin/siswa?angkatan=2024` | Status `200`, hasil sesuai filter |
| SIS-GET-006 | Filter tingkat kelas | Akses `GET /admin/siswa?tingkat_kelas=10` | Status `200`, hasil sesuai filter |
| SIS-GET-007 | Filter id nama kelas | Akses `GET /admin/siswa?id_nama_kelas=3` | Status `200`, hasil sesuai filter |
| SIS-GET-008 | Pagination siswa | Akses `GET /admin/siswa?limit=20&offset=0` | Status `200`, hasil sesuai limit-offset |
| SIS-GET-009 | Detail siswa dengan id tidak ada | Akses `GET /admin/siswa/999999` | Status `404`, code `NOT_FOUND` |
| SIS-GET-010 | Detail siswa dengan id bukan angka | Akses `GET /admin/siswa/abc` | Status `400`, code `BAD_REQUEST` |
| SIS-GET-011 | Filter siswa angkatan bukan angka | Akses `GET /admin/siswa?angkatan=abc` | Status `400`, code `INVALID_INPUT` |
| SIS-GET-012 | Filter siswa tingkat_kelas bukan angka | Akses `GET /admin/siswa?tingkat_kelas=abc` | Status `400`, code `INVALID_INPUT` |
| SIS-GET-013 | Filter siswa id_nama_kelas bukan angka | Akses `GET /admin/siswa?id_nama_kelas=abc` | Status `400`, code `INVALID_INPUT` |
| SIS-GET-014 | Guest mengakses detail siswa | Tanpa login akses endpoint | Status `401`, code `UNAUTHORIZED` |
| SIS-GET-015 | Non-admin mengakses list siswa | Login guru/siswa lalu akses endpoint admin | Status `403`, code `FORBIDDEN` |

## 8. Skenario Happy Path Tambah Guru

Field wajib guru:

- `username`
- `password`
- `nama_lengkap`
- `jenis_kelamin`
- `nip`
- `jabatan`
- `bidang_studi`

Field opsional guru:

- `email`
- `no_hp`
- `foto_profil`

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| GRU-CRT-001 | Tambah guru dengan semua field valid | Kirim multipart lengkap termasuk foto | Status `200`, data guru baru tersimpan |
| GRU-CRT-002 | Tambah guru tanpa email dan no_hp | Kirim hanya field wajib | Status `200`, akun tetap berhasil dibuat |
| GRU-CRT-003 | Tambah guru tanpa foto | Kirim multipart tanpa `foto_profil` | Status `200`, akun berhasil dibuat |
| GRU-CRT-004 | Tambah guru dengan username bertitik | Username `guru.test` valid | Status `200` |
| GRU-CRT-005 | Tambah guru dengan username underscore | Username `guru_test` valid | Status `200` |
| GRU-CRT-006 | Tambah guru dengan gender variasi teks | Isi `jenis_kelamin` seperti `Laki-laki` atau `PRIA` | Status `200`, nilai dinormalisasi |
| GRU-CRT-007 | Tambah guru dengan NIP `-` | Isi `nip=-` | Status `200` jika rule bisnis mengizinkan placeholder |

## 9. Skenario Validasi Tambah Guru

| ID | Skenario | Input | Ekspektasi |
|---|---|---|---|
| GRU-VAL-001 | Content-Type bukan multipart | JSON atau text/plain | Status `400`, code `BAD_REQUEST` |
| GRU-VAL-002 | Multipart rusak | Body multipart invalid | Status `400`, code `BAD_REQUEST` |
| GRU-VAL-003 | Username kosong | `username=""` | Status `400`, code `BAD_REQUEST` atau `INVALID_INPUT` |
| GRU-VAL-004 | Username kurang dari 5 | `adm` | Status `400`, code `INVALID_INPUT` atau `USERNAME_LENGTH_INVALID` |
| GRU-VAL-005 | Username lebih dari 20 | 21 karakter | Status `400` |
| GRU-VAL-006 | Username karakter terlarang | `guru-01` atau `guru@01` | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-007 | Password kosong | `password=""` | Status `400` |
| GRU-VAL-008 | Password kurang dari 8 | `short` | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-009 | Password lebih dari 72 bytes | password >72 bytes | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-010 | Nama lengkap invalid | mengandung tag/script | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-011 | Jenis kelamin invalid | `unknown` | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-012 | Email invalid | `not-an-email` | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-013 | No HP invalid | karakter berbahaya atau format salah | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-014 | NIP invalid panjang | bukan `18 digit` dan bukan `-` | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-015 | Jabatan invalid | string berisi html/script | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-016 | Bidang studi invalid | string berisi html/script | Status `400`, code `INVALID_INPUT` |
| GRU-VAL-017 | File foto terlalu besar | upload >10 MB atau melebihi batas file store | Status `400`, code `FILE_TOO_LARGE` |
| GRU-VAL-018 | Body kurang field wajib | hilang `nip` atau `jabatan` | Status `400`, code `BAD_REQUEST` |

## 10. Skenario Konflik Tambah Guru

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| GRU-CFL-001 | Username guru sudah dipakai | Tambah guru dengan username existing | Status `409`, code `USERNAME_TAKEN` |
| GRU-CFL-002 | Email guru sudah dipakai | Tambah guru dengan email existing | Status `409`, code `EMAIL_TAKEN` |
| GRU-CFL-003 | No HP guru sudah dipakai | Tambah guru dengan no_hp existing | Status `409`, code `NO_HP_TAKEN` |
| GRU-CFL-004 | NIP guru sudah dipakai | Tambah guru dengan NIP existing | Status `409`, code `NIP_TAKEN` |

## 11. Skenario Happy Path Tambah Siswa

Field wajib siswa:

- `username`
- `password`
- `nama_lengkap`
- `jenis_kelamin`
- `id_nama_kelas`
- `nisn`
- `no_absen`
- `angkatan`
- `tempat_lahir`
- `tanggal_lahir`

Field opsional siswa:

- `email`
- `no_hp`
- `foto_profil`

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SIS-CRT-001 | Tambah siswa dengan semua field valid | Kirim multipart lengkap termasuk foto | Status `200`, data siswa baru tersimpan |
| SIS-CRT-002 | Tambah siswa tanpa email dan no_hp | Kirim field wajib saja | Status `200` |
| SIS-CRT-003 | Tambah siswa tanpa foto | Kirim multipart tanpa file | Status `200` |
| SIS-CRT-004 | Tambah siswa dengan NISN `-` | Isi `nisn=-` | Status `200` jika placeholder diizinkan |
| SIS-CRT-005 | Tambah siswa dengan tanggal lahir valid | Format `YYYY-MM-DD` | Status `200` |

## 12. Skenario Validasi Tambah Siswa

| ID | Skenario | Input | Ekspektasi |
|---|---|---|---|
| SIS-VAL-001 | Content-Type bukan multipart | JSON atau text/plain | Status `400`, code `BAD_REQUEST` |
| SIS-VAL-002 | Username invalid | terlalu pendek, terlalu panjang, atau karakter terlarang | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-003 | Password invalid | kosong, pendek, >72 bytes | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-004 | Nama lengkap invalid | mengandung karakter berbahaya | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-005 | Jenis kelamin invalid | nilai tidak dikenali | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-006 | Email invalid | format salah | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-007 | No HP invalid | format salah | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-008 | `id_nama_kelas` bukan angka | `abc` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-009 | `nisn` invalid | bukan `10 digit` dan bukan `-` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-010 | `no_absen` bukan angka | `abc` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-011 | `no_absen <= 0` | `0` atau `-1` | Status `400`, code `BAD_REQUEST` atau `INVALID_INPUT` |
| SIS-VAL-012 | `angkatan` bukan angka | `abc` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-013 | `angkatan` di bawah batas domain | mis. `2018` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-014 | `angkatan` di atas tahun berjalan | mis. `2099` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-015 | `tempat_lahir` invalid | mengandung html/script | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-016 | `tanggal_lahir` format salah | `31-12-2008` | Status `400`, code `INVALID_INPUT` |
| SIS-VAL-017 | Field wajib tidak lengkap | hilang `id_nama_kelas` atau `tanggal_lahir` | Status `400`, code `BAD_REQUEST` |
| SIS-VAL-018 | File foto terlalu besar | upload >10 MB atau melebihi batas store | Status `400`, code `FILE_TOO_LARGE` |

## 13. Skenario Konflik Tambah Siswa

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SIS-CFL-001 | Username siswa sudah dipakai | Tambah siswa dengan username existing | Status `409`, code `USERNAME_TAKEN` |
| SIS-CFL-002 | Email siswa sudah dipakai | Tambah siswa dengan email existing | Status `409`, code `EMAIL_TAKEN` |
| SIS-CFL-003 | No HP siswa sudah dipakai | Tambah siswa dengan no_hp existing | Status `409`, code `NO_HP_TAKEN` |
| SIS-CFL-004 | NISN siswa sudah dipakai | Tambah siswa dengan NISN existing | Status `409`, code `NISN_TAKEN` |

## 14. Skenario Happy Path Update Guru

Field guru yang bisa diubah:

- `username`
- `email`
- `nama_lengkap`
- `jenis_kelamin`
- `no_hp`
- `nip`
- `jabatan`
- `bidang_studi`
- `status_akun`
- `role`
- `foto_profil` atau `foto`

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| GRU-UPD-001 | Update satu field guru | Ubah hanya `nama_lengkap` | Status `200`, perubahan tersimpan |
| GRU-UPD-002 | Update banyak field guru | Ubah username, email, jabatan, bidang_studi | Status `200`, semua field berubah |
| GRU-UPD-003 | Update status akun guru | Ubah `status_akun=NONAKTIF` | Status `200`, status akun berubah |
| GRU-UPD-004 | Update role guru | Ubah `role=ADMIN` atau `GURU` | Status `200` bila aturan bisnis mengizinkan |
| GRU-UPD-005 | Update foto guru | Upload foto baru | Status `200`, foto baru tersimpan |
| GRU-UPD-006 | Kosongkan email opsional | Kirim `email=""` | Status `200`, email menjadi null/kosong sesuai implementasi |
| GRU-UPD-007 | Kosongkan no_hp opsional | Kirim `no_hp=""` | Status `200`, no_hp menjadi null/kosong sesuai implementasi |

## 15. Skenario Validasi dan Bad Scenario Update Guru

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| GRU-UPV-001 | Update guru tanpa field apapun | Kirim multipart tanpa field update | Status `400`, code `NO_FIELD_TO_UPDATE` |
| GRU-UPV-002 | Id guru bukan angka | `PATCH /admin/guru/abc` | Status `400`, code `INVALID_INPUT` atau `BAD_REQUEST` |
| GRU-UPV-003 | Id guru <= 0 | `PATCH /admin/guru/0` | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-004 | Content-Type bukan multipart | Kirim JSON | Status `400`, code `BAD_REQUEST` |
| GRU-UPV-005 | Username update invalid | karakter terlarang atau panjang salah | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-006 | Email update invalid | format salah | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-007 | No HP update invalid | format salah | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-008 | NIP update invalid | bukan `18 digit` dan bukan `-` | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-009 | Role update invalid | `SUPERADMIN` | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-010 | Status akun invalid | `BLOKIR` | Status `400`, code `INVALID_INPUT` |
| GRU-UPV-011 | Foto update terlalu besar | upload file besar | Status `400`, code `FILE_TOO_LARGE` |
| GRU-UPV-012 | Update guru tidak ditemukan | id valid format tapi tidak ada di DB | Status `404`, code `NOT_FOUND` |
| GRU-UPV-013 | Konflik username saat update | ubah ke username existing | Status `409`, code `USERNAME_TAKEN` |
| GRU-UPV-014 | Konflik email saat update | ubah ke email existing | Status `409`, code `EMAIL_TAKEN` |
| GRU-UPV-015 | Konflik no_hp saat update | ubah ke no_hp existing | Status `409`, code `NO_HP_TAKEN` |
| GRU-UPV-016 | Konflik NIP saat update | ubah ke NIP existing | Status `409`, code `NIP_TAKEN` |

## 16. Skenario Happy Path Update Siswa

Field siswa yang bisa diubah:

- `username`
- `email`
- `nama_lengkap`
- `jenis_kelamin`
- `no_hp`
- `status_akun`
- `role`
- `id_tingkat_kelas`
- `id_nama_kelas`
- `nisn`
- `no_absen`
- `angkatan`
- `tempat_lahir`
- `tanggal_lahir`
- `foto_profil` atau `foto`

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SIS-UPD-001 | Update satu field siswa | Ubah `nama_lengkap` | Status `200`, perubahan tersimpan |
| SIS-UPD-002 | Update banyak field siswa | Ubah username, kelas, angkatan, tempat lahir | Status `200` |
| SIS-UPD-003 | Update status akun siswa | Ubah `status_akun=NONAKTIF` | Status `200` |
| SIS-UPD-004 | Update role siswa | Ubah `role=SISWA` atau role lain yang diizinkan | Status `200` bila aturan bisnis mengizinkan |
| SIS-UPD-005 | Update foto siswa | Upload foto baru | Status `200` |
| SIS-UPD-006 | Kosongkan email opsional | Kirim `email=""` | Status `200` |
| SIS-UPD-007 | Kosongkan no_hp opsional | Kirim `no_hp=""` | Status `200` |

## 17. Skenario Validasi dan Bad Scenario Update Siswa

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SIS-UPV-001 | Update siswa tanpa field apapun | Multipart kosong | Status `400`, code `NO_FIELD_TO_UPDATE` |
| SIS-UPV-002 | Id siswa bukan angka | `PATCH /admin/siswa/abc` | Status `400`, code `INVALID_INPUT` atau `BAD_REQUEST` |
| SIS-UPV-003 | Content-Type bukan multipart | Kirim JSON | Status `400`, code `BAD_REQUEST` |
| SIS-UPV-004 | Username invalid | format salah | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-005 | Email invalid | format salah | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-006 | No HP invalid | format salah | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-007 | NISN invalid | bukan `10 digit` dan bukan `-` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-008 | No absen bukan angka | `abc` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-009 | Angkatan bukan angka | `abc` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-010 | Tanggal lahir format salah | `12/31/2008` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-011 | Role invalid | `SUPERADMIN` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-012 | Status akun invalid | `BLOKIR` | Status `400`, code `INVALID_INPUT` |
| SIS-UPV-013 | File foto terlalu besar | upload file besar | Status `400`, code `FILE_TOO_LARGE` |
| SIS-UPV-014 | Update siswa tidak ditemukan | id tidak ada | Status `404`, code `NOT_FOUND` |
| SIS-UPV-015 | Konflik username saat update | ubah ke username existing | Status `409`, code `USERNAME_TAKEN` |
| SIS-UPV-016 | Konflik email saat update | ubah ke email existing | Status `409`, code `EMAIL_TAKEN` |
| SIS-UPV-017 | Konflik no_hp saat update | ubah ke no_hp existing | Status `409`, code `NO_HP_TAKEN` |
| SIS-UPV-018 | Konflik NISN saat update | ubah ke NISN existing | Status `409`, code `NISN_TAKEN` |

## 18. Skenario Efek Session Setelah Update

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SES-UPD-001 | Update guru yang sedang login | Login sebagai guru target, lalu admin update data guru tersebut | Update `200`, seluruh session guru target direvoke |
| SES-UPD-002 | Guru target akses API setelah diupdate | Setelah SES-UPD-001, guru target panggil endpoint protected | Status `401`, harus login ulang |
| SES-UPD-003 | Update siswa yang sedang login | Login sebagai siswa target, lalu admin update data siswa tersebut | Update `200`, seluruh session siswa target direvoke |
| SES-UPD-004 | Siswa target akses API setelah diupdate | Setelah SES-UPD-003, siswa target panggil endpoint protected | Status `401`, harus login ulang |

## 19. Skenario Hapus Satu Akun

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| DEL-001 | Admin hapus satu akun guru | `DELETE /admin/pengguna/:id` dengan id guru valid | Status `200`, akun terhapus |
| DEL-002 | Admin hapus satu akun siswa | `DELETE /admin/pengguna/:id` dengan id siswa valid | Status `200`, akun terhapus |
| DEL-003 | Hapus akun id tidak ditemukan | id valid format tapi tidak ada | Status `404`, code `NOT_FOUND` |
| DEL-004 | Hapus akun id bukan angka | `/admin/pengguna/abc` | Status `400`, code `BAD_REQUEST` |
| DEL-005 | Hapus akun id <= 0 | `/admin/pengguna/0` | Status `400`, code `BAD_REQUEST` |
| DEL-006 | Guest hapus akun | Tanpa login | Status `401`, code `UNAUTHORIZED` |
| DEL-007 | Non-admin hapus akun | Login guru/siswa lalu delete | Status `403`, code `FORBIDDEN` |

## 20. Skenario Hapus Banyak Akun

Format body:

```json
{
  "ids": [1, 2, 3]
}
```

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| DELM-001 | Admin hapus banyak akun valid | Kirim beberapa id valid | Status `200`, response berisi jumlah `deleted` |
| DELM-002 | Ada id duplikat dalam body | `ids:[5,5,6]` | Status `200`, duplikat diabaikan, delete tetap jalan |
| DELM-003 | Body ids kosong | `{"ids":[]}` | Status `400`, code `BAD_REQUEST` |
| DELM-004 | Body ids tidak ada | `{}` | Status `400`, code `BAD_REQUEST` |
| DELM-005 | Salah satu id tidak valid `<=0` | `ids:[1,0,2]` | Status `400`, code `BAD_REQUEST` |
| DELM-006 | Content-Type bukan JSON | multipart atau text/plain | Status `400`, code `BAD_REQUEST` |
| DELM-007 | Body JSON rusak | malformed JSON | Status `400`, code `BAD_REQUEST` |
| DELM-008 | Semua id tidak ditemukan | kirim id yang tidak ada semua | Status `404`, code `NOT_FOUND` |
| DELM-009 | Guest hapus banyak akun | tanpa login | Status `401`, code `UNAUTHORIZED` |
| DELM-010 | Non-admin hapus banyak akun | login guru/siswa | Status `403`, code `FORBIDDEN` |

## 21. Skenario Reset Password

Format body:

```json
{
  "password": "NewPassword123"
}
```

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| RST-001 | Admin reset password guru | PUT dengan password valid | Status `200`, password baru tersimpan |
| RST-002 | Admin reset password siswa | PUT dengan password valid | Status `200`, password baru tersimpan |
| RST-003 | Password reset kosong | `{"password":""}` | Status `400`, code `BAD_REQUEST` |
| RST-004 | Password reset kurang dari 8 | `short` | Status `400`, code `BAD_REQUEST` |
| RST-005 | Password reset lebih dari 72 bytes | password terlalu panjang | Status `400`, code `BAD_REQUEST` |
| RST-006 | Id pengguna bukan angka | `/admin/pengguna/abc/reset-password` | Status `400`, code `BAD_REQUEST` |
| RST-007 | Id pengguna <= 0 | `/admin/pengguna/0/reset-password` | Status `400`, code `BAD_REQUEST` |
| RST-008 | Content-Type bukan JSON | multipart atau text/plain | Status `400`, code `BAD_REQUEST` |
| RST-009 | JSON body rusak | malformed JSON | Status `400`, code `BAD_REQUEST` |
| RST-010 | Guest reset password | tanpa login | Status `401`, code `UNAUTHORIZED` |
| RST-011 | Non-admin reset password | login guru/siswa | Status `403`, code `FORBIDDEN` |

## 22. Skenario Otorisasi

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| AUTHZ-001 | Semua endpoint kelola akun dapat diakses admin | Login admin lalu coba semua endpoint utama | Request diproses normal |
| AUTHZ-002 | Guru tidak boleh akses area kelola akun | Login guru lalu akses list/create/update/delete/reset | Status `403`, code `FORBIDDEN` |
| AUTHZ-003 | Siswa tidak boleh akses area kelola akun | Login siswa lalu akses list/create/update/delete/reset | Status `403`, code `FORBIDDEN` |
| AUTHZ-004 | Guest tidak boleh akses area kelola akun | Tanpa login akses endpoint | Status `401`, code `UNAUTHORIZED` |

## 23. Skenario UI Admin

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| UI-001 | Halaman list guru tampil | Buka halaman manajemen guru sebagai admin | Tabel/daftar guru tampil |
| UI-002 | Halaman list siswa tampil | Buka halaman manajemen siswa sebagai admin | Tabel/daftar siswa tampil |
| UI-003 | Form tambah guru submit sukses | Isi form valid lalu submit | Data tersimpan, UI memberi notifikasi sukses |
| UI-004 | Form tambah siswa submit sukses | Isi form valid lalu submit | Data tersimpan, UI memberi notifikasi sukses |
| UI-005 | Error validasi tampil ke admin | Isi form dengan input invalid | UI menampilkan pesan error yang sesuai |
| UI-006 | Error konflik tampil ke admin | Submit username/email/NIP/NISN duplikat | UI menampilkan pesan unik/duplikat |
| UI-007 | Edit akun guru sukses | Ubah data guru dari UI | Data berubah dan list/detail ikut terbarui |
| UI-008 | Edit akun siswa sukses | Ubah data siswa dari UI | Data berubah dan list/detail ikut terbarui |
| UI-009 | Hapus satu akun dari UI | Klik hapus dan konfirmasi | Data hilang dari list |
| UI-010 | Hapus banyak akun dari UI | Pilih beberapa akun lalu bulk delete | Data terhapus dan count sesuai |
| UI-011 | Reset password dari UI | Isi password baru valid | UI menampilkan sukses |
| UI-012 | User non-admin tidak melihat menu kelola akun | Login guru/siswa | Menu/hak akses admin tidak tampil |

## 24. Skenario Keamanan dan Konsistensi Data

| ID | Skenario | Langkah Uji | Ekspektasi |
|---|---|---|---|
| SEC-001 | Email disimpan lowercase | Tambah/update email `USER@MAIL.COM` | Data tersimpan dalam lowercase |
| SEC-002 | Script injection pada field teks | Isi `nama_lengkap`, `jabatan`, `tempat_lahir` dengan script/html | Status `400`, code `INVALID_INPUT` |
| SEC-003 | Tidak bisa spoof role pada create guru | Kirim field tambahan `role=ADMIN` saat create guru | Akun tetap tersimpan sebagai `GURU` |
| SEC-004 | Tidak bisa spoof role pada create siswa | Kirim field tambahan `role=ADMIN` saat create siswa | Akun tetap tersimpan sebagai `SISWA` |
| SEC-005 | Session target terputus setelah update | Update user yang sedang login | Session lama tidak lagi valid |
| SEC-006 | Hapus akun membuat detail tidak dapat diakses | Hapus akun lalu buka detailnya | Status `404`, code `NOT_FOUND` |
| SEC-007 | Password lama tidak berlaku setelah reset | Reset password lalu coba login dengan password lama | Login gagal |
| SEC-008 | Password baru berlaku setelah reset | Reset password lalu login dengan password baru | Login berhasil |

## 25. Checklist Regression Minimum

- `GRU-GET-001`, `GRU-GET-002`, `GRU-GET-014`, `GRU-GET-015`
- `SIS-GET-001`, `SIS-GET-002`, `SIS-GET-014`, `SIS-GET-015`
- `GRU-CRT-001`, `GRU-VAL-004`, `GRU-VAL-012`, `GRU-CFL-001`, `GRU-CFL-004`
- `SIS-CRT-001`, `SIS-VAL-008`, `SIS-VAL-013`, `SIS-CFL-001`, `SIS-CFL-004`
- `GRU-UPD-001`, `GRU-UPV-001`, `GRU-UPV-013`
- `SIS-UPD-001`, `SIS-UPV-001`, `SIS-UPV-018`
- `SES-UPD-001`, `SES-UPD-004`
- `DEL-001`, `DEL-003`, `DELM-001`, `DELM-003`
- `RST-001`, `RST-003`, `RST-010`
- `AUTHZ-001`, `AUTHZ-002`, `AUTHZ-004`
- `SEC-007`, `SEC-008`

## 26. Kriteria Lulus

Fitur kelola akun siswa dan guru dinyatakan lulus jika:

- Admin dapat melihat, menambah, mengubah, menghapus, dan reset password akun siswa/guru.
- Guru dan siswa tidak dapat mengakses endpoint/admin page kelola akun.
- Validasi input berjalan konsisten untuk multipart dan JSON.
- Konflik data unik terdeteksi dengan response `409` yang tepat.
- Upload foto bekerja normal dan file oversized ditolak.
- Update data user membuat session user target dicabut.
- Reset password membuat password lama tidak berlaku dan password baru dapat dipakai login.
- Endpoint list, detail, delete, bulk delete, dan reset password menangani id invalid dan not found dengan benar.

## 27. Catatan Verifikasi Lanjutan

- Perlu verifikasi apakah perubahan `role` pada update memang diizinkan oleh aturan bisnis produk, karena endpoint teknis saat ini menerima field `role`.
- Perlu verifikasi apakah update akun memang diharapkan selalu memutus semua session aktif user target. Implementasi saat ini melakukan revoke semua session setelah update berhasil.
- Perlu verifikasi perilaku jika `id_nama_kelas` valid format tetapi referensi kelas tidak ada, karena hasil akhir bisa bergantung pada validasi repository/database.
- Perlu verifikasi apakah delete user juga harus memutus session aktif user yang dihapus. Implementasi handler delete tidak secara eksplisit merevoke session.
