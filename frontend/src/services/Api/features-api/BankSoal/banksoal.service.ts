import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import type { KelasOption,MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";

// =====================
// DUMMY DATA
// =====================
const DUMMY_KELAS: KelasOption[] = [
  { id: "kls-10", tingkat_kelas: 10, nama_kelas: "A" },
  { id: "kls-11", tingkat_kelas: 11, nama_kelas: "IPA 1" },
  { id: "kls-12", tingkat_kelas: 12, nama_kelas: "IPS 2" },
];

const DUMMY_MAPEL: MataPelajaranOption[] = [
  { id: "mp-bindo-10", label: "Bahasa Indonesia (Kelas 10)" },
  { id: "mp-fisika-10", label: "Fisika (Kelas 10)" },
  { id: "mp-mtk-11", label: "Matematika (Kelas 11)" },
  { id: "mp-bindo-11", label: "Bahasa Indonesia (Kelas 11)" },
  { id: "mp-ekonomi-12", label: "Ekonomi (Kelas 12)" },
];

// Catatan: BankSoalItem kamu saat ini punya `kelas: number` + `mata_pelajaran: string`.
// Untuk filter mapel-by-id yang rapi, idealnya ada `mapelId` di item.
// Sementara: kita simpan `__mapelId` sebagai field tambahan lokal (tidak wajib disimpan ke DB).
type BankSoalItemLocal = BankSoalItem & {
  __kelasId: string;
  __mapelId: string;
};

const DUMMY_BANKSOAL: BankSoalItemLocal[] = [
  {
    id: "bs-001",
    __kelasId: "kls-11",
    __mapelId: "mp-bindo-11",
    nama_banksoal: "Bank Soal Ujian Bahasa",
    kelas: 11,
    mata_pelajaran: "Bahasa Indonesia",
    materi: "Ejaan Bahasa",
    jumlah_soal_pg: 10,
    jumlah_soal_essay: 5,
    deskripsi: "Bank soal untuk latihan kelas XI. Fokus ejaan dan tata bahasa.",
    tgl_buat: "2025-03-10",
  },
  {
    id: "bs-002",
    __kelasId: "kls-10",
    __mapelId: "mp-fisika-10",
    nama_banksoal: "Bank Soal Ujian Fisika",
    kelas: 10,
    mata_pelajaran: "Fisika",
    materi: "Hukum Newton",
    jumlah_soal_pg: 10,
    jumlah_soal_essay: 5,
    deskripsi: "Latihan konsep gaya, massa, percepatan, dan penerapannya.",
    tgl_buat: "2025-03-12",
  },
  {
    id: "bs-003",
    __kelasId: "kls-10",
    __mapelId: "mp-bindo-10",
    nama_banksoal: "Bank Soal Bahasa Indonesia - Teks",
    kelas: 10,
    mata_pelajaran: "Bahasa Indonesia",
    materi: "Teks Eksposisi",
    jumlah_soal_pg: 15,
    jumlah_soal_essay: 3,
    deskripsi: "Pemahaman teks eksposisi, struktur, dan kaidah kebahasaan.",
    tgl_buat: "2025-04-02",
  },
  {
    id: "bs-004",
    __kelasId: "kls-11",
    __mapelId: "mp-mtk-11",
    nama_banksoal: "Bank Soal Matematika - Fungsi",
    kelas: 11,
    mata_pelajaran: "Matematika",
    materi: "Fungsi Kuadrat",
    jumlah_soal_pg: 20,
    jumlah_soal_essay: 0,
    deskripsi:
      "Latihan fungsi kuadrat, grafik, dan analisis nilai maksimum/minimum.",
    tgl_buat: "2025-05-11",
  },
  {
    id: "bs-005",
    __kelasId: "kls-12",
    __mapelId: "mp-ekonomi-12",
    nama_banksoal: "Bank Soal Ekonomi - Pasar",
    kelas: 12,
    mata_pelajaran: "Ekonomi",
    materi: "Permintaan & Penawaran",
    jumlah_soal_pg: 12,
    jumlah_soal_essay: 4,
    deskripsi: "Studi kasus pasar, kurva, elastisitas, dan keseimbangan.",
    tgl_buat: "2025-06-01",
  },
];

// =====================
// MOCK SERVICE
// (nanti tinggal ganti implementasi ke API beneran)
// =====================
function normalize(s: string) {
  return s.toLowerCase().trim();
}

export async function getKelasOptions(): Promise<KelasOption[]> {
  return new Promise((resolve) => setTimeout(() => resolve(DUMMY_KELAS), 250));
}

export async function getMataPelajaranOptions(params: {
  kelasId?: string;
}): Promise<MataPelajaranOption[]> {
  const { kelasId } = params;

  // contoh: kalau mode per kelas, batasi mapel sesuai kelas (dummy sederhana pakai label)
  const filtered = !kelasId
    ? DUMMY_MAPEL
    : DUMMY_MAPEL.filter((m) => {
        if (kelasId === "kls-10") return m.label.includes("(Kelas 10)");
        if (kelasId === "kls-11") return m.label.includes("(Kelas 11)");
        if (kelasId === "kls-12") return m.label.includes("(Kelas 12)");
        return true;
      });

  return new Promise((resolve) => setTimeout(() => resolve(filtered), 250));
}

export async function getBankSoalByKelas(params: {
  kelasId?: string;
  mapelId?: string;
  q?: string;
}): Promise<BankSoalItem[]> {
  const { kelasId, mapelId, q } = params;

  let result = [...DUMMY_BANKSOAL];

  if (kelasId) {
    result = result.filter((x) => x.__kelasId === kelasId);
  }

  if (mapelId) {
    result = result.filter((x) => x.__mapelId === mapelId);
  }

  if (q && normalize(q).length > 0) {
    const nq = normalize(q);
    result = result.filter((x) => {
      const haystack = normalize(
        [
          x.nama_banksoal,
          x.mata_pelajaran,
          x.materi,
          x.deskripsi,
          String(x.kelas),
          x.tgl_buat,
        ]
          .filter(Boolean)
          .join(" ")
      );
      return haystack.includes(nq);
    });
  }

  // sort terbaru dulu
  result.sort((a, b) => (b.tgl_buat ?? "").localeCompare(a.tgl_buat ?? ""));


  // balikin tanpa field internal
  const stripped: BankSoalItem[] = result.map(
    ({ __kelasId, __mapelId, ...rest }) => rest
  );

  return new Promise((resolve) => setTimeout(() => resolve(stripped), 350));
}
