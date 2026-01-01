export type KelasOption = {
  id: string;
  label: string;
};

export type MataPelajaranRow = {
  id: string;
  kelasId: string;
  namaMapel: string;
  deskripsiMapel: string;
};

export type MataPelajaranOption = {
  id: string;
  label: string;
};

export type MataPelajaranFormValues = {
  kelasId: string;
  namaMapel: string;
  deskripsiMapel: string;
};
