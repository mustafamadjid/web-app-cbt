import type { MataPelajaranOption } from "@/types/DataMaster/MataPelajaran";
import type { RuangUjianRow } from "@/types/DataMaster/RuangUjian";
import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import type {
  BankSoalOption,
  GuruPengawasOption,
  SesiUjianOption,
} from "@/types/Ujian/BuatUjian";
import { getBankSoal } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { getMapel } from "@/services/Api/features-api/DataMaster/mapel.service";
import { getRuangUjian } from "@/services/Api/features-api/DataMaster/ruang-ujian.service";
import { useFetch } from "@/hooks/fetch";

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

type MataPelajaranOptionsParams = {
  source?: "bankSoal" | "dataMaster";
  tingkatKelasId?: number;
};

export async function getMataPelajaranOptions(
  params: MataPelajaranOptionsParams = {},
): Promise<MataPelajaranOption[]> {
  const { source = "dataMaster", tingkatKelasId } = params;

  if (source === "bankSoal") {
    const data = await getMapel({
      tingkatKelas: tingkatKelasId,
      limit: 50,
      offset: 0,
    });
    return data.map((mapel) => ({
      id: mapel.id,
      label: mapel.namaMapel,
    }));
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
  const [bankSoalData, mapelData] = await Promise.all([
    getBankSoal({
      id_kelas: params.tingkatKelasId,
      limit: 50,
      offset: 0,
    }),
    getMapel({
      limit: 50,
      offset: 0,
    }),
  ]);

  const mapelById = new Map<number, string>(
    mapelData.map((mapel) => [mapel.id, mapel.namaMapel]),
  );

  return bankSoalData.map((item) => ({
    id: item.id_bank_soal,
    nama: item.nama_bank_soal || "-",
    mata_pelajaran: mapelById.get(item.id_mapel),
    materi: item.materi,
    kelas: item.id_kelas,
    total_soal: 0,
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
