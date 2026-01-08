export const TINGKAT_KELAS_BY_ID: Record<number, number> = {
  1: 10,
  2: 11,
  3: 12,
};

export const getTingkatKelasById = (
  id?: number | string | null
): number | undefined => {
  if (id == null || id === "") return undefined;
  const parsed = typeof id === "number" ? id : Number(id);
  if (!Number.isFinite(parsed)) return undefined;
  return TINGKAT_KELAS_BY_ID[parsed];
};
