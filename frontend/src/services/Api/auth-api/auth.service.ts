import { api } from "@/services/Api/api";
import type { LoginPayload, User } from "@/types/Auth/Auth";

export async function login(payload: LoginPayload) {
  await api<null>("/auth/login", { method: "POST", data: payload });
}

export async function logout() {
  await api<null>("/auth/logout", { method: "POST" });
}

export async function getAuthMe() {
  return api<User>("/auth/me", { method: "GET" });
}
