import type { RichContent } from "@/types/Content/RichContent";

export type AnalisisSoalOpsiResponse = {
  id_pilihan_ganda: number;
  isi_pilihan: string;
  isi_pilihan_content?: RichContent | null;
  is_benar: boolean;
};

export type AnalisisSoalItem = {
  id_soal: number;
  id_bank_soal_version: number;
  tipe_soal: string;
  pertanyaan: string;
  pertanyaan_content?: RichContent | null;
  gambar: string;
  bobot_soal: number;
  no_urut_soal: number;
  jumlah_jawaban_benar: number;
  jumlah_jawaban_salah: number;
  opsi_jawaban: AnalisisSoalOpsiResponse[];
};

export type AnalisisSoalResponse = {
  id_jadwal_ujian: number;
  analisis_soal: AnalisisSoalItem[];
};
