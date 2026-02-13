import type {
  RuangUjianFilterParams,
  RuangUjianFormValues,
  RuangUjianRow,
} from "@/types/DataMaster/RuangUjian";
import { api } from "@/services/Api/api";
import type { ApiEnvelope } from "@/services/Api/api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

const DUMMY_RUANG_UJIAN: RuangUjianRow[] = [
  { id: 1, namaRuangan: "Ruang Ujian 01" },
  { id: 2, namaRuangan: "Ruang Ujian 02" },
  { id: 3, namaRuangan: "Lab Komputer" },
  { id: 4, namaRuangan: "Aula Utama" },
  {id: 5, namaRuangan: "Lab IPA"},
  {id: 6, namaRuangan: "Lab IPS"}
];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getRuangUjian(
  params: RuangUjianFilterParams = {}
): Promise<RuangUjianRow[]> {
  await sleep(200);
  const q = params.q ? normalize(params.q) : "";

  const filtered = !q
    ? DUMMY_RUANG_UJIAN
    : DUMMY_RUANG_UJIAN.filter((ruang) =>
        ruang.namaRuangan.toLowerCase().includes(q)
      );

  const offset = params.offset ?? 0;
  const limit = params.limit ?? filtered.length;

  return filtered.slice(offset, offset + limit);
}

export async function getRuangUjianById(
  id: number
): Promise<RuangUjianFormValues | null> {
  const target = DUMMY_RUANG_UJIAN.find((ruang) => ruang.id === id);
  if (!target) return null;

  return {
    nama_ruangan_ujian: target.namaRuangan,
  };
}

export async function updateRuangUjian(
  id: number,
  values: RuangUjianFormValues
) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<{ id: number }>>(`/ruang-ujian/${id}`, {
    method: "PUT",
    data,
  });

  return res.data;
}
