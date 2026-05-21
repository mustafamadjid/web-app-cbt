import { registerSkippedScenarios } from "../helpers/spec.helper.js";

describe("F-08 Total Mata Pelajaran Dashboard", function () {
  registerSkippedScenarios(
    [
      "[F-08-S01] Total mata pelajaran tampil sesuai data",
      "[F-08-S02] Total berubah setelah mapel ditambah",
      "[F-08-S03] Total berubah setelah mapel dihapus/nonaktif",
    ],
    "perlu seed mapel E2E dan selector angka dashboard yang stabil"
  );
});

describe("F-10 Data Jurusan", function () {
  registerSkippedScenarios(
    [
      "[F-10-S01] Tambah jurusan valid",
      "[F-10-S02] Daftar jurusan ditampilkan",
      "[F-10-S03] Jurusan duplikat ditolak",
    ],
    "route jurusan belum ditemukan pada frontend"
  );
});

describe("F-11 Data Kelas", function () {
  registerSkippedScenarios(
    [
      "[F-11-S01] Tambah kelas valid",
      "[F-11-S02] Daftar kelas ditampilkan",
      "[F-11-S03] Kelas tanpa data wajib ditolak",
    ],
    "perlu selector stabil untuk form tingkat kelas dan nama kelas"
  );
});

describe("F-12 Data Ruangan", function () {
  registerSkippedScenarios(
    [
      "[F-12-S01] Tambah ruangan valid",
      "[F-12-S02] Daftar ruangan ditampilkan",
      "[F-12-S03] Kapasitas ruangan tidak valid ditolak",
    ],
    "perlu selector stabil pada form ruang ujian dan tabel"
  );
});

describe("F-13 Data Sesi", function () {
  registerSkippedScenarios(
    [
      "[F-13-S01] Tambah sesi valid",
      "[F-13-S02] Daftar sesi ditampilkan",
      "[F-13-S03] Rentang waktu sesi tidak valid ditolak",
    ],
    "perlu selector stabil pada form sesi dan validasi waktu"
  );
});

describe("F-35 Mata Pelajaran", function () {
  registerSkippedScenarios(
    [
      "[F-35-S01] Tambah mata pelajaran valid",
      "[F-35-S02] Daftar mata pelajaran ditampilkan",
      "[F-35-S03] Mata pelajaran duplikat ditolak",
    ],
    "perlu selector stabil pada form dan tabel mata pelajaran"
  );
});

describe("F-36 Modifikasi Mata Pelajaran", function () {
  registerSkippedScenarios(
    [
      "[F-36-S01] Modifikasi mata pelajaran valid",
      "[F-36-S02] Modifikasi mapel ditolak bila duplikat",
      "[F-36-S03] Perubahan mapel tercermin pada form ujian",
    ],
    "perlu data mapel disposable dan selector edit yang stabil"
  );
});
