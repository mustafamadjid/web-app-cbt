import { api } from "../../api";
import { buildFormData } from "@/helper/FormData/BuildFormData";

import type {
  DataGuru,
  DataUpdateGuru,
  TeacherRegisterFormValues,
  TeacherRegisterResponse,
  TeacherUpdatePayload,
} from "@/types/KelolaAkun/AkunGuru";
import type { StatusAkun } from "@/types/OpsiTypes/Option";
import type { ApiEnvelope } from "../../api";

export type GuruFilterParams = {
  q?: string;
  status?: StatusAkun;
};
// // Data Dummy
// export const daftarPengguna: DataGuru[] = [
//   {
//     id: 1,
//     namaLengkap: "Neil Sims",
//     email: "neil.sims@flowbite.com",
//     username: "neilsims",
//     noHp: "081234567890",
//     jenisKelamin: "LAKI_LAKI",
//     statusAkun: "aktif",
//     nip: "198701012010121001",
//     jabatan: "Guru Matematika",
//     bidangStudi: "Sains & Teknologi",
//     urlGambarProfil: "https://i.pravatar.cc/150?u=neil",
//     role: "ADMIN",
//   },
//   {
//     id: 2,
//     namaLengkap: "Bonnie Green",
//     email: "bonnie@flowbite.com",
//     username: "bonnieg",
//     noHp: "082233445566",
//     jenisKelamin: "PEREMPUAN",
//     statusAkun: "nonaktif",
//     nip: "199002022011112002",
//     jabatan: "Guru Bahasa Inggris",
//     bidangStudi: "Bahasa",
//     urlGambarProfil: "https://i.pravatar.cc/150?u=bonnie",
//     role: "GURU",
//   },
//   {
//     id: 3,
//     namaLengkap: "Jese Leos",
//     email: "jese@flowbite.com",
//     username: "jeseleos",
//     noHp: "085677889900",
//     jenisKelamin: "LAKI_LAKI",
//     statusAkun: "dibekukan",
//     nip: "199505052015051005",
//     jabatan: "Guru Olahraga",
//     bidangStudi: "PJOK",
//     urlGambarProfil: "https://i.pravatar.cc/150?u=jese",
//     role: "GURU",
//   },
//   {
//     id: 4,
//     namaLengkap: "Jese Leos",
//     email: "jese@flowbite.com",
//     username: "jeseleos",
//     noHp: "085677889900",
//     jenisKelamin: "LAKI_LAKI",
//     statusAkun: "dibekukan",
//     nip: "199505052015051005",
//     jabatan: "Guru Olahraga",
//     bidangStudi: "PJOK",
//     urlGambarProfil: "https://i.pravatar.cc/150?u=jese",
//     role: "ADMIN",
//   },
// ];


export async function submitTeacherRegister(values: TeacherRegisterFormValues) {
  const formData = buildFormData(values, {
    transform: (key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") {
        if (key === "email") return value.trim().toLowerCase();
        if (key === "password") return value;
        return value.trim();
      }
      return value as any;
    },
    
  });

  return api<TeacherRegisterFormValues>("/admin/guru",{
    method: "POST",
    data: formData
  })
}

export async function getGuruById(id: number): Promise<DataUpdateGuru | null> {
  return api<DataUpdateGuru | null>(`/admin/guru/${id}`,{
    method: "GET"
  });
}

export async function updateGuru(
  id: number,
  payload: TeacherUpdatePayload
) {
  const formData = buildFormData(payload, {
    transform: (key, value) => {
      if (value instanceof Blob) return value;
      if (typeof value === "string") {
        if (key === "email") return value.trim().toLowerCase();
        return value.trim();
      }
      return value as any;
    },
    skipNullish: true,
  });

  const res = await api<ApiEnvelope<TeacherRegisterResponse>>(
    `/admin/guru/${id}`,
    {
      method: "PATCH",
      data: formData,
    }
  );

  return res.data;
}

// /** === MOCK "API" (simulasikan network delay) === */
// const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

// type DataGuruIndexed = DataGuru & {
//   __searchText: string;
// };

// const normalize = (value: string) => value.toLowerCase().trim();

// const buildSearchableText = (guru: DataGuru) =>
//   normalize(
//     [guru.nama_lengkap, guru.email, guru.nip, guru.username]
//       .filter(Boolean)
//       .join(" ")
//   );

// const indexDummy = (data: DataGuru[]): DataGuruIndexed[] =>
//   data.map((item) => ({
//     ...item,
//     __searchText: buildSearchableText(item),
//   }));

// const filterBySearch = (data: DataGuruIndexed[], q?: string) => {
//   const query = q ? normalize(q) : "";
//   if (!query) return data;
//   return data.filter((guru) => guru.__searchText.includes(query));
// };

// const applyGuruFilters = (
//   data: DataGuruIndexed[],
//   params: GuruFilterParams = {}
// ) => {
//   let out = data;
//   out = filterBySearch(out, params.q);
//   return out;
// };

// const stripInternal = (items: DataGuruIndexed[]): DataGuru[] =>
//   items.map(({ __searchText, ...rest }) => rest);

// const USE_DUMMY = true;

export async function GetAllGuru(
  params: GuruFilterParams = {}
): Promise<DataGuru[]> {
  // if (USE_DUMMY) {
  //   await sleep(250);
  //   const filtered = applyGuruFilters(indexDummy(daftarPengguna), params);
  //   return stripInternal(filtered);
  // }

  const queryParams: Record<string, string | undefined> = {
    q: params.q || undefined,
    status: params.status || undefined,
  };

  
    const result = await api<DataGuru[]>("/admin/guru", {
      params: queryParams,
    });

 
    return result;
}
