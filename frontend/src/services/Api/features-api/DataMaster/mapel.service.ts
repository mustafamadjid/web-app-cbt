import type {
  KelasOption,
  MataPelajaranFilterParams,
  MataPelajaranOption,
  MataPelajaranRow,
} from "@/types/DataMaster/MataPelajaran";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";

const DUMMY_KELAS: KelasOption[] = [
  {
    id: 10,
    tingkat_kelas: 10,
    label: "Kelas 10",
  },
  {
    id: 11,
    tingkat_kelas: 11,
    label: "Kelas 11",
  },
  {
    id: 12,
    tingkat_kelas: 12,
    label: "Kelas 12",
  },
];

const DUMMY_MAPEL: MataPelajaranRow[] = [
  {
    id: 1,
    kelasId: 10,
    kodeMapel: "MAT-10-01",
    namaMapel: "Matematika",
    deskripsiMapel: "Aljabar dasar, geometri, dan statistika.",
  },
  {
    id: 2,
    kelasId: 10,
    kodeMapel: "EKO-10-01",
    namaMapel: "Ekonomi",
    deskripsiMapel: "Dasar-dasar ekonomi mikro dan makro.",
  },
  {
    id: 3,
    kelasId: 11,
    kodeMapel: "BIO-11-01",
    namaMapel: "Biologi",
    deskripsiMapel: "Sistem makhluk hidup dan genetika.",
  },
  {
    id: 4,
    kelasId: 11,
    kodeMapel: "GEO-11-01",
    namaMapel: "Geografi",
    deskripsiMapel: "Peta, lingkungan, dan dinamika wilayah.",
  },
  {
    id: 5,
    kelasId: 12,
    kodeMapel: "FIS-12-01",
    namaMapel: "Fisika",
    deskripsiMapel: "Listrik, gelombang, dan mekanika.",
  },
];

const DUMMY_MAPEL_OPTIONS: MataPelajaranOption[] = DUMMY_MAPEL.map((mapel) => ({
  id: mapel.id,
  label: mapel.namaMapel,
}));

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getMataPelajaranOptions(): Promise<MataPelajaranOption[]> {
  await sleep(150);
  return DUMMY_MAPEL_OPTIONS;
}

export async function getMataPelajaran(
  params: MataPelajaranFilterParams = {}
): Promise<MataPelajaranRow[]> {
  await sleep(250);
  const q = params.q ? normalize(params.q) : "";

  const kelasById = DUMMY_KELAS.reduce<Record<number, KelasOption>>(
    (acc, kelas) => {
      acc[kelas.id] = kelas;
      return acc;
    },
    {}
  );

  return DUMMY_MAPEL.filter((mapel) => {
    const tingkatKelas = kelasById[mapel.kelasId]?.tingkat_kelas;
    const tingkatKelasById = getTingkatKelasById(params.tingkatKelas);

    if (params.tingkatKelas && tingkatKelas !== tingkatKelasById) {
      return false;
    }

    if (params.mapelId && params.mapelId !== mapel.id) {
      return false;
    }

    if (!q) return true;

    return (
      mapel.kodeMapel.toLowerCase().includes(q) ||
      mapel.namaMapel.toLowerCase().includes(q) ||
      mapel.deskripsiMapel.toLowerCase().includes(q) ||
      (tingkatKelas !== undefined && String(tingkatKelas).includes(q))
    );
  });
}
