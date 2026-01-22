export type KelasFormValues = {
  tingkat_kelas: number | "";
  nama_kelas: string;
};

export type KelasRow = {
  id: number;
  id_tingkat_kelas: number;
  tingkat_kelas: number;
  nama_kelas: string;
};

export type TingkatKelas = {
  id_tingkat_kelas: number;
  tingkat_kelas: number;
};
export type NamaKelas = {
  id_nama_kelas: number;
  id_tingkat_kelas: number;
  nama_kelas: string;
};

export type KelasFilterParams = {
  q?: string;
  tingkatKelas?: number;
};

export type KelasSubmitResponse = {
  id: number;
};



