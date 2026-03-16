import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";

export type UjianSiswaFilterParams = {
  bulan?: number;
  tahun?: number;
  mapel?: string;
  search?: string;
};

export type UjianSiswaExamItem = JadwalUjianItem & {
  mapel: string;
  id_bank_soal: number;
};

export type UjianSiswaResultItem = UjianSiswaExamItem & {
  jumlah_benar: number;
  jumlah_salah: number;
  nilai: number;
};

export type UjianSiswaResponse = {
  upcoming: UjianSiswaExamItem[];
  ongoing: UjianSiswaExamItem[];
  completed: UjianSiswaResultItem[];
  mapelOptions: string[];
};

export type WaktuSelesaiUjian = {
  id_jadwal_ujian : number,
  waktu_selesai : string
}

export type ActiveAttemptUjian = {
  id_attempt: number;
  id_peserta_ujian: number;
  status_attempt: string;
  waktu_mulai: string | null;
  waktu_submit: string | null;
  deadline_at: string | null;
};

export type JawabanUjianSiswaItem = {
  id_jawaban: number;
  id_soal: number;
  id_pilihan: number | null;
  jawaban_essay: string | null;
  waktu_jawab: string | null;
};

export type JawabanUjianSiswaResponse = {
  id_attempt: number;
  jawaban: JawabanUjianSiswaItem[];
};

export type SaveJawabanUjianSiswaItemRequest = {
  id_soal: number;
  id_pilihan: number | null;
  jawaban_essay: string | null;
  waktu_jawab: string | null;
};

export type SaveJawabanUjianSiswaRequest = {
  id_attempt: number;
  jawaban: SaveJawabanUjianSiswaItemRequest[];
};
