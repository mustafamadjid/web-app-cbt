import { api, type ApiEnvelope } from "../../api";
import { useFetch, usePost } from "@/hooks/fetch";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";

import type {
  BuatUjianFormValues,
  BuatUjianSubmitResponse,
} from "@/types/Ujian/BuatUjian";
import type { DetailUjianItem } from "@/types/Ujian/DetailUjian";
import type {
  ActiveAttemptUjian,
  JawabanUjianSiswaResponse,
  SaveJawabanUjianSiswaRequest,
  WaktuSelesaiUjian,
} from "@/types/Ujian/ujianSiswa";



const UJIAN_DETAIL_ENDPOINT = "/ujian/detail";
const ACTIVE_ATTEMPT_UJIAN_ENDPOINT = "/siswa/ujian/attempt/active";
const SISWA_ATTEMPT_UJIAN_ENDPOINT = "/siswa/ujian/attempt";
const SISWA_SUBMIT_UJIAN_ENDPOINT = "/siswa/uijan/submit";
const SISWA_JAWABAN_UJIAN_ENDPOINT = "/siswa/ujian/jawaban";
const API_BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

const buildExpireAttemptPayload = () =>
  buildJsonData(
    {
      status_attempt: "expired",
    },
    { nullishToEmptyString: false },
  );


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

export async function getActiveAttemptUjian(
  idJadwalUjian: number,
): Promise<ActiveAttemptUjian> {
  return api<ActiveAttemptUjian>(ACTIVE_ATTEMPT_UJIAN_ENDPOINT, {
    method: "GET",
    params: {
      id_jadwal_ujian: idJadwalUjian,
    },
  });
}

export async function getJawabanUjianSiswaByAttemptId(
  idAttempt: number,
): Promise<JawabanUjianSiswaResponse> {
  return api<JawabanUjianSiswaResponse>(
    `${SISWA_JAWABAN_UJIAN_ENDPOINT}/${idAttempt}`,
    {
      method: "GET",
    },
  );
}

export async function saveJawabanUjianSiswa(
  payload: SaveJawabanUjianSiswaRequest,
): Promise<boolean> {
  return api<boolean>(SISWA_JAWABAN_UJIAN_ENDPOINT, {
    method: "POST",
    data: payload,
  });
}

export async function expireAttemptUjianSiswa(idAttempt: number): Promise<boolean> {
  return api<boolean>(`${SISWA_ATTEMPT_UJIAN_ENDPOINT}/${idAttempt}`, {
    method: "PATCH",
    data: buildExpireAttemptPayload(),
  });
}

export async function submitAttemptUjianSiswa(idAttempt: number): Promise<boolean> {
  return api<boolean>(`${SISWA_SUBMIT_UJIAN_ENDPOINT}/${idAttempt}`, {
    method: "PATCH",
  });
}

export function expireAttemptUjianSiswaOnPageLeave(idAttempt: number): void {
  if (!Number.isInteger(idAttempt) || idAttempt <= 0) {
    return;
  }

  void fetch(`${API_BASE_URL}${SISWA_ATTEMPT_UJIAN_ENDPOINT}/${idAttempt}`, {
    method: "PATCH",
    credentials: "include",
    keepalive: true,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(buildExpireAttemptPayload()),
  });
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

export function useGetActiveAttemptUjian(
  idJadwalUjian: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getActiveAttemptUjian(idJadwalUjian)
        : Promise.resolve(null as ActiveAttemptUjian | null),
    [idJadwalUjian, enabled],
  );
}

export function useGetJawabanUjianSiswaByAttemptId(
  idAttempt: number,
  enabled = true,
  deps: readonly unknown[] = [],
) {
  return useFetch(
    () =>
      enabled
        ? getJawabanUjianSiswaByAttemptId(idAttempt)
        : Promise.resolve(null as JawabanUjianSiswaResponse | null),
    [idAttempt, enabled, ...deps],
  );
}

export function useSaveJawabanUjianSiswa() {
  return usePost((payload: SaveJawabanUjianSiswaRequest) =>
    saveJawabanUjianSiswa(payload),
  );
}

export function useExpireAttemptUjianSiswa() {
  return usePost((idAttempt: number) => expireAttemptUjianSiswa(idAttempt));
}

export function useSubmitAttemptUjianSiswa() {
  return usePost((idAttempt: number) => submitAttemptUjianSiswa(idAttempt));
}
