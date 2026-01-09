import { formatTanggalToIso } from "@/helper/dateFormatting/formatToIso";
import { normalize } from "@/helper/normalizeString/normalizeString";
import type {
  JadwalUjianFilterParams,
  JadwalUjianItem,
} from "@/types/Ujian/jadwalUjian";

import { api } from "@/services/Api/api"; // ✅ aktifkan saat BE siap
import type { ApiEnvelope } from "@/services/Api/api"; // kalau pakai envelope

/**
 * Dummy source data.
 * Catatan: untuk konsisten dengan filter by ID,
 * dummy data harus memiliki tingkat_kelas_id dan ruang_ujian_id.
 */
export const dummyJadwalUjian: JadwalUjianItem[] = [
  {
    id: 1,
    nama_ujian: "Ujian Matematika",
    pengawas_ujian: "Pak Budi Santoso",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "07:30",
    sesi_ujian: 1,
    
    ruang_ujian: "Ruang 101",
    id_ruang: 1,
    status_ujian: "belum_dimulai",
   
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 1",
  },
  {
    id: 2,
    nama_ujian: "Ujian Bahasa Indonesia",
    pengawas_ujian: "Bu Siti Nuraini",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "10:00",
    sesi_ujian: 2,
  
    ruang_ujian: "Ruang 102",
    id_ruang: 2,
    status_ujian: "berlangsung",
    
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPA 2",
  },
  {
    id: 3,
    nama_ujian: "Ujian IPA",
    pengawas_ujian: "Pak Andi Pratama",
    tgl_ujian: "Selasa, 13 Februari 2026",
    waktu_mulai: "08:00",
    sesi_ujian: 1,
   
    ruang_ujian: "Lab IPA",
    id_ruang: 3,
    status_ujian: "selesai",
    
    tingkat_kelas: 12,
    tingkat_kelas_id: 3,
    nama_kelas: "12 IPS 1",
  },
  {
    id: 4,
    nama_ujian: "Ujian Bahasa Inggris",
    pengawas_ujian: "Bu Rina Oktavia",
    tgl_ujian: "Rabu, 14 Februari 2026",
    waktu_mulai: "09:30",
    
    ruang_ujian: "Ruang 203",
    id_ruang: 2,
    status_ujian: "belum_dimulai",
    
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPS 2",
  },
  {
    id: 5,
    nama_ujian: "Ujian Sejarah",
    pengawas_ujian: "Pak Dedi Kurniawan",
    tgl_ujian: "Kamis, 15 Februari 2026",
    waktu_mulai: "13:00",
    sesi_ujian: 3,
    
    ruang_ujian: "Ruang 301",
    id_ruang: 1,
    status_ujian: "belum_dimulai",
    
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 3",
  },
];

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

/**
 * --- Indexing (dummy mode only) ---
 * Precompute agar search & filter tanggal tidak rebuild/parse terus.
 */
type JadwalUjianItemIndexed = JadwalUjianItem & {
  __searchText: string;
  __tglIso?: string;
};

const buildSearchableText = (ujian: JadwalUjianItem) =>
  normalize(
    [
      ujian.nama_ujian,
      ujian.pengawas_ujian ?? "",
      ujian.tgl_ujian,
      ujian.waktu_mulai,
      ujian.sesi_ujian != null ? String(ujian.sesi_ujian) : "",
      ujian.ruang_ujian ?? "",
      ujian.status_ujian ?? "",
      ujian.tingkat_kelas_id != null ? String(ujian.tingkat_kelas_id) : "",
      ujian.tingkat_kelas != null ? String(ujian.tingkat_kelas) : "",
      ujian.nama_kelas ?? "",
    ]
      .filter(Boolean)
      .join(" ")
  );

const indexDummy = (data: JadwalUjianItem[]): JadwalUjianItemIndexed[] =>
  data.map((item) => ({
    ...item,
    __searchText: buildSearchableText(item),
    __tglIso: formatTanggalToIso(item.tgl_ujian) ?? undefined,
  }));


/**
 * --- Filter helpers (dummy mode only) ---
 * Semua filter di sini pakai ID agar konsisten dengan BE.
 */

const filterByTingkatKelasId = (
  data: JadwalUjianItemIndexed[],
  tingkatKelasId?: number
) => {
  if (tingkatKelasId == null) return data;
  return data.filter((ujian) => ujian.tingkat_kelas_id === tingkatKelasId);
};

const filterByRuangUjianId = (
  data: JadwalUjianItemIndexed[],
  ruangUjianId?: number
) => {
  if (ruangUjianId == null) return data;
  return data.filter((ujian) => ujian.id_ruang === ruangUjianId);
};

const filterByTanggalIso = (
  data: JadwalUjianItemIndexed[],
  tanggal?: string
) => {
  if (!tanggal) return data;
  return data.filter((ujian) => ujian.__tglIso === tanggal);
};

const filterBySearch = (data: JadwalUjianItemIndexed[], q?: string) => {
  const query = q ? normalize(q) : "";
  if (!query) return data;
  return data.filter((ujian) => ujian.__searchText.includes(query));
};

const applyJadwalUjianFilters = (
  data: JadwalUjianItemIndexed[],
  params: JadwalUjianFilterParams = {}
) => {
  let out = data;
  out = filterByTingkatKelasId(out, params.tingkatKelasId);
  out = filterByRuangUjianId(out, params.ruangUjianId);
  out = filterByTanggalIso(out, params.tanggal);
  out = filterBySearch(out, params.q);
  return out;
};

const stripInternal = (items: JadwalUjianItemIndexed[]): JadwalUjianItem[] =>
  items.map(({ __searchText, __tglIso, ...rest }) => rest);

/**
 * Single public API for UI layer.
 * - Dummy mode: filter lokal
 * - API mode: query params ke BE
 *
 * Ganti USE_DUMMY sesuai strategi kamu.
 */
const USE_DUMMY = true; // ✅ set false saat BE sudah siap

export async function getJadwalUjian(
  params: JadwalUjianFilterParams = {}
): Promise<JadwalUjianItem[]> {
  if (USE_DUMMY) {
    await sleep(250);
    const filtered = applyJadwalUjianFilters(
      indexDummy(dummyJadwalUjian),
      params
    );
    return stripInternal(filtered);
  }

  // ✅ PRODUCTION MODE (aktifkan saat BE siap)
  const queryParams: Record<string, string | number | undefined> = {
    q: params.q || undefined,
    tanggal: params.tanggal || undefined,
    tingkat_kelas_id: params.tingkatKelasId ?? undefined,
    ruang_ujian_id: params.ruangUjianId ?? undefined,
  };

  const res = await api<ApiEnvelope<JadwalUjianItem[]>>("/jadwal-ujian", {
    params: queryParams,
  });
  return res.data;
}
