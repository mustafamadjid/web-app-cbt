import { buildFormData } from "@/helper/FormData/BuildFormData";
import { api } from "@/services/Api/api";
import type {
  PengumumanCreatePayload,
  PengumumanFormValues,
  PengumumanGetResponse,
  PengumumanUpdatePayload,
} from "@/types/Widget/Pengumuman";

const PENGUMUMAN_ENDPOINT = "/pengumuman";

const toPartialUpdatePayload = (
  values: PengumumanFormValues,
  initialValues: PengumumanFormValues,
): PengumumanUpdatePayload => {
  const payload: PengumumanUpdatePayload = {};

  if (values.judul_pengumuman !== initialValues.judul_pengumuman) {
    payload.judul_pengumuman = values.judul_pengumuman;
  }

  if (values.isi_pengumuman !== initialValues.isi_pengumuman) {
    payload.isi_pengumuman = values.isi_pengumuman;
  }

  if (
    values.tanggal_rilis_pengumuman !== initialValues.tanggal_rilis_pengumuman
  ) {
    payload.tanggal_rilis_pengumuman = values.tanggal_rilis_pengumuman;
  }

  if (
    values.tanggal_selesai_pengumuman !==
    initialValues.tanggal_selesai_pengumuman
  ) {
    payload.tanggal_selesai_pengumuman = values.tanggal_selesai_pengumuman;
  }

  if (values.dokumen_pengumuman) {
    payload.dokumen_pengumuman = values.dokumen_pengumuman;
  }

  return payload;
};

export async function createPengumuman(payload: PengumumanCreatePayload) {
  const formData = buildFormData(payload, {
    transform: (_key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") return value.trim();
      return value as any;
    },
    skipNullish: true,
  });

  return await api<null>(PENGUMUMAN_ENDPOINT, {
    method: "POST",
    data: formData,
  });
}

export async function getPengumumanActive(): Promise<PengumumanGetResponse[]> {
  return await api<PengumumanGetResponse[]>(`${PENGUMUMAN_ENDPOINT}/active`, {
    method: "GET",
  });
}

export async function getPengumumanIncoming(): Promise<PengumumanGetResponse[]> {
  return await api<PengumumanGetResponse[]>(`${PENGUMUMAN_ENDPOINT}/incoming`, {
    method: "GET",
  });
}

export async function getPengumumanNonActive(): Promise<
  PengumumanGetResponse[]
> {
  return await api<PengumumanGetResponse[]>(
    `${PENGUMUMAN_ENDPOINT}/non-active`,
    {
      method: "GET",
    },
  );
}

export async function getPengumumanById(
  idPengumuman: number,
): Promise<PengumumanGetResponse> {
  return await api<PengumumanGetResponse>(
    `${PENGUMUMAN_ENDPOINT}/id/${idPengumuman}`,
    {
      method: "GET",
    },
  );
}

export async function updatePengumumanPartial(
  idPengumuman: number,
  values: PengumumanFormValues,
  initialValues: PengumumanFormValues,
) {
  const payload = toPartialUpdatePayload(values, initialValues);

  const formData = buildFormData(payload, {
    transform: (_key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") return value.trim();
      return value as any;
    },
    skipNullish: true,
  });

  return await api<null>(`${PENGUMUMAN_ENDPOINT}/${idPengumuman}`, {
    method: "PATCH",
    data: formData,
  });
}

export async function deletePengumuman(idPengumuman: number) {
  return await api<null>(`${PENGUMUMAN_ENDPOINT}/${idPengumuman}`, {
    method: "DELETE",
  });
}
