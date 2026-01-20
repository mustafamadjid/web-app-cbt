import { api, type ApiEnvelope } from "@/services/Api/api";
import { formatTanggalToIso } from "@/helper/dateFormatting/formatToIso";
import { normalize } from "@/helper/normalizeString/normalizeString";
import type {
  UjianSiswaExamItem,
  UjianSiswaFilterParams,
  UjianSiswaResponse,
  UjianSiswaResultItem,
} from "@/types/Ujian/ujianSiswa";

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const dummyUjianList: UjianSiswaExamItem[] = [
  {
    id: 201,
    nama_ujian: "Ujian Matematika",
    mapel: "Matematika",
    pengawas_ujian: "Pak Budi Santoso",
    tgl_ujian: "Senin, 12 Februari 2026",
    tanggal_ujian: "2026-02-12",
    waktu_mulai: "07:30",
    waktu_selesai: "09:00",
    sesi_ujian: 1,
    ruang_ujian: "Ruang 101",
    id_ruang: 1,
    status_ujian: "belum_dimulai",
    started: 0,
    pembuat_username: "guru_sri",
    pengawas_username: "guru_sri",
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 1",
  },
  {
    id: 202,
    nama_ujian: "Ujian Bahasa Indonesia",
    mapel: "Bahasa Indonesia",
    pengawas_ujian: "Bu Siti Nuraini",
    tgl_ujian: "Senin, 12 Februari 2026",
    tanggal_ujian: "2026-02-12",
    waktu_mulai: "10:00",
    waktu_selesai: "11:30",
    sesi_ujian: 2,
    ruang_ujian: "Ruang 102",
    id_ruang: 2,
    status_ujian: "berlangsung",
    started: 1,
    pembuat_username: "guru_sri",
    pengawas_username: "guru_sri",
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPA 2",
  },
  {
    id: 203,
    nama_ujian: "Ujian Bahasa Inggris",
    mapel: "Bahasa Inggris",
    pengawas_ujian: "Bu Rina Oktavia",
    tgl_ujian: "Rabu, 14 Februari 2026",
    tanggal_ujian: "2026-02-14",
    waktu_mulai: "09:30",
    waktu_selesai: "11:00",
    sesi_ujian: 2,
    ruang_ujian: "Ruang 203",
    id_ruang: 2,
    status_ujian: "belum_dimulai",
    started: 0,
    pembuat_username: "guru_sri",
    pengawas_username: "guru_rina",
    tingkat_kelas: 11,
    tingkat_kelas_id: 2,
    nama_kelas: "11 IPS 2",
  },
];

const dummyUjianCompleted: UjianSiswaResultItem[] = [
  {
    id: 301,
    nama_ujian: "Ujian IPA",
    mapel: "Ilmu Pengetahuan Alam",
    pengawas_ujian: "Pak Andi Pratama",
    tgl_ujian: "Selasa, 13 Februari 2026",
    tanggal_ujian: "2026-02-13",
    waktu_mulai: "08:00",
    waktu_selesai: "09:40",
    sesi_ujian: 1,
    ruang_ujian: "Lab IPA",
    id_ruang: 3,
    status_ujian: "selesai",
    started: 1,
    pembuat_username: "guru_raka",
    pengawas_username: "guru_raka",
    tingkat_kelas: 12,
    tingkat_kelas_id: 3,
    nama_kelas: "12 IPS 1",
    jumlah_benar: 16,
    jumlah_salah: 4,
    nilai: 84,
  },
  {
    id: 302,
    nama_ujian: "Ujian Sejarah",
    mapel: "Sejarah",
    pengawas_ujian: "Pak Dedi Kurniawan",
    tgl_ujian: "Kamis, 15 Februari 2026",
    tanggal_ujian: "2026-02-15",
    waktu_mulai: "13:00",
    waktu_selesai: "14:30",
    sesi_ujian: 3,
    ruang_ujian: "Ruang 301",
    id_ruang: 1,
    status_ujian: "selesai",
    started: 1,
    pembuat_username: "guru_raka",
    pengawas_username: "guru_dedi",
    tingkat_kelas: 10,
    tingkat_kelas_id: 1,
    nama_kelas: "10 IPA 3",
    jumlah_benar: 14,
    jumlah_salah: 6,
    nilai: 78,
  },
];

type IndexedItem = {
    mapel: string;
  __searchText: string;
  __month?: number;
  __year?: number;
};

type IndexedExamItem = UjianSiswaExamItem & IndexedItem;

type IndexedResultItem = UjianSiswaResultItem & IndexedItem;

const buildSearchableText = (item: UjianSiswaExamItem) =>
  normalize(
    [
      item.nama_ujian,
      item.mapel,
      item.pengawas_ujian,
      item.tgl_ujian,
      item.tanggal_ujian ?? "",
      item.ruang_ujian ?? "",
    ]
      .filter(Boolean)
      .join(" ")
  );

const getDateMeta = (item: UjianSiswaExamItem) => {
  const iso = item.tanggal_ujian ?? formatTanggalToIso(item.tgl_ujian ?? "");
  if (!iso) return { month: undefined, year: undefined };
  const [y, m] = iso.split("-").map(Number);
  return { month: m, year: y };
};

const indexExamItem = (item: UjianSiswaExamItem): IndexedExamItem => {
  const meta = getDateMeta(item);
  return {
    ...item,
    __searchText: buildSearchableText(item),
    __month: meta.month,
    __year: meta.year,
  };
};

const indexResultItem = (item: UjianSiswaResultItem): IndexedResultItem => {
  const meta = getDateMeta(item);
  return {
    ...item,
    __searchText: buildSearchableText(item),
    __month: meta.month,
    __year: meta.year,
  };
};

const applyFilters = <T extends IndexedItem>(
  data: T[],
  params: UjianSiswaFilterParams,
) => {
  const query = params.search ? normalize(params.search) : "";

  return data.filter((item) => {
    if (params.bulan && item.__month !== params.bulan) return false;
    if (params.tahun && item.__year !== params.tahun) return false;

    if (params.mapel) {
      if (normalize(item.mapel) !== normalize(params.mapel)) return false;
    }

    if (query && !item.__searchText.includes(query)) return false;
    return true;
  });
};


const stripInternalExam = (items: IndexedExamItem[]): UjianSiswaExamItem[] =>
  items.map(({ __searchText, __month, __year, ...rest }) => rest);

const stripInternalResult = (
  items: IndexedResultItem[]
): UjianSiswaResultItem[] =>
  items.map(({ __searchText, __month, __year, ...rest }) => rest);

const collectMapelOptions = (items: UjianSiswaExamItem[]) =>
  Array.from(new Set(items.map((item) => item.mapel))).sort((a, b) =>
    a.localeCompare(b)
  );

const USE_DUMMY = true;

export async function getUjianSiswaOverview(params: {
  siswaId: number;
  filter?: UjianSiswaFilterParams;
}): Promise<UjianSiswaResponse> {
  const filter = params.filter ?? {};

  if (USE_DUMMY) {
    await sleep(240);
    const indexedUpcoming = dummyUjianList.map(indexExamItem);
    const indexedCompleted = dummyUjianCompleted.map(indexResultItem);

    const filteredUpcoming = applyFilters(indexedUpcoming, filter);
    const filteredCompleted = applyFilters(indexedCompleted, filter);

    const upcoming = stripInternalExam(
      filteredUpcoming.filter((item) => item.status_ujian === "belum_dimulai")
    );
    const ongoing = stripInternalExam(
      filteredUpcoming.filter((item) => item.status_ujian === "berlangsung")
    );
    const completed = stripInternalResult(
      filteredCompleted.filter((item) => item.status_ujian === "selesai")
    );

    const mapelOptions = collectMapelOptions([
      ...dummyUjianList,
      ...dummyUjianCompleted,
    ]);

    return {
      upcoming,
      ongoing,
      completed,
      mapelOptions,
    };
  }


 const queryParams: Record<string, string | number | undefined> = {
   bulan: filter.bulan,
   tahun: filter.tahun,
   mapel: filter.mapel,
   search: filter.search?.trim() || undefined,
 };

 const res = await api<ApiEnvelope<UjianSiswaResponse>>(
   `siswa/${params.siswaId}/ujian`,
   { params: queryParams },
 );

 return res.data;
}

export async function getUjianSiswaResultDetail(
  id: number
): Promise<UjianSiswaResultItem> {
  if (USE_DUMMY) {
    await sleep(180);
    const detail = dummyUjianCompleted.find((item) => item.id === id);
    if (!detail) {
      throw new Error("Detail hasil ujian tidak ditemukan.");
    }
    return detail;
  }

  const res = await api<ApiEnvelope<UjianSiswaResultItem>>(
    `/ujian-siswa/hasil/${id}`
  );
  return res.data;
}
