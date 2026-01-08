import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";

export const dummyJadwalUjian: JadwalUjianItem[] = [
  {
    id: 1,
    nama_ujian: "Ujian Matematika",
    pengawas_ujian: "Pak Budi Santoso",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "07:30",
    sesi_ujian: 1,
    ruang_ujian: "Ruang 101",
    status_ujian: "belum_dimulai",
    tingkat_kelas: 10,
    nama_kelas: "10 IPA 1",
  },
  {
    id: 2,
    nama_ujian: "Ujian Bahasa Indonesia",
    pengawas_ujian: "Bu Siti Nuraini",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "10:00",
    sesi_ujian: 2,
    ruang_ujian: "Ruang 102",
    status_ujian: "berlangsung",
    tingkat_kelas: 11,
    nama_kelas: "11 IPA 2",
  },
  {
    id: 3,
    nama_ujian: "Ujian IPA",
    pengawas_ujian: "Pak Andi Pratama",
    tgl_ujian: "Selasa, 13 Februari 2026",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
    ruang_ujian: "Lab IPA",
    status_ujian: "selesai",
    tingkat_kelas: 12,
    nama_kelas: "12 IPS 1",
  },
  {
    id: 4,
    nama_ujian: "Ujian Bahasa Inggris",
    pengawas_ujian: "Bu Rina Oktavia",
    tgl_ujian: "Rabu, 14 Februari 2026",
    waktu_mulai: "09:30",
    ruang_ujian: "Ruang 203",
    status_ujian: "belum_dimulai",
    tingkat_kelas: 11,
    nama_kelas: "11 IPS 2",
  },
  {
    id: 5,
    nama_ujian: "Ujian Sejarah",
    pengawas_ujian: "Pak Dedi Kurniawan",
    tgl_ujian: "Kamis, 15 Februari 2026",
    waktu_mulai: "13:00",
    sesi_ujian: 3,
    tingkat_kelas: 10,
    nama_kelas: "10 IPA 3",
  },
];

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
export async function getJadwalUjian (): Promise<JadwalUjianItem[]> {
  await sleep(250);
  return dummyJadwalUjian;
}