import { api } from "../../api";
import { usePost } from "@/hooks/fetch";

export type ImportSoalResponse = {
  id_job: number;
};


export async function uploadImportSoal(
  idBankSoal: number,
  file: File,
): Promise<ImportSoalResponse> {
  const formData = new FormData();
  formData.append("file", file);

  return api<ImportSoalResponse>(
    `/admin/bank-soal/import/${idBankSoal}`,
    {
      method: "POST",
      data: formData,
      headers: { "Content-Type": "multipart/form-data" },
    },
  );
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
