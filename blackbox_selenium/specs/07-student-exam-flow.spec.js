import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-41 Resume Sesi Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-41-S01] Siswa dapat melanjutkan ujian setelah keluar sesi",
      "[F-41-S02] Timer tetap mengikuti sisa waktu",
      "[F-41-S03] Ujian tidak dapat dilanjutkan setelah sudah submit",
    ],
    "perlu ujian aktif berdurasi pendek dan fixture jawaban siswa"
  );
});

describe("F-42 Pencegahan Tab Baru Saat Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-42-S01] Pembukaan tab baru saat ujian dicegah",
      "[F-42-S02] Pelanggaran tab baru dicatat/ditangani",
      "[F-42-S03] Fitur tidak mengganggu halaman non-ujian",
    ],
    "perlu mode browser non-headless untuk validasi perilaku tab baru secara stabil"
  );
});

describe("F-43 Timer Ujian", function () {
  registerSkippedScenarios(
    [
      "[F-43-S01] Timer ujian tampil sesuai durasi",
      "[F-43-S02] Timer berkurang real-time",
      "[F-43-S03] Ujian otomatis berakhir saat timer habis",
    ],
    "perlu ujian E2E dengan durasi sangat pendek"
  );
});

describe("F-44 Ujian Siswa Berdasarkan Jadwal", function () {
  registerSkippedScenarios(
    [
      "[F-44-S01] Siswa melihat ujian sesuai kelas dan sesi",
      "[F-44-S02] Ujian di luar jadwal tidak tampil/tidak bisa diikuti",
      "[F-44-S03] Ujian kelas lain tidak tampil",
    ],
    "perlu beberapa akun siswa pada kelas berbeda dan jadwal berbeda"
  );
});

describe("F-45 Indikator Jawaban", function () {
  registerSkippedScenarios(
    [
      "[F-45-S01] Indikator soal belum dijawab tampil",
      "[F-45-S02] Indikator berubah setelah soal dijawab",
      "[F-45-S03] Indikator tetap setelah pindah soal/refresh",
    ],
    "perlu selector stabil pada navigator soal ujian"
  );
});

describe("F-48 Review Sebelum Submit", function () {
  registerSkippedScenarios(
    [
      "[F-48-S01] Review jawaban muncul sebelum submit akhir",
      "[F-48-S02] Submit akhir berhasil setelah konfirmasi review",
      "[F-48-S03] Batal submit kembali ke pengerjaan",
    ],
    "perlu ujian aktif disposable dan selector dialog submit"
  );
});
