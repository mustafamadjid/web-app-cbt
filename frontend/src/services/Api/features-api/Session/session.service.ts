import { buildJsonData } from "@/helper/FormData/BuildJsonData";
import { useFetch, usePut } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type {
  ActiveSessionRow,
  AdminRevokeSessionPayload,
  ListActiveSessionsResponse,
} from "@/types/Session/Session";

const ACTIVE_LOGIN_SESSION_ENDPOINT = "/admin/sesi/login-aktif";
const ADMIN_REVOKE_SESSION_ENDPOINT = "/admin/auth/revoke-session";

export async function getActiveLoginSessions(): Promise<ActiveSessionRow[]> {
  const response = await api<ListActiveSessionsResponse>(ACTIVE_LOGIN_SESSION_ENDPOINT, {
    method: "GET",
  });

  return response.items;
}

export async function adminRevokeSession(payload: AdminRevokeSessionPayload) {
  return await api<null>(ADMIN_REVOKE_SESSION_ENDPOINT, {
    method: "PUT",
    data: buildJsonData(payload),
  });
}

export function useGetActiveLoginSessions() {
  return useFetch(() => getActiveLoginSessions(), []);
}

export function useAdminRevokeSession() {
  return usePut((payload: AdminRevokeSessionPayload) => adminRevokeSession(payload));
}
