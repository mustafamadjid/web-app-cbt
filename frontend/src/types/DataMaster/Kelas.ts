export type KelasFormValues = {
  tingkat_kelas: number | "";
  nama_kelas: string;
};

export type KelasRow = {
  id: string;
  id_tingkat_kelas: number;
  tingkat_kelas: number;
  nama_kelas: string;
};

export type TingkatKelasOption = {
  id_tingkat_kelas: number;
  tingkat_kelas: number;
};

export type KelasFilterParams = {
  q?: string;
  tingkatKelas?: number;
};

export type KelasSubmitResponse = {
  id: number;
};
