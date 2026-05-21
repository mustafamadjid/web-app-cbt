import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-06 Total Bank Soal", function () {
  registerSkippedScenarios(
    [
      "[F-06-S01] Jumlah bank soal tampil sesuai data",
      "[F-06-S02] Jumlah bank soal bertambah setelah import/tambah",
      "[F-06-S03] Jumlah bank soal berkurang setelah penghapusan",
    ],
    "perlu seed bank soal disposable dan selector angka dashboard"
  );
});

describe("F-23 Import Bank Soal", function () {
  registerSkippedScenarios(
    [
      "[F-23-S01] Import bank soal pilihan ganda dari docx valid",
      "[F-23-S02] Import bank soal essay dari docx valid",
      "[F-23-S03] Import ditolak untuk format file tidak valid",
    ],
    "perlu fixture docx valid final sesuai parser backend"
  );
});

describe("F-24 Hapus Bank Soal", function () {
  registerSkippedScenarios(
    [
      "[F-24-S01] Hapus bank soal berhasil",
      "[F-24-S02] Batal hapus bank soal",
      "[F-24-S03] Hapus ditolak untuk pengguna tidak berwenang",
    ],
    "perlu data bank soal disposable dan selector dialog konfirmasi"
  );
});

describe("F-25 Modifikasi Bank Soal", function () {
  registerSkippedScenarios(
    [
      "[F-25-S01] Modifikasi bank soal valid",
      "[F-25-S02] Modifikasi ditolak bila field wajib kosong",
      "[F-25-S03] Perubahan bank soal tidak disimpan saat batal",
    ],
    "perlu selector stabil editor bank soal"
  );
});

describe("F-26 Pengacakan Soal", function () {
  registerSkippedScenarios(
    [
      "[F-26-S01] Nomor soal diacak saat ujian dimulai",
      "[F-26-S02] Urutan pilihan jawaban diacak",
      "[F-26-S03] Urutan tetap saat fitur acak nonaktif",
    ],
    "perlu dua akun siswa dan ujian E2E khusus dengan konfigurasi acak"
  );
});

describe("F-29 Bobot Soal", function () {
  registerSkippedScenarios(
    [
      "[F-29-S01] Simpan bobot soal valid",
      "[F-29-S02] Bobot tidak valid ditolak",
      "[F-29-S03] Nilai mengikuti bobot soal",
    ],
    "perlu ujian fixture dengan bobot soal deterministik"
  );
});
