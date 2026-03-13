import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { SoalUjian, SoalUjianSiswa } from "@/types/Ujian/SoalUjian";

export async function getSoalUjian(
  bankSoalId: number,
  acakSoal = false,
): Promise<SoalUjian[]> {
  return api<SoalUjian[]>(`/ujian/soal/bank-soal/${bankSoalId}`, {
    method: "GET",
    params: {
      acak_soal: String(acakSoal),
    },
  });
}

export async function GetSoalUjianForSiswa(
  idJadwalUjian: number,
): Promise<SoalUjianSiswa[]> {
  return api<SoalUjianSiswa[]>(`/siswa/soal-ujian/${idJadwalUjian}`, {
    method: "GET",
  });
}

export function useGetSoalUjian(
  bankSoalId: number,
  acakSoal = false,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getSoalUjian(bankSoalId, acakSoal)
        : Promise.resolve([] as SoalUjian[]),
    [bankSoalId, acakSoal, enabled],
  );
}

export function useGetSoalUjianForSiswa(
  idJadwalUjian: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? GetSoalUjianForSiswa(idJadwalUjian)
        : Promise.resolve([] as SoalUjianSiswa[]),
    [idJadwalUjian, enabled],
  );
}
