import { api } from "@/services/Api/api";
import type { LogAktivitasApiItem } from "@/types/Log/LogAktivitas";

export async function getLogAktivitas() {
  return api<LogAktivitasApiItem[]>("/admin/aktivitas-user", {
    method: "GET",
  });
}
