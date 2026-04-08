import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";

export type HasilUjianListApiItem = Omit<
  JadwalUjianItem,
  "status_ujian" | "started"
> & {
  status_ujian: string;
  started: number;
};
