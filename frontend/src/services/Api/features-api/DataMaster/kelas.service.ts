import { api } from "../../api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type {
  KelasFilterParams,
  KelasRow,
  TingkatKelas,
} from "@/types/DataMaster/Kelas";
import type {
  KelasSubmitResponse,
  KelasFormValues,
} from "@/types/DataMaster/Kelas";
import type { ApiEnvelope } from "../../api";

const DUMMY_KELAS: KelasRow[] = [
  {
    id: 1,
    id_tingkat_kelas: 1,
    tingkat_kelas: 10,
    nama_kelas: "X IPA 1",
  },
  {
    id: 2,
    id_tingkat_kelas: 1,
    tingkat_kelas: 10,
    nama_kelas: "X IPS 1",
  },
  {
    id: 3,
    id_tingkat_kelas: 2,
    tingkat_kelas: 11,
    nama_kelas: "XI IPA 1",
  },
  {
    id: 4,
    id_tingkat_kelas: 2,
    tingkat_kelas: 11,
    nama_kelas: "XI IPS 1",
  },
  {
    id: 5,
    id_tingkat_kelas: 3,
    tingkat_kelas: 12,
    nama_kelas: "XII IPA 1",
  },
  {
    id: 6,
    id_tingkat_kelas: 3,
    tingkat_kelas: 12,
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

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getKelas(
  params: KelasFilterParams = {},
): Promise<KelasRow[]> {
  await sleep(250);
  const q = params.q ? normalize(params.q) : "";

  const filtered = DUMMY_KELAS.filter((kelas) => {
    if (params.tingkatKelas && kelas.id_tingkat_kelas !== params.tingkatKelas) {
      return false;
    }

    if (!q) return true;

    return (
      kelas.nama_kelas.toLowerCase().includes(q) ||
      String(kelas.id_tingkat_kelas).includes(q) ||
      String(kelas.tingkat_kelas).includes(q)
    );
  });

  return [...filtered].sort((a, b) => {
    if (a.tingkat_kelas !== b.tingkat_kelas) {
      return a.tingkat_kelas - b.tingkat_kelas;
    }
    return a.nama_kelas.localeCompare(b.nama_kelas, "id", {
      sensitivity: "base",
      numeric: true,
    });
  });
}

export const getTingkatKelasById = (id?: number | null): number | undefined => {
  if (id == null) return undefined;
  const dataId = DUMMY_KELAS.reduce<Record<number, TingkatKelas>>(
    (acc, kelas) => {
      if (!acc[kelas.id_tingkat_kelas]) {
        acc[kelas.id_tingkat_kelas] = {
          id_tingkat_kelas: kelas.id_tingkat_kelas,
          tingkat_kelas: kelas.tingkat_kelas,
        };
      }
      return acc;
    },
    {},
  );
  return dataId[id]?.tingkat_kelas;
};
