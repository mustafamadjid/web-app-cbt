# Performance Test Ujian Siswa

Performance test ini menggunakan k6 untuk menjalankan alur ujian siswa pada backend CBT.

## Prasyarat

- k6 sudah terpasang.
- Akun siswa di `k6/data/students.json` valid.
- Siswa memiliki jadwal ujian yang sama.
- Token ujian aktif tersedia.

## Konfigurasi

Environment variable yang dipakai:

| Variable | Wajib | Default | Keterangan |
|---|---|---|---|
| `BASE_URL` | Ya | `https://staging-srv.smafi.my.id/` | Base URL backend API |
| `TOKEN_UJIAN` | Ya | kosong | Token ujian dari pengawas |
| `TEST_TYPE` | Tidak | `load_100` | `baseline_5`, `load_25`, `load_50`, atau `load_100` |
| `STUDENTS_FILE` | Tidak | `./data/students.json` | Lokasi fixture akun siswa relatif dari file k6 |
| `ID_JADWAL_UJIAN` | Tidak | kosong | Pakai jadwal spesifik jika ada banyak jadwal |
| `ANSWER_LIMIT` | Tidak | `10` | Jumlah maksimal jawaban yang disimpan per flow |

## Menjalankan Test

Semua skenario memakai executor `constant-vus`.
`baseline_5` menjalankan 5 VU selama 5 menit.
`load_25`, `load_50`, dan `load_100` menjalankan VU sesuai nama skenario selama 10 menit.
Setiap iterasi mengambil akun siswa yang belum pernah dipakai pada run tersebut.
Jika akun di `k6/data/students.json` habis, test akan berhenti dengan error `Not enough student accounts`.

Baseline 5 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_5 \
  --summary-export results/baseline-5.json
```

Load 25 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_25 \
  --summary-export results/load-25.json
```

Load 50 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_50 \
  --summary-export results/load-50.json
```

Load 100 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=load_100 \
  --summary-export results/load-100.json
```

## Metrik Utama

- `http_req_duration` P95 harus kurang dari 4000 ms.
- `http_req_failed` harus kurang dari 1%.
- `checks` harus lebih dari 99%.
- `exam_flow_failed` harus kurang dari 1%.

`exam_flow_duration` dicatat sebagai durasi total alur ujian, tetapi batas 4 detik hanya berlaku untuk response time API.

## Rotating User Account

Data akun siswa dibaca dari `k6/data/students.json`.
Setiap iterasi memakai akun siswa berdasarkan indeks iterasi global k6, sehingga akun yang sudah digunakan tidak dipakai ulang pada run yang sama.

Jumlah akun siswa harus lebih besar atau sama dengan total iterasi yang terjadi selama durasi skenario.
Contoh: jika `load_100` menyelesaikan 2.500 iterasi total selama run, maka minimal dibutuhkan 2.500 akun unik.
Script tidak akan mengulang akun saat akun habis.

## Catatan Endpoint

Script memakai route aktual backend saat ini. Route submit mengikuti route yang ada di backend, yaitu `/siswa/uijan/submit/:idAttempt`.
