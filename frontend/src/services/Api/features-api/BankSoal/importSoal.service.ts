import { api } from "../../api";
import { buildFormData } from "@/helper/FormData/BuildFormData";
import { usePost } from "@/hooks/fetch";

export type ImportSoalResponse = {
  id_job: number;
};

export type ImportSoalJobResponse = {
  id_job: number;
  id_bank_soal: number;
  status: string;
  error_msg?: string;
  warning_msg?: string;
  total_soal: number;
  created_at: string;
  updated_at: string;
};

export async function uploadImportSoal(
  idBankSoal: number,
  file: File,
): Promise<ImportSoalResponse> {
  const data = buildFormData({ file });

  return api<ImportSoalResponse>(
    `/admin/bank-soal/import/${idBankSoal}`,
    {
      method: "POST",
      data,
    },
  );
}

export async function getImportSoalJob(
  idJob: number,
): Promise<ImportSoalJobResponse> {
  return api<ImportSoalJobResponse>(`/admin/bank-soal/import-job/${idJob}`, {
    method: "GET",
  });
}

// =====================
// Hook Wrappers
// =====================

export function useUploadImportSoal() {
  return usePost(
    (payload: { idBankSoal: number; file: File }) =>
      uploadImportSoal(payload.idBankSoal, payload.file),
  );
}
