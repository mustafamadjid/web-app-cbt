import { api, type ApiEnvelope } from "@/services/Api/api";
import type { AktivitasLogItem } from "@/types/Log/LogAktivitas";

export async function getLogAktivitas(): Promise<AktivitasLogItem[]> {
  const res = await api<ApiEnvelope<AktivitasLogItem[]>>(
    "/admin/aktivitas-user",
    {
      method: "GET",
    },
  );

  return res.data;
}
