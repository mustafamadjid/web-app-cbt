import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { StatistikUjian } from "@/types/Ujian/StatistikUjian";

const STATISTIK_UJIAN_ENDPOINT = "/ujian/statistik";

export async function getStatistikUjian(
  idJadwalUjian: number,
): Promise<StatistikUjian> {
  return api<StatistikUjian>(`${STATISTIK_UJIAN_ENDPOINT}/${idJadwalUjian}`, {
    method: "GET",
  });
}

export function useGetStatistikUjian(
  idJadwalUjian: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getStatistikUjian(idJadwalUjian)
        : Promise.resolve(null as StatistikUjian | null),
    [idJadwalUjian, enabled],
  );
}
