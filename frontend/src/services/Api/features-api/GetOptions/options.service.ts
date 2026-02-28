import type { BankSoalItem } from "@/types/BankSoal/BankSoal";

import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import type {
  BankSoalOption,
  GuruPengawasOption,
  SesiUjianOption,
} from "@/types/Ujian/BuatUjian";
import { getBankSoalByKelas } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";
import { getMapel } from "@/services/Api/features-api/DataMaster/mapel.service";
import { getRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { useFetch } from "@/hooks/fetch";

const DUMMY_BANKSOAL_MAPEL: MataPelajaranOption[] = [
  { id: 101, label: "Bahasa Indonesia (Kelas 10)" },
  { id: 102, label: "Fisika (Kelas 10)" },
  { id: 103, label: "Matematika (Kelas 11)" },
  { id: 104, label: "Bahasa Indonesia (Kelas 11)" },
  { id: 105, label: "Ekonomi (Kelas 12)" },
];

const DUMMY_GURU_PENGAWAS: GuruPengawasOption[] = [
  {
    id: 1,
    nama: "Budi Santoso",
    nip: "198305232010011001",
    mapel: "Matematika",
  },
  {
    id: 2,
    nama: "Siti Rahmawati",
    nip: "198912112012012002",
    mapel: "Bahasa Indonesia",
  },
  { id: 3, nama: "Dedi Pratama", nip: "198706142011021003", mapel: "Fisika" },
];

const DUMMY_SESI: SesiUjianOption[] = [
  { id: 1, kode: "S1", nama: "Sesi 1 (Pagi)" },
  { id: 2, kode: "S2", nama: "Sesi 2 (Siang)" },
  { id: 3, kode: "S3", nama: "Sesi 3 (Sore)" },
];

const DUMMY_JENIS_KELAMIN: Array<{
  value: JenisKelamin;
  label: string;
}> = [
  { value: "LAKI_LAKI", label: "Laki-laki" },
  { value: "PEREMPUAN", label: "Perempuan" },
];

const DUMMY_ANGKATAN: number[] = [2023, 2024, 2025];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const getTotalSoal = (item: BankSoalItem) => {
  if (typeof item.total_soal === "number") return item.total_soal;
  const pg = item.jumlah_soal_pg ?? 0;
  const essay = item.jumlah_soal_essay ?? 0;
  return pg + essay;
};

type MataPelajaranOptionsParams = {
  source?: "bankSoal" | "dataMaster";
  tingkatKelasId?: number;
};

export async function getMataPelajaranOptions(
  params: MataPelajaranOptionsParams = {},
): Promise<MataPelajaranOption[]> {
  const { source = "dataMaster", tingkatKelasId } = params;

  if (source === "bankSoal") {
    const tingkatKelas = getTingkatKelasById(tingkatKelasId);
    const kelasId = tingkatKelas ?? undefined;
    const filtered = !kelasId
      ? DUMMY_BANKSOAL_MAPEL
      : DUMMY_BANKSOAL_MAPEL.filter((mapel) => {
          if (kelasId === 10) return mapel.label.includes("(Kelas 10)");
          if (kelasId === 11) return mapel.label.includes("(Kelas 11)");
          if (kelasId === 12) return mapel.label.includes("(Kelas 12)");
          return true;
        });

    await sleep(250);
    return filtered;
  }

  const data = await getMapel();
  return data.map((mapel) => ({
    id: mapel.id,
    label: mapel.namaMapel,
  }));
}

export async function getRuangUjianOptions(): Promise<RuangUjianRow[]> {
  const data = await getRuangUjian();
  return [...data].sort((a, b) => a.id_ruangan - b.id_ruangan);
}

export async function getAngkatanOptions(): Promise<number[]> {
  await sleep(250);
  return DUMMY_ANGKATAN;
}

export async function getJenisKelaminOptions(): Promise<
  Array<{ value: JenisKelamin; label: string }>
> {
  await sleep(150);
  return DUMMY_JENIS_KELAMIN;
}

export async function getUjianGuruPengawasOptions(): Promise<
  GuruPengawasOption[]
> {
  await sleep(180);
  return [...DUMMY_GURU_PENGAWAS];
}

export async function getUjianSesiOptions(): Promise<SesiUjianOption[]> {
  await sleep(180);
  return [...DUMMY_SESI];
}

export async function getUjianBankSoalOptions(params: {
  tingkatKelasId?: number;
}): Promise<BankSoalOption[]> {
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  const kelasId = tingkatKelas ?? undefined;
  const data = await getBankSoalByKelas({
    tingkatKelasId: params.tingkatKelasId,
    kelasId,
  });

  return data.map((item) => ({
    id: item.id,
    nama: item.nama_banksoal ?? "-",
    mata_pelajaran: item.mata_pelajaran,
    materi: item.materi,
    kelas: item.kelas,
    jumlah_soal_pg: item.jumlah_soal_pg,
    jumlah_soal_essay: item.jumlah_soal_essay,
    total_soal: getTotalSoal(item),
  }));
}

// =====================
// Hook Wrappers
// =====================

export function useGetMataPelajaranOptions(params: MataPelajaranOptionsParams = {}) {
  return useFetch(
    () => getMataPelajaranOptions(params),
    [params.source, params.tingkatKelasId],
  );
}

export function useGetRuangUjianOptions() {
  return useFetch(() => getRuangUjianOptions(), []);
}

export function useGetAngkatanOptions() {
  return useFetch(() => getAngkatanOptions(), []);
}

export function useGetJenisKelaminOptions() {
  return useFetch(() => getJenisKelaminOptions(), []);
}

export function useGetUjianGuruPengawasOptions() {
  return useFetch(() => getUjianGuruPengawasOptions(), []);
}

export function useGetUjianSesiOptions() {
  return useFetch(() => getUjianSesiOptions(), []);
}

export function useGetUjianBankSoalOptions(params: { tingkatKelasId?: number }) {
  return useFetch(
    () => getUjianBankSoalOptions(params),
    [params.tingkatKelasId],
  );
}
