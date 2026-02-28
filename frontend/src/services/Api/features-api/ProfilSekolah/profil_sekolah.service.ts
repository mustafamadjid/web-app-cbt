import { api } from "@/services/Api/api";
import { buildFormData } from "@/helper/FormData/BuildFormData";
import { useFetch, usePut } from "@/hooks/fetch";

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
  payload: Partial<ProfilSekolahUpdatePayload>,
  logoFile?: File | null,
): Promise<ProfilSekolahUpdatePayload> {
  const formData = buildFormData(payload, {
    transform: (_key, value) => {
      if (typeof value === "string") return value.trim();
      return value as any;
    },
  });

  if (logoFile) {
    formData.append("logo_sekolah", logoFile);
  }

  return api<ProfilSekolahUpdatePayload>(PROFIL_SEKOLAH_ENDPOINT, {
    method: "PATCH",
    data: formData,
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetProfilSekolah() {
  return useFetch(() => getProfilSekolah(), []);
}

export function useUpdateProfilSekolah() {
  return usePut(
    (payload: { data: Partial<ProfilSekolahUpdatePayload>; logoFile?: File | null }) =>
      updateProfilSekolah(payload.data, payload.logoFile),
  );
}
