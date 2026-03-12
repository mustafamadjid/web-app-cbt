import { api } from "../../api";
import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useFetch, usePost, useDelete, usePut } from "@/hooks/fetch";

import type {
  KelasFilterParams,
  FullDataKelas,
  DataKelas
} from "@/types/DataMaster/Kelas";
import type {
  KelasSubmitResponse,
  KelasFormValues,
} from "@/types/DataMaster/Kelas";


// Real Request
type CreateTingkatKelasPayload = {
  tingkat_kelas: number;
};

type CreateNamaKelasPayload = {
  id_tingkat_kelas: number;
  nama_kelas: string;
};

type UpdateNamaKelasPayload = {
  id_tingkat_kelas?: number;
  nama_kelas?: string;
};


export async function createTingkatKelas(values: CreateTingkatKelasPayload) {
  const data = buildJsonData(values);
  const res = await api<boolean>("/admin/kelas/tingkat-kelas", {
    method: "POST",
    data,
  });

  return res;
}

export async function createNamaKelas(values: CreateNamaKelasPayload) {
  const data = buildJsonData(values);
  const res = await api<boolean>("/admin/kelas/nama-kelas", {
    method: "POST",
    data,
  });

  return res;
}
export async function GetDataKelasFull(
  params: KelasFilterParams = {},
): Promise<FullDataKelas> {

  const queryParams : Record<string, string | undefined> = {
    search: params.search || undefined,
    tingkat_kelas: params.tingkatKelas != null ? String(params.tingkatKelas) : undefined,
    limit: params.limit != null ? String(params.limit) : undefined,
    offset: params.offset != null ? String(params.offset) : undefined,
  };

  return api<FullDataKelas>("/admin/kelas",{
    method: "GET",
    params : queryParams
  })
}



export async function getKelasByIdsRequest(
  idTingkatKelas: number,
  idNamaKelas: number,
): Promise<DataKelas | null> {
    return api<DataKelas>(`/admin/kelas/${idTingkatKelas}/${idNamaKelas}`, {
      method: "GET",
    });  
}


export async function deleteNamaKelas(idNamaKelas: number) {
  return api<boolean>(`/admin/kelas/nama-kelas/${idNamaKelas}`, {
    method: "DELETE",
  });
}

export async function updateNamaKelasPartial(
  idNamaKelas: number,
  values: UpdateNamaKelasPayload,
) {
  const data = buildJsonData(values);
  return await api<UpdateNamaKelasPayload>(
    `/admin/kelas/nama-kelas/${idNamaKelas}`,
    {
      method:"PATCH",
      data
    }
  )
}

// Submit Handler
export async function submitKelasResponse(values: KelasFormValues) {
  const data = buildJsonData(values);
  const res = await api<KelasSubmitResponse>("/kelas", {
    method: "POST",
    data,
  });

  return res;
}



// =====================
// Hook Wrappers
// =====================

export function useGetDataKelasFull(params: KelasFilterParams = {}) {
  return useFetch(
    () => GetDataKelasFull(params),
    [params.search, params.tingkatKelas, params.limit, params.offset],
  );
}


export function useGetKelasByIds(idTingkatKelas: number, idNamaKelas: number) {
  return useFetch(
    () => getKelasByIdsRequest(idTingkatKelas, idNamaKelas),
    [idTingkatKelas, idNamaKelas],
  );
}


export function useCreateTingkatKelas() {
  return usePost((values: CreateTingkatKelasPayload) => createTingkatKelas(values));
}

export function useCreateNamaKelas() {
  return usePost((values: CreateNamaKelasPayload) => createNamaKelas(values));
}

export function useDeleteNamaKelas() {
  return useDelete((id: number) => deleteNamaKelas(id));
}

export function useUpdateNamaKelas() {
  return usePut(
    (payload: { id: number; values: UpdateNamaKelasPayload }) =>
      updateNamaKelasPartial(payload.id, payload.values),
  );
}
