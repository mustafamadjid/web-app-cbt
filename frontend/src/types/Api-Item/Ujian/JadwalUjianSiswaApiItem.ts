import type { JadwalUjianSiswaItem } from "@/types/Ujian/jadwalUjian";

export type JadwalUjianSiswaApiItem = Omit<
  JadwalUjianSiswaItem,
  "status_ujian" | "pengawas_nama_lengkap"
> & {
  status_ujian: string;
  pengawas_nama_lengkap?: string | null;
};
