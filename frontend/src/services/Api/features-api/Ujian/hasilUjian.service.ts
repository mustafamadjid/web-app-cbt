import { api, type ApiEnvelope } from "@/services/Api/api";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type {
  HasilUjianSiswa,
  StatistikUjian,
} from "@/types/Ujian/HasilUjian";

export type HasilUjianDetailResponse = {
  statistik: StatistikUjian;
  siswa: HasilUjianSiswa[];
};

export async function getHasilUjianList(): Promise<JadwalUjianItem[]> {
  const res = await api<ApiEnvelope<JadwalUjianItem[]>>("/ujian/hasil", {
    method: "GET",
  });
  return res.data;
}

export async function getHasilUjianDetail(
  ujianId: number
): Promise<HasilUjianDetailResponse> {
  const res = await api<ApiEnvelope<HasilUjianDetailResponse>>(
    `/ujian/hasil/${ujianId}`,
    {
      method: "GET",
    }
  );
  return res.data;
}
