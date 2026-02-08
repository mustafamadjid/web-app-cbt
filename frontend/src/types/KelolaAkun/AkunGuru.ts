import type { StatusAkun, JenisKelamin } from "../OpsiTypes/Option";

export type TeacherRegisterFormValues = {
  username: string;
  email: string;
  nama_lengkap: string;
  password: string;
  no_hp: string;
  jenis_kelamin: JenisKelamin;
  nip: string;
  jabatan: string;
  bidang_studi: string;
  foto_profil: File | null;
};

export type GuruFilterParams = {
  q?: string;
  status?: StatusAkun;
  limit?: number;
  offset?: number;
};

export type DataGuru = Omit<TeacherRegisterFormValues, "foto_profil" | "password"> &{
  id_pengguna:number;
  foto_profil:string;
  role:string
  status_akun:StatusAkun
}

export type DataUpdateGuruComplete = Omit<
  TeacherRegisterFormValues,
  "password"
> & {
  role: string;
};

export type DataUpdateGuru = Omit<DataGuru,"password">

export type TeacherUpdateFormValues = Omit<TeacherRegisterFormValues, "password"> & {
  role: string;
  status_akun: StatusAkun;
  id_pengguna?: number;
};

export type TeacherUpdatePayload = Partial<TeacherUpdateFormValues>;

export type TeacherRegisterResponse = { id: number };
