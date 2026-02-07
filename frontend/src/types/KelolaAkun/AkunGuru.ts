import type { StatusAkun,JenisKelamin } from "../OpsiTypes/Option";

export type TeacherRegisterFormValues = {
  username: string;
  email: string;
  // role?: string;
  nama_lengkap: string;
  password: string;
  no_hp: string;
  jenis_kelamin: JenisKelamin;
  // status_akun: StatusAkun;
  nip: string;
  jabatan: string;
  bidang_studi: string;
  foto_profil: File | null;
};

export type DataGuru = Omit<TeacherRegisterFormValues, "fotoProfil" | "password"> &{
  id_pengguna:number;
  foto_profil:string;
  role:string
  status_akun:StatusAkun
}

export type TeacherRegisterResponse = { id: number };