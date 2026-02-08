import { api } from "@/services/Api/api";
import type { DataAkunSiswa } from "@/types/KelolaAkun/AkunSiswa";
import type { DataGuru } from "@/types/KelolaAkun/AkunGuru";

export type ProfileData = DataGuru | DataAkunSiswa;

export async function getProfileByUserId(userId: number) {
  const queryParams: Record<string, string | number | undefined> = {
    id: userId,
  };

  return api<ProfileData>("/profile", {
    method: "GET",
    params: queryParams,
  });
}
