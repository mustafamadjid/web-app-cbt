import type {
  KelasOption,
  MataPelajaranFilterParams,
  MataPelajaranOption,
  MataPelajaranRow,
} from "@/types/DataMaster/MataPelajaran";

const DUMMY_KELAS: KelasOption[] = [
  {
    id: "kelas-10",
    tingkat_kelas: 10,
    label: "Kelas 10",
  },
  {
    id: "kelas-11",
    tingkat_kelas: 11,
    label: "Kelas 11",
  },
  {
    id: "kelas-12",
    tingkat_kelas: 12,
    label: "Kelas 12",
  },
];

const DUMMY_MAPEL: MataPelajaranRow[] = [
  {
    id: "mapel-1",
    kelasId: "kelas-10",
    kodeMapel: "MAT-10-01",
    namaMapel: "Matematika",
    deskripsiMapel: "Aljabar dasar, geometri, dan statistika.",
  },
  {
    id: "mapel-2",
    kelasId: "kelas-10",
    kodeMapel: "EKO-10-01",
    namaMapel: "Ekonomi",
    deskripsiMapel: "Dasar-dasar ekonomi mikro dan makro.",
  },
  {
    id: "mapel-3",
    kelasId: "kelas-11",
    kodeMapel: "BIO-11-01",
    namaMapel: "Biologi",
    deskripsiMapel: "Sistem makhluk hidup dan genetika.",
  },
  {
    id: "mapel-4",
    kelasId: "kelas-11",
    kodeMapel: "GEO-11-01",
    namaMapel: "Geografi",
    deskripsiMapel: "Peta, lingkungan, dan dinamika wilayah.",
  },
  {
    id: "mapel-5",
    kelasId: "kelas-12",
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

  const kelasById = DUMMY_KELAS.reduce<Record<string, KelasOption>>(
    (acc, kelas) => {
      acc[kelas.id] = kelas;
      return acc;
    },
    {}
  );

  return DUMMY_MAPEL.filter((mapel) => {
    const tingkatKelas = kelasById[mapel.kelasId]?.tingkat_kelas;

    if (params.tingkatKelas && tingkatKelas !== params.tingkatKelas) {
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
