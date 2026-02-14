import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { api } from "@/services/Api/api";
import type {
  ListRuangUjianResponse,
  RuangUjianFilterParams,
  RuangUjianFormValues,
  RuangUjianRow,
  UpdateRuangUjianPayload,
} from "@/types/DataMaster/RuangUjian";

const RUANG_UJIAN_ENDPOINT = "/admin/ruang-ujian";

const toPartialPayload = (
  values: RuangUjianFormValues,
  initialValues: RuangUjianFormValues,
): UpdateRuangUjianPayload => {
  const payload: UpdateRuangUjianPayload = {};

  if (values.kode_ruang !== initialValues.kode_ruang) {
    payload.kode_ruang = values.kode_ruang;
  }

  if (values.nama_ruangan !== initialValues.nama_ruangan) {
    payload.nama_ruangan = values.nama_ruangan;
  }

  return payload;
};

export async function createRuangUjian(values: RuangUjianFormValues) {
  const data = buildJsonData(values);

  return api<null>(RUANG_UJIAN_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function getRuangUjian(
  params: RuangUjianFilterParams = {},
): Promise<RuangUjianRow[]> {
  return api<ListRuangUjianResponse>(RUANG_UJIAN_ENDPOINT, {
    method: "GET",
    params: {
      q: params.q?.trim() || undefined,
      search: params.search?.trim() || undefined,
      limit: params.limit != null ? String(params.limit) : undefined,
      offset: params.offset != null ? String(params.offset) : undefined,
    },
  });
}

export async function getRuangUjianById(
  idRuangan: number,
): Promise<RuangUjianFormValues> {
  const response = await api<RuangUjianRow>(`${RUANG_UJIAN_ENDPOINT}/id/${idRuangan}`, {
    method: "GET",
  });

  return {
    kode_ruang: response.kode_ruang,
    nama_ruangan: response.nama_ruangan,
  };
}

export async function updateRuangUjianPartial(
  idRuangan: number,
  values: RuangUjianFormValues,
  initialValues: RuangUjianFormValues,
) {
  const payload = toPartialPayload(values, initialValues);
  const data = buildJsonData(payload);

  return api<null>(`${RUANG_UJIAN_ENDPOINT}/${idRuangan}`, {
    method: "PATCH",
    data,
  });
}

export async function deleteRuangUjian(idRuangan: number) {
  return api<null>(`${RUANG_UJIAN_ENDPOINT}/${idRuangan}`, {
    method: "DELETE",
  });
}
