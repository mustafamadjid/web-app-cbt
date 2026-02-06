import { api } from "@/services/Api/api";

import type { ApiEnvelope } from "@/services/Api/api";
import type {
  ProfilSekolahResponse,
  ProfilSekolahUpdatePayload,
} from "@/types/ProfilSekolah/ProfilSekolah";

const PROFIL_SEKOLAH_ENDPOINT = "/admin/profil-sekolah";

export async function getProfilSekolah(): Promise<ProfilSekolahResponse> {
  const res = await api<ApiEnvelope<ProfilSekolahResponse>>(
    PROFIL_SEKOLAH_ENDPOINT,
    {
      method: "GET",
    }
  );

  return res.data;
}

export async function updateProfilSekolah(
  formData: FormData
): Promise<ProfilSekolahUpdatePayload> {
  const res = await api<ApiEnvelope<ProfilSekolahUpdatePayload>>(
    PROFIL_SEKOLAH_ENDPOINT,
    {
      method: "PATCH",
      data: formData,
    }
  );

  return res.data;
}
