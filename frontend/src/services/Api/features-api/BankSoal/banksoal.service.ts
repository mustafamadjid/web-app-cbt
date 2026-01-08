import type { BankSoalItem } from "@/types/DataMaster/BankSoal";
import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";
import { getTingkatKelasById } from "@/helper/tingkatKelas/tingkatKelas";

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
    __kelasId: "kelas-11",
    __mapelId: "mp-bindo-11",
    nama_banksoal: "Bank Soal Ujian Bahasa",
    guru: "Budianto",
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
    __kelasId: "kelas-10",
    __mapelId: "mp-fisika-10",
    nama_banksoal: "Bank Soal Ujian Fisika",
    guru: "Andi",
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
    __kelasId: "kelas-10",
    __mapelId: "mp-bindo-10",
    nama_banksoal: "Bank Soal Bahasa Indonesia - Teks",
    guru: "Budi",
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
    __kelasId: "kelas-11",
    __mapelId: "mp-mtk-11",
    nama_banksoal: "Bank Soal Matematika - Fungsi",
    guru: "Budi",
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
    __kelasId: "kelas-12",
    __mapelId: "mp-ekonomi-12",
    nama_banksoal: "Bank Soal Ekonomi - Pasar",
    guru: "Andi",
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

export async function getMataPelajaranOptions(params: {
  tingkatKelasId?: number;
}): Promise<MataPelajaranOption[]> {
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  const kelasId = tingkatKelas ? `kelas-${tingkatKelas}` : undefined;

  // contoh: kalau mode per kelas, batasi mapel sesuai kelas (dummy sederhana pakai label)
  const filtered = !kelasId
    ? DUMMY_MAPEL
    : DUMMY_MAPEL.filter((m) => {
        if (kelasId === "kelas-10") return m.label.includes("(Kelas 10)");
        if (kelasId === "kelas-11") return m.label.includes("(Kelas 11)");
        if (kelasId === "kelas-12") return m.label.includes("(Kelas 12)");
        return true;
      });

  return new Promise((resolve) => setTimeout(() => resolve(filtered), 250));
}

export async function getBankSoalByKelas(params: {
  tingkatKelasId?: number;
  kelasId?: string;
  mapelId?: string;
  q?: string;
}): Promise<BankSoalItem[]> {
  const { tingkatKelasId, kelasId, mapelId, q } = params;
  const tingkatKelas = getTingkatKelasById(tingkatKelasId);
  const resolvedKelasId = kelasId ?? (tingkatKelas ? `kelas-${tingkatKelas}` : undefined);

  let result = [...DUMMY_BANKSOAL];

  if (resolvedKelasId) {
    result = result.filter((x) => x.__kelasId === resolvedKelasId);
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
