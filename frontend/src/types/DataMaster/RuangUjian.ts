export type RuangUjianFormValues = {
  nama_ruangan_ujian: string;
};

export type RuangUjianRow = {
  id: number;
  namaRuangan: string;
};

export type RuangUjianFilterParams = {
  q?: string;
  limit?: number;
  offset?: number;
};
