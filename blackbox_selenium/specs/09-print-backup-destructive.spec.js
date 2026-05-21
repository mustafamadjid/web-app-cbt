import { env } from "../config/env.js";
import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-15 Backup Data", function () {
  registerSkippedScenarios(
    [
      "[F-15-S01] Backup data berhasil dibuat",
      "[F-15-S02] Backup dapat diunduh",
      "[F-15-S03] Backup gagal ditangani dengan pesan",
    ],
    "route backup belum ditemukan pada frontend atau perlu mode backend khusus untuk simulasi gagal"
  );
});

describe("F-20 Pengosongan Data Destructive", function () {
  before(function () {
    if (!env.allowDestructive) this.skip();
  });

  registerSkippedScenarios(
    [
      "[F-20-S01] Kosongkan data yang dipilih berhasil",
      "[F-20-S02] Batal konfirmasi tidak menghapus data",
      "[F-20-S03] Pengosongan tanpa pilihan ditolak",
    ],
    "fitur pengosongan data perlu route UI dan database E2E khusus"
  );
});

describe("F-22 Cetak Daftar Hadir Destructive/Download", function () {
  registerSkippedScenarios(
    [
      "[F-22-S01] Cetak daftar hadir dengan filter valid",
      "[F-22-S02] Daftar hadir hanya berisi peserta sesuai filter",
      "[F-22-S03] Cetak ditolak bila filter wajib kosong",
    ],
    "dicakup di spec manajemen ujian setelah seed jadwal tersedia"
  );
});
