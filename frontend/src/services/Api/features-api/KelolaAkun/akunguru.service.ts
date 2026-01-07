// // services/users.service.ts
// import { api } from "./api";

// export type User = {
//   id: string;
//   name: string;
//   email: string;
// };

// // untuk update (biasanya partial)
// export type UpdateUserPayload = {
//   name?: string;
//   email?: string;
// };

// export async function getUserById(userId: string, token?: string | null) {
//   return api<User>(`/users/${encodeURIComponent(userId)}`, {
//     method: "GET",
//     token,
//   });
// }

// export async function updateUser(
//   userId: string,
//   payload: UpdateUserPayload,
//   token?: string | null
// ) {
//   return api<User>(`/users/${encodeURIComponent(userId)}`, {
//     method: "PATCH", // atau "PUT" tergantung backend
//     data: payload,
//     token,
//   });
// }

// export async function listUsers(
//   params: { q?: string; page?: number; limit?: number },
//   token?: string | null
// ) {
//   return api<User[]>("/users", {
//     method: "GET",
//     params,
//     token,
//   });
// }

// export async function deleteUser(userId: string, token?: string | null) {
//   return api<void>(`/users/${encodeURIComponent(userId)}`, {
//     method: "DELETE",
//     token,
//   });
// }

import { api } from "../../api";
import { buildFormData } from "@/helper/FormData/BuildFormData";

import type { TeacherRegisterFormValues,TeacherRegisterResponse } from "@/types/KelolaAkun/AkunGuru";
import type { ApiEnvelope } from "../../api";


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

  // Penting: jangan set Content-Type manual untuk FormData
  const res = await api<ApiEnvelope<TeacherRegisterResponse>>(
    "/teachers/register",
    {
      method: "POST",
      data: formData,
      
    }
  );

  return res.data;
}