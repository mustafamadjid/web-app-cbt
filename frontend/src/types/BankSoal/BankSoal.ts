import type { TipeUjian } from "../Ujian/BuatUjian";

export type BankSoalItem = {
    id: number;
    guru?: string;
    nama_banksoal?: string;
    mata_pelajaran?: string;
    materi?: string;
    kelas?: number;
    deskripsi?: string;
    tgl_buat?: string;
    jumlah_soal_pg?: number;
    jumlah_soal_essay?: number;
    total_soal?: number;
};

export type SoalUjianItem = {
    id_soal: number;
    tipe_soal:TipeUjian;
    pertanyaan: string;
    urlGambar?:string;
    bobot?:string;
    opsi_a?:string;
    opsi_b?:string;
    opsi_c?:string;
    opsi_d?:string;
    opsi_e?:string;
    jawaban?:string;
    
}

export type SoalUjianResponse = {
  id_bank_soal: number;
  nama_ujian: string;
  mata_pelajaran?: string;
  jumlah_soal: number;
  sisa_waktu?: string;
  soal: SoalUjianItem[];
};