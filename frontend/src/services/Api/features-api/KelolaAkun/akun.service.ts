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