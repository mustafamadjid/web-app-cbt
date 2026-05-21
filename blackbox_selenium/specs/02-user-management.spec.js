import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-01 Registrasi Akun", function () {
  registerSkippedScenarios(
    [
      "[F-01-S01] Registrasi akun valid untuk setiap role",
      "[F-01-S02] Registrasi ditolak saat data wajib kosong",
      "[F-01-S03] Registrasi ditolak saat username/email duplikat",
    ],
    "perlu finalisasi selector form tambah akun admin/guru/siswa dan seed akun E2E"
  );
});

describe("F-02 Pengelolaan Data Pengguna", function () {
  registerSkippedScenarios(
    [
      "[F-02-S01] Tambah data pengguna baru",
      "[F-02-S02] Lihat daftar dan detail pengguna",
      "[F-02-S03] Hapus/nonaktifkan pengguna",
    ],
    "perlu data-testid atau selector stabil pada tabel dan dialog hapus akun"
  );
});

describe("F-03 Perubahan Data Guru dan Siswa", function () {
  registerSkippedScenarios(
    [
      "[F-03-S01] Ubah data guru dengan input valid",
      "[F-03-S02] Ubah data siswa dengan input valid",
      "[F-03-S03] Perubahan ditolak saat data tidak valid",
    ],
    "perlu data akun disposable dan selector stabil pada halaman edit"
  );
});
