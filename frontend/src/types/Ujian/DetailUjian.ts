import type { JadwalUjianStatusClient } from "@/types/Ujian/jadwalUjian";
import type { ReactNode } from "react";

export type PrintJenis = "daftar-hadir" | "berita-acara" | "kartu-peserta";

export type DetailUjianItem = {
  id: number;
  id_ujian: number;
  id_guru: number;
  id_pengawas: number;
  nama_ujian: string;
  pengawas_ujian: string;
  tgl_ujian: string;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  sesi_ujian: number;
  ruang_ujian: string;
  id_ruang: number;
  status_ujian: JadwalUjianStatusClient;
  started: 0 | 1;
  tingkat_kelas: number;
  tingkat_kelas_id: number;
  nama_kelas: string;
  pembuat_username: string;
  pengawas_username: string;
  deskripsi_ujian: string;
  token: string;
};

export type InfoItem = {
  label: string;
  value: string;
  icon?: ReactNode;
};
