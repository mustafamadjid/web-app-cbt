export type SoalPreviewOption = {
  id: number;
  label: string;
  text: string;
};

export type SoalPreviewItem = {
  id: number;
  nomor: number;
  tipe: string;
  pertanyaan: string;
  gambar_url?: string;
  opsi: SoalPreviewOption[];
};

export type SoalPreviewData = {
  title: string;
  sisa_waktu: string;
  soal: SoalPreviewItem[];
};
