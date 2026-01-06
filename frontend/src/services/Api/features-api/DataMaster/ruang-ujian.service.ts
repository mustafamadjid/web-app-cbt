import type {
  RuangUjianFilterParams,
  RuangUjianRow,
} from "@/types/DataMaster/RuangUjian";

const DUMMY_RUANG_UJIAN: RuangUjianRow[] = [
  { id: "ruang-1", namaRuangan: "Ruang Ujian 01" },
  { id: "ruang-2", namaRuangan: "Ruang Ujian 02" },
  { id: "ruang-3", namaRuangan: "Lab Komputer" },
  { id: "ruang-4", namaRuangan: "Aula Utama" },
];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getRuangUjian(
  params: RuangUjianFilterParams = {}
): Promise<RuangUjianRow[]> {
  await sleep(200);
  const q = params.q ? normalize(params.q) : "";

  if (!q) return DUMMY_RUANG_UJIAN;

  return DUMMY_RUANG_UJIAN.filter((ruang) =>
    ruang.namaRuangan.toLowerCase().includes(q)
  );
}
