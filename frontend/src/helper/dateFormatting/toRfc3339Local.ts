const toNumber = (value: string): number => Number.parseInt(value, 10);

const pad = (value: number) => String(value).padStart(2, "0");

const WIB_OFFSET = "+07:00";

export const toRfc3339Local = (tanggal: string, waktu: string) => {
  const [yearRaw, monthRaw, dayRaw] = tanggal.split("-").map(toNumber);
  const [hourRaw, minuteRaw] = waktu.split(":").map(toNumber);

  const year = yearRaw || 1970;
  const month = pad(monthRaw || 1);
  const day = pad(dayRaw || 1);
  const hour = pad(hourRaw || 0);
  const minute = pad(minuteRaw || 0);
  const second = "00";

  return `${year}-${month}-${day}T${hour}:${minute}:${second}${WIB_OFFSET}`;
};
