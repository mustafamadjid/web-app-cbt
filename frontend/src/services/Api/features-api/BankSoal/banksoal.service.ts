import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useDelete, useFetch, usePost, usePut } from "@/hooks/fetch";
import type {
  BankSoalItem,
  CreateBankSoalPayload,
  GetBankSoalParams,
  UpdateBankSoalPayload,
} from "@/types/BankSoal/BankSoal";
import { api } from "../../api";

const BANK_SOAL_ENDPOINT = "/admin-guru/bank-soal";
const BANK_SOAL_UPLOADED_ENDPOINT = "/admin-guru/bank-soal-uploaded";

const buildBankSoalQueryParams = (
  params: GetBankSoalParams = {},
): Record<string, string | undefined> => ({
  search: params.search?.trim() || undefined,
  id_kelas: params.id_kelas != null ? String(params.id_kelas) : undefined,
  id_mapel: params.id_mapel != null ? String(params.id_mapel) : undefined,
  limit: params.limit != null ? String(params.limit) : undefined,
  offset: params.offset != null ? String(params.offset) : undefined,
});

export async function getBankSoal(
  params: GetBankSoalParams = {},
): Promise<BankSoalItem[]> {
  const queryParams = buildBankSoalQueryParams(params);

  return api<BankSoalItem[]>(BANK_SOAL_ENDPOINT, {
    method: "GET",
    params: queryParams,
  });
}

export async function getBankSoalUploaded(
  params: GetBankSoalParams = {},
): Promise<BankSoalItem[]> {
  const queryParams = buildBankSoalQueryParams(params);

  return api<BankSoalItem[]>(BANK_SOAL_UPLOADED_ENDPOINT, {
    method: "GET",
    params: queryParams,
  });
}

export async function getBankSoalById(idBankSoal: number): Promise<BankSoalItem> {
  return api<BankSoalItem>(`${BANK_SOAL_ENDPOINT}/${idBankSoal}`, {
    method: "GET",
  });
}

export async function getBankSoalByGuru(
  idPengguna: number,
): Promise<BankSoalItem[]> {
  if (!idPengguna || idPengguna <= 0) {
    return [];
  }

  return api<BankSoalItem[]>(`/admin-guru/guru/bank-soal/${idPengguna}`, {
    method: "GET",
  });
}

export async function createBankSoal(
  values: CreateBankSoalPayload,
): Promise<boolean> {
  const data = buildJsonData(values);
  return api<boolean>(BANK_SOAL_ENDPOINT, {
    method: "POST",
    data,
  });
}

export async function updateBankSoalPartial(
  idBankSoal: number,
  values: UpdateBankSoalPayload,
): Promise<boolean> {
  const data = buildJsonData(values);
  return api<boolean>(`${BANK_SOAL_ENDPOINT}/${idBankSoal}`, {
    method: "PATCH",
    data,
  });
}

export async function deleteBankSoal(idBankSoal: number): Promise<boolean> {
  return api<boolean>(`${BANK_SOAL_ENDPOINT}/${idBankSoal}`, {
    method: "DELETE",
  });
}

// =====================
// Hook Wrappers
// =====================

export function useGetBankSoal(params: GetBankSoalParams = {}, enabled = true) {
  return useFetch(
    () => (enabled ? getBankSoal(params) : Promise.resolve([])),
    [params.search, params.id_kelas, params.id_mapel, params.limit, params.offset, enabled],
  );
}

export function useGetBankSoalUploaded(
  params: GetBankSoalParams = {},
  enabled = true,
) {
  return useFetch(
    () => (enabled ? getBankSoalUploaded(params) : Promise.resolve([])),
    [params.search, params.id_kelas, params.id_mapel, params.limit, params.offset, enabled],
  );
}

export function useGetBankSoalById(idBankSoal: number, enabled = true) {
  return useFetch(
    () =>
      enabled
        ? getBankSoalById(idBankSoal)
        : Promise.resolve(null as BankSoalItem | null),
    [idBankSoal, enabled],
  );
}

export function useGetBankSoalByGuru(idPengguna: number) {
  return useFetch(() => getBankSoalByGuru(idPengguna), [idPengguna]);
}

export function useCreateBankSoal() {
  return usePost((values: CreateBankSoalPayload) => createBankSoal(values));
}

export function useUpdateBankSoal() {
  return usePut(
    (payload: { id: number; values: UpdateBankSoalPayload }) =>
      updateBankSoalPartial(payload.id, payload.values),
  );
}

export function useDeleteBankSoal() {
  return useDelete((idBankSoal: number) => deleteBankSoal(idBankSoal));
}
