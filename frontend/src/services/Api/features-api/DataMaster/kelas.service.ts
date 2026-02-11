import { api } from "../../api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type {
  KelasFilterParams,
  NamaKelas,
  TingkatKelas,
  FullDataKelas,
  KelasByIdData
} from "@/types/DataMaster/Kelas";
import type {
  KelasSubmitResponse,
  KelasFormValues,
} from "@/types/DataMaster/Kelas";


import type { ApiEnvelope } from "../../api";

const USE_DUMMY = true; // ✅ set false saat BE sudah siap

const DUMMY_TINGKAT_KELAS: TingkatKelas[] = [
  { id_tingkat_kelas: 1, tingkat_kelas: 10 },
  { id_tingkat_kelas: 2, tingkat_kelas: 11 },
  { id_tingkat_kelas: 3, tingkat_kelas: 12 },
];

const DUMMY_NAMA_KELAS: NamaKelas[] = [
  {
    id_nama_kelas: 1,
    id_tingkat_kelas: 1,
    nama_kelas: "X IPA 1",
  },
  {
    id_nama_kelas: 2,
    id_tingkat_kelas: 1,
    nama_kelas: "X IPS 1",
  },
  {
    id_nama_kelas: 3,
    id_tingkat_kelas: 2,
    nama_kelas: "XI IPA 1",
  },
  {
    id_nama_kelas: 4,
    id_tingkat_kelas: 2,
    nama_kelas: "XI IPS 1",
  },
  {
    id_nama_kelas: 5,
    id_tingkat_kelas: 3,
    nama_kelas: "XII IPA 1",
  },
  {
    id_nama_kelas: 6,
    id_tingkat_kelas: 3,
    nama_kelas: "XII IPS 2",
  },
];

// Submit Handler
export async function submitKelasResponse(values: KelasFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<KelasSubmitResponse>>("/kelas", {
    method: "POST",
    data,
  });

  return res.data;
}

type CreateTingkatKelasPayload = {
  tingkat_kelas: number;
};

type CreateNamaKelasPayload = {
  id_tingkat_kelas: number;
  nama_kelas: string;
};

export async function createTingkatKelas(values: CreateTingkatKelasPayload) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<null>>("/admin/kelas/tingkat-kelas", {
    method: "POST",
    data,
  });

  return res;
}

export async function createNamaKelas(values: CreateNamaKelasPayload) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<null>>("/admin/kelas/nama-kelas", {
    method: "POST",
    data,
  });

  return res;
}

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

let cachedTingkatKelas: TingkatKelas[] | null = null;

export async function GetDataKelasFull(
  params: KelasFilterParams = {},
): Promise<FullDataKelas> {

  const queryParams : Record<string, string | undefined> = {
    search: params.search || undefined,
    tingkat_kelas: params.tingkatKelas != null ? String(params.tingkatKelas) : undefined
  };

  return api<FullDataKelas>("/admin/kelas",{
    method: "GET",
    params : queryParams
  })
}






// Pertimbangkn untuk dihapus nanti
export async function getNamaKelas(
  params: KelasFilterParams = {},
): Promise<NamaKelas[]> {
  if (!USE_DUMMY) {
    const res = await api<ApiEnvelope<NamaKelas[]>>("/kelas/nama", {
      method: "GET",
      params,
    });
    return res.data;
  }

  await sleep(250);
  const q = params.q ? normalize(params.q) : "";

  const filtered = DUMMY_NAMA_KELAS.filter((kelas) => {
    if (params.tingkatKelas && kelas.id_tingkat_kelas !== params.tingkatKelas) {
      return false;
    }

    if (!q) return true;

    const tingkatLabel = getTingkatKelasById(kelas.id_tingkat_kelas);

    return (
      kelas.nama_kelas.toLowerCase().includes(q) ||
      String(kelas.id_tingkat_kelas).includes(q) ||
      (tingkatLabel != null && String(tingkatLabel).includes(q))
    );
  });

  return [...filtered].sort((a, b) => {
    if (a.id_tingkat_kelas !== b.id_tingkat_kelas) {
      return a.id_tingkat_kelas - b.id_tingkat_kelas;
    }
    return a.nama_kelas.localeCompare(b.nama_kelas, "id", {
      sensitivity: "base",
      numeric: true,
    });
  });
}

export async function getKelasById(id: number): Promise<KelasByIdData | null> {
  if (!USE_DUMMY) {
    return api<KelasByIdData>(`/admin/kelas/${id}`, {
      method: "GET",
    });
  }

  await sleep(150);
  const namaKelas = DUMMY_NAMA_KELAS.find((kelas) => kelas.id_nama_kelas === id);
  if (!namaKelas) return null;

  const tingkatKelas =
    DUMMY_TINGKAT_KELAS.find(
      (tingkat) => tingkat.id_tingkat_kelas === namaKelas.id_tingkat_kelas,
    ) ?? null;

  if (!tingkatKelas) return null;

  return {
    item_tingkat_kelas: tingkatKelas,
    item_nama_kelas: namaKelas,
  };
}

export async function getTingkatKelasByIdRequest(
  id: number,
): Promise<TingkatKelas | null> {
  if (!USE_DUMMY) {
    const res = await api<ApiEnvelope<TingkatKelas>>(
      `/admin/kelas/tingkat-kelas/${id}`,
      {
        method: "GET",
      },
    );
    return res.data;
  }

  await sleep(150);
  return DUMMY_TINGKAT_KELAS.find((item) => item.id_tingkat_kelas === id) ?? null;
}

export async function getNamaKelasById(
  id: number,
): Promise<NamaKelas | null> {
  if (!USE_DUMMY) {
    const res = await api<ApiEnvelope<NamaKelas>>(`/admin/kelas/nama-kelas/${id}`, {
      method: "GET",
    });
    return res.data;
  }

  await sleep(150);
  return DUMMY_NAMA_KELAS.find((item) => item.id_nama_kelas === id) ?? null;
}

export async function updateKelas(id: number, values: KelasFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<KelasSubmitResponse>>(`/kelas/${id}`, {
    method: "PUT",
    data,
  });

  return res.data;
}

export async function getTingkatKelas(): Promise<TingkatKelas[]> {
  if (!USE_DUMMY) {
    const res = await api<ApiEnvelope<TingkatKelas[]>>("/kelas/tingkat", {
      method: "GET",
    });
    cachedTingkatKelas = res.data;
    return res.data;
  }

  await sleep(150);
  cachedTingkatKelas = [...DUMMY_TINGKAT_KELAS];
  return [...DUMMY_TINGKAT_KELAS].sort(
    (a, b) => a.id_tingkat_kelas - b.id_tingkat_kelas,
  );
}

export const getTingkatKelasById = (id?: number | null): number | undefined => {
  if (id == null) return undefined;
  const data = cachedTingkatKelas ?? (USE_DUMMY ? DUMMY_TINGKAT_KELAS : null);
  return data?.find((item) => item.id_tingkat_kelas === id)?.tingkat_kelas;
};
