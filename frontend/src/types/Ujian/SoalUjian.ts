import type { RichContent } from "@/types/Content/RichContent";

export type OpsiJawabanSoalUjian = {
  id_pilihan_ganda: number;
  id_soal: number;
  isi_pilihan: string;
  isi_pilihan_content?: RichContent | null;
  is_benar: boolean;
};

export type SoalUjian = {
  id_soal: number;
  id_bank_soal_version: number;
  tipe_soal: string;
  pertanyaan: string;
  pertanyaan_content?: RichContent | null;
  gambar: string;
  bobot_soal: number;
  no_urut_soal: number;
  opsi_jawaban: OpsiJawabanSoalUjian[];
};

export type OpsiJawabanSoalUjianSiswa = {
  id_pilihan_ganda: number;
  isi_pilihan: string;
  isi_pilihan_content?: RichContent | null;
};

export type SoalUjianSiswa = {
  id_soal: number;
  tipe_soal: string;
  pertanyaan: string;
  pertanyaan_content?: RichContent | null;
  gambar: string;
  bobot_soal: number;
  no_urut_soal: number;
  opsi_jawaban: OpsiJawabanSoalUjianSiswa[];
};
