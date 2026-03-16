import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type { SoalUjianSiswa } from "@/types/Ujian/SoalUjian";

const OPTION_LABELS = ["A", "B", "C", "D", "E", "F", "G", "H"];

export function buildSoalPreview(
  soalRows: SoalUjianSiswa[] | null | undefined,
): SoalPreviewItem[] {
  return (soalRows ?? []).map((soal, index) => ({
    id: soal.id_soal,
    nomor: index + 1,
    tipe: soal.tipe_soal,
    pertanyaan: soal.pertanyaan,
    gambar_url: resolveImageUrl(soal.gambar) || undefined,
    opsi: soal.opsi_jawaban.map((opsi, optionIndex) => ({
      id: opsi.id_pilihan_ganda,
      label: OPTION_LABELS[optionIndex] ?? String(optionIndex + 1),
      text: opsi.isi_pilihan,
    })),
  }));
}
