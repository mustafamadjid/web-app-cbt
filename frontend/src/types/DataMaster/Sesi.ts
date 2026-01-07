export type SesiFormValues = {
  kode_sesi: string;
  nama_sesi: string;
};

export type SesiRow = {
  id: string;
  kodeSesi: string;
  namaSesi: string;
};

export type SesiFilterParams = {
  q?: string;
};
