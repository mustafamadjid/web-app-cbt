import type { RichContent } from "@/types/Content/RichContent";

export type SoalPreviewOption = {
  id: number;
  label: string;
  text: string;
  content?: RichContent | null;
};

export type SoalPreviewItem = {
  id: number;
  nomor: number;
  tipe: string;
  pertanyaan: string;
  pertanyaan_content?: RichContent | null;
  gambar_url?: string;
  opsi: SoalPreviewOption[];
};

export type SoalPreviewData = {
  title: string;
  sisa_waktu: string;
  soal: SoalPreviewItem[];
};
