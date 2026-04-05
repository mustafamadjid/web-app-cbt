export type AttemptUjianRequest = {
  id_siswa: number;
  id_jadwal_ujian: number;
  token_ujian: string;
  waktu_mulai: string;
};

export type PesertaUjianSubmittedItem = {
  id_peserta_ujian: number;
  id_attempt: number;
  id_siswa: number;
  tingkat_kelas: number;
  nama_kelas: string;
  nama_lengkap: string;
  no_absen: number;
  nilai_akhir: number | null;
  waktu_mulai: string | null;
  waktu_submit: string | null;
};

export type ListPesertaUjianSubmittedResponse = PesertaUjianSubmittedItem[];
