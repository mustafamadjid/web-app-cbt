import type { TipeUjian } from "../Ujian/BuatUjian";

export type BankSoalApiItem = {
  id_bank_soal: number;
  id_mapel: number;
  id_kelas: number;
  id_pengguna: number;
  nama_bank_soal: string;
  deskripsi: string;
  materi: string;
};

export type BankSoalItem = BankSoalApiItem;

export type GetBankSoalParams = {
  search?: string;
  id_kelas?: number;
  id_mapel?: number;
  limit?: number;
  offset?: number;
};

export type SoalUjianItem = {
  id_soal: number;
  nomor_urut_soal: number;
  tipe_soal: TipeUjian;
  pertanyaan: string;
  urlGambar?: string;
  bobot?: string;
  opsi_a?: string;
  opsi_b?: string;
  opsi_c?: string;
  opsi_d?: string;
  opsi_e?: string;
  jawaban?: string;
};

export type BankSoalFormValues = {
  namaBankSoal: string;
  kelasId: number | "";
  mapelId: number | "";
  deskripsi: string;
};

export type CreateBankSoalPayload = {
  id_pengguna?: number;
  id_mapel: number;
  id_kelas: number;
  nama_bank_soal: string;
  deskripsi: string;
  materi: string;
};

export type UpdateBankSoalPayload = Partial<CreateBankSoalPayload>;

export type SoalUjianResponse = {
  id_bank_soal: number;
  nama_ujian: string;
  mata_pelajaran?: string;
  jumlah_soal: number;
  sisa_waktu?: string;
  soal: SoalUjianItem[];
};
