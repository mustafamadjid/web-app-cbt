import type { RichContent } from "@/types/Content/RichContent";

export type HasilJawabanUjianOpsi = {
  id_pilihan_ganda: number;
  isi_pilihan: string;
  isi_pilihan_content?: RichContent | null;
  is_benar: boolean;
};

export type HasilJawabanUjianJawabanSiswa = {
  id_jawaban: number | null;
  id_pilihan: number | null;
  jawaban_essay: string | null;
  waktu_jawab: string | null;
  essay_is_benar: boolean | null;
};

export type HasilJawabanUjianItem = {
  id_soal: number;
  id_bank_soal_version: number;
  tipe_soal: string;
  pertanyaan: string;
  pertanyaan_content?: RichContent | null;
  gambar: string;
  bobot_soal: number;
  no_urut_soal: number;
  opsi_jawaban: HasilJawabanUjianOpsi[];
  jawaban_siswa: HasilJawabanUjianJawabanSiswa;
};

export type HasilJawabanUjianResponse = {
  id_attempt: number;
  nilai_akhir: number | null;
  hasil_jawaban: HasilJawabanUjianItem[];
};
