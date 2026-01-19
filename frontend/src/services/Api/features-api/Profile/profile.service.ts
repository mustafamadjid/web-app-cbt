import { api, type ApiEnvelope } from "@/services/Api/api";
import type { DataAkunSiswa } from "@/types/KelolaAkun/AkunSiswa";
import type { DataGuru } from "@/types/KelolaAkun/AkunGuru";

export type ProfileData = DataGuru | DataAkunSiswa;

export async function getProfileByUserId(userId: number) {
  const res = await api<ApiEnvelope<ProfileData>>("/profile", {
    method: "GET",
    params: {
      id: userId,
    },
  });

  return res.data;
}
