export type KelasOption = {
  id: number;
  tingkat_kelas: number;
  label: string; // label tingkat kelas
};
export type MataPelajaranRow = {
  id: number;
  kelasId: number;
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranOption = {
  id: number;
  label: string;
};

export type MataPelajaranFormValues = {
  kelasId: number | "";
  kodeMapel: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranFilterParams = {
  search?: string;
  tingkatKelas?: number;
  namaMapel?: string;
  limit?: number;
  offset?: number;
};

export type MapelItemResponse = {
  id_mapel: number;
  id_kelas: number;
  kode_mapel: string;
  nama_mapel: string;
  deskripsi: string;
};

export type ListMapelResponse = {
  items: MapelItemResponse[];
};

export type CreateMapelPayload = {
  id_kelas: number;
  kode_mapel: string;
  nama_mapel: string;
  deskripsi: string;
};

export type UpdateMapelPayload = Partial<CreateMapelPayload>;
