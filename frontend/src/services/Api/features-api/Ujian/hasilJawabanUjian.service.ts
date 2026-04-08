import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useFetch, usePut } from "@/hooks/fetch";
import type {
  HasilJawabanUjianResponse,
} from "@/types/Ujian/HasilJawabanUjian";
import type { SubmitKoreksiEssayRequest } from "@/types/Ujian/SubmitKoreksiEssay";
import { api } from "../../api";

export async function getHasilJawabanUjian(
  idAttempt: number,
): Promise<HasilJawabanUjianResponse> {
  return api<HasilJawabanUjianResponse>(
    `/ujian/jawaban/hasil/${idAttempt}`,
    { method: "GET" },
  );
}

export async function submitKoreksiEssay(
  payload: SubmitKoreksiEssayRequest,
): Promise<boolean> {
  const data = buildJsonData(payload, { nullishToEmptyString: false });

  return api<boolean>("/ujian/koreksi-essay", {
    method: "PATCH",
    data,
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetHasilJawabanUjian(
  idAttempt: number,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getHasilJawabanUjian(idAttempt)
        : Promise.resolve({
            id_attempt: idAttempt,
            nilai_akhir: null,
            hasil_jawaban: [],
          } as HasilJawabanUjianResponse),
    [idAttempt, enabled],
  );
}

export function useSubmitKoreksiEssay() {
  return usePut((payload: SubmitKoreksiEssayRequest) => submitKoreksiEssay(payload));
}
