import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useDelete, useFetch, usePost } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { CreatePenjadwalanUjianPayload } from "@/types/Ujian/BuatUjian";
import type {
  JadwalUjianFilterParams,
  JadwalUjianItem,
  JadwalUjianStatusClient,
} from "@/types/Ujian/jadwalUjian";

const JADWAL_UJIAN_ENDPOINT = "/jadwal-ujian";
const UJIAN_DETAIL_ENDPOINT = "/ujian/detail";

type JadwalUjianApiItem = {
  id: number;
  id_ujian: number;
  id_guru: number;
  id_pengawas: number;
  nama_ujian: string;
  pengawas_ujian: string;
  tgl_ujian: string;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  sesi_ujian: number;
  ruang_ujian: string;
  id_ruang: number;
  status_ujian: string;
  started: number;
  tingkat_kelas: number;
  tingkat_kelas_id: number;
  nama_kelas: string;
  pembuat_username: string;
  pengawas_username: string;
};

const toClientStatus = (status: string): JadwalUjianStatusClient => {
  const normalized = status.trim().toUpperCase();
  switch (normalized) {
    case "BELUM_DIMULAI":
    case "BELUM_MULAI":
      return "belum_dimulai";
    case "BERLANGSUNG":
    case "MULAI":
      return "berlangsung";
    case "SELESAI":
      return "selesai";
    case "DIBATALKAN":
      return "dibatalkan";
    default:
      return "belum_dimulai";
  }
};

const mapJadwalUjianItem = (item: JadwalUjianApiItem): JadwalUjianItem => ({
  id: item.id,
  id_ujian: item.id_ujian,
  id_guru: item.id_guru,
  id_pengawas: item.id_pengawas,
  nama_ujian: item.nama_ujian,
  pengawas_ujian: item.pengawas_ujian,
  tgl_ujian: item.tgl_ujian,
  tanggal_ujian: item.tanggal_ujian,
  waktu_mulai: item.waktu_mulai,
  waktu_selesai: item.waktu_selesai,
  sesi_ujian: item.sesi_ujian,
  ruang_ujian: item.ruang_ujian,
  id_ruang: item.id_ruang,
  status_ujian: toClientStatus(item.status_ujian),
  started: item.started === 1 ? 1 : 0,
  pembuat_username: item.pembuat_username,
  pengawas_username: item.pengawas_username,
  tingkat_kelas: item.tingkat_kelas,
  tingkat_kelas_id: item.tingkat_kelas_id,
  nama_kelas: item.nama_kelas,
});

export async function getJadwalUjian(
  params: JadwalUjianFilterParams = {},
): Promise<JadwalUjianItem[]> {
  const queryParams: Record<string, string | number | undefined> = {
    search: params.search?.trim() || undefined,
    tanggal: params.tanggal || undefined,
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    ruang_ujian_id: params.ruangUjianId ?? undefined,
    tahun: params.tahun ?? undefined,
  };

  const response = await api<JadwalUjianApiItem[]>(JADWAL_UJIAN_ENDPOINT, {
    method: "GET",
    params: queryParams,
  });

  return response.map(mapJadwalUjianItem);
}

export async function createJadwalUjian(payload: CreatePenjadwalanUjianPayload) {
  const data = buildJsonData(payload, { nullishToEmptyString: false });
  return api<boolean>("/ujian", {
    method: "POST",
    data,
  });
}

export async function deleteJadwalUjian(idUjian: number) {
  return api<boolean>(`${UJIAN_DETAIL_ENDPOINT}/${idUjian}`, {
    method: "DELETE",
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetJadwalUjian(params: JadwalUjianFilterParams = {}) {
  return useFetch(
    () => getJadwalUjian(params),
    [
      params.search,
      params.tanggal,
      params.tingkatKelasId,
      params.ruangUjianId,
      params.tahun,
    ],
  );
}

export function useCreateJadwalUjian() {
  return usePost((payload: CreatePenjadwalanUjianPayload) =>
    createJadwalUjian(payload),
  );
}

export function useDeleteJadwalUjian() {
  return useDelete((idUjian: number) => deleteJadwalUjian(idUjian));
}
