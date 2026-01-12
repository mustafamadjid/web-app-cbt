import type { BuatUjianFormValues } from "@/types/Ujian/BuatUjian";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type { ReactNode } from "react";


export type PrintJenis = "daftar-hadir" | "berita-acara" | "kartu-peserta";
export type DetailUjianItem = JadwalUjianItem &
  BuatUjianFormValues & {
    id: number;
  };

  export type InfoItem = {
    label: string;
    value: string;
    icon?: ReactNode;
  };

