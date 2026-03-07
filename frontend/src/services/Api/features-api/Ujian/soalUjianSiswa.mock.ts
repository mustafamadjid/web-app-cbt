import { useFetch } from "@/hooks/fetch";
import type { SoalPreviewData } from "@/types/Ujian/SoalPreview";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const buildMockSoalUjianSiswa = (bankSoalId: number): SoalPreviewData => ({
  title: `Ujian ${bankSoalId}`,
  sisa_waktu: "01:35:00",
  soal: [
    {
      id: bankSoalId * 100 + 1,
      nomor: 1,
      tipe: "PILIHAN_GANDA",
      pertanyaan:
        "Fenomena banjir di kota besar disebabkan oleh berkurangnya daerah resapan air, meningkatnya pembangunan, serta kebiasaan masyarakat membuang sampah sembarangan. Apa gagasan utama paragraf tersebut?",
      opsi: [
        {
          id: bankSoalId * 1000 + 11,
          label: "A",
          text: "Banjir dipicu berkurangnya resapan air dan perilaku masyarakat.",
        },
        {
          id: bankSoalId * 1000 + 12,
          label: "B",
          text: "Banjir hanya terjadi karena curah hujan tinggi.",
        },
        {
          id: bankSoalId * 1000 + 13,
          label: "C",
          text: "Permukiman selalu tergenang ketika hujan.",
        },
        {
          id: bankSoalId * 1000 + 14,
          label: "D",
          text: "Pembangunan kota besar tidak memiliki dampak apa pun.",
        },
      ],
    },
    {
      id: bankSoalId * 100 + 2,
      nomor: 2,
      tipe: "PILIHAN_GANDA",
      pertanyaan:
        'Perhatikan kalimat berikut: "Para siswa diminta membuat ringkasan bacaan untuk menilai pemahaman isi teks." Apa tujuan kegiatan tersebut?',
      opsi: [
        {
          id: bankSoalId * 1000 + 21,
          label: "A",
          text: "Menilai pemahaman siswa terhadap isi bacaan.",
        },
        {
          id: bankSoalId * 1000 + 22,
          label: "B",
          text: "Menghafal isi bacaan secara cepat.",
        },
        {
          id: bankSoalId * 1000 + 23,
          label: "C",
          text: "Menyusun daftar pustaka dari bacaan.",
        },
        {
          id: bankSoalId * 1000 + 24,
          label: "D",
          text: "Menyalin kembali seluruh isi teks.",
        },
      ],
    },
    {
      id: bankSoalId * 100 + 3,
      nomor: 3,
      tipe: "ESSAY",
      pertanyaan:
        "Jelaskan mengapa ringkasan dapat membantu guru menilai pemahaman siswa terhadap sebuah teks.",
      opsi: [],
    },
  ],
});

export async function getMockSoalUjianSiswa(
  bankSoalId: number,
): Promise<SoalPreviewData> {
  if (!Number.isInteger(bankSoalId) || bankSoalId <= 0) {
    throw new Error("Bank soal tidak ditemukan.");
  }

  await sleep(150);
  return buildMockSoalUjianSiswa(bankSoalId);
}

export function useGetMockSoalUjianSiswa(
  bankSoalId: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getMockSoalUjianSiswa(bankSoalId)
        : Promise.resolve(null as SoalPreviewData | null),
    [bankSoalId, enabled],
  );
}
