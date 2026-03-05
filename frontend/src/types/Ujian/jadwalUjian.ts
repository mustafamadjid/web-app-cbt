export type JadwalUjianStatusClient =
  | "belum_dimulai"
  | "berlangsung"
  | "selesai"
  | "dibatalkan";

export type JadwalUjianItem = {
  id: number;
  id_ujian?: number;
  id_guru?: number;
  id_pengawas?: number;
  nama_ujian: string;
  pengawas_ujian: string;
  tgl_ujian: string;
  tanggal_ujian?: string;
  waktu_mulai: string;
  waktu_selesai?: string;
  sesi_ujian: number;
  ruang_ujian: string;
  id_ruang?: number;
  status_ujian: JadwalUjianStatusClient;
  started: 0 | 1;
  pembuat_username?: string;
  pengawas_username?: string;
  tingkat_kelas?: number;
  tingkat_kelas_id?: number;
  nama_kelas?: string;
};

export type JadwalUjianFilterParams = {
  search?: string;
  tanggal?: string; // "YYYY-MM-DD"
  tingkatKelasId?: number;
  ruangUjianId?: number;
  tahun?: string | number;
};
