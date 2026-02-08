export type KelasFormValues = {
  tingkat_kelas: number | "";
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

export type FullDataKelas = {
  item_tingkat_kelas: TingkatKelas[];
  item_nama_kelas: NamaKelas[];
};

export type KelasFilterParams = {
  search?: string;
  tingkatKelas?: number;
};

export type KelasSubmitResponse = {
  id: number;
};

