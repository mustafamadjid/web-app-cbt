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
| `TEST_TYPE` | Ya | `baseline_10` | `baseline_10`, `baseline_20`, `load_25`, `load_50`, atau `load_100` |
| `STUDENTS_FILE` | Tidak | `./data/students.json` | Lokasi fixture akun siswa relatif dari file k6 |
| `ID_JADWAL_UJIAN` | Tidak | kosong | Pakai jadwal spesifik jika ada banyak jadwal |
| `ANSWER_LIMIT` | Tidak | `10` | Jumlah maksimal jawaban yang disimpan per flow |

## Menjalankan Test

Semua skenario memakai executor `per-vu-iterations` dengan `iterations: 1`.
Artinya setiap VU login satu kali, mengerjakan satu flow ujian, lalu berhenti.
`maxDuration` hanya menjadi batas waktu maksimum penyelesaian flow, bukan durasi loop test.

Baseline 10 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_10
```

Baseline 20 user:

```bash
k6 run k6/exam-flow.test.js \
  -e BASE_URL=https://staging-srv.smafi.my.id/ \
  -e TOKEN_UJIAN=ABC123 \
  -e TEST_TYPE=baseline_20 \
  --summary-export results/baseline-20.json
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

## Catatan Endpoint

Script memakai route aktual backend saat ini. Route submit mengikuti route yang ada di backend, yaitu `/siswa/uijan/submit/:idAttempt`.
