import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type {
  HasilUjianSiswa,
  StatistikUjian,
} from "@/types/Ujian/HasilUjian";
import type {
  JadwalUjianItem,
  JadwalUjianStatusClient,
} from "@/types/Ujian/jadwalUjian";

export type HasilUjianDetailResponse = {
  statistik: StatistikUjian;
  siswa: HasilUjianSiswa[];
};

type HasilUjianListApiItem = Omit<JadwalUjianItem, "status_ujian" | "started"> & {
  status_ujian: string;
  started: number;
};

type UjianEssayUngradedApiItem = {
  id: number;
  id_ujian: number;
  id_bank_soal: number;
  id_guru: number;
  id_pengawas: number;
  nama_ujian: string;
  pengawas_ujian: string;
  tgl_ujian: string;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  id_sesi: number;
  nama_sesi: string;
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

const mapHasilUjianItem = (item: HasilUjianListApiItem): JadwalUjianItem => ({
  ...item,
  status_ujian: toClientStatus(item.status_ujian),
  started: item.started === 1 ? 1 : 0,
});

const mapUjianEssayUngradedItem = (
  item: UjianEssayUngradedApiItem,
): JadwalUjianItem => ({
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
  sesi_ujian: item.id_sesi,
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

export type HasilUjianFilterParams = {
  tingkatKelasId?: number;
  tahun?: string;
};

export type UjianEssayUngradedFilterParams = {
  search?: string;
  limit?: number;
  offset?: number;
  tanggal?: string;
  tahun?: string | number;
  bulan?: string | number;
  tingkatKelasId?: number;
  namaKelasId?: number;
  idMapel?: number;
  sesiId?: number;
};

export async function getHasilUjianList(
  params: HasilUjianFilterParams = {},
): Promise<JadwalUjianItem[]> {
  const queryParams: Record<string, string | number | undefined> = {
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    tahun: params.tahun ?? undefined,
  };

  const response = await api<HasilUjianListApiItem[]>("/ujian/hasil", {
    method: "GET",
    params: queryParams,
  });

  return response.map(mapHasilUjianItem);
}

export async function getUjianEssayUngradedList(
  params: UjianEssayUngradedFilterParams = {},
): Promise<JadwalUjianItem[]> {
  const queryParams: Record<string, string | number | undefined> = {
    search: params.search?.trim() || undefined,
    limit: params.limit ?? undefined,
    offset: params.offset ?? undefined,
    tanggal: params.tanggal || undefined,
    tahun: params.tahun ?? undefined,
    bulan: params.bulan ?? undefined,
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    nama_kelas_id: params.namaKelasId ?? undefined,
    id_mapel: params.idMapel ?? undefined,
    sesi_id: params.sesiId ?? undefined,
  };

  const response = await api<UjianEssayUngradedApiItem[]>(
    "/ujian/koreksi-essay/list",
    {
      method: "GET",
      params: queryParams,
    },
  );

  return response.map(mapUjianEssayUngradedItem);
}

export async function getHasilUjianDetail(
  ujianId: number,
): Promise<HasilUjianDetailResponse> {
  return api<HasilUjianDetailResponse>(`/ujian/hasil/${ujianId}`, {
    method: "GET",
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetHasilUjianList(
  params: HasilUjianFilterParams = {},
  enabled = true,
) {
  return useFetch(
    () => (enabled ? getHasilUjianList(params) : Promise.resolve([])),
    [params.tingkatKelasId, params.tahun, enabled],
  );
}

export function useGetUjianEssayUngradedList(
  params: UjianEssayUngradedFilterParams = {},
  enabled = true,
) {
  return useFetch(
    () => (enabled ? getUjianEssayUngradedList(params) : Promise.resolve([])),
    [
      params.search,
      params.limit,
      params.offset,
      params.tanggal,
      params.tahun,
      params.bulan,
      params.tingkatKelasId,
      params.namaKelasId,
      params.idMapel,
      params.sesiId,
      enabled,
    ],
  );
}

export function useGetHasilUjianDetail(ujianId: number, enabled = true) {
  return useFetch(
    () => (enabled ? getHasilUjianDetail(ujianId) : Promise.resolve(null)),
    [ujianId, enabled],
  );
}
