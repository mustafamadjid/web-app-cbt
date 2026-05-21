import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-05 Total Peserta/Siswa", function () {
  registerSkippedScenarios(
    [
      "[F-05-S01] Jumlah siswa tampil sesuai data",
      "[F-05-S02] Jumlah siswa berubah setelah data ditambah",
      "[F-05-S03] Jumlah siswa nol saat belum ada data",
    ],
    "perlu database E2E terkontrol untuk validasi angka dashboard"
  );
});

describe("F-07 Total Ujian Terlaksana", function () {
  registerSkippedScenarios(
    [
      "[F-07-S01] Total ujian terlaksana tampil sesuai riwayat",
      "[F-07-S02] Ujian belum berlangsung tidak dihitung",
      "[F-07-S03] Total nol saat belum ada ujian selesai",
    ],
    "perlu seed ujian selesai dan belum mulai yang deterministik"
  );
});

describe("F-09 Profil Sekolah", function () {
  registerSkippedScenarios(
    [
      "[F-09-S01] Simpan profil sekolah valid",
      "[F-09-S02] Validasi profil sekolah wajib",
      "[F-09-S03] Tampilan profil tetap setelah refresh",
    ],
    "perlu selector stabil pada form pengaturan profil sekolah"
  );
});

describe("F-14 Data Umum Aplikasi", function () {
  registerSkippedScenarios(
    [
      "[F-14-S01] Simpan data umum aplikasi valid",
      "[F-14-S02] Validasi data umum aplikasi",
      "[F-14-S03] Data umum tetap setelah login ulang",
    ],
    "fitur data umum selain profil sekolah perlu dipastikan di UI"
  );
});

describe("F-17 Log Aktivitas", function () {
  registerSkippedScenarios(
    [
      "[F-17-S01] Log aktivitas siswa tersimpan",
      "[F-17-S02] Log aktivitas guru tersimpan",
      "[F-17-S03] Log aktivitas operator/admin tersimpan",
    ],
    "route log aktivitas admin belum terlihat di router frontend"
  );
});

describe("F-18 Audit Log", function () {
  registerSkippedScenarios(
    [
      "[F-18-S01] Audit log tercatat saat login berhasil",
      "[F-18-S02] Audit log tercatat saat login gagal",
      "[F-18-S03] Audit log tidak dapat diubah oleh pengguna biasa",
    ],
    "route audit log belum terlihat di router frontend"
  );
});

describe("F-19 Tampilan Log User", function () {
  registerSkippedScenarios(
    [
      "[F-19-S01] Administrator melihat daftar log user",
      "[F-19-S02] Filter log berdasarkan pengguna/tanggal",
      "[F-19-S03] Akses log ditolak untuk non-admin",
    ],
    "route log user belum terlihat di router frontend"
  );
});

describe("F-20 Pengosongan Data", function () {
  registerSkippedScenarios(
    [
      "[F-20-S01] Kosongkan data yang dipilih berhasil",
      "[F-20-S02] Batal konfirmasi tidak menghapus data",
      "[F-20-S03] Pengosongan tanpa pilihan ditolak",
    ],
    "destructive test hanya boleh diaktifkan pada database E2E khusus"
  );
});

describe("F-21 Riwayat Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-21-S01] Riwayat ujian tersimpan setelah ujian selesai",
      "[F-21-S02] Riwayat ujian ditampilkan sesuai pengguna/role",
      "[F-21-S03] Riwayat kosong ditampilkan dengan aman",
    ],
    "perlu seed attempt ujian selesai dan akun siswa tanpa riwayat"
  );
});
