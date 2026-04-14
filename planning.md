# Planning Unit Test API HAndler 

Dokumen ini berisi rencana unit testing untuk semua API (Request-Response) dan non-API (Helper) di level Handler. Pendekatan yang digunakan adalah **White Box** dengan teknik **Branch Coverage**.

## Aturan & Panduan Umum

1. **Struktur Folder Test**
   Seluruh file test diletakkan di dalam folder `test` khusus pada masing-masing fitur.
   Contoh: `internal/adapter/handler/http/features/bank_soal/test/bank_soal_test.go`
2. **Library Testing**
   Menggunakan build-in `testing` Go dan package eksternal (contoh: `stretchr/testify/assert`) untuk mempermudah _assertion_.
3. **Consistency (Aturan Data)**
   Untuk memastikan _state_ yang konsisten, setiap skenario test wajib diawali dengan tahapan membersihkan mock state atau menghapus data-data _dummy_ (tear down/set up).
4. **Scope (Cakupan Uji)**
   - **Fungsi API Handler**: Tes HTTP Status Code, format JSON Response, serta berbagai skenario error handling (berdasarkan percabangan respon).
   - **Fungsi Helper**: Diuji secara independen layaknya standard function I/O (Input/Output mapping, error return) tanpa HTTP Request.

---

## Daftar Fitur & Skenario Test

Berikut adalah daftar skenario yang perlu diimplementasikan oleh _programmer_. Perhatikan direktori `routes/` untuk memastikan kesesuaian URL jika diperlukan.

### 1. Fitur Auth
*Lokasi: `internal/adapter/handler/http/auth`*

- **Handler Login**:
  - Melakukan login berhasil dengan kredensial valid (Branch: Success `200`).
  - Gagal login karena akun/email tidak ditemukan (Branch: Not Found `404` / Unauthorized `401`).
  - Gagal login karena password salah (Branch: Err Invalid Password).
  - Request login tidak lengkap/format salah (Branch: Err Validation Payload `400`).
- **Handler Logout**:
  - Logout berhasil dan mereset session/cookie.
  - Logout dengan tidak ada indikasi session terdaftar.

### 2. Fitur Users (Pengguna, Guru, Siswa)
*Lokasi: `internal/adapter/handler/http/features/user/`*

- **Handler Create (Guru/Siswa)**:
  - Berhasil membuat user dengan data yang lengkap.
  - Gagal karena payload kurang atau id/email bentrok (Branch: Duplicate `409` atau Validation Error).
- **Handler Get & List**:
  - Berhasil mengambil daftar (List) seluruh siswa/guru.
  - Berhasil mengambil detail berdasarkan ID (GetByID).
  - Gagal GetByID karena ID yang dikirim tidak ditemukan di sistem.
- **Handler Update & Delete**:
  - Berhasil memperbarui data pengguna yang ada, hapus data sukses.
  - Gagal jika ID tidak valid, payload tidak valid.
- **Handler Reset Password**:
  - Berhasil reset password user, dan gagal apabila user ID tidak ditemukan.
- **Helper Function**: 
  - Uji seluruh `get_helper.go` dan `update_helper.go` untuk memvalidasi _mapping_ logika Request ke Entity/DTO tanpa bergantung ke HTTP.

### 3. Fitur Bank Soal & Import Soal
*Lokasi: `internal/adapter/handler/http/features/bank_soal/` & `import_soal/`*

- **Handler Get Bank Soal**:
  - Berhasil mengambil daftar seluruh bank soal, beserta skenario Get List by Guru yang valid.
  - Get By ID berhasil mengembalikan struktur bank soal, dan error 404 ketika ID salah.
- **Handler Create, Update & Delete Bank Soal**:
  - Berhasil manipulasi data apabila diakses oleh Admin/Guru.
  - Mengeksekusi branch validasi _bad request_.
- **Handler Import Soal (Docx, dll)**:
  - Berhasil proses file dengan ekstensi dan MIME type yang diizinkan (Upload berhasil).
  - Branch handling: upload gagal karena format file bukan docx, file size melampaui batas, atau form error.
- **Helper Function**: 
  - Mapping logika parsing model ujian dan pembentukan DTO bank soal.

### 4. Fitur Ujian (Jadwal, Pertanyaan, Koreksi, Hasil)
*Lokasi: `internal/adapter/handler/http/features/ujian/`*

- **Handler Skenario Siswa (Attempt Ujian & Jawaban)**:
  - Sukses memulai ujian (`attempt`), gagal bila jadwal kadaluarsa/belum waktunya/password salah.
  - Sukses simpan jawaban (branch menyimpan pilihan ganda/essay).
  - Sukses _submit_ akhir ujian, tolak jika _attempt_ sudah pernah disubmit sebelumnya.
- **Handler Skenario Admin/Guru (Create, Koreksi, Statistik)**:
  - Membuat/Update ujian sukses dengan konfigurasi waktu, target bank soal dan kelas valid.
  - Get List ujian terkirim, get Statistik ujian dengan response data rekap.
  - Memperbarui poin (koreksi essay) sukses, menolak input skor diluar batas yang ada.
  - Memaksa _expire_ (Selesai paksa) _attempt_ siswa dari sisi _admin/guru_.

### 5. Fitur Kurikulum Tambahan (Kelas, Mata Pelajaran)
*Lokasi: `internal/adapter/handler/http/features/kelas/` & `mata_pelajaran/`*

- **Semua API (CRUD) pada Kelas & Mapel**:
  - Menjalankan test _CREATE_ data master berhasil dan menguji jika data tidak berpasangan (error invalid relation).
  - Menjalankan _LIST_ dan _GET By ID_.
  - Menjalankan _UPDATE_ dan _DELETE_, dipadukan skenario kegagalan bila _constraints rule / duplicate code_ terjadi di Database (direpresentasikan lewat mock).

### 6. Fitur Pengumuman
*Lokasi: `internal/adapter/handler/http/features/pengumuman/`*

- **Handler CRUD Pengumuman**:
  - List pengumuman aktif & list semua pengumuman berjalan sukses.
  - Error branch pada saat param ID atau Payload tidak valid untuk update.
- **Helper**: 
  - Pastikan parser tipe *Attachment* berjalan benar pada data lampiran/dokumen.

### 7. Fitur Profil Sekolah
*Lokasi: `internal/adapter/handler/http/features/profil_sekolah/`*

- **Handler View & Update**:
  - Mengambil Profil sekolah mengembalikan JSON format profil.
  - Memperbarui Identitas sukses (Nama Sekolah, Alamat, dsb), Branch error jika request struct tidak lolos tag validation.

### 8. Ruang Ujian & Sesi (Pelaksanaan)
*Lokasi: `internal/adapter/handler/http/features/ruang_ujian/` & `sesi/`*

- **Hanlder CRUD Sesi & Ruangan**:
  - Get List dan Get by ID, Get by Kode sukses di sisi Controller respon.
  - Skenario Create berhasil. Skenario *Failed Validator*, gagal menyimpan karena data duplikat kode ruang/Sesi.
  - Branch khusus mencari ruangan dan di-_map_ dengan sesi.

### 9. Fitur Aktivitas User
*Lokasi: `internal/adapter/handler/http/features/aktivitas_user/`*

- **Handler List & Create Log**:
  - Merekam Log melalui Payload dan memberikan `HTTP 201 Created`.
  - Memberikan Listing _Activity_ sukses tanpa error parsing pagination.

---

> _Setiap skenario wajib menyertakan assertion untuk Status Code (contoh HTTP 200, 201, 400, 404, 500), kembalian struktur Message / Body, serta memastikan bahwa skenario GAGAL (Branch Exception) tercapai pada logika handler tersebut._
