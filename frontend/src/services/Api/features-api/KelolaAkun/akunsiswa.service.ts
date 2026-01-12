import { buildFormData } from "@/helper/FormData/BuildFormData";

import type { JenisKelamin } from "@/types/OpsiTypes/Option";
import type { StudentRegisterFormValues, StudentRegisterResponse } from "@/types/KelolaAkun/AkunSiswa";
import { api, type ApiEnvelope } from "../../api";
import { getTingkatKelasById } from "@/services/Api/features-api/DataMaster/kelas.service";

import type { DataAkunSiswa } from "@/types/KelolaAkun/AkunSiswa";

// export type BarisSiswa = {
//   id: number;
//   namaLengkap: string;
//   username: string;
//   email?: string;
//   noHp?: string;
//   jenisKelamin: JenisKelamin;
//   statusAkun: "aktif" | "nonaktif" | "dibekukan";
//   noAbsen: number;
//   angkatan: number;
//   tempatLahir: string;
//   tanggalLahir: string; // yyyy-mm-dd
//   kelas: string;
//   urlGambarProfil: string;
// };

export type SiswaFilterParams = {
  q?: string;
  angkatan?: number;
  tingkatKelasId?: number;
  jenisKelamin?: JenisKelamin;
};


export const DUMMY_JENIS_KELAMIN: Array<{
  value: JenisKelamin;
  label: string;
}> = [
  { value: "LAKI_LAKI", label: "Laki-laki" },
  { value: "PEREMPUAN", label: "Perempuan" },
];

export const DUMMY_SISWA: DataAkunSiswa[] = [
  {
    id: 1,
    kelasId: 11,
    namaLengkap: "Siti Aminah",
    username: "siti.aminah",
    email: "siti.aminah@gmail.com",
    noHp: "081234567890",
    jenisKelamin: "PEREMPUAN",
    statusAkun: "aktif",
    noAbsen: 12,
    angkatan: 2025,
    tempatLahir: "Bandung",
    tanggalLahir: "2008-01-31",
    kelas: "XI IPA 1",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0001",
  },
  {
    id: 2,
    __kelasId: 10,
    namaLengkap: "Raka Pratama",
    username: "raka.pratama",
    email: "",
    noHp: "",
    jenisKelamin: "LAKI_LAKI",
    statusAkun: "nonaktif",
    noAbsen: 7,
    angkatan: 2024,
    tempatLahir: "Jakarta",
    tanggalLahir: "2009-08-12",
    kelas: "X IPS 1",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0002",
  },
  {
    id: 3,
    __kelasId: 10,
    namaLengkap: "Dimas Saputra",
    username: "dimas.saputra",
    email: "dimas.saputra@mail.com",
    noHp: "082198765432",
    jenisKelamin: "LAKI_LAKI",
    statusAkun: "aktif",
    noAbsen: 3,
    angkatan: 2025,
    tempatLahir: "Surabaya",
    tanggalLahir: "2009-02-20",
    kelas: "X IPA 1",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0003",
  },
  {
    id: 4,
    __kelasId: 11,
    namaLengkap: "Nadya Putri",
    username: "nadya.putri",
    email: "nadya.putri@gmail.com",
    noHp: "081355500011",
    jenisKelamin: "PEREMPUAN",
    statusAkun: "dibekukan",
    noAbsen: 18,
    angkatan: 2023,
    tempatLahir: "Semarang",
    tanggalLahir: "2008-11-05",
    kelas: "XI IPS 1",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0004",
  },
  {
    id: 5,
    __kelasId: 12,
    namaLengkap: "Bagas Wiratama",
    username: "bagas.wiratama",
    email: "bagas.wiratama@school.id",
    noHp: "081200011122",
    jenisKelamin: "LAKI_LAKI",
    statusAkun: "aktif",
    noAbsen: 9,
    angkatan: 2024,
    tempatLahir: "Depok",
    tanggalLahir: "2007-06-14",
    kelas: "XII IPA 2",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0005",
  },
  {
    id: 6,
    __kelasId: 10,
    namaLengkap: "Alya Maharani",
    username: "alya.maharani",
    email: "alya.maharani@mail.com",
    noHp: "085700099988",
    jenisKelamin: "PEREMPUAN",
    statusAkun: "aktif",
    noAbsen: 21,
    angkatan: 2025,
    tempatLahir: "Bogor",
    tanggalLahir: "2009-09-01",
    kelas: "X IPS 1",
    urlGambarProfil: "https://i.pravatar.cc/150?u=s-0006",
  },
];
/** === DUMMY OPTIONS === */
export const DUMMY_ANGKATAN: number[] = [2023, 2024, 2025];

// Submit Data
export async function submitStudentRegister(values: StudentRegisterFormValues) {
  const formData = buildFormData(values, {
    transform: (key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") {
        if (key === "email") return value.trim().toLowerCase();
        if (key === "password") return value;
        return value.trim();
      }
      return value as any;
    },
  });

  const res = await api<ApiEnvelope<StudentRegisterResponse>>(
    "/students/register",
    {
      method: "POST",
      data: formData,
    }
  );

  return res.data;
}


/** === MOCK "API" (simulasikan network delay) === */
const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

function normalize(s: string) {
  return s.toLowerCase().trim();
}

export async function getAngkatanOptions(): Promise<number[]> {
  await sleep(250);
  return DUMMY_ANGKATAN;
}

export async function getJenisKelaminOptions(): Promise<
  Array<{ value: JenisKelamin; label: string }>
> {
  await sleep(150);
  return DUMMY_JENIS_KELAMIN;
}

export async function getSiswa(params: SiswaFilterParams): Promise<BarisSiswa[]> {
  await sleep(350);

  let data = [...DUMMY_SISWA];
  const tingkatKelas = getTingkatKelasById(params.tingkatKelasId);
  const kelasId = tingkatKelas ?? undefined;

  if (params.angkatan) {
    data = data.filter((s) => s.angkatan === params.angkatan);
  }

  if (kelasId != null) {
    data = data.filter((s) => s.__kelasId === kelasId);
  }

  if (params.jenisKelamin) {
    data = data.filter((s) => s.jenisKelamin === params.jenisKelamin);
  }

  if (params.q) {
    const q = normalize(params.q);
    data = data.filter((s) => {
      const hay = normalize(
        [
          s.namaLengkap,
          s.username,
          s.email ?? "",
          s.noHp ?? "",
          s.kelas,
          String(s.noAbsen),
          String(s.angkatan),
          s.tempatLahir,
          s.tanggalLahir,
          s.statusAkun,
        ].join(" ")
      );
      return hay.includes(q);
    });
  }

  // strip internal field sebelum return
  return data.map(({ __kelasId, ...rest }) => rest);
}
