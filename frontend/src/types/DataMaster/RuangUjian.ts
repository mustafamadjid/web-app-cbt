export type RuangUjianFormValues = {
  kode_ruang: string;
  nama_ruangan: string;
};

export type RuangUjianRow = {
  id_ruangan: number;
  kode_ruang: string;
  nama_ruangan: string;
};

export type RuangUjianFilterParams = {
  q?: string;
  search?: string;
  limit?: number;
  offset?: number;
};

export type ListRuangUjianResponse = RuangUjianRow[];

export type UpdateRuangUjianPayload = Partial<RuangUjianFormValues>;
