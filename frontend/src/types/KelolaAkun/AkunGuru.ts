import type { O } from "node_modules/react-router/dist/development/router-5fbeEIMQ.d.mts";
import type { StatusAkun,JenisKelamin } from "../OpsiTypes/Option";

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

export type TeacherRegisterResponse = { id: number };