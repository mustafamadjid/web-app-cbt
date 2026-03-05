import type { JadwalUjianStatusClient } from "@/types/Ujian/jadwalUjian";
import type { ReactNode } from "react";

export type PrintJenis = "daftar-hadir" | "berita-acara" | "kartu-peserta";

export type DetailUjianItem = {
  id: number;
  id_ujian: number;
  id_bank_soal: number;
  id_kelas: number;
  id_nama_kelas?: number;
  id_guru: number;
  id_sesi: number;
  id_ruangan: number;
  id_pengawas: number;
  nama_ujian: string;
  deskripsi_ujian: string;
  acak_soal: boolean;
  tanggal_ujian: string;
  tgl_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  durasi_menit: number;
  status_ujian: JadwalUjianStatusClient;
  token: string;
  tingkat_kelas: number;
  nama_kelas: string;
  pengawas_ujian: string;
  ruang_ujian: string;
  sesi_ujian: number;
};

export type InfoItem = {
  label: string;
  value: string;
  icon?: ReactNode;
};
