import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type { PengumumanGetResponse } from "@/types/Widget/Pengumuman";
import type {
  SiswaDashboardSummary,
  SiswaProfile,
  SiswaSemesterAverage,
} from "@/types/Widget/SiswaDashboard";
import { getJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";
import { getPengumumanActive } from "@/services/Api/features-api/pengumuman/pengumuman.service";
import { useFetch } from "@/hooks/fetch";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const dummyProfile: SiswaProfile = {
  id: 14,
  nama: "Nabila Putri",
  kelas: "XI IPS 2",
  tingkat_kelas_id: 2,
};

const dummyRataRata: SiswaSemesterAverage[] = [
  { semester: "Ganjil 2024/2025", rata_rata: 84.5, target: 85 },
  { semester: "Genap 2024/2025", rata_rata: 88.2, target: 87 },
];

export async function getSiswaDashboardSummary(): Promise<SiswaDashboardSummary> {
  await sleep(200);
  const jadwal = await getJadwalUjian({
    tingkatKelasId: dummyProfile.tingkat_kelas_id,
  });
  const ujianSelesai = jadwal.filter(
    (item) => item.status_ujian === "selesai"
  ).length;

  return {
    ujian_selesai: ujianSelesai,
    total_ujian: jadwal.length,
    rata_rata_semester: dummyRataRata,
  };
}

export async function getSiswaJadwalUjian(): Promise<JadwalUjianItem[]> {
  await sleep(150);
  const jadwal = await getJadwalUjian({
    tingkatKelasId: dummyProfile.tingkat_kelas_id,
  });
  return jadwal.filter((item) => item.status_ujian !== "selesai");
}

export async function getSiswaPengumuman(): Promise<PengumumanGetResponse[]> {
  await sleep(120);
  return getPengumumanActive();
}

export async function getSiswaRataRataSemester(): Promise<
  SiswaSemesterAverage[]
> {
  await sleep(80);
  return dummyRataRata;
}

// =====================
// Hook Wrappers
// =====================

export function useGetSiswaDashboardSummary() {
  return useFetch(() => getSiswaDashboardSummary(), []);
}

export function useGetSiswaJadwalUjian() {
  return useFetch(() => getSiswaJadwalUjian(), []);
}

export function useGetSiswaPengumuman() {
  return useFetch(() => getSiswaPengumuman(), []);
}

export function useGetSiswaRataRataSemester() {
  return useFetch(() => getSiswaRataRataSemester(), []);
}
