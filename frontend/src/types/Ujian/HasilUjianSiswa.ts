import type { JadwalUjianStatusClient } from "./jadwalUjian";

export type HasilUjianSiswaItem = {
  id: number;
  id_attempt: number;
  id_ujian: number;
  id_bank_soal: number;
  id_guru: number;
  id_pengawas: number;
  nama_ujian: string;
  pengawas_ujian: string;
  tgl_ujian: string;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  sesi_ujian: number;
  nama_sesi: string;
  ruang_ujian: string;
  id_ruang: number;
  status_ujian: JadwalUjianStatusClient;
  started: 0 | 1;
  tingkat_kelas: number;
  tingkat_kelas_id: number;
  nama_kelas: string;
  pengawas_nama_lengkap: string;
  deskripsi_ujian: string;
  acak_soal: boolean;
};

export type HasilUjianSiswaResponse = HasilUjianSiswaItem[];
