import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useFetch, usePost } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type {
  AttemptUjianRequest,
  ListPesertaUjianSubmittedResponse,
} from "@/types/Ujian/AttemptUjian";

const ATTEMPT_UJIAN_ENDPOINT = "/siswa/ujian/attempt";
const LIST_PESERTA_UJIAN_SUBMITTED_ENDPOINT = "/ujian/peserta-submitted";

export async function attemptUjian(payload: AttemptUjianRequest): Promise<boolean> {
  const data = buildJsonData(payload, { nullishToEmptyString: false });
  return api<boolean>(ATTEMPT_UJIAN_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function getListPesertaUjianSubmitted(
  idJadwalUjian: number,
): Promise<ListPesertaUjianSubmittedResponse> {
  return api<ListPesertaUjianSubmittedResponse>(
    `${LIST_PESERTA_UJIAN_SUBMITTED_ENDPOINT}/${idJadwalUjian}`,
    {
      method: "GET",
    },
  );
}

// =====================
// Hook Wrappers
// =====================

export function useAttemptUjian() {
  return usePost((payload: AttemptUjianRequest) => attemptUjian(payload));
}

export function useGetListPesertaUjianSubmitted(
  idJadwalUjian: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getListPesertaUjianSubmitted(idJadwalUjian)
        : Promise.resolve([] as ListPesertaUjianSubmittedResponse),
    [idJadwalUjian, enabled],
  );
}
