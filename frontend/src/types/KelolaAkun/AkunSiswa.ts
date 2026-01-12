import type { JenisKelamin, StatusAkun } from "../OpsiTypes/Option";

export type StudentRegisterFormValues = {
  role : "SISWA";
  namaLengkap: string;
  username: string;
  password: string;

  jenisKelamin: JenisKelamin;

  email?: string;
  noHp?: string;

  noAbsen: number;
  angkatan: number;

  tempatLahir: string;
  tanggalLahir: string;

  id_tingkat_kelas: number | "";
  id_nama_kelas: string | "";

  fotoProfil: File | null;

  statusAkun: StatusAkun;
};

export type DataAkunSiswa = Omit<StudentRegisterFormValues, "fotoProfil" | "password"> &{
  id:number;
  urlGambarProfil:string;
}


export type StudentRegisterResponse = {
    id: number;
}
