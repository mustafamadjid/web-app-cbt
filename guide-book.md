# Guide Book Sistem Ujian CBT

Dokumen ini berisi panduan penggunaan fitur-fitur yang tersedia di dalam proyek sistem ujian CBT berdasarkan implementasi aplikasi saat ini. Panduan dibagi menjadi 3 bagian utama:

1. Panduan Admin
2. Panduan Guru
3. Panduan Siswa

Dokumen ini ditulis dari sudut pandang pengguna akhir. Beberapa menu yang tampil di UI tetapi belum selesai diimplementasikan tetap dicantumkan dengan catatan status agar tidak menimbulkan salah persepsi.

## Cara Menggunakan Guide Book

- Ikuti langkah sesuai peran akun Anda.
- Gunakan urutan kerja yang disarankan agar data saling terhubung dengan benar.
- Jika sebuah fitur diberi catatan `Status implementasi terbatas` atau `Placeholder`, artinya menu tersebut ada tetapi perilakunya belum sepenuhnya final.

## Ringkasan Hak Akses

| Peran | Akses Utama                                                                                              |
| ----- | -------------------------------------------------------------------------------------------------------- |
| Admin | Kelola akun, sesi login aktif, data master, bank soal, ujian, hasil ujian, pengumuman, profil sekolah    |
| Guru  | Kelola bank soal, membuat ujian, melihat jadwal, menilai hasil essay, membuat pengumuman, melihat profil |
| Siswa | Melihat pengumuman, mengikuti ujian, mengisi token, submit jawaban, melihat hasil ujian, melihat profil  |

## Urutan Kerja yang Disarankan

Sebelum ujian dapat dipakai penuh, urutan kerja yang paling aman adalah:

1. Admin mengisi `Profil Sekolah`.
2. Admin membuat `Data Kelas`, `Mata Pelajaran`, `Ruang Ujian`, dan `Sesi`.
3. Admin membuat akun guru dan akun siswa.
4. Admin atau guru membuat `Bank Soal`.
5. Admin atau guru membuat `Ujian`.
6. Siswa mengerjakan ujian sesuai jadwal dan token.
7. Guru atau admin memeriksa `Hasil Ujian`, terutama essay yang belum dinilai.

## Catatan Penting Tentang Status Fitur

- Menu `Backup dan Restore` di halaman `Pengaturan` masih berupa placeholder.
- Aksi `Arsipkan` yang muncul di beberapa tabel belum menunjukkan implementasi aktif.
- Fitur `Reset Password` tersedia di layer service/backend, tetapi pada UI yang saya audit belum terlihat tombol/halaman aktif untuk menggunakannya.
- Halaman `Cetak` hanya berfungsi sebagai pengarah ke `Detail Ujian`; proses cetak dilakukan dari halaman detail ujian.

---

# Bagian 1: Panduan Admin

## 1. Login Admin

### Tujuan

Masuk ke dashboard administrator untuk mengelola seluruh sistem ujian.

### Langkah

1. Buka halaman login aplikasi.
2. Isi `Username`.
3. Isi `Password`.
4. Klik tombol `Masuk`.
5. Jika login berhasil, sistem mengarahkan Anda ke `Dashboard Admin`.

### Catatan

- Jika username atau password salah, aplikasi menampilkan pesan error di halaman login.
- Jika sesi sudah habis, Anda perlu login ulang.

## 2. Mengenal Dashboard Admin

### Fitur yang tampil

- Ringkasan total siswa dan guru.
- Statistik total ujian terlaksana.
- Statistik total bank soal.
- Statistik total mata pelajaran aktif.
- Papan pengumuman aktif.
- Jadwal ujian mendatang.
- Log aktivitas pengguna.

### Cara menggunakan

1. Setelah login, baca kartu statistik di bagian atas untuk melihat kondisi umum sistem.
2. Gunakan panel `Papan Pengumuman` untuk melihat pengumuman aktif.
3. Gunakan panel `Jadwal Ujian` untuk memantau ujian yang belum dimulai.
4. Gunakan panel `Log Aktivitas` untuk melihat aktivitas penting pengguna.

### Catatan

- Log aktivitas hanya muncul untuk admin.
- Jika data gagal dimuat, akan muncul notifikasi peringatan pada dashboard.

## 3. Kelola Pengguna

## 3.1 Akun Guru

### Tujuan

Menambah, melihat, memfilter, mengedit, dan menghapus akun guru.

### Cara membuka

Masuk ke menu `Kelola Pengguna` → `Akun Guru`.

### Fitur yang tersedia

- Pencarian nama, NIP, atau email.
- Filter status akun `Aktif` dan `Nonaktif`.
- Toggle samarkan/tampilkan data sensitif.
- Tambah akun guru.
- Edit akun guru.
- Hapus satu akun.
- Hapus banyak akun sekaligus.

### Langkah menambah akun guru

1. Klik `Tambah Akun Guru`.
2. Isi bagian `Informasi Dasar`:
   - Nama lengkap
   - Email
   - Username
   - Password
   - Jenis kelamin
   - Nomor HP
3. Isi bagian `Data Kepegawaian`:
   - NIP
   - Jabatan
   - Bidang studi
4. Unggah `Foto Profil` jika diperlukan.
5. Klik `Daftarkan Guru`.

### Langkah melihat dan memfilter daftar guru

1. Gunakan kolom pencarian untuk mencari guru tertentu.
2. Gunakan chip `Aktif` atau `Nonaktif` untuk memfilter status.
3. Gunakan tombol `Tampilkan Data` atau `Samarkan Data` jika ingin menampilkan/menyembunyikan data sensitif seperti NIP dan nomor HP.

### Langkah mengedit akun guru

1. Cari guru yang ingin diubah.
2. Klik ikon `Edit`.
3. Ubah data yang diperlukan.
4. Simpan perubahan.

### Langkah menghapus akun guru

1. Cari guru yang akan dihapus.
2. Klik ikon `Hapus`.
3. Baca pesan konfirmasi dengan teliti.
4. Konfirmasi penghapusan.

### Langkah menghapus banyak akun guru

1. Centang beberapa baris guru.
2. Klik tombol jumlah data terpilih.
3. Pilih `Hapus Data`.
4. Konfirmasi penghapusan.

### Catatan

- Penghapusan akun guru berpotensi ikut memengaruhi data terkait, seperti bank soal, aktivitas, ujian, dan hasil ujian.
- Tombol `Arsipkan` terlihat di menu aksi, tetapi belum menunjukkan perilaku final.

## 3.2 Akun Siswa

### Tujuan

Menambah, melihat, memfilter, mengedit, dan menghapus akun siswa.

### Cara membuka

Masuk ke menu `Kelola Pengguna` → `Akun Siswa`.

### Fitur yang tersedia

- Pencarian nama, kelas, email, atau nomor absen.
- Filter angkatan.
- Filter tingkat kelas.
- Filter jenis kelamin.
- Toggle samarkan/tampilkan data sensitif.
- Tambah akun siswa.
- Edit akun siswa.
- Hapus satu akun.
- Hapus banyak akun sekaligus.

### Langkah menambah akun siswa

1. Klik `Tambah Akun Siswa`.
2. Isi bagian `Informasi Dasar`:
   - Nama lengkap
   - Username
   - Password
   - Jenis kelamin
   - NISN
   - Email opsional
   - Nomor HP opsional
3. Isi bagian `Data Akademik`:
   - Nomor absen
   - Angkatan
   - Nama kelas
4. Isi bagian `Data Kelahiran`:
   - Tempat lahir
   - Tanggal lahir
5. Unggah `Foto Profil` jika diperlukan.
6. Klik `Tambah Siswa`.

### Langkah memfilter daftar siswa

1. Ketik kata kunci pada kolom pencarian.
2. Pilih `Angkatan`.
3. Pilih `Tingkat Kelas`.
4. Pilih `Jenis Kelamin`.
5. Gunakan `Reset Filter` jika ingin menghapus semua filter.

### Langkah mengedit akun siswa

1. Cari siswa yang akan diubah.
2. Klik ikon `Edit`.
3. Perbarui data yang diperlukan.
4. Simpan perubahan.

### Langkah menghapus akun siswa

1. Klik ikon `Hapus` pada baris siswa.
2. Konfirmasi penghapusan.

### Langkah menghapus banyak akun siswa

1. Centang beberapa siswa.
2. Klik menu aksi jumlah data terpilih.
3. Pilih `Hapus Data`.
4. Konfirmasi penghapusan.

### Catatan

- Penghapusan akun siswa berpotensi ikut menghapus data ujian, hasil ujian, dan data kartu ujian yang terkait.
- Tombol reset password belum terlihat aktif di UI yang diaudit, walaupun service backend sudah tersedia.

## 4. Kelola Sesi Login Aktif

### Tujuan

Memantau sesi login aktif dan memaksa logout akun lain jika diperlukan.

### Cara membuka

Masuk ke menu `Kelola Sesi`.

### Fitur yang tersedia

- Melihat daftar sesi login aktif.
- Melihat nama pengguna, kontak, role, status akun, dan masa berlaku sesi.
- Logout pengguna lain dari sesi aktifnya.

### Langkah menggunakan

1. Buka halaman `Kelola Sesi`.
2. Periksa daftar sesi aktif.
3. Cari akun yang ingin dikeluarkan.
4. Klik tombol `Logout`.
5. Konfirmasi tindakan.

### Catatan

- Admin tidak dapat merevoke sesi miliknya sendiri. Pada sesi milik sendiri akan muncul label `Sesi Anda`.
- Data sesi akan diperbarui setelah halaman dimuat ulang atau setelah proses revoke berhasil.

## 5. Data Master

Data master harus diisi lebih dulu sebelum pembuatan ujian berjalan lancar.

## 5.1 Data Kelas

### Tujuan

Menyusun struktur tingkat kelas dan nama kelas.

### Cara membuka

Masuk ke menu `Data Master` → `Data Kelas`.

### Fitur yang tersedia

- Pencarian kelas.
- Filter tingkat kelas.
- Tambah tingkat kelas.
- Tambah nama kelas.
- Edit data kelas.
- Hapus satu atau banyak nama kelas.

### Langkah menambah kelas

1. Klik `Tambah Kelas`.
2. Pada bagian `Tambah Tingkat Kelas`, isi tingkat kelas dalam bentuk angka, lalu simpan.
3. Pada bagian `Tambah Nama Kelas`, pilih tingkat kelas yang sudah ada.
4. Isi nama kelas.
5. Simpan data.

### Langkah mengedit kelas

1. Di daftar kelas, klik ikon `Edit`.
2. Ubah `Tingkat Kelas` bila perlu.
3. Ubah `Nama Kelas` bila perlu.
4. Klik `Simpan Perubahan`.

### Langkah menghapus kelas

1. Klik ikon `Hapus` pada kelas yang dipilih.
2. Konfirmasi penghapusan.

### Catatan

- Struktur kelas dipisah menjadi dua komponen: tingkat kelas dan nama kelas.
- Tombol `Arsipkan` belum menunjukkan implementasi final.

## 5.2 Data Mata Pelajaran

### Tujuan

Mencatat daftar mata pelajaran per tingkat kelas.

### Cara membuka

Masuk ke menu `Data Master` → `Data Mata Pelajaran`.

### Fitur yang tersedia

- Pencarian kode mapel, nama mapel, atau deskripsi.
- Filter tingkat kelas.
- Filter nama mapel.
- Tambah mata pelajaran.
- Edit mata pelajaran.
- Hapus satu atau banyak mata pelajaran.

### Langkah menambah mata pelajaran

1. Klik `Tambah Mata Pelajaran`.
2. Pilih `Tingkat Kelas`.
3. Isi `Kode Mapel`.
4. Isi `Nama Mata Pelajaran`.
5. Isi `Deskripsi Mata Pelajaran`.
6. Klik `Simpan Mata Pelajaran`.

### Langkah mengedit mata pelajaran

1. Cari mapel yang ingin diubah.
2. Klik ikon `Edit`.
3. Ubah data yang diperlukan.
4. Simpan perubahan.

### Langkah menghapus mata pelajaran

1. Klik ikon `Hapus` pada baris mapel.
2. Konfirmasi penghapusan.

## 5.3 Data Ruang Ujian

### Tujuan

Mencatat ruangan ujian yang dipakai saat penjadwalan.

### Cara membuka

Masuk ke menu `Data Master` → `Data Ruang Ujian`.

### Fitur yang tersedia

- Pencarian nama ruang atau kode ruang.
- Tambah ruang.
- Edit ruang.
- Hapus ruang.

### Langkah menambah ruang

1. Klik `Tambah Ruang`.
2. Isi `Kode Ruang`.
3. Isi `Nama Ruangan`.
4. Klik `Simpan Ruangan Ujian`.

### Langkah mengedit ruang

1. Klik ikon `Edit`.
2. Perbarui kode atau nama ruangan.
3. Simpan perubahan.

### Langkah menghapus ruang

1. Klik ikon `Hapus`.
2. Konfirmasi penghapusan.

## 5.4 Data Sesi

### Tujuan

Menyimpan sesi ujian seperti sesi pagi, siang, atau sesi khusus lain.

### Cara membuka

Masuk ke menu `Data Master` → `Data Sesi`.

### Fitur yang tersedia

- Pencarian kode atau nama sesi.
- Tambah sesi.
- Edit sesi.
- Hapus sesi.

### Langkah menambah sesi

1. Klik `Tambah Sesi`.
2. Isi `Kode Sesi`.
3. Isi `Nama Sesi`.
4. Klik `Simpan Sesi`.

### Langkah mengedit sesi

1. Klik ikon `Edit`.
2. Perbarui data sesi.
3. Simpan perubahan.

### Langkah menghapus sesi

1. Klik ikon `Hapus`.
2. Konfirmasi penghapusan.

## 6. Bank Soal

### Tujuan

Membuat, memfilter, melihat preview, mengimpor, dan menghapus bank soal.

### Cara membuka

Masuk ke menu `Bank Soal`.

### Fitur yang tersedia

- Tab `Semua Bank Soal`.
- Tab `Soal Saya`.
- Filter kelas.
- Filter mata pelajaran.
- Filter status upload.
- Pencarian nama bank soal, materi, atau deskripsi.
- Buat bank soal.
- Preview bank soal.
- Upload/import soal dari file DOCX.
- Hapus bank soal.

## 6.1 Membuat bank soal

1. Klik `Buat Bank Soal`.
2. Isi:
   - Nama bank soal
   - Tingkat kelas
   - Mata pelajaran
   - Materi
   - Deskripsi
3. Klik `Simpan Bank Soal`.

## 6.2 Melihat dan memfilter bank soal

1. Gunakan tab `Semua Bank Soal` untuk melihat seluruh bank soal.
2. Gunakan tab `Soal Saya` untuk melihat bank soal milik pembuat saat ini.
3. Gunakan pencarian untuk mencari bank soal tertentu.
4. Gunakan filter `Kelas`, `Mapel`, dan `Status Upload`.
5. Gunakan `Reset Filter` untuk mengembalikan tampilan awal.

## 6.3 Preview isi bank soal

1. Klik aksi `Kelola` atau `Preview` pada kartu bank soal.
2. Sistem membuka halaman preview soal.
3. Gunakan navigator nomor soal untuk berpindah antar soal.
4. Gunakan opsi `Acak Soal` jika ingin melihat urutan acak pada preview.

### Catatan

- Preview mendukung soal pilihan ganda dan essay.
- Jika bank soal belum memiliki pertanyaan, akan muncul informasi bahwa bank soal masih kosong.

## 6.4 Upload/import soal dari file DOCX

1. Dari daftar bank soal, klik aksi `Upload` atau `Tambah Soal`.
2. Pilih file `.docx`.
3. Pastikan ukuran file tidak lebih dari 20 MB.
4. Klik `Submit Soal`.
5. Tunggu proses upload dan import.
6. Sistem akan melakukan polling status import secara otomatis.
7. Jika import selesai, sistem menampilkan jumlah soal atau peringatan jika ada.

### Catatan

- Halaman upload menerima file DOCX saja.
- Jika job import gagal, aplikasi menampilkan pesan error.
- Status upload dapat dipakai sebagai filter pada halaman daftar bank soal.

## 6.5 Menghapus bank soal

1. Pilih bank soal yang ingin dihapus.
2. Klik aksi `Hapus`.
3. Konfirmasi penghapusan.

## 7. Ujian

Menu `Ujian` pada admin terdiri dari:

- `Buat Ujian`
- `Jadwal Ujian`
- `Hasil Ujian`

## 7.1 Membuat Ujian

### Tujuan

Membuat jadwal ujian dari bank soal yang sudah tersedia.

### Prasyarat

- Tingkat kelas dan nama kelas sudah ada.
- Mata pelajaran sudah ada.
- Bank soal sudah ada.
- Ruang ujian sudah ada.
- Sesi sudah ada.
- Guru pengawas sudah ada.

### Langkah

1. Buka menu `Ujian` → `Buat Ujian`.
2. Isi bagian `Informasi Ujian`:
   - Nama ujian
   - Deskripsi ujian
3. Isi bagian `Kelas & Bank Soal`:
   - Pilih tingkat kelas
   - Pilih cakupan kelas:
     - `Semua kelas di tingkat ini`
     - `Spesifik nama kelas`
   - Jika memilih spesifik, pilih nama kelas
   - Pilih mapel
   - Pilih bank soal
4. Isi bagian `Jadwal Ujian`:
   - Tanggal ujian
   - Waktu mulai
   - Waktu selesai
   - Periksa durasi otomatis
5. Isi bagian `Ruang, Sesi, Pengawas`:
   - Ruang ujian
   - Sesi ujian
   - Guru pengawas
6. Isi bagian `Keamanan & Token`:
   - Tentukan apakah soal diacak
   - Isi token ujian
7. Periksa `Preview Daftar Siswa` untuk memastikan peserta sudah sesuai.
8. Klik `Simpan Ujian`.

### Catatan

- Token ujian maksimal 30 karakter.
- Preview siswa sangat penting untuk memastikan ujian tidak salah sasaran.
- Jika kelas tidak cocok atau siswa belum ada, preview bisa kosong.

## 7.2 Melihat Jadwal Ujian

### Tujuan

Melihat dan memantau ujian berdasarkan status pelaksanaannya.

### Langkah

1. Buka menu `Ujian` → `Jadwal Ujian`.
2. Gunakan filter:
   - Cari jadwal
   - Tanggal
   - Tahun
   - Kelas
   - Ruangan
3. Gunakan tab status:
   - `Berlangsung`
   - `Belum Mulai`
   - `Selesai`
   - `Dibatalkan`
4. Klik kartu jadwal untuk membuka detail ujian.

## 7.3 Membuka Detail Ujian

### Informasi yang dapat dilihat

- Nama ujian
- Status ujian
- Deskripsi ujian
- Tanggal
- Jam mulai dan selesai
- Durasi
- Ruangan
- Sesi
- Pengawas
- Target kelas
- Token akses

### Langkah mencetak dokumen ujian

1. Buka `Detail Ujian`.
2. Klik salah satu tombol cetak:
   - `Daftar Hadir`
   - `Berita Acara`
   - `Kartu Peserta`
3. Sistem memanggil fungsi cetak browser.
4. Pilih printer atau simpan sebagai PDF sesuai kebutuhan.

### Catatan

- Menu `Cetak` di sidebar hanya mengarahkan admin ke area jadwal/detail ujian.
- Fungsi cetak saat ini terpusat di halaman detail ujian.

## 7.4 Menghapus Jadwal Ujian

1. Buka halaman `Jadwal Ujian`.
2. Cari ujian yang ingin dihapus.
3. Klik aksi hapus pada kartu jadwal.
4. Konfirmasi penghapusan.

## 8. Hasil Ujian

### Tujuan

Memantau hasil ujian yang selesai dan memproses penilaian essay.

### Cara membuka

Masuk ke menu `Ujian` → `Hasil Ujian`.

### Tab yang tersedia

- `Ujian Selesai`
- `Essay Belum Dinilai`

## 8.1 Melihat daftar hasil ujian selesai

1. Buka tab `Ujian Selesai`.
2. Gunakan filter:
   - Kelas
   - Tahun
   - Bulan
   - Tanggal
3. Klik kartu hasil untuk membuka detail hasil ujian.

## 8.2 Melihat daftar essay belum dinilai

1. Buka tab `Essay Belum Dinilai`.
2. Gunakan filter yang sama jika perlu.
3. Gunakan daftar ini untuk memprioritaskan koreksi manual essay.

## 8.3 Membuka Detail Hasil Ujian

### Informasi yang tersedia

- Nilai tertinggi
- Rata-rata kelas
- Nilai terendah
- Total peserta
- Daftar peserta yang sudah submit
- Waktu mulai dan waktu submit peserta
- Nilai akhir tiap peserta

### Langkah

1. Klik salah satu kartu hasil ujian.
2. Tinjau statistik di bagian atas.
3. Pada tabel peserta submit, pilih peserta tertentu.
4. Klik `Lihat Jawaban`.

## 8.4 Melihat Jawaban Siswa dan Koreksi Essay

### Langkah

1. Dari detail hasil ujian, klik `Lihat Jawaban`.
2. Tinjau semua jawaban siswa.
3. Untuk soal essay yang belum dinilai, isi atau sesuaikan koreksi.
4. Simpan koreksi essay.
5. Setelah berhasil, muat ulang data untuk memastikan nilai terbarui.

### Catatan

- Soal pilihan ganda dinilai otomatis.
- Soal essay memerlukan koreksi manual dari guru atau admin.

## 9. Pengumuman

### Tujuan

Membuat dan mengelola pengumuman sekolah yang akan dilihat pengguna.

### Cara membuka

Masuk ke menu `Pengumuman`.

### Status daftar pengumuman

- `Akan Rilis`
- `Sedang Rilis`
- `Sudah Rilis`

### Fitur yang tersedia

- Tambah pengumuman
- Edit pengumuman pada status yang masih bisa diubah
- Hapus pengumuman
- Melihat periode pengumuman
- Membuka dokumen lampiran

## 9.1 Membuat pengumuman

1. Klik `Tambah Pengumuman`.
2. Isi:
   - Judul pengumuman
   - Isi pengumuman
   - Tanggal rilis
   - Tanggal selesai
3. Unggah dokumen jika diperlukan.
4. Klik `Simpan Pengumuman`.

### Aturan lampiran

- Format yang diterima: `PDF` dan `DOCX`
- Ukuran maksimal: `10 MB`

## 9.2 Mengelola daftar pengumuman

1. Pindah antar tab status sesuai kebutuhan.
2. Klik nama dokumen untuk membuka lampiran.
3. Klik ikon `Edit` untuk memperbarui pengumuman.
4. Klik ikon `Hapus` untuk menghapus pengumuman.

### Catatan

- Pengumuman yang sudah berada pada kategori `Sudah Rilis` tidak dapat diedit dari UI yang diaudit.
- Pada implementasi saat ini, isi pengumuman ditulis menggunakan textarea standar, bukan editor rich text penuh.

## 10. Pengaturan

## 10.1 Profil Sekolah

### Tujuan

Mengisi identitas sekolah yang digunakan sistem.

### Cara membuka

Masuk ke menu `Pengaturan` → tab `Profil Sekolah`.

### Langkah

1. Isi atau ubah:
   - Nama sekolah
   - Alamat sekolah
   - Nomor telepon sekolah
   - Email sekolah
   - Kepala sekolah
   - Waka sekolah
2. Unggah logo sekolah jika diperlukan.
3. Klik `Simpan Perubahan`.

### Aturan logo

- File harus berupa gambar.
- Ukuran maksimal 2 MB.

## 10.2 Backup dan Restore

### Status implementasi

`Placeholder`

### Catatan

- Tab ini sudah tampil di menu `Pengaturan`, tetapi pada implementasi saat ini baru menampilkan teks placeholder dan belum memiliki alur kerja aktif.

## 11. Profil Admin

### Tujuan

Melihat informasi profil akun admin.

### Data yang dapat dilihat

- Nama lengkap
- Username
- Email
- Nomor HP
- Jenis kelamin
- Status akun
- Informasi role
- Foto profil jika tersedia

### Catatan

- Halaman profil yang diaudit bersifat tampilan data, bukan form edit.

---

# Bagian 2: Panduan Guru

## 1. Login Guru

### Langkah

1. Buka halaman login.
2. Masukkan username dan password guru.
3. Klik `Masuk`.
4. Sistem mengarahkan Anda ke dashboard guru.

## 2. Dashboard Guru

### Fitur yang tampil

- Statistik total ujian terlaksana
- Total bank soal
- Total mata pelajaran aktif
- Papan pengumuman
- Jadwal ujian mendatang

### Catatan

- Guru tidak melihat log aktivitas pengguna.

## 3. Bank Soal Guru

### Cara membuka

Masuk ke menu `Bank Soal`.

### Fitur yang tersedia

- Melihat `Semua Bank Soal`
- Melihat `Soal Saya`
- Filter kelas, mapel, status upload
- Pencarian
- Buat bank soal
- Preview bank soal
- Import soal DOCX
- Hapus bank soal

### Langkah kerja yang disarankan

1. Buat bank soal baru terlebih dahulu.
2. Import atau lengkapi isi soal.
3. Gunakan tab `Soal Saya` untuk fokus pada bank soal buatan Anda.
4. Gunakan preview untuk memeriksa urutan dan tampilan soal.

### Cara membuat bank soal

1. Klik `Buat Bank Soal`.
2. Isi nama bank soal, kelas, mapel, materi, dan deskripsi.
3. Simpan.

### Cara import soal DOCX

1. Pada bank soal yang sudah dibuat, klik aksi `Upload`.
2. Pilih file DOCX.
3. Klik `Submit Soal`.
4. Tunggu sampai sistem menyatakan import selesai.

## 4. Ujian Guru

Menu `Ujian` untuk guru terdiri dari:

- `Buat Ujian`
- `Jadwal Ujian`
- `Hasil Ujian`

## 4.1 Membuat ujian

Langkahnya sama dengan admin:

1. Isi informasi ujian.
2. Pilih kelas, mapel, dan bank soal.
3. Tentukan jadwal.
4. Pilih ruang, sesi, dan pengawas.
5. Atur token dan opsi acak soal.
6. Periksa preview daftar siswa.
7. Simpan ujian.

### Catatan

- Guru hanya bisa bekerja dengan akun dan data yang sesuai hak aksesnya.

## 4.2 Melihat jadwal ujian

### Langkah

1. Buka `Jadwal Ujian`.
2. Gunakan filter pencarian, tanggal, tahun, kelas, dan ruang.
3. Pindah antar tab status ujian.
4. Klik jadwal untuk membuka detail.

### Catatan

- Guru dapat mengontrol ujian yang dibuatnya atau ujian yang diawasinya.

## 4.3 Melihat detail ujian dan mencetak dokumen

Guru dapat melihat:

- Deskripsi ujian
- Waktu dan durasi
- Ruang dan sesi
- Pengawas
- Target kelas
- Token akses

Guru juga dapat mencetak:

- Daftar hadir
- Berita acara
- Kartu peserta

## 4.4 Menghapus ujian

1. Buka `Jadwal Ujian`.
2. Pilih jadwal yang bisa Anda kontrol.
3. Klik aksi hapus.
4. Konfirmasi penghapusan.

## 5. Hasil Ujian Guru

### Fitur yang tersedia

- Melihat ujian selesai
- Melihat daftar essay yang belum dinilai
- Melihat statistik hasil ujian
- Membuka jawaban siswa
- Menilai essay

### Langkah menilai essay

1. Masuk ke `Hasil Ujian`.
2. Buka tab `Essay Belum Dinilai`.
3. Pilih ujian yang ingin diperiksa.
4. Klik `Lihat Jawaban` pada peserta.
5. Koreksi jawaban essay.
6. Simpan koreksi.

## 6. Pengumuman Guru

### Tujuan

Guru dapat membuat dan mengelola pengumuman dari menu `Pengumuman`.

### Langkah

1. Buka menu `Pengumuman`.
2. Pilih status daftar yang ingin dilihat.
3. Klik `Tambah Pengumuman` untuk membuat pengumuman baru.
4. Isi judul, isi, periode, dan lampiran jika perlu.
5. Simpan.
6. Gunakan ikon edit atau hapus untuk pengelolaan lanjutan.

## 7. Profil Guru

Guru dapat membuka halaman profil untuk melihat:

- Identitas pribadi
- Username
- Email
- Nomor HP
- Jenis kelamin
- Status akun
- NIP
- Jabatan
- Bidang studi

### Catatan

- Halaman profil yang diaudit bersifat tampilan data.

## 8. Batasan Hak Akses Guru

- Guru tidak memiliki akses ke kelola akun pengguna.
- Guru tidak memiliki akses ke data master.
- Guru tidak memiliki akses ke kelola sesi login aktif.
- Guru tidak memiliki akses ke pengaturan profil sekolah.

---

# Bagian 3: Panduan Siswa

## 1. Login Siswa

### Langkah

1. Buka halaman login.
2. Masukkan username dan password siswa.
3. Klik `Masuk`.
4. Sistem mengarahkan Anda ke dashboard siswa.

## 2. Dashboard Siswa

### Fitur yang tersedia

- Sapaan pengguna
- Pengumuman terbaru

### Cara menggunakan

1. Setelah login, baca pengumuman yang tampil di dashboard.
2. Periksa informasi yang berkaitan dengan jadwal atau pelaksanaan ujian.

## 3. Melihat Daftar Ujian

### Cara membuka

Masuk ke menu `Ujian`.

### Kategori yang tersedia

- `Ujian Mendatang`
- `Ujian Berlangsung`

### Langkah

1. Pilih kategori `Ujian Mendatang` untuk melihat ujian yang belum dimulai.
2. Pilih kategori `Ujian Berlangsung` untuk melihat ujian yang bisa dikerjakan saat ini.
3. Baca informasi pada kartu ujian:
   - Nama ujian
   - Mapel
   - Tanggal
   - Waktu
   - Ruang
   - Pengawas

## 4. Memulai Ujian

### Alur umum

1. Dari kategori `Ujian Berlangsung`, klik `Mulai Sekarang`.
2. Sistem mencoba membuka sesi ujian aktif.
3. Jika belum ada attempt aktif, sistem akan mengarahkan ke halaman token.
4. Jika attempt aktif sudah ada, sistem bisa langsung membuka halaman pengerjaan ujian.

## 5. Memasukkan Token Ujian

### Tujuan

Memverifikasi hak siswa untuk mulai mengerjakan ujian.

### Langkah

1. Masuk ke halaman token.
2. Isi `Token Ujian`.
3. Klik `Mulai Ujian`.

### Kemungkinan pesan error

- Token ujian salah
- Anda tidak diizinkan mengikuti ujian ini
- Ujian telah selesai
- Anda telah mengikuti ujian ini

## 6. Mengerjakan Ujian

### Fitur yang tersedia

- Soal pilihan ganda
- Soal essay
- Navigator nomor soal
- Penyimpanan jawaban
- Countdown waktu
- Preview sebelum submit
- Submit final
- Penanganan waktu habis

### Langkah mengerjakan

1. Baca soal pada panel utama.
2. Untuk pilihan ganda, pilih satu jawaban yang tersedia.
3. Untuk essay, ketik jawaban pada area essay.
4. Gunakan navigator nomor soal untuk berpindah antar soal.
5. Gunakan tombol navigasi berikutnya atau sebelumnya jika tersedia.
6. Perhatikan `Sisa Waktu` di halaman ujian.
7. Jika sudah selesai, klik tombol submit ujian untuk membuka preview submit.

### Preview submit

1. Pada mode preview, periksa kembali soal yang sudah dan belum dijawab.
2. Jika masih ingin mengubah jawaban, kembali ke mode soal.
3. Jika sudah yakin, klik submit final.

### Saat keluar dari halaman ujian

Jika Anda mencoba meninggalkan halaman ujian:

1. Sistem menampilkan konfirmasi.
2. Jawaban saat ini akan dicoba disimpan terlebih dahulu.
3. Jika Anda tetap keluar, sesi ujian bisa berakhir.

### Saat waktu habis

1. Sistem menampilkan peringatan bahwa waktu ujian habis.
2. Jawaban terakhir akan dicoba disimpan.
3. Anda perlu menekan `Submit`.
4. Sistem mengakhiri attempt ujian.

## 7. Melihat Hasil Ujian

### Cara membuka

Masuk ke menu `Hasil Ujian`.

### Fitur yang tersedia

- Daftar ujian yang sudah disubmit
- Detail hasil ujian
- Nilai akhir
- Detail semua jawaban

### Langkah

1. Buka menu `Hasil Ujian`.
2. Pilih ujian yang ingin dilihat.
3. Klik `Lihat Detail`.
4. Baca informasi nilai dan hasil jawaban.

### Catatan

- Halaman ini menampilkan ujian yang sudah berhasil disubmit.
- Jika belum pernah submit ujian, daftar hasil akan kosong.

## 8. Melihat Detail Hasil Jawaban

### Informasi yang tersedia

- Nilai akhir
- Daftar soal
- Jawaban siswa
- Status penilaian tiap soal

### Langkah

1. Dari daftar hasil ujian, klik `Lihat Detail`.
2. Tinjau satu per satu hasil jawaban.
3. Gunakan tombol `Refresh` jika ingin memuat ulang hasil terbaru.

## 9. Profil Siswa

Siswa dapat membuka halaman profil untuk melihat:

- Nama lengkap
- Username
- Email
- Nomor HP
- Jenis kelamin
- Status akun
- Nomor absen
- Angkatan
- Tempat lahir
- Tanggal lahir
- Tingkat kelas
- Nama kelas

### Catatan

- Halaman profil saat ini bersifat tampilan data.

---

# Ringkasan Keterbatasan Implementasi Saat Ini

## Fitur yang ada tetapi belum final

- `Backup dan Restore`: masih placeholder.
- `Arsipkan` pada beberapa tabel: belum terlihat memiliki efek kerja final.
- `Reset Password` pengguna: service ada, tetapi tombol UI aktif belum terlihat saat audit.

## Fitur yang berjalan melalui halaman lain

- `Cetak`: tidak dikelola dari halaman cetak khusus; proses cetak dilakukan dari `Detail Ujian`.

## Fitur yang berbeda dari deskripsi umum proyek

- Pengumuman pada UI yang diaudit memakai input teks biasa dan lampiran dokumen opsional.
- Import soal yang benar-benar tampak di UI saat ini menggunakan file `DOCX`.

---

# Penutup

Jika sistem akan dipakai secara operasional, peran admin sebaiknya selalu memastikan data master, akun, bank soal, dan jadwal ujian sudah benar sebelum hari pelaksanaan. Guru sebaiknya fokus pada kualitas bank soal, validitas jadwal, dan koreksi essay. Siswa sebaiknya memeriksa pengumuman, memastikan token benar, dan menjaga kestabilan koneksi saat mengerjakan ujian.

Dokumen ini dapat diperbarui kembali jika ada penambahan fitur, perubahan alur UI, atau aktivasi modul yang saat ini masih placeholder.
