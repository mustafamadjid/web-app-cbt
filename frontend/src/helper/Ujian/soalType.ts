export type NormalizedSoalType = "essay" | "pilihan_ganda" | "unknown";

const canonicalizeSoalType = (value: string | null | undefined) =>
  (value ?? "")
    .trim()
    .toLowerCase()
    .replaceAll("-", "_")
    .replaceAll(" ", "_");

export const normalizeSoalType = (
  value: string | null | undefined,
): NormalizedSoalType => {
  const normalized = canonicalizeSoalType(value);

  if (normalized === "essay") {
    return "essay";
  }

  if (normalized === "pilihan_ganda") {
    return "pilihan_ganda";
  }

  return "unknown";
};

export const isEssaySoal = (value: string | null | undefined) =>
  normalizeSoalType(value) === "essay";

export const isPilihanGandaSoal = (value: string | null | undefined) =>
  normalizeSoalType(value) === "pilihan_ganda";

export const formatSoalTypeLabel = (value: string | null | undefined) => {
  const normalized = normalizeSoalType(value);

  if (normalized === "essay") {
    return "Essay";
  }

  if (normalized === "pilihan_ganda") {
    return "Pilihan Ganda";
  }

  const fallback = canonicalizeSoalType(value);
  if (!fallback) {
    return "Soal";
  }

  return fallback
    .split("_")
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
};
