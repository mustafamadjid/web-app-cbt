import type { StatusAkun,JenisKelamin } from "../OpsiTypes/Option";

export type TeacherRegisterFormValues = {
  namaLengkap: string;
  email: string;
  username: string;
  password: string;
  noHp: string;
  jenisKelamin: JenisKelamin;
  statusAkun: StatusAkun;
  nip: string;
  jabatan: string;
  bidangStudi: string;
  fotoProfil: File | null;
};

export type TeacherRegisterResponse = { id: number };