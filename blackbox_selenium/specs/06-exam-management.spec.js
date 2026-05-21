import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-22 Cetak Daftar Hadir", function () {
  registerSkippedScenarios(
    [
      "[F-22-S01] Cetak daftar hadir dengan filter valid",
      "[F-22-S02] Daftar hadir hanya berisi peserta sesuai filter",
      "[F-22-S03] Cetak ditolak bila filter wajib kosong",
    ],
    "perlu seed jadwal ujian dan validasi file download"
  );
});

describe("F-30 Buat Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-30-S01] Buat ujian baru valid",
      "[F-30-S02] Buat ujian ditolak bila mapel/kelas kosong",
      "[F-30-S03] Ujian baru hanya memakai data tersedia",
    ],
    "perlu data master dan bank soal E2E yang deterministik"
  );
});

describe("F-31 Modifikasi Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-31-S01] Modifikasi ujian terjadwal valid",
      "[F-31-S02] Modifikasi ujian telah dilaksanakan sesuai aturan",
      "[F-31-S03] Modifikasi ditolak bila jadwal tidak valid",
    ],
    "perlu ujian terjadwal dan selesai yang disposable"
  );
});

describe("F-32 Hapus Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-32-S01] Hapus ujian berhasil",
      "[F-32-S02] Batal hapus ujian",
      "[F-32-S03] Hapus ujian ditolak untuk pengguna tidak berwenang",
    ],
    "perlu ujian disposable dan selector dialog konfirmasi"
  );
});

describe("F-33 Filter Ujian Berdasarkan Mata Pelajaran", function () {
  registerSkippedScenarios(
    [
      "[F-33-S01] Tampilkan ujian berdasarkan mata pelajaran",
      "[F-33-S02] Daftar memuat ujian akan dan telah berlangsung",
      "[F-33-S03] Filter tanpa hasil menampilkan keadaan kosong",
    ],
    "perlu seed beberapa ujian lintas status dan mapel"
  );
});

describe("F-34 Token Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-34-S01] Buat token ujian valid",
      "[F-34-S02] Token kosong atau tidak valid ditolak",
      "[F-34-S03] Token baru menggantikan token lama",
    ],
    "perlu selector stabil untuk field token pada form ujian"
  );
});

describe("F-37 Pengumuman", function () {
  registerSkippedScenarios(
    [
      "[F-37-S01] Buat pengumuman valid",
      "[F-37-S02] Pengumuman wajib isi ditolak bila kosong",
      "[F-37-S03] Pengumuman hanya tampil pada target penerima",
    ],
    "perlu selector stabil editor pengumuman dan data target penerima"
  );
});
