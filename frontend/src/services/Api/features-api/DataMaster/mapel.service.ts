import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { api } from "@/services/Api/api";
import type {
  CreateMapelPayload,
  ListMapelResponse,
  MapelItemResponse,
  MataPelajaranFilterParams,
  MataPelajaranFormValues,
  MataPelajaranRow,
  UpdateMapelPayload,
} from "@/types/DataMaster/MataPelajaran";

const MAPEL_ENDPOINT = "/admin/mata-pelajaran";

const toMataPelajaranRow = (item: MapelItemResponse): MataPelajaranRow => ({
  id: item.id_mapel,
  kelasId: item.id_kelas,
  kodeMapel: item.kode_mapel,
  namaMapel: item.nama_mapel,
  deskripsiMapel: item.deskripsi,
});

const toCreatePayload = (
  values: MataPelajaranFormValues,
): CreateMapelPayload => ({
  id_kelas: Number(values.kelasId),
  kode_mapel: values.kodeMapel,
  nama_mapel: values.namaMapel,
  deskripsi: values.deskripsiMapel,
});

const toPartialPayload = (
  values: MataPelajaranFormValues,
  initialValues: MataPelajaranFormValues,
): UpdateMapelPayload => {
  const payload: UpdateMapelPayload = {};

  if (Number(values.kelasId) !== Number(initialValues.kelasId)) {
    payload.id_kelas = Number(values.kelasId);
  }
  if (values.kodeMapel !== initialValues.kodeMapel) {
    payload.kode_mapel = values.kodeMapel;
  }
  if (values.namaMapel !== initialValues.namaMapel) {
    payload.nama_mapel = values.namaMapel;
  }
  if (values.deskripsiMapel !== initialValues.deskripsiMapel) {
    payload.deskripsi = values.deskripsiMapel;
  }

  return payload;
};

export async function createMataPelajaran(values: MataPelajaranFormValues) {
  const payload = toCreatePayload(values);
  const data = buildJsonData(payload);

  return api<null>(MAPEL_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function getMapel(
  params: MataPelajaranFilterParams = {},
): Promise<MataPelajaranRow[]> {
  const queryParams: Record<string, string | undefined> = {
    search: params.search?.trim() || undefined,
    tingkat_kelas:
      params.tingkatKelas != null ? String(params.tingkatKelas) : undefined,
    nama_mapel: params.namaMapel?.trim() || undefined,
    limit: params.limit != null ? String(params.limit) : undefined,
    offset: params.offset != null ? String(params.offset) : undefined,
  };

  const response = await api<ListMapelResponse>(MAPEL_ENDPOINT, {
    method: "GET",
    params: queryParams,
  });

  return response.items.map(toMataPelajaranRow);
}

export async function getMapelById(
  idMapel: number,
): Promise<MataPelajaranFormValues> {
  const response = await api<MapelItemResponse>(
    `${MAPEL_ENDPOINT}/${idMapel}`,
    {
      method: "GET",
    },
  );

  return {
    kelasId: response.id_kelas,
    kodeMapel: response.kode_mapel,
    namaMapel: response.nama_mapel,
    deskripsiMapel: response.deskripsi,
  };
}

export async function updateMataPelajaranPartial(
  idMapel: number,
  values: MataPelajaranFormValues,
  initialValues: MataPelajaranFormValues,
) {
  const payload = toPartialPayload(values, initialValues);
  const data = buildJsonData(payload);

  return api<null>(`${MAPEL_ENDPOINT}/${idMapel}`, {
    method: "PATCH",
    data,
  });
}

export async function deleteMataPelajaran(idMapel: number) {
  const data = buildJsonData({ id_mapel: idMapel });

  return api<null>(`${MAPEL_ENDPOINT}/${idMapel}`, {
    method: "DELETE",
    data,
  });
}
