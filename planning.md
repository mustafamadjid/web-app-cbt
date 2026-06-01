# Planning: Performance Test Fitur Ujian Siswa dengan k6

## 1. Tujuan

Dokumen ini berisi rencana implementasi performance test menggunakan **k6 JavaScript** untuk fitur ujian siswa pada sistem ujian sekolah berbasis website.

Pengujian difokuskan pada fitur ujian karena fitur ini merupakan fitur utama yang digunakan secara bersamaan oleh banyak siswa dan menjadi skenario kritis dalam sistem.

Jenis pengujian yang perlu diimplementasikan:

1. **Baseline Test**
2. **Load Test**

Stress test tidak perlu dibuat karena fokus pengujian adalah memastikan sistem memenuhi kebutuhan non-fungsional, bukan mencari breaking point sistem.

---

## 2. Kebutuhan Non-Fungsional

Kebutuhan non-fungsional utama:

> Waktu respons sistem tidak lebih dari 4 detik.

Kriteria ini diterapkan pada response time request API, bukan total durasi satu alur ujian dari login sampai submit.

Metrik utama yang perlu diamati:

| Metrik | Target |
|---|---|
| `http_req_duration` P95 | `< 4000 ms` |
| `http_req_failed` | `< 1%` |
| Status Code 503 | Tidak boleh muncul |
| Error Rate | Serendah mungkin, idealnya 0% |
| RPS | Stabil selama durasi pengujian |

---

## 3. Alur Skenario Ujian Siswa

Setiap virtual user di k6 merepresentasikan satu siswa.

Satu iterasi merepresentasikan satu alur siswa mengikuti ujian.

Alur yang harus diimplementasikan:

1. Siswa login.
2. Siswa masuk ke halaman ujian.
3. Sistem mengambil jadwal atau detail ujian.
4. Siswa memasukkan token ujian.
5. Sistem mengambil daftar soal.
6. Siswa menyimpan jawaban beberapa kali.
7. Siswa melakukan submit ujian.
8. Sistem menyelesaikan proses penilaian atau penyimpanan hasil.

Catatan penting:

- Gunakan akun siswa yang berbeda untuk setiap virtual user.
- Jangan menggunakan satu akun siswa untuk banyak virtual user karena sistem memiliki mekanisme pembatasan login satu perangkat.
- Token ujian dapat dikirim melalui environment variable.
- Endpoint aktual harus disesuaikan dengan route backend yang digunakan pada project.

---

## 4. Rencana Skenario Test

### 4.1 Baseline Test

Baseline test digunakan untuk memperoleh nilai awal performa sistem dan menjadi pembanding terhadap load test.

| No | Jenis Test | Beban Pengguna | Durasi | Tujuan |
|---:|---|---:|---|---|
| 1 | Baseline Test | 10 concurrent user | 5 menit | Mengukur performa awal sistem pada beban ringan |
| 2 | Baseline Test | 20 concurrent user | 5 menit | Membandingkan performa sistem dengan kondisi masalah sistem lama yang gagal saat pengguna lebih dari 20 siswa |

### 4.2 Load Test

Load test digunakan untuk menguji apakah sistem memenuhi kebutuhan non-fungsional pada beban pengguna yang lebih tinggi.

| No | Jenis Test | Beban Pengguna | Durasi | Tujuan |
|---:|---|---:|---|---|
| 1 | Load Test | 50 concurrent user | 10 menit | Menguji kemampuan sistem pada beban sedang |
| 2 | Load Test | 100 concurrent user | 10 menit | Menguji kemampuan sistem pada beban target yang lebih tinggi |

---

## 5. Struktur File yang Direkomendasikan

Codex agent perlu membuat struktur file berikut:

```txt
performance-test/
├── planning.md
├── k6/
│   ├── exam-flow.test.js
│   ├── config.js
│   └── data/
│       └── students.json
└── README.md
```

Keterangan:

| File | Fungsi |
|---|---|
| `planning.md` | Dokumen rencana implementasi performance test |
| `k6/exam-flow.test.js` | Script utama k6 untuk alur ujian siswa |
| `k6/config.js` | Konfigurasi endpoint, environment variable, dan helper konfigurasi |
| `k6/data/students.json` | Data akun siswa untuk virtual user |
| `README.md` | Panduan menjalankan test |

---

## 6. Environment Variable

Script k6 harus mendukung konfigurasi melalui environment variable.

Contoh variable:

| Variable | Wajib | Contoh | Keterangan |
|---|---|---|---|
| `BASE_URL` | Ya | `https://staging-srv.smafi.my.id/` | Base URL backend API |
| `TOKEN_UJIAN` | Ya | `ABC123` | Token ujian yang akan digunakan siswa |
| `TEST_TYPE` | Ya | `baseline_10` | Jenis skenario test yang dijalankan |
| `STUDENTS_FILE` | Tidak | `./data/students.json` | Lokasi file data siswa jika ingin dibuat fleksibel |

Contoh perintah:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_10
```

---

## 7. Data Uji Siswa

Buat file:

```txt
k6/data/students.json
```

Format data:

```json
[
  {
    "username": "siswa001",
    "password": "siswa12345"
  },
  {
    "username": "siswa002",
    "password": "siswa12345"
  }
]
```

Ketentuan:

- Minimal sediakan jumlah akun sesuai jumlah virtual user terbesar.
- Untuk load test 100 concurrent user, sediakan minimal 100 akun siswa.
- Semua akun siswa harus valid dan memiliki akses ke jadwal ujian yang sama.
- Semua akun siswa sebaiknya berada pada kelas, sesi, dan ruangan yang sesuai dengan jadwal ujian yang diuji.

---

## 8. Desain Script k6

### 8.1 Prinsip Implementasi

Script k6 harus dibuat modular dan mudah dibaca.

Minimal pisahkan logic menjadi beberapa fungsi:

```txt
loginStudent()
getActiveExamSchedule()
startExamAttempt()
getExamQuestions()
saveAnswers()
submitExam()
```

Setiap fungsi harus:

- menerima parameter yang jelas,
- mengembalikan data yang diperlukan untuk step berikutnya,
- melakukan `check()` terhadap status response,
- menangani kegagalan dengan aman,
- tidak melanjutkan flow jika step penting gagal.

---

## 9. Skenario k6 Options

Codex agent perlu membuat konfigurasi `options` berdasarkan `TEST_TYPE`.

### 9.1 Baseline 10 User

```js
{
  scenarios: {
    baseline_10: {
      executor: "constant-vus",
      vus: 10,
      duration: "5m"
    }
  }
}
```

### 9.2 Baseline 20 User

```js
{
  scenarios: {
    baseline_20: {
      executor: "constant-vus",
      vus: 20,
      duration: "5m"
    }
  }
}
```

### 9.3 Load 50 User

```js
{
  scenarios: {
    load_50: {
      executor: "constant-vus",
      vus: 50,
      duration: "10m"
    }
  }
}
```

### 9.4 Load 100 User

```js
{
  scenarios: {
    load_100: {
      executor: "constant-vus",
      vus: 100,
      duration: "10m"
    }
  }
}
```

---

## 10. Threshold k6

Gunakan threshold berikut:

```js
thresholds: {
  http_req_duration: ["p(95)<4000"],
  http_req_failed: ["rate<0.01"],
  checks: ["rate>0.99"]
}
```

Tambahkan custom metric jika diperlukan:

```js
import { Rate, Trend } from "k6/metrics";

export const examFlowFailed = new Rate("exam_flow_failed");
export const examFlowDuration = new Trend("exam_flow_duration", true);
```

Custom threshold:

```js
thresholds: {
  http_req_duration: ["p(95)<4000"],
  http_req_failed: ["rate<0.01"],
  checks: ["rate>0.99"],
  exam_flow_failed: ["rate<0.01"]
}
```

Catatan:

- `http_req_duration` digunakan untuk mengukur response time API.
- `exam_flow_duration` boleh dicatat, tetapi jangan dijadikan acuan batas 4 detik karena total flow mencakup waktu simulasi siswa membaca atau mengerjakan soal.
- Batas 4 detik berlaku pada response API, bukan total waktu pengerjaan ujian.

---

## 11. Rencana Endpoint Placeholder

Codex agent harus menyesuaikan endpoint dengan route aktual pada backend.

Gunakan placeholder berikut jika route belum diketahui:

| Step | Method | Endpoint Placeholder | Keterangan |
|---|---|---|---|
| Login siswa | `POST` | `/auth/login` | Menghasilkan access token |
| Masuk halaman ujian / ambil ujian aktif | `GET` | `/siswa/ujian/aktif` | Mengambil daftar ujian aktif siswa |
| Mulai attempt / validasi token | `POST` | `/siswa/ujian/attempt` | Mengirim token dan ID jadwal ujian |
| Ambil daftar soal | `GET` | `/siswa/ujian/{attemptId}/soal` | Mengambil soal berdasarkan attempt |
| Simpan jawaban | `POST` | `/siswa/ujian/jawaban` | Menyimpan jawaban siswa |
| Submit ujian | `POST` | `/siswa/ujian/submit` | Menyelesaikan ujian |
| Ambil hasil / validasi selesai | `GET` atau internal response | `/siswa/ujian/{attemptId}/hasil` | Opsional, jika sistem menyediakan endpoint hasil |

Jika sistem menggunakan cookie untuk autentikasi:

- Pastikan k6 menyimpan cookie dari response login.
- Jika access token dikirim melalui response body, gunakan header `Authorization: Bearer <token>`.
- Jika refresh token disimpan dalam cookie, tidak perlu diambil manual selama flow cukup menggunakan access token.

---

## 12. Detail Flow Implementasi

### 12.1 Login Siswa

Input:

```json
{
  "username": "siswa001",
  "password": "siswa12345"
}
```

Validasi:

- Status response `200`.
- Response memiliki access token atau cookie session.
- Tidak ada status `401`, `403`, atau `503`.

Output:

```txt
accessToken
```

---

### 12.2 Ambil Jadwal atau Detail Ujian

Request:

```txt
GET /siswa/ujian/aktif
```

Validasi:

- Status response `200`.
- Data ujian aktif tersedia.
- ID jadwal ujian tersedia.

Output:

```txt
idJadwalUjian
```

---

### 12.3 Masukkan Token Ujian / Mulai Attempt

Request body:

```json
{
  "idJadwalUjian": 1,
  "token": "ABC123"
}
```

Validasi:

- Status response `200` atau `201`.
- Response memiliki `idAttempt`.
- Tidak ada status `400`, `401`, `403`, atau `503`.

Output:

```txt
idAttempt
```

---

### 12.4 Ambil Daftar Soal

Request:

```txt
GET /siswa/ujian/{attemptId}/soal
```

Validasi:

- Status response `200`.
- Daftar soal tidak kosong.
- Setiap soal memiliki ID soal.
- Untuk pilihan ganda, opsi jawaban tersedia.

Output:

```txt
questions[]
```

---

### 12.5 Simpan Jawaban Beberapa Kali

Simulasikan siswa menjawab beberapa soal.

Ketentuan:

- Jangan harus menjawab semua soal jika ingin test lebih cepat.
- Minimal simpan 5 sampai 10 jawaban.
- Gunakan jawaban statis atau acak sesuai struktur data soal.
- Beri `sleep()` kecil antar request untuk meniru perilaku siswa.

Contoh body:

```json
{
  "idAttempt": 10,
  "idSoal": 99,
  "jawaban": "A"
}
```

Validasi:

- Status response `200` atau `201`.
- Tidak ada status `500` atau `503`.

---

### 12.6 Submit Ujian

Request body:

```json
{
  "idAttempt": 10
}
```

Validasi:

- Status response `200` atau `201`.
- Response menyatakan ujian berhasil disubmit.
- Jika sistem melakukan penilaian otomatis, pastikan proses tidak menghasilkan error.

---

### 12.7 Penyelesaian Penilaian atau Penyimpanan Hasil

Jika submit ujian langsung menjalankan proses penilaian otomatis, cukup validasi response submit.

Jika hasil ujian tersedia melalui endpoint terpisah, tambahkan step opsional:

```txt
GET /siswa/ujian/{attemptId}/hasil
```

Validasi:

- Status response `200`.
- Data hasil tersedia.
- Nilai atau status selesai tersedia.

---

## 13. Handling Failure

Jika salah satu step penting gagal, flow harus dihentikan untuk virtual user tersebut.

Contoh aturan:

- Jika login gagal, jangan lanjut ambil jadwal.
- Jika ambil jadwal gagal, jangan lanjut attempt.
- Jika attempt gagal, jangan lanjut ambil soal.
- Jika ambil soal gagal, jangan lanjut simpan jawaban.
- Jika simpan jawaban gagal, tandai flow gagal.
- Jika submit gagal, tandai flow gagal.

Catat kegagalan menggunakan custom metric:

```js
examFlowFailed.add(1);
```

Jika flow berhasil:

```js
examFlowFailed.add(0);
```

---

## 14. Simulasi Think Time

Gunakan `sleep()` untuk meniru perilaku siswa.

Rekomendasi:

| Step | Sleep |
|---|---:|
| Setelah login | 1 detik |
| Setelah ambil jadwal | 1 detik |
| Setelah mulai ujian | 1 detik |
| Antar simpan jawaban | 1 sampai 3 detik |
| Sebelum submit | 1 detik |

Catatan:

- `sleep()` tidak dihitung sebagai response time API.
- `sleep()` memengaruhi total durasi flow dan RPS.
- Jangan gunakan sleep terlalu panjang jika tujuan test adalah mengamati beban API dalam durasi terbatas.

---

## 15. Output dan Laporan

Codex agent perlu memastikan hasil k6 dapat diekspor.

Contoh perintah output JSON:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_20 \
  --summary-export results-baseline-20.json
```

File hasil yang disarankan:

```txt
performance-test/results/
├── baseline-10.json
├── baseline-20.json
├── load-50.json
└── load-100.json
```

---

## 16. README yang Perlu Dibuat

Codex agent perlu membuat `README.md` berisi:

1. Tujuan performance test.
2. Prasyarat menjalankan test.
3. Cara menyiapkan data siswa.
4. Cara menjalankan baseline test.
5. Cara menjalankan load test.
6. Penjelasan metrik k6.
7. Penjelasan threshold.
8. Contoh interpretasi hasil.

---

## 17. Acceptance Criteria

Implementasi dianggap selesai jika memenuhi kriteria berikut:

- Tersedia script `k6/exam-flow.test.js`.
- Script dapat menjalankan baseline 10 user.
- Script dapat menjalankan baseline 20 user.
- Script dapat menjalankan load 50 user.
- Script dapat menjalankan load 100 user.
- Script menjalankan alur ujian siswa secara lengkap:
  - login,
  - ambil jadwal/detail ujian,
  - input token,
  - mulai attempt,
  - ambil soal,
  - simpan jawaban beberapa kali,
  - submit ujian,
  - validasi hasil atau status selesai jika tersedia.
- Script menggunakan data siswa dari `students.json`.
- Script mendukung `BASE_URL`, `TOKEN_UJIAN`, dan `TEST_TYPE`.
- Script memiliki threshold `p(95)<4000`.
- Script mencatat failure flow menggunakan custom metric.
- Dokumentasi menjalankan test tersedia di `README.md`.

---

## 18. Catatan untuk Codex Agent

Saat mengimplementasikan:

1. Jangan hardcode endpoint jika route aktual dapat dibaca dari source code backend.
2. Cari route aktual pada project backend terlebih dahulu.
3. Sesuaikan parsing response JSON dengan struktur response aktual.
4. Gunakan helper function agar script mudah dirawat.
5. Hindari duplikasi logic request.
6. Gunakan nama fungsi yang deskriptif.
7. Pastikan error message di `check()` mudah dibaca.
8. Jangan menggunakan satu akun siswa untuk semua VU.
9. Pastikan data ujian, jadwal, token, kelas, sesi, dan akun siswa sudah valid sebelum menjalankan test.
10. Jangan membuat stress test karena tidak termasuk scope.
