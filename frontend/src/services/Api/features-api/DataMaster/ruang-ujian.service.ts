import type {
  RuangUjianFilterParams,
  RuangUjianRow,
} from "@/types/DataMaster/RuangUjian";

const DUMMY_RUANG_UJIAN: RuangUjianRow[] = [
  { id: 1, namaRuangan: "Ruang Ujian 01" },
  { id: 2, namaRuangan: "Ruang Ujian 02" },
  { id: 3, namaRuangan: "Lab Komputer" },
  { id: 4, namaRuangan: "Aula Utama" },
  {id: 5, namaRuangan: "Lab IPA"},
  {id: 6, namaRuangan: "Lab IPS"}
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

export async function getRuangUjianOptions(): Promise<RuangUjianRow[]> {
  await sleep(180);
  

  // Ini harusnya nanit kalau sudah fetch API ga perlu sorting id
  return [...DUMMY_RUANG_UJIAN].sort((a, b) => a.id - b.id);
}
