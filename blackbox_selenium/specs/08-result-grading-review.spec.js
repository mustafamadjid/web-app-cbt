import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-27 Analisis Soal Paling Banyak Salah", function () {
  registerSkippedScenarios(
    [
      "[F-27-S01] Identifikasi soal paling banyak salah",
      "[F-27-S02] Perhitungan salah diperbarui setelah koreksi",
      "[F-27-S03] Analisis kosong saat belum ada jawaban",
    ],
    "perlu attempt siswa dan hasil koreksi deterministik"
  );
});

describe("F-28 Analisis Soal Paling Banyak Benar", function () {
  registerSkippedScenarios(
    [
      "[F-28-S01] Identifikasi soal paling banyak benar",
      "[F-28-S02] Perhitungan benar diperbarui setelah koreksi",
      "[F-28-S03] Analisis benar kosong saat belum ada jawaban",
    ],
    "perlu attempt siswa dan hasil koreksi deterministik"
  );
});

describe("F-38 Penilaian Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-38-S01] Penilaian otomatis pilihan ganda",
      "[F-38-S02] Penilaian gabungan PG dan essay",
      "[F-38-S03] Penilaian belum final saat essay belum dikoreksi",
    ],
    "perlu ujian PG dan campuran PG/essay yang sudah dikerjakan"
  );
});

describe("F-39 Hasil Nilai Siswa", function () {
  registerSkippedScenarios(
    [
      "[F-39-S01] Siswa melihat nilai ujian yang diikuti",
      "[F-39-S02] Siswa tidak melihat nilai ujian yang tidak diikuti",
      "[F-39-S03] Nilai belum tersedia ditampilkan dengan status",
    ],
    "perlu akun siswa dengan riwayat ujian berbeda"
  );
});

describe("F-40 Koreksi Essay Manual", function () {
  registerSkippedScenarios(
    [
      "[F-40-S01] Guru menandai jawaban essay benar",
      "[F-40-S02] Guru menandai jawaban essay salah",
      "[F-40-S03] Koreksi essay ditolak untuk guru tidak berwenang",
    ],
    "perlu ujian essay submitted dan selector stabil halaman koreksi"
  );
});

describe("F-46 Review Jawaban Benar dan Salah", function () {
  registerSkippedScenarios(
    [
      "[F-46-S01] Review jawaban benar dan salah ditampilkan",
      "[F-46-S02] Jawaban benar tampil pada soal yang salah",
      "[F-46-S03] Review tidak tampil sebelum ujian dikoreksi/final",
    ],
    "perlu hasil final dan hasil belum final untuk siswa"
  );
});

describe("F-47 Nilai Berdasarkan Hasil Koreksi", function () {
  registerSkippedScenarios(
    [
      "[F-47-S01] Nilai tampil setelah koreksi selesai",
      "[F-47-S02] Nilai berubah setelah koreksi diperbarui",
      "[F-47-S03] Nilai tidak tampil bila koreksi belum lengkap",
    ],
    "perlu ujian essay yang bisa dikoreksi ulang secara aman"
  );
});
