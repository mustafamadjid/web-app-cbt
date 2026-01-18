import { api, type ApiEnvelope } from "@/services/Api/api";
import type { SoalUjianItem } from "@/types/BankSoal/BankSoal";

type SoalUjianResponse = {
  id_bank_soal: number;
  nama_ujian: string;
  mata_pelajaran?: string;
  jumlah_soal: number;
  sisa_waktu?: string;
  soal: SoalUjianItem[];
};

const dummySoalUjian: Record<number, SoalUjianResponse> = {
  1: {
    id_bank_soal: 1,
    nama_ujian: "Soal Bahasa Indonesia",
    mata_pelajaran: "Bahasa Indonesia",
    jumlah_soal: 30,
    sisa_waktu: "01:35:00",
    soal: [
      {
        id_soal: 101,
        tipe_soal: "Teks pendek",
        pertanyaan:
          "Fenomena banjir di kota besar disebabkan oleh berkurangnya daerah resapan air, meningkatnya pembangunan, serta kebiasaan masyarakat membuang sampah sembarangan. Ketika curah hujan tinggi, air tidak dapat terserap dengan baik sehingga meluap dan menggenangi permukaan jalan maupun permukiman. Apa gagasan utama paragraf tersebut?",
        opsi_a:
          "Fenomena banjir di kota besar disebabkan oleh berkurangnya daerah resapan air",
        opsi_b:
          "Meningkatnya pembangunan, serta kebiasaan masyarakat membuang sampah sembarangan",
        opsi_c:
          "Ketika curah hujan tinggi, air tidak dapat terserap dengan baik sehingga meluap",
        opsi_d:
          "Air tidak dapat terserap dengan baik sehingga meluap dan menggenangi permukaan jalan maupun permukiman",
        jawaban: "A",
      },
    ],
  },
};

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

const USE_DUMMY = true;

export async function getSoalUjian(
  bankSoalId: number
): Promise<SoalUjianResponse> {
  if (USE_DUMMY) {
    await sleep(200);
    const data = dummySoalUjian[bankSoalId];
    if (!data) {
      throw new Error("Soal ujian tidak ditemukan.");
    }
    return data;
  }

  const res = await api<ApiEnvelope<SoalUjianResponse>>(
    `/bank-soal/${bankSoalId}/soal`,
    {
      method: "GET",
    }
  );
  return res.data;
}
