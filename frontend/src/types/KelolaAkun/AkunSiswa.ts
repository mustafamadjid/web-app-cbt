import type { JenisKelamin } from "../OpsiTypes/Option";

export type StudentRegisterFormValues = {
  namaLengkap: string;
  username: string;
  password: string;

  jenisKelamin: JenisKelamin;

  email?: string;
  noHp?: string;

  noAbsen: string;
  angkatan: string;

  tempatLahir: string;
  tanggalLahir: string;

  kelasId: number | "";

  fotoProfil: File | null;
};

export type StudentRegisterResponse = {
    id: number;
}
