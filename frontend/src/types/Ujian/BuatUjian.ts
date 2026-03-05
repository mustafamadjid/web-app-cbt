import type { StatusAkun } from "../OpsiTypes/Option";

export type TipeUjian = "PILIHAN_GANDA" | "ESSAY" | "CAMPURAN";
export type KelasScope = "SEMUA" | "SPESIFIK";
export type StatusUjianServer =
  | "BELUM_MULAI"
  | "MULAI"
  | "SELESAI"
  | "DIBATALKAN";

export type BuatUjianFormValues = {
  nama_ujian: string;
  deskripsi_ujian: string;
  id_kelas: number;
  kelas_scope: KelasScope;
  id_nama_kelas: number;
  id_bank_soal: number;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  id_ruangan: number;
  acak_soal: boolean;
  id_pengawas: number;
  id_sesi: number;
  token: string;
};

export type CreatePenjadwalanUjianPayload = {
  id_bank_soal: number;
  id_kelas: number;
  id_nama_kelas?: number;
  id_guru: number;
  nama_ujian: string;
  deskripsi_ujian?: string;
  acak_soal: boolean;
  id_sesi: number;
  id_ruangan: number;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  status_ujian: StatusUjianServer;
  token: string;
  id_pengawas: number;
};

export type UpdatePenjadwalanUjianPayload = Partial<CreatePenjadwalanUjianPayload>;

export type BankSoalOption = {
  id: number;
  nama: string;
  mata_pelajaran?: string;
  materi?: string;
  kelas?: number;
  total_soal: number;
  jumlah_soal_pg?: number;
  jumlah_soal_essay?: number;
};

export type GuruPengawasOption = {
  id: number;
  nama: string;
  nip?: string;
  mapel?: string;
};

export type SesiUjianOption = {
  id: number;
  kode: string;
  nama: string;
};



export type SiswaPreviewItem = {
  id: number;
  nama: string;
  username: string;
  no_absen: number;
  kelas: string;
  status_akun: StatusAkun;
};

export type BuatUjianSubmitResponse = {
  id: number;
};
