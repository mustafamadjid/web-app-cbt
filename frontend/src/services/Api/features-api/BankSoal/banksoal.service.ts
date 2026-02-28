import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import type {
  BankSoalItem,
  CreateBankSoalPayload,
} from "@/types/BankSoal/BankSoal";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";
import { api } from "../../api";
import { useFetch, usePost } from "@/hooks/fetch";


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

const toBankSoalItem = (item: BankSoalItemLocal): BankSoalItem => ({
  id: item.id,
  id_guru: item.id_guru,
  guru: item.guru,
  nama_banksoal: item.nama_banksoal,
  mata_pelajaran: item.mata_pelajaran,
  materi: item.materi,
  kelas: item.kelas,
  deskripsi: item.deskripsi,
  tgl_buat: item.tgl_buat,
  jumlah_soal_pg: item.jumlah_soal_pg,
  jumlah_soal_essay: item.jumlah_soal_essay,
  total_soal: item.total_soal,
});

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

export async function getBankSoalByKelas(params: {
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
}): Promise<BankSoalItem[]> {
  const result = filterBankSoal(params);
  const stripped: BankSoalItem[] = result.map(toBankSoalItem);

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
    const stripped: BankSoalItem[] = result.map(toBankSoalItem);

    return new Promise((resolve) => setTimeout(() => resolve(stripped), 350));
  }

  const queryParams: Record<string, string | number | undefined> = {
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    kelas_id: params.kelasId ?? undefined,
    mapel_id: params.mapelId ?? undefined,
    q: params.q ?? undefined,
  };

  const res = await api<BankSoalItem[]>(`/guru/${params.idGuru}/bank-soal`, {
    params: queryParams,
  });
  return res;
}

export async function createBankSoal(
  values: CreateBankSoalPayload,
): Promise<BankSoalItem> {
  if (USE_DUMMY) {
    const nextId =
      DUMMY_BANKSOAL.reduce((maxId, item) => Math.max(maxId, item.id), 0) + 1;

    const item: BankSoalItemLocal = {
      id: nextId,
      id_guru: values.id_guru ?? 0,
      __kelasId: values.kelas,
      __mapelId: values.mapel_id,
      nama_banksoal: values.nama_banksoal,
      kelas: values.kelas,
      mata_pelajaran: values.mata_pelajaran,
      materi: values.materi ?? "-",
      deskripsi: values.deskripsi,
      tgl_buat: new Date().toISOString(),
      jumlah_soal_pg: 0,
      jumlah_soal_essay: 0,
    };

    DUMMY_BANKSOAL.unshift(item);
    return new Promise((resolve) =>
      setTimeout(() => resolve(toBankSoalItem(item)), 250),
    );
  }

  const data = buildJsonData(values);
  const res = await api<BankSoalItem>("/admin/bank-soal", {
    method: "POST",
    data,
  });

  return res;
}

// =====================
// Hook Wrappers
// =====================

export function useGetBankSoalByKelas(params: {
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
}) {
  return useFetch(
    () => getBankSoalByKelas(params),
    [params.tingkatKelasId, params.kelasId, params.mapelId, params.q],
  );
}

export function useGetBankSoalByGuru(params: {
  idGuru?: number;
  tingkatKelasId?: number;
  kelasId?: number;
  mapelId?: number;
  q?: string;
}) {
  return useFetch(
    () => getBankSoalByGuru(params),
    [params.idGuru, params.tingkatKelasId, params.kelasId, params.mapelId, params.q],
  );
}

export function useCreateBankSoal() {
  return usePost((values: CreateBankSoalPayload) => createBankSoal(values));
}
