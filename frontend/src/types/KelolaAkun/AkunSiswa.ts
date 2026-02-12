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

  id_nama_kelas: string | "";

  foto_profil: File | null;
};

export type StudentUpdateFormValues = Omit<StudentRegisterFormValues, "password"> & {
  role: string;
  status_akun: StatusAkun;
  id_pengguna?: number;
};

export type StudentUpdatePayload = Partial<StudentUpdateFormValues>;

export type StudentListResponseItem = {
  id_pengguna: number;
  username: string;
  email: string;
  nama_lengkap: string;
  jenis_kelamin: JenisKelamin;
  no_hp: string;
  foto_profil: string;
  status_akun: StatusAkun;
  nama_kelas: string;
  tingkat_kelas: number;
  angkatan: number;
  no_absen: number;
  tempat_lahir: string;
  tanggal_lahir: string;
  role: string;
};

export type StudentDetailResponse = StudentListResponseItem & {
  nisn: string;
  kelas: string;
};

export type DataAkunSiswa = StudentDetailResponse;

export type StudentRegisterResponse = {
  id: number;
};
