import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useDelete, useFetch, usePost } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { JadwalUjianApiItem } from "@/types/Api-Item/Ujian/JadwalUjianApiItem";
import type { JadwalUjianSiswaApiItem } from "@/types/Api-Item/Ujian/JadwalUjianSiswaApiItem";
import type { CreatePenjadwalanUjianPayload } from "@/types/Ujian/BuatUjian";
import type {
  JadwalUjianFilterParams,
  JadwalUjianItem,
  JadwalUjianSiswaFilterParams,
  JadwalUjianSiswaItem,
  JadwalUjianStatusClient,
} from "@/types/Ujian/jadwalUjian";

const JADWAL_UJIAN_ENDPOINT = "/jadwal-ujian";
const SISWA_JADWAL_UJIAN_ENDPOINT = "/siswa/ujian/list";
const UJIAN_DETAIL_ENDPOINT = "/ujian/detail";

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

const mapJadwalUjianSiswaItem = (
  item: JadwalUjianSiswaApiItem,
): JadwalUjianSiswaItem => ({
  ...item,
  status_ujian: toClientStatus(item.status_ujian),
  pengawas_nama_lengkap:
    item.pengawas_nama_lengkap?.trim() || item.pengawas_ujian,
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
    kategori_ujian: params.kategoriUjian ?? undefined,
  };

  const response = await api<JadwalUjianApiItem[]>(JADWAL_UJIAN_ENDPOINT, {
    method: "GET",
    params: queryParams,
  });

  return response.map(mapJadwalUjianItem);
}

export async function GetAllJadwalUjianForSiswa(
  params: JadwalUjianSiswaFilterParams = {},
): Promise<JadwalUjianSiswaItem[]> {
  const queryParams: Record<string, string | number | undefined> = {
    search: params.search?.trim() || undefined,
    limit: params.limit ?? undefined,
    offset: params.offset ?? undefined,
    tanggal: params.tanggal || undefined,
    tahun: params.tahun ?? undefined,
    bulan: params.bulan ?? undefined,
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    tingkat_kelas: params.tingkatKelas ?? undefined,
    ruang_ujian_id: params.ruangUjianId ?? undefined,
    id_mapel: params.idMapel ?? undefined,
    kategori_ujian: params.kategoriUjian ?? undefined,
  };

  const response = await api<JadwalUjianSiswaApiItem[]>(
    SISWA_JADWAL_UJIAN_ENDPOINT,
    {
      method: "GET",
      params: queryParams,
    },
  );

  return response.map(mapJadwalUjianSiswaItem);
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
      params.kategoriUjian,
    ],
  );
}

export function useGetAllJadwalUjianForSiswa(
  params: JadwalUjianSiswaFilterParams = {},
) {
  return useFetch(
    () => GetAllJadwalUjianForSiswa(params),
    [
      params.search,
      params.limit,
      params.offset,
      params.tanggal,
      params.tahun,
      params.bulan,
      params.tingkatKelasId,
      params.tingkatKelas,
      params.ruangUjianId,
      params.idMapel,
      params.kategoriUjian,
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
