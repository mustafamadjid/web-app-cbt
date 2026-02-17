export type SesiFormValues = {
  kode_sesi: string;
  nama_sesi: string;
};

export type SesiRow = {
  id_sesi: number;
  kode_sesi: string;
  nama_sesi: string;
};

export type SesiFilterParams = {
  q?: string;
  search?: string;
  limit?: number;
  offset?: number;
};

export type ListSesiResponse = {
  items: SesiRow[];
};

export type UpdateSesiPayload = Partial<SesiFormValues>;
