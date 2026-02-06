import { api } from "@/services/Api/api";

import type {
  ProfilSekolahResponse,
  ProfilSekolahUpdatePayload,
} from "@/types/ProfilSekolah/ProfilSekolah";

const PROFIL_SEKOLAH_ENDPOINT = "/admin/profil-sekolah";

export async function getProfilSekolah(): Promise<ProfilSekolahResponse> {
  return api<ProfilSekolahResponse>(PROFIL_SEKOLAH_ENDPOINT, {
    method: "GET",
  });
}


export async function updateProfilSekolah(
  formData: FormData,
): Promise<ProfilSekolahUpdatePayload> {
  return api<ProfilSekolahUpdatePayload>(PROFIL_SEKOLAH_ENDPOINT, {
    method: "PATCH",
    data: formData,
  });
}
