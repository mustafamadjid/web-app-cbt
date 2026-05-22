import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { AnalisisSoalResponse } from "@/types/Ujian/AnalisisSoal";

const ANALISIS_SOAL_ENDPOINT = "/ujian/analisis-soal";

export async function getAnalisisSoal(
  idJadwalUjian: number,
): Promise<AnalisisSoalResponse> {
  return api<AnalisisSoalResponse>(
    `${ANALISIS_SOAL_ENDPOINT}/${idJadwalUjian}`,
    {
      method: "GET",
    },
  );
}

export function useGetAnalisisSoal(
  idJadwalUjian: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getAnalisisSoal(idJadwalUjian)
        : Promise.resolve({
            id_jadwal_ujian: idJadwalUjian,
            analisis_soal: [],
          } as AnalisisSoalResponse),
    [idJadwalUjian, enabled],
  );
}
