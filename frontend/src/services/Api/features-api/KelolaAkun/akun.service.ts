import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import type { ResetPasswordRequest } from "@/types/KelolaAkun/ResetPassword";
import { useDelete, usePut } from "@/hooks/fetch";
import { api } from "../../api";

export async function DeletePengguna(id: number) {
    return api<any>(`/admin/pengguna/${id}`, {
        method: "DELETE",
    })
}

export async function DeletePenggunaBulk(ids: number[]) {
  return api<any>("/admin/pengguna", {
    method: "DELETE",
    data: { ids },
  });
}

export function resetPasswordPengguna(idPengguna: number, payload: ResetPasswordRequest) {
  return api<ResetPasswordRequest>(`/admin/pengguna/${idPengguna}/reset-password`, {
    method: "PUT",
    data: buildJsonData(payload),
  });
}

// =====================
// Hook Wrappers
// =====================

export function useDeletePengguna() {
  return useDelete((id: number) => DeletePengguna(id));
}

export function useDeletePenggunaBulk() {
  return useDelete((ids: number[]) => DeletePenggunaBulk(ids));
}

export function useResetPasswordPengguna() {
  return usePut(
    (payload: { id: number; values: ResetPasswordRequest }) =>
      resetPasswordPengguna(payload.id, payload.values),
  );
}
