export type TipeUjian = "PILIHAN_GANDA" | "ESSAY" | "CAMPURAN";

export type BuatUjianFormValues = {
  nama_ujian: string;
  deskripsi_ujian: string;
  tipe_ujian: TipeUjian;
  kelas_id: number | "";
  kelas_scope?: "SEMUA" | "SPESIFIK";
  kelas_detail_id: number | "";
  bank_soal_id: number | "";
  jumlah_soal: number;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  durasi_menit: number;
  ruang_ujian_id: number | "";
  acak_soal: boolean;
  guru_pengawas_id: number | "";
  sesi_id: number | "";
  token_ujian: string;
  // Nanti status dibuat default tergantung jam jadwal
  // nanti started dibuat default dengan nilai 0
};

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
  status_akun: "aktif" | "nonaktif" | "dibekukan";
};

export type BuatUjianSubmitResponse = {
  id: number;
};
