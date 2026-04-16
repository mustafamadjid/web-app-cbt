import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { HasilUjianSiswaApiItem } from "@/types/Api-Item/Ujian/HasilUjianSiswaApiItem";
import type {
  HasilUjianSiswaItem,
  HasilUjianSiswaResponse,
} from "@/types/Ujian/HasilUjianSiswa";

const HasilUjianSiswaEndpoint = "/siswa/ujian/list-selesai";

const mapHasilUjianSiswaItem = (
  item: HasilUjianSiswaApiItem,
): HasilUjianSiswaItem => ({
  ...item,
  started: item.started === 1 ? 1 : 0,
});

export async function getHasilUjianSiswaList(
  siswaId: number,
): Promise<HasilUjianSiswaResponse> {
  void siswaId;

  const response = await api<HasilUjianSiswaApiItem[]>(
    HasilUjianSiswaEndpoint,
    {
      method: "GET",
    },
  );

  return response.map(mapHasilUjianSiswaItem);
}

export function useGetHasilUjianSiswaList(
  siswaId: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled && siswaId > 0
        ? getHasilUjianSiswaList(siswaId)
        : Promise.resolve([] as HasilUjianSiswaResponse),
    [siswaId, enabled],
  );
}
