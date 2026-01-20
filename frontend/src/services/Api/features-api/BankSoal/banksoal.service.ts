import type { BankSoalItem } from "@/types/BankSoal/BankSoal";
import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";
import { api, type ApiEnvelope } from "../../api";

const DUMMY_MAPEL: MataPelajaranOption[] = [
  { id: 101, label: "Bahasa Indonesia (Kelas 10)" },
  { id: 102, label: "Fisika (Kelas 10)" },
  { id: 103, label: "Matematika (Kelas 11)" },
  { id: 104, label: "Bahasa Indonesia (Kelas 11)" },
  { id: 105, label: "Ekonomi (Kelas 12)" },
];

// Catatan: BankSoalItem kamu saat ini punya `kelas: number` + `mata_pelajaran: string`.
// Untuk filter mapel-by-id yang rapi, idealnya ada `mapelId` di item.
// Sementara: kita simpan `__mapelId` sebagai field tambahan lokal (tidak wajib disimpan ke DB).
type BankSoalItemLocal = BankSoalItem & {
  __kelasId: number;
  __mapelId: number;
};

const DUMMY_BANKSOAL: BankSoalItemLocal[] = [
  {
    id: 1,
    id_guru: 101,
    __kelasId: 11,
    __mapelId: 104,
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
    id: 2,
    id_guru: 102,
    __kelasId: 10,
    __mapelId: 102,
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
    id: 3,
    id_guru: 103,
    __kelasId: 10,
    __mapelId: 101,
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
    id: 4,
    id_guru: 103,
    __kelasId: 11,
    __mapelId: 103,
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
    id: 5,
    id_guru: 102,
    __kelasId: 12,
    __mapelId: 105,
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

const USE_DUMMY = true; // ✅ set false saat BE sudah siap

// =====================
// MOCK SERVICE
// (nanti tinggal ganti implementasi ke API beneran)
// =====================
function normalize(s: string) {
  return s.toLowerCase().trim();
}

type BankSoalFilterParams = {
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
  idGuru?: number;
};

function filterBankSoal(params: BankSoalFilterParams) {
  const { tingkatKelasId, kelasId, mapelId, q, idGuru } = params;
  const tingkatKelas = getTingkatKelasById(tingkatKelasId);
  const resolvedKelasId = kelasId ?? tingkatKelas ?? undefined;

  let result = [...DUMMY_BANKSOAL];

  if (resolvedKelasId) {
    result = result.filter((x) => x.__kelasId === resolvedKelasId);
  }

  if (mapelId) {
    result = result.filter((x) => x.__mapelId === mapelId);
  }

  if (idGuru) {
    result = result.filter((x) => x.id_guru === idGuru);
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
          .join(" "),
      );
      return haystack.includes(nq);
    });
  }

  // sort terbaru dulu
  result.sort((a, b) => (b.tgl_buat ?? "").localeCompare(a.tgl_buat ?? ""));

  return result;
}

export async function getMataPelajaranOptions(params: {
  tingkatKelasId?: number;
}): Promise<MataPelajaranOption[]> {
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  const kelasId = tingkatKelas ?? undefined;

  // contoh: kalau mode per kelas, batasi mapel sesuai kelas (dummy sederhana pakai label)
  const filtered = !kelasId
    ? DUMMY_MAPEL
    : DUMMY_MAPEL.filter((m) => {
        if (kelasId === 10) return m.label.includes("(Kelas 10)");
        if (kelasId === 11) return m.label.includes("(Kelas 11)");
        if (kelasId === 12) return m.label.includes("(Kelas 12)");
        return true;
      });

  return new Promise((resolve) => setTimeout(() => resolve(filtered), 250));
}

export async function getBankSoalByKelas(params: {
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
}): Promise<BankSoalItem[]> {
  const result = filterBankSoal(params);

  // balikin tanpa field internal
  const stripped: BankSoalItem[] = result.map(
    ({ __kelasId, __mapelId, ...rest }) => rest,
  );

  return new Promise((resolve) => setTimeout(() => resolve(stripped), 350));
}

export async function getBankSoalByGuru(params: {
  idGuru?: number;
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
}): Promise<BankSoalItem[]> {
  if (!params.idGuru) {
    throw new Error("id guru wajib diisi.");
  }

  if (USE_DUMMY) {
    const result = filterBankSoal({
      ...params,
      idGuru: params.idGuru,
    });

    const stripped: BankSoalItem[] = result.map(
      ({ __kelasId, __mapelId, ...rest }) => rest,
    );

    return new Promise((resolve) => setTimeout(() => resolve(stripped), 350));
  }

  const queryParams: Record<string, string | number | undefined> = {
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    kelas_id: params.kelasId ?? undefined,
    mapel_id: params.mapelId ?? undefined,
    q: params.q ?? undefined,
  };

  const res = await api<ApiEnvelope<BankSoalItem[]>>(`/guru/${params.idGuru}/bank-soal`, {
    params: queryParams,
  });
  return res.data;
}
