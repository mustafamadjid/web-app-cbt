
import type { AktivitasLogItem } from "@/types/Log/LogAktivitas";
import { api } from "../../api";
import { useFetch } from "@/hooks/fetch";

export async function getLogAktivitas(): Promise<AktivitasLogItem[]> {
  return await api<AktivitasLogItem[]>(
    "/admin/aktivitas-user",
    {
      method: "GET",
    },
  );
}

// =====================
// Hook Wrappers
// =====================

export function useGetLogAktivitas() {
  return useFetch(() => getLogAktivitas(), []);
}
