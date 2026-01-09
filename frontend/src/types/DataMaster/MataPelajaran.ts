export type KelasOption = {
  id: number;
  tingkat_kelas: number;
  label: string; // label tingkat kelas
};
export type MataPelajaranRow = {
  id: number;
  kelasId: number;
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranOption = {
  id: number;
  label: string;
};

export type MataPelajaranFormValues = {
  kelasId: number | "";
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranFilterParams = {
  q?: string;
  tingkatKelas?: number;
  mapelId?: number;
};
