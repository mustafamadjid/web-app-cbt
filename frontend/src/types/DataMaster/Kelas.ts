export type KelasFormValues = {
  tingkat_kelas: number | "";
  nama_kelas: string;
};

export type KelasRow = {
  id: string;
  tingkat_kelas: number;
  nama_kelas: string;
};

export type KelasFilterParams = {
  q?: string;
  tingkatKelas?: number;
};

export type KelasSubmitResponse = {
  id: number;
};