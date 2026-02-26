import { api } from "../../api";

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
