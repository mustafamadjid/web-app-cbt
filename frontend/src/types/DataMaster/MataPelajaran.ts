export type KelasOption = {
  id: string;
  tingkat_kelas: number;
  label: string; // label tingkat kelas
};
export type MataPelajaranRow = {
  id: string;
  kelasId: string;
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranOption = {
  id: string;
  label: string;
};

export type MataPelajaranFormValues = {
  kelasId: string;
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranFilterParams = {
  q?: string;
  tingkatKelas?: number;
  mapelId?: string;
};
