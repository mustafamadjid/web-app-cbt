import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { usePost } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { AttemptUjianRequest } from "@/types/Ujian/AttemptUjian";

const ATTEMPT_UJIAN_ENDPOINT = "/siswa/ujian/attempt";

export async function attemptUjian(payload: AttemptUjianRequest): Promise<boolean> {
  const data = buildJsonData(payload, { nullishToEmptyString: false });
  return api<boolean>(ATTEMPT_UJIAN_ENDPOINT, {
    method: "POST",
    data,
  });
}

// =====================
// Hook Wrappers
// =====================

export function useAttemptUjian() {
  return usePost((payload: AttemptUjianRequest) => attemptUjian(payload));
}
