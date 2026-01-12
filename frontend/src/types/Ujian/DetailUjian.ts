import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";

export type DetailUjianItem = JadwalUjianItem &
  BuatUjianFormValues & {
    id: number;
  };
