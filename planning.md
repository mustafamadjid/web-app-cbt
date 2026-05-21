# Planning Automation Testing Selenium JavaScript

Dokumen ini adalah rencana kerja untuk membuat script automation testing Selenium memakai bahasa JavaScript untuk sistem SMAFI CBT.

Acuan wajib:

- Semua skenario test mengikuti `dokumentasi_black_box_testing.xlsx`.
- Total acuan dari dokumen Excel: 48 kebutuhan fungsional dan 144 skenario black box.
- Bahasa script test: JavaScript, bukan TypeScript.
- Target aplikasi: frontend React/Vite di folder `frontend`.
- Test runner yang disarankan: Mocha.
- Browser automation: Selenium WebDriver.

Dokumen ini ditulis sebagai instruksi implementasi yang eksplisit agar mudah diterapkan oleh model dengan reasoning lebih rendah.

## 1. Tujuan Akhir

Hasil akhir yang harus dibuat pada tahap implementasi:

- Folder test E2E Selenium di `frontend/tests/e2e`.
- Script npm untuk menjalankan test Selenium.
- File konfigurasi environment untuk base URL, kredensial, browser, timeout, dan mode headless.
- Helper Selenium untuk login, logout, navigasi, isi form, assert toast/error, table assertion, upload file, download file, dan screenshot saat gagal.
- Page Object untuk halaman utama yang dipakai di test.
- Spec test JavaScript yang mencakup seluruh skenario dari `dokumentasi_black_box_testing.xlsx`.
- Dokumentasi cara menjalankan test secara lokal.

## 2. Stack Testing yang Dipakai

Gunakan dependency berikut di folder `frontend`:

```bash
npm install --save-dev selenium-webdriver mocha chai dotenv chromedriver
```

Tambahkan script berikut ke `frontend/package.json`:

```json
{
  "scripts": {
    "test:e2e": "mocha \"tests/e2e/specs/**/*.spec.js\" --timeout 90000 --exit",
    "test:e2e:auth": "mocha \"tests/e2e/specs/01-auth.spec.js\" --timeout 90000 --exit",
    "test:e2e:master": "mocha \"tests/e2e/specs/03-master-data.spec.js\" --timeout 90000 --exit",
    "test:e2e:exam": "mocha \"tests/e2e/specs/06-exam-management.spec.js\" --timeout 120000 --exit"
  }
}
```

Catatan:

- Project frontend memakai `"type": "module"`, jadi semua file test JavaScript harus memakai `import` dan `export`.
- Gunakan `chromedriver` agar setup Selenium di Windows lebih mudah.
- Jalankan frontend sebelum test: `npm run dev`.
- Jalankan backend dan database sebelum test.

## 3. Struktur Folder yang Harus Dibuat

Buat struktur berikut:

```text
frontend/
  tests/
    e2e/
      config/
        env.js
        browser.js
        routes.js
      data/
        accounts.js
        fixtures.js
        expected-text.js
      helpers/
        auth.helper.js
        navigation.helper.js
        form.helper.js
        table.helper.js
        assertion.helper.js
        api.helper.js
        download.helper.js
        screenshot.helper.js
        wait.helper.js
      pages/
        LoginPage.js
        DashboardPage.js
        UserManagementPage.js
        MasterDataPage.js
        ProfileSettingsPage.js
        BankSoalPage.js
        ExamSchedulePage.js
        ExamStudentPage.js
        ResultPage.js
        AnnouncementPage.js
        PrintPage.js
        ActivityLogPage.js
      specs/
        01-auth.spec.js
        02-user-management.spec.js
        03-master-data.spec.js
        04-dashboard-settings-log.spec.js
        05-bank-soal.spec.js
        06-exam-management.spec.js
        07-student-exam-flow.spec.js
        08-result-grading-review.spec.js
        09-print-backup-destructive.spec.js
      artifacts/
        screenshots/
        downloads/
```

Jangan campur test Selenium dengan source aplikasi di `frontend/src`.

## 4. File Environment

Buat file contoh `frontend/.env.e2e.example`.

Isi minimal:

```env
E2E_BASE_URL=http://localhost:5173
E2E_API_URL=http://localhost:8080
E2E_BROWSER=chrome
E2E_HEADLESS=true
E2E_TIMEOUT=15000

E2E_ADMIN_USERNAME=admin_e2e
E2E_ADMIN_PASSWORD=Password123!

E2E_GURU_USERNAME=guru_e2e
E2E_GURU_PASSWORD=Password123!

E2E_SISWA_USERNAME=siswa_e2e
E2E_SISWA_PASSWORD=Password123!

E2E_SISWA_2_USERNAME=siswa2_e2e
E2E_SISWA_2_PASSWORD=Password123!
```

Instruksi implementasi:

- Jangan commit file `.env.e2e` yang berisi kredensial nyata.
- Commit `.env.e2e.example`.
- `env.js` harus load `.env.e2e` memakai `dotenv`.
- Jika value tidak ada, pakai default aman untuk local.

## 5. Route Aplikasi yang Sudah Terlihat di Frontend

Gunakan route berikut sebagai baseline di `tests/e2e/config/routes.js`:

```js
export const routes = {
  login: "/login",
  adminDashboard: "/dashboard/administrator",
  guruDashboard: "/dashboard/guru",
  siswaDashboard: "/dashboard/siswa",

  adminGuru: "/dashboard/administrator/kelola-akun/guru",
  adminSiswa: "/dashboard/administrator/kelola-akun/siswa",
  tambahGuru: "/dashboard/administrator/kelola-akun/tambah-guru",
  tambahSiswa: "/dashboard/administrator/kelola-akun/tambah-siswa",

  mapel: "/dashboard/administrator/data-master/mapel",
  kelas: "/dashboard/administrator/data-master/kelas",
  ruang: "/dashboard/administrator/data-master/ruang",
  sesi: "/dashboard/administrator/data-master/sesi",

  bankSoalAdmin: "/dashboard/administrator/bank-soal",
  bankSoalGuru: "/dashboard/guru/bank-soal",

  jadwalUjianAdmin: "/dashboard/administrator/ujian/jadwal",
  jadwalUjianGuru: "/dashboard/guru/ujian/jadwal",
  buatUjianAdmin: "/dashboard/administrator/ujian/buat-ujian",
  buatUjianGuru: "/dashboard/guru/ujian/buat-ujian",
  hasilUjianAdmin: "/dashboard/administrator/ujian/hasil",
  hasilUjianGuru: "/dashboard/guru/ujian/hasil",

  pengumumanAdmin: "/dashboard/administrator/pengumuman",
  pengumumanGuru: "/dashboard/guru/pengumuman",
  pengumumanSiswa: "/dashboard/siswa/pengumuman",

  cetakAdmin: "/dashboard/administrator/cetak",
  cetakGuru: "/dashboard/guru/cetak",
  pengaturan: "/dashboard/administrator/pengaturan",

  ujianSiswa: "/dashboard/siswa/ujian",
  hasilUjianSiswa: "/dashboard/siswa/ujian/hasil"
};
```

Jika ada skenario yang membutuhkan route yang belum terlihat, lakukan salah satu:

- Cari route di `frontend/src/routes/paths.ts`.
- Cari menu di `frontend/src/layouts/MainLayout/Sidebar/sidebarMenuItems.tsx`.
- Jika fitur memang belum ada di UI, tulis test sebagai `it.skip(...)` dengan alasan jelas.

## 6. Strategi Selector

Prioritas selector Selenium:

1. `By.css('[data-testid="..."]')`
2. `By.id("...")`
3. `By.name("...")`
4. `By.css('input[placeholder="..."]')`
5. XPath berbasis label atau text button sebagai pilihan terakhir.

Instruksi penting:

- Jika element penting belum punya `data-testid`, tambahkan `data-testid` kecil pada source UI.
- Jangan memakai selector CSS yang bergantung pada class Tailwind panjang.
- Jangan memakai XPath rapuh yang bergantung pada posisi seperti `(//button)[3]` kecuali tidak ada pilihan lain.

Konvensi `data-testid` yang disarankan:

```text
login-username-input
login-password-input
login-submit-button
logout-button
page-title
search-input
filter-select
add-button
save-button
cancel-button
delete-button
confirm-button
toast-success
toast-error
validation-message
data-table
table-row
empty-state
```

## 7. Strategi Data Test

Gunakan data unik dengan timestamp agar test dapat dijalankan berulang.

Contoh helper di `fixtures.js`:

```js
export function unique(prefix) {
  return `${prefix}_${Date.now()}`;
}
```

Data dasar yang perlu tersedia sebelum semua test:

- 1 akun admin aktif.
- 1 akun guru aktif.
- 2 akun siswa aktif pada kelas yang sama untuk skenario pengacakan dan login multi perangkat.
- 1 akun siswa tambahan pada kelas berbeda untuk skenario ujian kelas lain.
- Minimal 2 mata pelajaran.
- Minimal 2 kelas.
- Minimal 2 ruang ujian.
- Minimal 2 sesi.
- Minimal 2 bank soal:
  - bank soal pilihan ganda.
  - bank soal campuran pilihan ganda dan essay.
- Minimal 3 ujian:
  - ujian terjadwal belum mulai.
  - ujian aktif yang bisa diikuti siswa.
  - ujian selesai dengan hasil.

Aturan data setup:

- Untuk test black box UI, aksi utama tetap dilakukan dari UI.
- Untuk mempercepat prasyarat, boleh memakai API helper pada `before()` selama bukan aksi utama yang sedang diuji.
- Jika API seed dipakai, beri nama fungsi jelas, misalnya `ensureMapelExists()`, `ensureExamWithEssayExists()`.
- Jangan memakai data produksi.

## 8. Pola File Test

Setiap spec mengikuti pola ini:

```js
import { expect } from "chai";
import { createDriver } from "../config/browser.js";
import { loginAs, logout } from "../helpers/auth.helper.js";

describe("F-04 Autentikasi dan Otorisasi", function () {
  let driver;

  before(async function () {
    driver = await createDriver();
  });

  after(async function () {
    if (driver) await driver.quit();
  });

  afterEach(async function () {
    if (this.currentTest.state === "failed") {
      await takeScreenshot(driver, this.currentTest.title);
    }
  });

  it("[F-04-S01] Login berhasil sesuai role", async function () {
    await loginAs(driver, "admin");
    await expectCurrentUrlContains(driver, "/dashboard/administrator");
  });
});
```

Aturan nama test:

- Semua `it()` wajib diawali ID skenario dari Excel.
- Format: `it("[F-04-S01] Login berhasil sesuai role", async function () { ... })`.
- Jika test belum bisa diimplementasikan karena fitur UI tidak ditemukan, gunakan `it.skip("[F-15-S01] ... - route backup belum ditemukan", async function () {})`.

## 9. Helper yang Wajib Dibuat

### 9.1 `browser.js`

Tanggung jawab:

- Membuat instance Chrome.
- Mengaktifkan headless berdasarkan `.env.e2e`.
- Menentukan ukuran window, misalnya `1366x768`.
- Menentukan folder download ke `tests/e2e/artifacts/downloads`.
- Mengatur timeout default.

### 9.2 `auth.helper.js`

Fungsi wajib:

- `loginAs(driver, role)`
- `loginWithCredentials(driver, username, password)`
- `logout(driver)`
- `expectLoggedInAsRole(driver, role)`
- `expectAccessDenied(driver)`

Role yang wajib didukung:

- `admin`
- `guru`
- `siswa`
- `siswa2`

### 9.3 `form.helper.js`

Fungsi wajib:

- `typeByTestId(driver, testId, value)`
- `clickByTestId(driver, testId)`
- `selectByText(driver, testId, visibleText)`
- `clearAndType(driver, locator, value)`
- `submitAndExpectSuccess(driver)`
- `submitAndExpectValidation(driver, expectedText)`

### 9.4 `table.helper.js`

Fungsi wajib:

- `expectTableContains(driver, text)`
- `expectTableNotContains(driver, text)`
- `getTableRowCount(driver)`
- `searchTable(driver, keyword)`
- `openFirstRowDetail(driver)`
- `deleteRowByText(driver, text)`

### 9.5 `assertion.helper.js`

Fungsi wajib:

- `expectCurrentUrlContains(driver, path)`
- `expectTextVisible(driver, text)`
- `expectTextNotVisible(driver, text)`
- `expectToastSuccess(driver)`
- `expectToastError(driver)`
- `expectValidationMessage(driver, text)`
- `expectEmptyState(driver)`

### 9.6 `api.helper.js`

Fungsi opsional tetapi disarankan:

- `apiLogin(role)`
- `apiGet(path)`
- `apiPost(path, body)`
- `apiDelete(path)`
- `ensureBaseData()`
- `cleanupE2EData()`

Gunakan API helper untuk menyiapkan prasyarat, bukan menggantikan aksi UI utama.

## 10. Urutan Implementasi

Kerjakan bertahap agar mudah diverifikasi.

### Tahap 1: Setup Dasar

1. Install dependency Selenium.
2. Tambahkan npm script.
3. Buat `.env.e2e.example`.
4. Buat `config/env.js`.
5. Buat `config/browser.js`.
6. Buat `config/routes.js`.
7. Buat helper screenshot.
8. Buat spec smoke login admin.
9. Jalankan `npm run test:e2e:auth`.

Selesai jika login admin berhasil dan screenshot gagal bekerja.

### Tahap 2: Auth dan Role

Implementasikan `01-auth.spec.js`.

Target skenario:

- F-04-S01
- F-04-S02
- F-04-S03
- F-16-S01
- F-16-S02
- F-16-S03

Selesai jika test bisa membuktikan login valid, login gagal, role restriction, dan single device session siswa.

### Tahap 3: Manajemen User

Implementasikan `02-user-management.spec.js`.

Target skenario:

- F-01-S01 sampai F-01-S03
- F-02-S01 sampai F-02-S03
- F-03-S01 sampai F-03-S03

Catatan:

- Jika tidak ada halaman registrasi publik, map F-01 ke flow tambah akun oleh admin.
- Jika benar-benar ada route registrasi, buat `RegistrationPage.js`.
- Jangan hapus akun utama admin/guru/siswa yang dipakai test.

### Tahap 4: Master Data

Implementasikan `03-master-data.spec.js`.

Target skenario:

- F-10 sampai F-13
- F-35 sampai F-36
- F-08 bila total mapel terlihat dari dashboard atau menu mapel

Catatan:

- Route jurusan belum terlihat jelas dari route frontend. Jika tidak ada UI jurusan, buat `it.skip` untuk F-10 dengan alasan `route jurusan belum ditemukan`.
- Mata pelajaran memakai route `data-master/mapel`.
- Kelas memakai route `data-master/kelas`.
- Ruang memakai route `data-master/ruang`.
- Sesi memakai route `data-master/sesi`.

### Tahap 5: Dashboard, Pengaturan, Log

Implementasikan `04-dashboard-settings-log.spec.js`.

Target skenario:

- F-05 sampai F-09
- F-14
- F-17 sampai F-21

Catatan:

- Jika menu log aktivitas belum ada di UI, gunakan `it.skip` untuk F-17, F-18, F-19 dengan alasan jelas.
- F-20 adalah pengosongan data dan bersifat destructive. Jalankan hanya pada database E2E.

### Tahap 6: Bank Soal

Implementasikan `05-bank-soal.spec.js`.

Target skenario:

- F-06
- F-23 sampai F-29

Catatan:

- Siapkan file fixture `.docx` valid PG, `.docx` valid essay, dan file invalid.
- Simpan fixture di `frontend/tests/e2e/data/files`.
- Untuk pengacakan soal F-26, siapkan dua siswa dan satu ujian yang sama.

### Tahap 7: Manajemen Ujian

Implementasikan `06-exam-management.spec.js`.

Target skenario:

- F-30 sampai F-34
- F-22 bila cetak daftar hadir tergantung jadwal ujian

Catatan:

- Ujian valid membutuhkan mapel, kelas, ruang, sesi, bank soal, tanggal, durasi, dan token.
- Untuk test jadwal invalid, isi jam selesai lebih awal dari jam mulai atau durasi invalid.

### Tahap 8: Flow Ujian Siswa

Implementasikan `07-student-exam-flow.spec.js`.

Target skenario:

- F-41 sampai F-45
- F-48
- F-44

Catatan:

- Jangan gunakan ujian berdurasi lama untuk E2E otomatis.
- Buat fixture ujian singkat khusus E2E, misalnya 1 sampai 3 menit.
- Untuk timer habis F-43-S03, lebih baik setup ujian dengan sisa waktu sangat pendek daripada menunggu lama.
- Untuk tab baru F-42, Selenium dapat membuka tab baru dengan `driver.switchTo().newWindow("tab")`.

### Tahap 9: Hasil, Koreksi, Review

Implementasikan `08-result-grading-review.spec.js`.

Target skenario:

- F-27
- F-28
- F-38 sampai F-40
- F-46 sampai F-47
- F-39

Catatan:

- Buat satu ujian PG untuk koreksi otomatis.
- Buat satu ujian campuran PG dan essay untuk koreksi manual.
- Login sebagai guru untuk koreksi essay.
- Login sebagai siswa untuk melihat hasil.

### Tahap 10: Cetak, Backup, Destructive

Implementasikan `09-print-backup-destructive.spec.js`.

Target skenario:

- F-15
- F-20
- F-22

Catatan:

- Test download harus mengecek file muncul di folder `artifacts/downloads` dan ukuran file lebih dari 0 byte.
- Test backup gagal biasanya perlu dukungan backend atau environment khusus. Jika tidak bisa disimulasikan dari UI, tandai `it.skip` dengan alasan.
- Test pengosongan data hanya boleh berjalan jika `E2E_ALLOW_DESTRUCTIVE=true`.

## 11. Mapping Lengkap Skenario dari Excel

Gunakan tabel ini sebagai daftar wajib. Jangan menghapus ID skenario.

| ID | Kebutuhan Fungsional | Skenario yang Harus Dibuat |
|---|---|---|
| F-01 | Sistem dapat melakukan registrasi akun untuk administrator, guru, dan siswa | F-01-S01 Registrasi akun valid untuk setiap role<br>F-01-S02 Registrasi ditolak saat data wajib kosong<br>F-01-S03 Registrasi ditolak saat username/email duplikat |
| F-02 | Sistem dapat melakukan pengelolaan terhadap data administrator, guru, dan siswa | F-02-S01 Tambah data pengguna baru<br>F-02-S02 Lihat daftar dan detail pengguna<br>F-02-S03 Hapus/nonaktifkan pengguna |
| F-03 | Sistem dapat melakukan perubahan pada masing-masing data guru dan siswa | F-03-S01 Ubah data guru dengan input valid<br>F-03-S02 Ubah data siswa dengan input valid<br>F-03-S03 Perubahan ditolak saat data tidak valid |
| F-04 | Sistem dapat menjalankan aksi autentikasi dan otorisasi untuk akses ke dalam sistem ujian | F-04-S01 Login berhasil sesuai role<br>F-04-S02 Login gagal dengan kredensial salah<br>F-04-S03 Akses menu dibatasi oleh role |
| F-05 | Sistem dapat menampilkan jumlah data peserta/siswa | F-05-S01 Jumlah siswa tampil sesuai data<br>F-05-S02 Jumlah siswa berubah setelah data ditambah<br>F-05-S03 Jumlah siswa nol saat belum ada data |
| F-06 | Sistem dapat menampilkan jumlah bank soal yang tersedia | F-06-S01 Jumlah bank soal tampil sesuai data<br>F-06-S02 Jumlah bank soal bertambah setelah import/tambah<br>F-06-S03 Jumlah bank soal berkurang setelah penghapusan |
| F-07 | Sistem dapat menampilkan total ujian yang telah dilaksanakan | F-07-S01 Total ujian terlaksana tampil sesuai riwayat<br>F-07-S02 Ujian belum berlangsung tidak dihitung<br>F-07-S03 Total nol saat belum ada ujian selesai |
| F-08 | Sistem dapat menampilkan total mata pelajaran tersedia | F-08-S01 Total mata pelajaran tampil sesuai data<br>F-08-S02 Total berubah setelah mapel ditambah<br>F-08-S03 Total berubah setelah mapel dihapus/nonaktif |
| F-09 | Sebagai administrator, saya ingin mengelola pengaturan umum profil sekolah | F-09-S01 Simpan profil sekolah valid<br>F-09-S02 Validasi profil sekolah wajib<br>F-09-S03 Tampilan profil tetap setelah refresh |
| F-10 | Sistem dapat menyimpan dan menampilkan data jurusan yang tersedia | F-10-S01 Tambah jurusan valid<br>F-10-S02 Daftar jurusan ditampilkan<br>F-10-S03 Jurusan duplikat ditolak |
| F-11 | Sistem dapat menyimpan dan menampilkan data kelas yang tersedia | F-11-S01 Tambah kelas valid<br>F-11-S02 Daftar kelas ditampilkan<br>F-11-S03 Kelas tanpa data wajib ditolak |
| F-12 | Sistem dapat menyimpan dan menampilkan data ruangan yang tersedia | F-12-S01 Tambah ruangan valid<br>F-12-S02 Daftar ruangan ditampilkan<br>F-12-S03 Kapasitas ruangan tidak valid ditolak |
| F-13 | Sistem dapat menyimpan dan menampilkan data sesi yang tersedia | F-13-S01 Tambah sesi valid<br>F-13-S02 Daftar sesi ditampilkan<br>F-13-S03 Rentang waktu sesi tidak valid ditolak |
| F-14 | Sistem dapat melakukan pengelolaan terhadap data umum aplikasi | F-14-S01 Simpan data umum aplikasi valid<br>F-14-S02 Validasi data umum aplikasi<br>F-14-S03 Data umum tetap setelah login ulang |
| F-15 | Sistem dapat melakukan backup terhadap keseluruhan data aplikasi | F-15-S01 Backup data berhasil dibuat<br>F-15-S02 Backup dapat diunduh<br>F-15-S03 Backup gagal ditangani dengan pesan |
| F-16 | Sistem dapat membatasi login hanya dapat dilakukan di satu perangkat bagi siswa | F-16-S01 Login siswa pertama berhasil<br>F-16-S02 Login siswa kedua ditolak saat sesi masih aktif<br>F-16-S03 Login kembali berhasil setelah sesi lama logout/reset |
| F-17 | Sistem dapat menyimpan dan menampilkan seluruh log aktivitas siswa, guru, dan operator | F-17-S01 Log aktivitas siswa tersimpan<br>F-17-S02 Log aktivitas guru tersimpan<br>F-17-S03 Log aktivitas operator/admin tersimpan |
| F-18 | Sistem harus menyimpan log aktivitas pengguna untuk mendukung audit dan keamanan | F-18-S01 Audit log tercatat saat login berhasil<br>F-18-S02 Audit log tercatat saat login gagal<br>F-18-S03 Audit log tidak dapat diubah oleh pengguna biasa |
| F-19 | Sistem dapat menampilkan log aktivitas user kepada administrator | F-19-S01 Administrator melihat daftar log user<br>F-19-S02 Filter log berdasarkan pengguna/tanggal<br>F-19-S03 Akses log ditolak untuk non-admin |
| F-20 | Sistem dapat mengosongkan data keseluruhan yang dipilih oleh administrator | F-20-S01 Kosongkan data yang dipilih berhasil<br>F-20-S02 Batal konfirmasi tidak menghapus data<br>F-20-S03 Pengosongan tanpa pilihan ditolak |
| F-21 | Sistem dapat menyimpan dan menampilkan riwayat ujian | F-21-S01 Riwayat ujian tersimpan setelah ujian selesai<br>F-21-S02 Riwayat ujian ditampilkan sesuai pengguna/role<br>F-21-S03 Riwayat kosong ditampilkan dengan aman |
| F-22 | Sistem dapat membuat dokumen cetak untuk daftar hadir peserta ujian berdasarkan kelas, ruangan, mata pelajaran, dan sesi ujian | F-22-S01 Cetak daftar hadir dengan filter valid<br>F-22-S02 Daftar hadir hanya berisi peserta sesuai filter<br>F-22-S03 Cetak ditolak bila filter wajib kosong |
| F-23 | Sistem dapat melakukan penginputan bank soal berjenis pilihan ganda maupun essay melalui dokumen berformat .docx | F-23-S01 Import bank soal pilihan ganda dari docx valid<br>F-23-S02 Import bank soal essay dari docx valid<br>F-23-S03 Import ditolak untuk format file tidak valid |
| F-24 | Sistem dapat melakukan penghapusan terhadap bank soal yang telah tersimpan | F-24-S01 Hapus bank soal berhasil<br>F-24-S02 Batal hapus bank soal<br>F-24-S03 Hapus ditolak untuk pengguna tidak berwenang |
| F-25 | Sistem dapat melakukan modifikasi terhadap bank soal yang telah tersimpan | F-25-S01 Modifikasi bank soal valid<br>F-25-S02 Modifikasi ditolak bila field wajib kosong<br>F-25-S03 Perubahan bank soal tidak disimpan saat batal |
| F-26 | Sistem dapat melakukan pengacakan nomor dan urutan soal ujian | F-26-S01 Nomor soal diacak saat ujian dimulai<br>F-26-S02 Urutan pilihan jawaban diacak<br>F-26-S03 Urutan tetap saat fitur acak nonaktif |
| F-27 | Sistem dapat identifikasi soal yang paling banyak dijawab salah dalam satu ujian | F-27-S01 Identifikasi soal paling banyak salah<br>F-27-S02 Perhitungan salah diperbarui setelah koreksi<br>F-27-S03 Analisis kosong saat belum ada jawaban |
| F-28 | Sistem dapat identifikasi soal yang paling banyak dijawab benar dalam satu ujian | F-28-S01 Identifikasi soal paling banyak benar<br>F-28-S02 Perhitungan benar diperbarui setelah koreksi<br>F-28-S03 Analisis benar kosong saat belum ada jawaban |
| F-29 | Sistem dapat melakukan pembobotan terhadap masing-masing soal | F-29-S01 Simpan bobot soal valid<br>F-29-S02 Bobot tidak valid ditolak<br>F-29-S03 Nilai mengikuti bobot soal |
| F-30 | Sistem dapat membuat ujian baru sesuai dengan mata pelajaran dan kelas yang tersedia | F-30-S01 Buat ujian baru valid<br>F-30-S02 Buat ujian ditolak bila mapel/kelas kosong<br>F-30-S03 Ujian baru hanya memakai data tersedia |
| F-31 | Sistem dapat melakukan modifikasi terhadap ujian yang akan atau telah dilaksanakan | F-31-S01 Modifikasi ujian terjadwal valid<br>F-31-S02 Modifikasi ujian telah dilaksanakan sesuai aturan<br>F-31-S03 Modifikasi ditolak bila jadwal tidak valid |
| F-32 | Sistem dapat melakukan penghapusan terhadap ujian yang akan atau telah dilaksanakan | F-32-S01 Hapus ujian berhasil<br>F-32-S02 Batal hapus ujian<br>F-32-S03 Hapus ujian ditolak untuk pengguna tidak berwenang |
| F-33 | Sistem dapat menampilkan ujian-ujian yang akan atau telah berlangsung berdasarkan mata pelajaran | F-33-S01 Tampilkan ujian berdasarkan mata pelajaran<br>F-33-S02 Daftar memuat ujian akan dan telah berlangsung<br>F-33-S03 Filter tanpa hasil menampilkan keadaan kosong |
| F-34 | Sistem dapat melakukan pembuatan token berdasarkan masukan dari guru untuk setiap ujian | F-34-S01 Buat token ujian valid<br>F-34-S02 Token kosong atau tidak valid ditolak<br>F-34-S03 Token baru menggantikan token lama |
| F-35 | Sistem dapat menyimpan dan menampilkan mata pelajaran yang tersedia | F-35-S01 Tambah mata pelajaran valid<br>F-35-S02 Daftar mata pelajaran ditampilkan<br>F-35-S03 Mata pelajaran duplikat ditolak |
| F-36 | Sistem dapat memodifikasi mata pelajaran yang sebelumnya telah tersimpan | F-36-S01 Modifikasi mata pelajaran valid<br>F-36-S02 Modifikasi mapel ditolak bila duplikat<br>F-36-S03 Perubahan mapel tercermin pada form ujian |
| F-37 | Sistem dapat membuat pengumuman untuk ditampilkan kepada siswa dan guru | F-37-S01 Buat pengumuman valid<br>F-37-S02 Pengumuman wajib isi ditolak bila kosong<br>F-37-S03 Pengumuman hanya tampil pada target penerima |
| F-38 | Sistem dapat melakukan penilaian terhadap ujian yang telah dikerjakan | F-38-S01 Penilaian otomatis pilihan ganda<br>F-38-S02 Penilaian gabungan PG dan essay<br>F-38-S03 Penilaian belum final saat essay belum dikoreksi |
| F-39 | Sistem dapat menampilkan hasil nilai ujian siswa berdasarkan ujian yang diikuti siswa | F-39-S01 Siswa melihat nilai ujian yang diikuti<br>F-39-S02 Siswa tidak melihat nilai ujian yang tidak diikuti<br>F-39-S03 Nilai belum tersedia ditampilkan dengan status |
| F-40 | Sistem dapat melakukan koreksi essay secara manual berdasarkan benar/salah yang dimasukkan oleh guru | F-40-S01 Guru menandai jawaban essay benar<br>F-40-S02 Guru menandai jawaban essay salah<br>F-40-S03 Koreksi essay ditolak untuk guru tidak berwenang |
| F-41 | Sistem dapat menerapkan fitur ketika keluar dari sesi ujian, maka ujian dapat diikuti kembali tanpa reset | F-41-S01 Siswa dapat melanjutkan ujian setelah keluar sesi<br>F-41-S02 Timer tetap mengikuti sisa waktu<br>F-41-S03 Ujian tidak dapat dilanjutkan setelah sudah submit |
| F-42 | Sistem dapat mencegah siswa untuk membuka tab baru saat ujian | F-42-S01 Pembukaan tab baru saat ujian dicegah<br>F-42-S02 Pelanggaran tab baru dicatat/ditangani<br>F-42-S03 Fitur tidak mengganggu halaman non-ujian |
| F-43 | Sistem dapat mengatur dan menampilkan timer terhadap masing-masing ujian | F-43-S01 Timer ujian tampil sesuai durasi<br>F-43-S02 Timer berkurang real-time<br>F-43-S03 Ujian otomatis berakhir saat timer habis |
| F-44 | Sistem dapat menampilkan ujian-ujian yang dapat diikuti oleh siswa berdasarkan kelas, sesi, mata pelajaran, dan jadwal siswa tersebut | F-44-S01 Siswa melihat ujian sesuai kelas dan sesi<br>F-44-S02 Ujian di luar jadwal tidak tampil/tidak bisa diikuti<br>F-44-S03 Ujian kelas lain tidak tampil |
| F-45 | Sistem dapat menampilkan soal yang belum dijawab dan sudah dijawab | F-45-S01 Indikator soal belum dijawab tampil<br>F-45-S02 Indikator berubah setelah soal dijawab<br>F-45-S03 Indikator tetap setelah pindah soal/refresh |
| F-46 | Sistem dapat menampilkan soal yang dijawab salah dan dijawab benar serta menampilkan jawaban benar pada soal yang dijawab salah | F-46-S01 Review jawaban benar dan salah ditampilkan<br>F-46-S02 Jawaban benar tampil pada soal yang salah<br>F-46-S03 Review tidak tampil sebelum ujian dikoreksi/final |
| F-47 | Sistem dapat menampilkan nilai ujian siswa berdasarkan hasil koreksi | F-47-S01 Nilai tampil setelah koreksi selesai<br>F-47-S02 Nilai berubah setelah koreksi diperbarui<br>F-47-S03 Nilai tidak tampil bila koreksi belum lengkap |
| F-48 | Sistem dapat menampilkan review sebelum siswa dapat benar-benar submit jawaban | F-48-S01 Review jawaban muncul sebelum submit akhir<br>F-48-S02 Submit akhir berhasil setelah konfirmasi review<br>F-48-S03 Batal submit kembali ke pengerjaan |

## 12. Pembagian Spec Berdasarkan Skenario

Gunakan pembagian ini saat membuat file test.

| File spec | Skenario |
|---|---|
| `01-auth.spec.js` | F-04, F-16 |
| `02-user-management.spec.js` | F-01, F-02, F-03 |
| `03-master-data.spec.js` | F-08, F-10, F-11, F-12, F-13, F-35, F-36 |
| `04-dashboard-settings-log.spec.js` | F-05, F-07, F-09, F-14, F-17, F-18, F-19, F-20, F-21 |
| `05-bank-soal.spec.js` | F-06, F-23, F-24, F-25, F-26, F-29 |
| `06-exam-management.spec.js` | F-22, F-30, F-31, F-32, F-33, F-34, F-37 |
| `07-student-exam-flow.spec.js` | F-41, F-42, F-43, F-44, F-45, F-48 |
| `08-result-grading-review.spec.js` | F-27, F-28, F-38, F-39, F-40, F-46, F-47 |
| `09-print-backup-destructive.spec.js` | F-15 dan destructive subset F-20 |

Jika satu skenario lebih cocok berada di file lain karena flow data, boleh dipindah, tetapi ID skenario tetap wajib dipertahankan.

## 13. Detail Implementasi per Spec

### 13.1 `01-auth.spec.js`

Buat test:

- `[F-04-S01] Login berhasil sesuai role`
  - Login admin, expect URL admin dashboard.
  - Logout.
  - Login guru, expect URL guru dashboard.
  - Logout.
  - Login siswa, expect URL siswa dashboard.
- `[F-04-S02] Login gagal dengan kredensial salah`
  - Buka login.
  - Isi username valid dan password salah.
  - Klik masuk.
  - Expect pesan error.
- `[F-04-S03] Akses menu dibatasi oleh role`
  - Login siswa.
  - Buka URL admin.
  - Expect ditolak, redirect, atau not found.
- `[F-16-S01] Login siswa pertama berhasil`
  - Login siswa pada browser pertama.
  - Expect dashboard siswa.
- `[F-16-S02] Login siswa kedua ditolak saat sesi masih aktif`
  - Tanpa logout browser pertama, buka browser kedua.
  - Login dengan akun siswa yang sama.
  - Expect pesan sesi masih aktif.
- `[F-16-S03] Login kembali berhasil setelah sesi lama logout/reset`
  - Logout dari browser pertama atau reset session via API helper.
  - Login ulang siswa.
  - Expect berhasil.

### 13.2 `02-user-management.spec.js`

Buat test:

- F-01: registrasi/tambah akun untuk admin, guru, siswa.
- F-02: tambah, list/detail, hapus/nonaktifkan user.
- F-03: edit guru, edit siswa, validasi edit invalid.

Instruksi:

- Login sebagai admin di `beforeEach`.
- Gunakan username unik.
- Untuk test hapus, buat user khusus lalu hapus user itu.
- Jangan menghapus akun utama dari `.env.e2e`.

### 13.3 `03-master-data.spec.js`

Buat test:

- F-08: total mata pelajaran tampil dan berubah.
- F-10: jurusan jika UI tersedia.
- F-11: kelas.
- F-12: ruang ujian.
- F-13: sesi.
- F-35: tambah dan list mapel.
- F-36: edit mapel dan validasi duplikat.

Instruksi:

- Gunakan data unik.
- Setelah tambah data, assert data muncul di tabel.
- Untuk duplikat, tambahkan data pertama lalu submit data kedua dengan kode/nama yang sama.
- Untuk validasi, kosongkan field wajib lalu assert pesan validasi.

### 13.4 `04-dashboard-settings-log.spec.js`

Buat test:

- F-05: total siswa.
- F-07: total ujian terlaksana.
- F-09: profil sekolah.
- F-14: pengaturan umum aplikasi jika UI tersedia.
- F-17, F-18, F-19: log aktivitas/audit.
- F-20: pengosongan data, hanya jika destructive diizinkan.
- F-21: riwayat ujian.

Instruksi:

- Untuk counter dashboard, ambil angka awal, lakukan aksi tambah/hapus, lalu expect angka berubah.
- Untuk empty state, siapkan data kosong via API helper atau pakai user baru tanpa riwayat.
- Untuk log, lakukan aksi yang pasti menghasilkan log lalu buka menu log.

### 13.5 `05-bank-soal.spec.js`

Buat test:

- F-06: total bank soal.
- F-23: import `.docx` PG, essay, invalid.
- F-24: hapus bank soal, batal hapus, role tidak berwenang.
- F-25: edit bank soal valid, invalid, batal.
- F-26: acak soal dan opsi.
- F-29: bobot soal valid, invalid, nilai mengikuti bobot.

Instruksi:

- Siapkan fixture file.
- Untuk upload file, pakai Selenium `sendKeys` ke input file.
- Untuk acak soal, mulai ujian dengan dua akun siswa berbeda lalu bandingkan urutan soal.
- Jika random kadang menghasilkan urutan sama, ulang maksimal 2 kali atau assert bahwa fitur acak aktif lewat indikator konfigurasi.

### 13.6 `06-exam-management.spec.js`

Buat test:

- F-22: cetak daftar hadir.
- F-30: buat ujian.
- F-31: edit ujian.
- F-32: hapus ujian.
- F-33: filter ujian berdasarkan mapel.
- F-34: buat token ujian.
- F-37: buat pengumuman.

Instruksi:

- Ujian baru harus memakai data master yang sudah disiapkan.
- Token ujian harus disimpan dalam fixture runtime agar bisa dipakai siswa.
- Untuk cetak, pastikan file download muncul dan tidak kosong.
- Untuk pengumuman target penerima, login sebagai target dan non-target.

### 13.7 `07-student-exam-flow.spec.js`

Buat test:

- F-41: resume ujian setelah keluar.
- F-42: pembatasan tab baru saat ujian.
- F-43: timer ujian.
- F-44: ujian siswa sesuai kelas/sesi/jadwal.
- F-45: indikator soal dijawab/belum dijawab.
- F-48: review sebelum submit, submit, batal submit.

Instruksi:

- Gunakan ujian khusus E2E dengan jumlah soal sedikit.
- Simpan jawaban satu soal, keluar halaman, masuk lagi, assert jawaban tetap ada.
- Untuk timer, catat nilai timer awal dan nilai setelah beberapa detik.
- Untuk submit akhir, test batal dulu sebelum test konfirmasi submit.

### 13.8 `08-result-grading-review.spec.js`

Buat test:

- F-27: soal paling banyak salah.
- F-28: soal paling banyak benar.
- F-38: penilaian otomatis dan gabungan.
- F-39: siswa melihat nilai sesuai ujian yang diikuti.
- F-40: koreksi essay manual.
- F-46: review jawaban benar/salah.
- F-47: nilai final dan nilai berubah setelah koreksi.

Instruksi:

- Buat minimal dua attempt siswa agar statistik benar/salah bisa diuji.
- Untuk essay, guru memberi koreksi benar dan salah pada dua jawaban berbeda.
- Setelah koreksi, refresh hasil siswa dan assert nilai berubah.

### 13.9 `09-print-backup-destructive.spec.js`

Buat test:

- F-15: backup berhasil, download backup, backup gagal.
- F-20: kosongkan data terpilih, batal, tanpa pilihan.
- F-22: cetak daftar hadir jika belum dicakup di spec ujian.

Instruksi:

- Jalankan spec ini hanya pada database khusus E2E.
- Tambahkan guard:

```js
if (process.env.E2E_ALLOW_DESTRUCTIVE !== "true") {
  this.skip();
}
```

## 14. Strategi Skip yang Diperbolehkan

Boleh memakai `it.skip` hanya jika:

- Route UI tidak ditemukan.
- Fitur belum tersedia di frontend.
- Skenario membutuhkan kondisi infrastruktur yang tidak bisa disimulasikan lokal, misalnya storage backup gagal.
- Skenario destructive tidak diizinkan oleh environment.

Format alasan skip:

```js
it.skip("[F-15-S03] Backup gagal ditangani dengan pesan - perlu mode backend untuk simulasi storage gagal", async function () {});
```

Jangan skip karena selector sulit. Jika selector sulit, tambahkan `data-testid`.

## 15. Kriteria Lulus

Planning dianggap berhasil diimplementasikan jika:

- Semua 144 ID skenario dari Excel muncul di file spec.
- Skenario yang bisa dijalankan memiliki assertion nyata.
- Skenario yang belum bisa dijalankan memakai `it.skip` dengan alasan jelas.
- Test memakai JavaScript ESM.
- Test dapat dijalankan dari folder `frontend` dengan `npm run test:e2e`.
- Screenshot tersimpan otomatis saat test gagal.
- Test tidak bergantung pada data produksi.
- Test dapat dijalankan ulang tanpa konflik data karena memakai data unik atau cleanup.

## 16. Checklist Implementasi untuk Model Selanjutnya

Ikuti checklist ini secara berurutan:

1. Buka `dokumentasi_black_box_testing.xlsx` dan pastikan jumlah skenario tetap 144.
2. Buka `frontend/package.json`.
3. Install dependency Selenium.
4. Tambahkan script npm.
5. Buat struktur folder `frontend/tests/e2e`.
6. Buat `.env.e2e.example`.
7. Buat `env.js`, `browser.js`, dan `routes.js`.
8. Buat helper dasar: auth, wait, assertion, screenshot.
9. Buat `LoginPage.js`.
10. Implementasikan `01-auth.spec.js`.
11. Jalankan `npm run test:e2e:auth`.
12. Jika auth stabil, lanjut page object lain.
13. Implementasikan spec sesuai urutan tahap.
14. Pastikan semua `it()` diawali ID skenario.
15. Tambahkan `data-testid` bila selector tidak stabil.
16. Jalankan spec per modul.
17. Jalankan seluruh test.
18. Perbaiki test gagal yang disebabkan timing dengan explicit wait, bukan `sleep` panjang.
19. Tulis ringkasan skenario pass, fail, skip.
20. Jangan mengubah skenario Excel tanpa instruksi.

## 17. Catatan Risiko

- Beberapa fitur pada Excel tidak terlihat jelas pada route frontend, misalnya registrasi publik, jurusan, backup, audit log, dan pengosongan data. Implementasi harus memverifikasi UI terlebih dahulu.
- Test ujian dengan timer bisa lambat jika memakai durasi asli. Gunakan data E2E berdurasi pendek.
- Test randomisasi soal bisa flaky. Gunakan data lebih dari 3 soal dan bandingkan dua siswa.
- Test download perlu path download khusus di Chrome options.
- Test single-device login membutuhkan dua browser session atau reset session backend.
- Test destructive harus dipisah agar tidak merusak data test lain.

## 18. Rekomendasi Prioritas Awal

Prioritas implementasi pertama:

1. F-04 autentikasi dan otorisasi.
2. F-02/F-03 manajemen akun.
3. F-11/F-12/F-13/F-35/F-36 master data.
4. F-30/F-34 pembuatan ujian dan token.
5. F-44/F-48 flow siswa mengikuti ujian.
6. F-38/F-40/F-47 hasil dan koreksi.

Alasan:

- Flow ini membentuk jalur utama sistem CBT dari admin/guru menyiapkan data sampai siswa mengerjakan dan melihat hasil.
- Setelah jalur utama stabil, skenario dashboard, log, backup, cetak, dan destructive lebih mudah dilengkapi.
