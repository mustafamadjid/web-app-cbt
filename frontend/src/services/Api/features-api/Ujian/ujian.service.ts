import { api, type ApiEnvelope } from "../../api";
import { useFetch, usePost } from "@/hooks/fetch";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type {
  BuatUjianFormValues,
  BuatUjianSubmitResponse,
} from "@/types/Ujian/BuatUjian";
import type { DetailUjianItem } from "@/types/Ujian/DetailUjian";
import type { WaktuSelesaiUjian } from "@/types/Ujian/ujianSiswa";



const UJIAN_DETAIL_ENDPOINT = "/ujian/detail";


export async function submitBuatUjian(values: BuatUjianFormValues) {
  const data = buildJsonData(values);
  const res = await api<ApiEnvelope<BuatUjianSubmitResponse>>("/ujian", {
    method: "POST",
    data,
  });

  return res.data;
}

export async function getDetailUjianById(idUjian: number): Promise<DetailUjianItem> {
  return api<DetailUjianItem>(`${UJIAN_DETAIL_ENDPOINT}/${idUjian}`, {
    method: "GET",
  });
}

export async function GetWaktuSelesaiUjian(idJadwalUjian: number): Promise<WaktuSelesaiUjian> {
  return api<WaktuSelesaiUjian>(
    `/siswa/ujian/waktu-selesai/${idJadwalUjian}`,
    {
      method : "GET"
    },
  );
}

// =====================
// Hook Wrappers
// =====================


export function useSubmitBuatUjian() {
  return usePost((values: BuatUjianFormValues) => submitBuatUjian(values));
}

export function useGetDetailUjianById(idUjian: number, enabled = true) {
  return useFetch(
    () =>
      enabled
        ? getDetailUjianById(idUjian)
        : Promise.resolve(null as DetailUjianItem | null),
    [idUjian, enabled],
  );
}

export function useGetWaktuSelesaiUjian(idJadwalUjian: number, enabled = true) {
  return useFetch(
    () =>
      enabled
        ? GetWaktuSelesaiUjian(idJadwalUjian)
        : Promise.resolve(null as WaktuSelesaiUjian | null),
    [idJadwalUjian, enabled],
  );
}
