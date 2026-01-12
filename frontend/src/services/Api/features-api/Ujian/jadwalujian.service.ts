import { formatTanggalToIso } from "@/helper/dateFormatting/formatToIso";
import { normalize } from "@/helper/normalizeString/normalizeString";
import type {
  JadwalUjianFilterParams,
  JadwalUjianItem,
} from "@/types/Ujian/jadwalUjian";
import type { DetailUjianItem } from "@/types/Ujian/DetailUjian";

import { api } from "@/services/Api/api"; // ✅ aktifkan saat BE siap
import type { ApiEnvelope } from "@/services/Api/api"; // kalau pakai envelope

/**
 * Dummy source data.
 * Catatan: untuk konsisten dengan filter by ID,
 * dummy data harus memiliki tingkat_kelas_id dan ruang_ujian_id.
 */
const dummyJadwalUjianDetail: DetailUjianItem[] = [
  {
    id: 1,
    nama_ujian: "Ujian Matematika",
    deskripsi_ujian: "Evaluasi tengah semester untuk materi fungsi dan aljabar.",
    tipe_ujian: "PILIHAN_GANDA",
    kelas_id: 1,
    bank_soal_id: 4,
    jumlah_soal: 20,
    tanggal_ujian: "2026-02-12",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "07:30",
    waktu_selesai: "09:00",
    durasi_menit: 90,
    ruang_ujian_id: 1,
    ruang_ujian: "Ruang 101",
    id_ruang: 1,
    acak_soal: true,
    guru_pengawas_id: 1,
    pengawas_ujian: "Pak Budi Santoso",
    sesi_id: 1,
    sesi_ujian: 1,
    token_ujian: "MAT-UTS-2026",
    status_ujian: "belum_dimulai",
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 1",
  },
  {
    id: 2,
    nama_ujian: "Ujian Bahasa Indonesia",
    deskripsi_ujian: "Ujian pemahaman teks eksposisi dan struktur bahasa.",
    tipe_ujian: "ESSAY",
    kelas_id: 2,
    bank_soal_id: 1,
    jumlah_soal: 15,
    tanggal_ujian: "2026-02-12",
    tgl_ujian: "Senin, 12 Februari 2026",
    waktu_mulai: "10:00",
    waktu_selesai: "11:30",
    durasi_menit: 90,
    ruang_ujian_id: 2,
    ruang_ujian: "Ruang 102",
    id_ruang: 2,
    acak_soal: false,
    guru_pengawas_id: 2,
    pengawas_ujian: "Bu Siti Nuraini",
    sesi_id: 2,
    sesi_ujian: 2,
    token_ujian: "BIN-2026-02",
    status_ujian: "berlangsung",
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPA 2",
  },
  {
    id: 3,
    nama_ujian: "Ujian IPA",
    deskripsi_ujian: "Evaluasi akhir materi eksperimen dan praktikum.",
    tipe_ujian: "CAMPURAN",
    kelas_id: 3,
    bank_soal_id: 5,
    jumlah_soal: 16,
    tanggal_ujian: "2026-02-13",
    tgl_ujian: "Selasa, 13 Februari 2026",
    waktu_mulai: "08:00",
    waktu_selesai: "09:40",
    durasi_menit: 100,
    ruang_ujian_id: 5,
    ruang_ujian: "Lab IPA",
    id_ruang: 3,
    acak_soal: true,
    guru_pengawas_id: 3,
    pengawas_ujian: "Pak Andi Pratama",
    sesi_id: 1,
    sesi_ujian: 1,
    token_ujian: "IPA-EX-13",
    status_ujian: "selesai",
    tingkat_kelas: 12,
    tingkat_kelas_id: 3,
    nama_kelas: "12 IPS 1",

  },
  {
    id: 4,
    nama_ujian: "Ujian Bahasa Inggris",
    deskripsi_ujian: "Ujian akhir materi reading comprehension.",
    tipe_ujian: "PILIHAN_GANDA",
    kelas_id: 2,
    bank_soal_id: 1,
    jumlah_soal: 20,
    tanggal_ujian: "2026-02-14",
    tgl_ujian: "Rabu, 14 Februari 2026",
    waktu_mulai: "09:30",
    waktu_selesai: "11:00",
    durasi_menit: 90,
    ruang_ujian_id: 2,
    ruang_ujian: "Ruang 203",
    id_ruang: 2,
    acak_soal: true,
    guru_pengawas_id: 2,
    pengawas_ujian: "Bu Rina Oktavia",
    sesi_id: 2,
    token_ujian: "BING-2026",
    status_ujian: "belum_dimulai",
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPS 2",
  },
  {
    id: 5,
    nama_ujian: "Ujian Sejarah",
    deskripsi_ujian: "Ujian akhir bab nasionalisme dan pergerakan.",
    tipe_ujian: "ESSAY",
    kelas_id: 1,
    bank_soal_id: 3,
    jumlah_soal: 18,
    tanggal_ujian: "2026-02-15",
    tgl_ujian: "Kamis, 15 Februari 2026",
    waktu_mulai: "13:00",
    waktu_selesai: "14:30",
    durasi_menit: 90,
    ruang_ujian_id: 1,
    ruang_ujian: "Ruang 301",
    id_ruang: 1,
    acak_soal: false,
    guru_pengawas_id: 1,
    pengawas_ujian: "Pak Dedi Kurniawan",
    sesi_id: 3,
    sesi_ujian: 3,
    token_ujian: "SEJ-2026",
    status_ujian: "belum_dimulai",
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 3",
  },
];

export const dummyJadwalUjian: JadwalUjianItem[] = dummyJadwalUjianDetail.map(
  (item) => ({
    id: item.id,
    nama_ujian: item.nama_ujian,
    pengawas_ujian: item.pengawas_ujian,
    tgl_ujian: item.tgl_ujian,
    waktu_mulai: item.waktu_mulai,
    sesi_ujian: item.sesi_ujian,
    ruang_ujian: item.ruang_ujian,
    id_ruang: item.id_ruang,
    status_ujian: item.status_ujian,
    tingkat_kelas: item.tingkat_kelas,
    tingkat_kelas_id: item.tingkat_kelas_id,
    nama_kelas: item.nama_kelas,
  })
);

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

export async function getJadwalUjianDetail(
  id: number
): Promise<DetailUjianItem> {
  if (USE_DUMMY) {
    await sleep(200);
    const detail = dummyJadwalUjianDetail.find((item) => item.id === id);
    if (!detail) {
      throw new Error("Detail ujian tidak ditemukan.");
    }
    return detail;
  }

  const res = await api<ApiEnvelope<DetailUjianItem>>(`/jadwal-ujian/${id}`);
  return res.data;
}
