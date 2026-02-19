import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import type { ResetPasswordRequest } from "@/types/KelolaAkun/ResetPassword";
// export async function GetAllGuru(
//   params: GuruFilterParams = {}
// ): Promise<DataGuru[]> {
//   // if (USE_DUMMY) {
//   //   await sleep(250);
//   //   const filtered = applyGuruFilters(indexDummy(daftarPengguna), params);
//   //   return stripInternal(filtered);
//   // }

import { api } from "../../api";

//   const queryParams: Record<string, string | undefined> = {
//     q: params.q || undefined,
//   };

  
//     const result = await api<DataGuru[]>("/admin/guru", {
//       params: queryParams,
//     });

 
//     return result;
// }

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
