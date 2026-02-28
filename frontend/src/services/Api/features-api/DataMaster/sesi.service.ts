import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { api } from "@/services/Api/api";
import { useFetch, usePost, usePut, useDelete } from "@/hooks/fetch";
import type {
  ListSesiResponse,
  SesiFilterParams,
  SesiFormValues,
  SesiRow,
  UpdateSesiPayload,
} from "@/types/DataMaster/Sesi";

const SESI_ENDPOINT = "/admin/sesi";

const toPartialPayload = (
  values: SesiFormValues,
  initialValues: SesiFormValues,
): UpdateSesiPayload => {
  const payload: UpdateSesiPayload = {};

  if (values.kode_sesi !== initialValues.kode_sesi) {
    payload.kode_sesi = values.kode_sesi;
  }

  if (values.nama_sesi !== initialValues.nama_sesi) {
    payload.nama_sesi = values.nama_sesi;
  }

  return payload;
};

export async function createSesi(values: SesiFormValues) {
  const data = buildJsonData(values);

  return await api<null>(SESI_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function getSesi(
  params: SesiFilterParams = {},
): Promise<SesiRow[]> {
  const response = await api<ListSesiResponse>(SESI_ENDPOINT, {
    method: "GET",
    params: {
      q: params.q?.trim() || undefined,
      search: params.search?.trim() || undefined,
      limit: params.limit != null ? String(params.limit) : undefined,
      offset: params.offset != null ? String(params.offset) : undefined,
    },
  });

  return response.items;
}

export async function getSesiById(idSesi: number): Promise<SesiFormValues> {
  const response = await api<SesiRow>(`${SESI_ENDPOINT}/sesi-id/${idSesi}`, {
    method: "GET",
  });

  return {
    kode_sesi: response.kode_sesi,
    nama_sesi: response.nama_sesi,
  };
}

export async function updateSesiPartial(
  idSesi: number,
  values: SesiFormValues,
  initialValues: SesiFormValues,
) {
  const payload = toPartialPayload(values, initialValues);
  const data = buildJsonData(payload);

  return await api<null>(`${SESI_ENDPOINT}/${idSesi}`, {
    method: "PATCH",
    data,
  });
}

export async function deleteSesi(idSesi: number) {
  return await api<null>(`${SESI_ENDPOINT}/${idSesi}`, {
    method: "DELETE",
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetSesi(params: SesiFilterParams = {}) {
  return useFetch(
    () => getSesi(params),
    [params.q, params.search, params.limit, params.offset],
  );
}

export function useGetSesiById(idSesi: number) {
  return useFetch(() => getSesiById(idSesi), [idSesi]);
}

export function useCreateSesi() {
  return usePost((values: SesiFormValues) => createSesi(values));
}

export function useUpdateSesi() {
  return usePut(
    (payload: { id: number; values: SesiFormValues; initialValues: SesiFormValues }) =>
      updateSesiPartial(payload.id, payload.values, payload.initialValues),
  );
}

export function useDeleteSesi() {
  return useDelete((idSesi: number) => deleteSesi(idSesi));
}
