import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { api } from "@/services/Api/api";
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
