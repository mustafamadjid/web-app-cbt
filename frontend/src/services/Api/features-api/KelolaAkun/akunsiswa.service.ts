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
};

const mapJenisKelaminFilter = (jenisKelamin?: JenisKelamin) => {
  if (!jenisKelamin) return undefined;
  return jenisKelamin === "LAKI_LAKI" ? "1" : "2";
};



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
        if (key === "email") return value.trim().toLowerCase();
        return value.trim();
      }
      return value as any;
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
    jenis_kelamin: mapJenisKelaminFilter(params.jenisKelamin),
  };

  return await api<DataAkunSiswa[]>("/admin/siswa", {
    params: queryParams,
  });

 
}
