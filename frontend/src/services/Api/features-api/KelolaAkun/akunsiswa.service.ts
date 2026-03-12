import { buildFormData } from "@/helper/FormData/BuildFormData";

import type { JenisKelamin, StatusAkun } from "@/types/OpsiTypes/Option";
import type {
  DataAkunSiswa,
  StudentDetailResponse,
  StudentRegisterFormValues,
  StudentRegisterResponse,
  StudentUpdatePayload,
} from "@/types/KelolaAkun/AkunSiswa";
import { api } from "../../api";
import { useFetch, usePost, usePut } from "@/hooks/fetch";

export type BarisSiswa = DataAkunSiswa & {
  kelas: string;
};

export type SiswaFilterParams = {
  q?: string;
  status?: StatusAkun;
  limit?: number;
  offset?: number;
  angkatan?: number;
  tingkatKelas?: number;
  jenisKelamin?: JenisKelamin;
  idKelas?: number;
  idNamaKelas?: number;
};

const mapJenisKelaminFilter = (jenisKelamin?: JenisKelamin) => {
  if (!jenisKelamin) return undefined;
  return jenisKelamin === "LAKI_LAKI" ? "1" : "2";
};

export const DUMMY_SISWA: DataAkunSiswa[] = [
  {
    id_pengguna: 1,
    role: "SISWA",
    nama_lengkap: "Siti Aminah",
    username: "siti.aminah",
    email: "siti.aminah@gmail.com",
    no_hp: "081234567890",
    jenis_kelamin: "PEREMPUAN",
    status_akun: "AKTIF",
    nisn: "1234567890",
    no_absen: 12,
    angkatan: 2025,
    tempat_lahir: "Bandung",
    tanggal_lahir: "2008-01-31",
    tingkat_kelas: 10,
    nama_kelas: "X IPA 1",
    kelas: "X IPA 1",
    foto_profil: "https://i.pravatar.cc/150?u=s-0001",
  },
  {
    id_pengguna: 2,
    role: "SISWA",
    nama_lengkap: "Raka Pratama",
    username: "raka.pratama",
    email: "",
    no_hp: "",
    jenis_kelamin: "LAKI_LAKI",
    status_akun: "NONAKTIF",
    nisn: "2234567890",
    no_absen: 7,
    angkatan: 2024,
    tempat_lahir: "Jakarta",
    tanggal_lahir: "2009-08-12",
    tingkat_kelas: 10,
    nama_kelas: "X IPS 1",
    kelas: "X IPS 1",
    foto_profil: "https://i.pravatar.cc/150?u=s-0002",
  },
];

export async function submitStudentRegister(values: StudentRegisterFormValues) {
  const formData = buildFormData(values, {
    transform: (key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") {
        if (key === "email") return value.trim().toLowerCase();
        if (key === "password") return value;
        return value.trim();
      }
      return value as string | number | boolean | Blob | null | undefined;
    },
  });

  return api<StudentRegisterResponse>("/admin/siswa", {
    method: "POST",
    data: formData,
  });
}

export async function getSiswaById(
  id: number,
): Promise<StudentDetailResponse | null> {
  if (!id) return null;
  return api<StudentDetailResponse>(`/admin/siswa/${id}`, {
    method: "GET",
  });
}

export async function updateSiswa(id: number, payload: StudentUpdatePayload) {
  const formData = buildFormData(payload, {
    transform: (key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") {
        const trimmed = value.trim();
        if ((key === "email" || key === "no_hp") && trimmed === "") return null;
        if (key === "email") return trimmed.toLowerCase();
        return trimmed;
      }
      return value as string | number | boolean | Blob | null | undefined;
    },
    skipNullish: true,
  });

  return api<StudentRegisterResponse>(`/admin/siswa/${id}`, {
    method: "PATCH",
    data: formData,
  });
}

export async function GetListSiswa(
  params: SiswaFilterParams = {},
): Promise<DataAkunSiswa[]> {
  const queryParams: Record<string, string | undefined> = {
    q: params.q || undefined,
    status: params.status || undefined,
    limit: params.limit ? String(params.limit) : undefined,
    offset: params.offset ? String(params.offset) : undefined,
    angkatan: params.angkatan ? String(params.angkatan) : undefined,
    tingkat_kelas: params.tingkatKelas ? String(params.tingkatKelas) : undefined,
    id_tingkat_kelas: params.idKelas ? String(params.idKelas) : undefined,
    id_nama_kelas: params.idNamaKelas ? String(params.idNamaKelas) : undefined,
    jenis_kelamin: mapJenisKelaminFilter(params.jenisKelamin),
  };

  return await api<DataAkunSiswa[]>("/admin/siswa", {
    params: queryParams,
  });

 
}

// =====================
// Hook Wrappers
// =====================

export function useGetListSiswa(params: SiswaFilterParams = {}, enabled = true) {
  return useFetch(
    () => (enabled ? GetListSiswa(params) : Promise.resolve([])),
    [params.q, params.status, params.limit, params.offset, params.angkatan, params.tingkatKelas, params.jenisKelamin, params.idKelas, params.idNamaKelas, enabled],
  );
}

export function useGetSiswaById(id: number) {
  return useFetch(() => getSiswaById(id), [id]);
}

export function useSubmitStudentRegister() {
  return usePost((values: StudentRegisterFormValues) => submitStudentRegister(values));
}

export function useUpdateSiswa() {
  return usePut(
    (payload: { id: number; values: StudentUpdatePayload }) =>
      updateSiswa(payload.id, payload.values),
  );
}
