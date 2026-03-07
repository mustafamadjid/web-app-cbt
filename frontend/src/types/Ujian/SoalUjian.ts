export type OpsiJawabanSoalUjian = {
  id_pilihan_ganda: number;
  id_soal: number;
  isi_pilihan: string;
  is_benar: boolean;
};

export type SoalUjian = {
  id_soal: number;
  id_bank_soal_version: number;
  tipe_soal: string;
  pertanyaan: string;
  gambar: string;
  bobot_soal: number;
  no_urut_soal: number;
  opsi_jawaban: OpsiJawabanSoalUjian[];
};
