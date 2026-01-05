import type {
  KelasOption,
  MataPelajaranFilterParams,
  MataPelajaranOption,
  MataPelajaranRow,
} from "@/types/DataMaster/MataPelajaran";

const DUMMY_KELAS: KelasOption[] = [
  {
    id: "kelas-10-ipa-1",
    tingkat_kelas: 10,
    nama_kelas: "X IPA 1",
    label: "Kelas 10 - IPA 1",
  },
  {
    id: "kelas-10-ips-1",
    tingkat_kelas: 10,
    nama_kelas: "X IPS 1",
    label: "Kelas 10 - IPS 1",
  },
  {
    id: "kelas-11-ipa-1",
    tingkat_kelas: 11,
    nama_kelas: "XI IPA 1",
    label: "Kelas 11 - IPA 1",
  },
  {
    id: "kelas-11-ips-1",
    tingkat_kelas: 11,
    nama_kelas: "XI IPS 1",
    label: "Kelas 11 - IPS 1",
  },
  {
    id: "kelas-12-ipa-1",
    tingkat_kelas: 12,
    nama_kelas: "XII IPA 1",
    label: "Kelas 12 - IPA 1",
  },
  {
    id: "kelas-12-ips-2",
    tingkat_kelas: 12,
    nama_kelas: "XII IPS 2",
    label: "Kelas 12 - IPS 2",
  },
];

const DUMMY_MAPEL: MataPelajaranRow[] = [
  {
    id: "mapel-1",
    kelasId: "kelas-10-ipa-1",
    kodeMapel: "MAT-10-01",
    namaMapel: "Matematika",
    deskripsiMapel: "Aljabar dasar, geometri, dan statistika.",
  },
  {
    id: "mapel-2",
    kelasId: "kelas-10-ips-1",
    kodeMapel: "EKO-10-01",
    namaMapel: "Ekonomi",
    deskripsiMapel: "Dasar-dasar ekonomi mikro dan makro.",
  },
  {
    id: "mapel-3",
    kelasId: "kelas-11-ipa-1",
    kodeMapel: "BIO-11-01",
    namaMapel: "Biologi",
    deskripsiMapel: "Sistem makhluk hidup dan genetika.",
  },
  {
    id: "mapel-4",
    kelasId: "kelas-11-ips-1",
    kodeMapel: "GEO-11-01",
    namaMapel: "Geografi",
    deskripsiMapel: "Peta, lingkungan, dan dinamika wilayah.",
  },
  {
    id: "mapel-5",
    kelasId: "kelas-12-ipa-1",
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

export async function getKelasOptions(): Promise<KelasOption[]> {
  await sleep(150);
  return DUMMY_KELAS;
}

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
