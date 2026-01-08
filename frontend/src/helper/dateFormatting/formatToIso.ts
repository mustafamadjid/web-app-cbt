import { monthMap } from "./monthMap";
import { normalize } from "../normalizeString/normalizeString";

export const formatTanggalToIso = (value?: string) => {
  if (!value) return null;
  const parts = value.split(",").map((item) => item.trim());
  if (parts.length < 2) return null;

  const dateParts = parts[1].split(" ").filter(Boolean);
  if (dateParts.length < 3) return null;

  const [day, monthName, year] = dateParts;
  const month = monthMap[normalize(monthName)];
  if (!month) return null;

  const paddedDay = day.padStart(2, "0");
  return `${year}-${month}-${paddedDay}`;
};
