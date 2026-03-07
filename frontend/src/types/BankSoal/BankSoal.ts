export type BankSoalApiItem = {
  id_bank_soal: number;
  id_mapel: number;
  id_kelas: number;
  id_pengguna: number;
  mapel: string;
  guru_pembuat: string;
  kelas: string;
  nama_bank_soal: string;
  deskripsi: string;
  materi: string;
  tanggal_dibuat: string;
  soal_uploaded: boolean;
};

export type BankSoalItem = BankSoalApiItem;

export type GetBankSoalParams = {
  search?: string;
  id_kelas?: number;
  id_mapel?: number;
  limit?: number;
  offset?: number;
};

export type BankSoalFormValues = {
  namaBankSoal: string;
  kelasId: number | "";
  mapelId: number | "";
  materi: string;
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
