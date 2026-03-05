import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { toRfc3339Local } from "@/helper/dateFormatting/toRfc3339Local";
import { useDelete, useFetch, usePost, usePut } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type {
  CreatePenjadwalanUjianPayload,
  UpdatePenjadwalanUjianPayload,
} from "@/types/Ujian/BuatUjian";
import type { DetailUjianItem } from "@/types/Ujian/DetailUjian";
import type {
  JadwalUjianFilterParams,
  JadwalUjianItem,
  JadwalUjianStatusClient,
} from "@/types/Ujian/jadwalUjian";

const JADWAL_UJIAN_ENDPOINT = "/jadwal-ujian";
const UJIAN_ENDPOINT = "/ujian";

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

type JadwalUjianDetailApiResponse = {
  id_jadwal_ujian: number;
  id_ujian: number;
  id_sesi: number;
  id_ruangan: number;
  id_pengawas: number;
  token: string;
  tanggal_ujian: string;
  waktu_mulai: string;
  waktu_selesai: string;
  status_ujian: string;
  created_at: string;
  updated_at?: string;
};

type UjianDetailApiResponse = {
  id_ujian: number;
  id_bank_soal: number;
  id_kelas: number;
  id_nama_kelas?: number;
  id_guru: number;
  nama_ujian: string;
  deskripsi_ujian?: string;
  acak_soal: boolean;
  created_at: string;
  updated_at?: string;
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

const toServerStatus = (
  status: JadwalUjianStatusClient,
): UpdatePenjadwalanUjianPayload["status_ujian"] => {
  switch (status) {
    case "berlangsung":
      return "MULAI";
    case "selesai":
      return "SELESAI";
    case "dibatalkan":
      return "DIBATALKAN";
    default:
      return "BELUM_MULAI";
  }
};

const toClientTime = (value: string) => value.trim().slice(0, 5);

const formatTanggalIndonesia = (tanggalIso: string) => {
  if (!tanggalIso) return "-";
  const date = new Date(`${tanggalIso}T00:00:00`);
  return date.toLocaleDateString("id-ID", {
    weekday: "long",
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
};

const calculateDuration = (mulai: string, selesai: string) => {
  const [mulaiHour = 0, mulaiMinute = 0] = mulai.split(":").map(Number);
  const [selesaiHour = 0, selesaiMinute = 0] = selesai.split(":").map(Number);
  const mulaiMinuteTotal = mulaiHour * 60 + mulaiMinute;
  const selesaiMinuteTotal = selesaiHour * 60 + selesaiMinute;
  return Math.max(0, selesaiMinuteTotal - mulaiMinuteTotal);
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

export async function getJadwalUjianDetail(id: number): Promise<DetailUjianItem> {
  const jadwal = await api<JadwalUjianDetailApiResponse>(
    `${JADWAL_UJIAN_ENDPOINT}/${id}`,
    { method: "GET" },
  );

  const ujian = await api<UjianDetailApiResponse>(`${UJIAN_ENDPOINT}/${jadwal.id_ujian}`, {
    method: "GET",
  });

  const waktuMulai = toClientTime(jadwal.waktu_mulai);
  const waktuSelesai = toClientTime(jadwal.waktu_selesai);

  return {
    id: jadwal.id_jadwal_ujian,
    id_ujian: ujian.id_ujian,
    id_bank_soal: ujian.id_bank_soal,
    id_kelas: ujian.id_kelas,
    id_nama_kelas: ujian.id_nama_kelas,
    id_guru: ujian.id_guru,
    id_sesi: jadwal.id_sesi,
    id_ruangan: jadwal.id_ruangan,
    id_pengawas: jadwal.id_pengawas,
    nama_ujian: ujian.nama_ujian,
    deskripsi_ujian: ujian.deskripsi_ujian ?? "",
    acak_soal: ujian.acak_soal,
    tanggal_ujian: jadwal.tanggal_ujian,
    tgl_ujian: formatTanggalIndonesia(jadwal.tanggal_ujian),
    waktu_mulai: waktuMulai,
    waktu_selesai: waktuSelesai,
    durasi_menit: calculateDuration(waktuMulai, waktuSelesai),
    status_ujian: toClientStatus(jadwal.status_ujian),
    token: jadwal.token,
    tingkat_kelas: 0,
    nama_kelas: "-",
    pengawas_ujian: "-",
    ruang_ujian: "-",
    sesi_ujian: jadwal.id_sesi,
  };
}

export async function createJadwalUjian(payload: CreatePenjadwalanUjianPayload) {
  const data = buildJsonData(payload, { nullishToEmptyString: false });
  return api<boolean>(UJIAN_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function updateJadwalUjian(
  idUjian: number,
  payload: UpdatePenjadwalanUjianPayload,
) {
  const data = buildJsonData(payload, { nullishToEmptyString: false });
  return api<boolean>(`${UJIAN_ENDPOINT}/${idUjian}`, {
    method: "PATCH",
    data,
  });
}

export async function deleteJadwalUjian(idUjian: number) {
  return api<boolean>(`${UJIAN_ENDPOINT}/${idUjian}`, {
    method: "DELETE",
  });
}

export async function updateStatusUjian(
  idUjian: number,
  nextStatus: JadwalUjianStatusClient,
) {
  return updateJadwalUjian(idUjian, {
    status_ujian: toServerStatus(nextStatus),
  });
}

export const buildStatusPayloadWithLocalTime = (
  tanggal: string,
  waktuMulai: string,
  waktuSelesai: string,
) => ({
  tanggal_ujian: toRfc3339Local(tanggal, "00:00"),
  waktu_mulai: toRfc3339Local(tanggal, waktuMulai),
  waktu_selesai: toRfc3339Local(tanggal, waktuSelesai),
});

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

export function useGetJadwalUjianDetail(id: number) {
  return useFetch(() => getJadwalUjianDetail(id), [id]);
}

export function useCreateJadwalUjian() {
  return usePost((payload: CreatePenjadwalanUjianPayload) =>
    createJadwalUjian(payload),
  );
}

export function useUpdateJadwalUjian() {
  return usePut(
    (payload: { idUjian: number; values: UpdatePenjadwalanUjianPayload }) =>
      updateJadwalUjian(payload.idUjian, payload.values),
  );
}

export function useDeleteJadwalUjian() {
  return useDelete((idUjian: number) => deleteJadwalUjian(idUjian));
}

export function useUpdateStatusUjian() {
  return usePut(
    (payload: { idUjian: number; status: JadwalUjianStatusClient }) =>
      updateStatusUjian(payload.idUjian, payload.status),
  );
}
