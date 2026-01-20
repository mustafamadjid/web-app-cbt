import { api, type ApiEnvelope } from "@/services/Api/api";
import type { SoalUjianResponse } from "@/types/BankSoal/BankSoal";



const dummySoalUjian: SoalUjianResponse[] = [
  {
    id_bank_soal: 1,
    nama_ujian: "Soal Bahasa Indonesia ",
    mata_pelajaran: "Bahasa Indonesia",
    jumlah_soal: 3,
    sisa_waktu: "01:35:00",
    soal: [
      {
        id_soal: 101,
        nomor_urut_soal: 1,
        tipe_soal: "PILIHAN_GANDA",
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
      {
        id_soal: 102,
        nomor_urut_soal: 2,
        tipe_soal: "PILIHAN_GANDA",
        pertanyaan:
          "Perhatikan kalimat berikut: \"Para siswa diminta membuat ringkasan bacaan untuk menilai pemahaman isi teks.\" Apa tujuan kegiatan pada kalimat tersebut?",
        opsi_a: "Menilai pemahaman siswa terhadap isi bacaan",
        opsi_b: "Menghafal isi bacaan secara cepat",
        opsi_c: "Membiasakan siswa menyalin teks bacaan",
        opsi_d: "Menyusun daftar pustaka dari bacaan",
        jawaban: "A",
      },
      {
        id_soal: 103,
        nomor_urut_soal: 3,
        tipe_soal: "PILIHAN_GANDA",
        pertanyaan:
          "Kalimat yang menggunakan kata berimbuhan me- dengan tepat adalah ....",
        opsi_a: "Ibu memakaikan adik seragam baru.",
        opsi_b: "Ayah menulisakan surat untuk kerabat.",
        opsi_c: "Dina meminumkan obat agar cepat sembuh.",
        opsi_d: "Rani menyapukan debu di meja belajar.",
        jawaban: "A",
      },
    ],
  },
];

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

const USE_DUMMY = true;

export async function getSoalUjian(
  bankSoalId: number
): Promise<SoalUjianResponse> {
  if (USE_DUMMY) {
    await sleep(200);
    const data = dummySoalUjian.find((item)=>item.id_bank_soal === bankSoalId);
    console.log(data);
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
