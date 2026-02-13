export type SesiFormValues = {
  kode_sesi: string;
  nama_sesi: string;
};

export type SesiRow = {
  id: number;
  kodeSesi: string;
  namaSesi: string;
};

export type SesiFilterParams = {
  q?: string;
  limit?: number;
  offset?: number;
};
