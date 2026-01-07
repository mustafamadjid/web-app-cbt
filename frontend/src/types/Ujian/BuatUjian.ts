export type TipeUjian = "PILIHAN_GANDA" | "ESSAY" | "CAMPURAN";

export type BuatUjianFormValues = {
  nama_ujian: string;
  deskripsi_ujian: string;
  tipe_ujian: TipeUjian;
  kelas_id: number | "";
  bank_soal_id: string;
  jumlah_soal: number;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  durasi_menit: number;
  ruang_ujian_id: string;
  acak_soal: boolean;
  guru_pengawas_id: string;
  sesi_id: string;
  token_ujian: string;
};

export type BankSoalOption = {
  id: string;
  nama: string;
  mata_pelajaran?: string;
  materi?: string;
  kelas?: number;
  total_soal: number;
  jumlah_soal_pg?: number;
  jumlah_soal_essay?: number;
};

export type GuruPengawasOption = {
  id: string;
  nama: string;
  nip?: string;
  mapel?: string;
};

export type SesiUjianOption = {
  id: string;
  kode: string;
  nama: string;
};

export type RuangUjianOption = {
  id: string;
  nama: string;
};

export type SiswaPreviewItem = {
  id: string;
  nama: string;
  username: string;
  no_absen: number;
  kelas: string;
  status_akun: "aktif" | "nonaktif" | "dibekukan";
};

export type BuatUjianSubmitResponse = {
  id: number;
};
