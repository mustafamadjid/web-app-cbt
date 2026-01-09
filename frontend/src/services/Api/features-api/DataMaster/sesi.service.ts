import type { SesiFilterParams, SesiRow } from "@/types/DataMaster/Sesi";

const DUMMY_SESI: SesiRow[] = [
  { id: 1, kodeSesi: "SESI-01", namaSesi: "Sesi Pagi" },
  { id: 2, kodeSesi: "SESI-02", namaSesi: "Sesi Siang" },
  { id: 3, kodeSesi: "SESI-03", namaSesi: "Sesi Sore" },
  { id: 4, kodeSesi: "SESI-04", namaSesi: "Sesi Malam" },
];

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const normalize = (value: string) => value.toLowerCase().trim();

export async function getSesi(
  params: SesiFilterParams = {}
): Promise<SesiRow[]> {
  await sleep(200);
  const q = params.q ? normalize(params.q) : "";

  if (!q) return DUMMY_SESI;

  return DUMMY_SESI.filter((sesi) => {
    return (
      sesi.kodeSesi.toLowerCase().includes(q) ||
      sesi.namaSesi.toLowerCase().includes(q)
    );
  });
}
