import type { SesiFilterParams, SesiFormValues, SesiRow } from "@/types/DataMaster/Sesi";
import { api } from "@/services/Api/api";
import type { ApiEnvelope } from "@/services/Api/api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

const DUMMY_SESI: SesiRow[] = [
  { id: 1, kodeSesi: "SESI-01", namaSesi: "Sesi Pagi" },
  { id: 2, kodeSesi: "SESI-02", namaSesi: "Sesi Siang" },
  { id: 3, kodeSesi: "SESI-03", namaSesi: "Sesi Sore" },
  { id: 4, kodeSesi: "SESI-04", namaSesi: "Sesi Malam" },
];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getSesi(
  params: SesiFilterParams = {}
): Promise<SesiRow[]> {
  await sleep(200);
  const q = params.q ? normalize(params.q) : "";

  if (!q) return DUMMY_SESI;

  return DUMMY_SESI.filter((sesi) => {
    return (
      sesi.kodeSesi.toLowerCase().includes(q) ||
      sesi.namaSesi.toLowerCase().includes(q)
    );
  });
}

export async function getSesiById(id: number): Promise<SesiFormValues | null> {
  const target = DUMMY_SESI.find((sesi) => sesi.id === id);
  if (!target) return null;

  return {
    kode_sesi: target.kodeSesi,
    nama_sesi: target.namaSesi,
  };
}

export async function updateSesi(id: number, values: SesiFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<{ id: number }>>(`/sesi/${id}`, {
    method: "PUT",
    data,
  });

  return res.data;
}
