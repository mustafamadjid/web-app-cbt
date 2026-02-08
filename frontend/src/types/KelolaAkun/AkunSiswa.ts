import type { JenisKelamin, StatusAkun } from "../OpsiTypes/Option";

export type StudentRegisterFormValues = {
  nama_lengkap: string;
  username: string;
  password: string;

  jenis_kelamin: JenisKelamin;

  email?: string;
  no_hp?: string;

  no_absen: number;
  angkatan: number;
  nisn: string;

  tempat_lahir: string;
  tanggal_lahir: string;

  id_tingkat_kelas: number | "";
  id_nama_kelas: string | "";

  foto_profil: File | null;
};

export type DataAkunSiswa = Omit<
  StudentRegisterFormValues,
  "fotoProfil" | "password"
> & {
  id: number;
  urlGambarProfil: string;
};

export type StudentRegisterResponse = {
  id: number;
};
