
import type { AktivitasLogItem } from "@/types/Log/LogAktivitas";
import { api } from "../../api";

export async function getLogAktivitas(): Promise<AktivitasLogItem[]> {
  return await api<AktivitasLogItem[]>(
    "/admin/aktivitas-user",
    {
      method: "GET",
    },
  );
}
