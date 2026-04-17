import { useFetch } from "@/hooks/fetch";
import { api } from "@/services/Api/api";
import type { Role } from "@/types/Sidebar/SidebarMenu";
import type { DashboardStatistik } from "@/types/Dashboard/DashboardStatistik";

type DashboardRole = Extract<Role, "ADMIN" | "GURU">;

const DASHBOARD_ENDPOINT_BY_ROLE: Record<DashboardRole, string> = {
  ADMIN: "/admin/dashboard",
  GURU: "/guru/dashboard",
};

function resolveDashboardEndpoint(role: DashboardRole): string {
  return DASHBOARD_ENDPOINT_BY_ROLE[role];
}

export async function getDashboardStatistik(
  role: DashboardRole,
): Promise<DashboardStatistik> {
  return api<DashboardStatistik>(resolveDashboardEndpoint(role), {
    method: "GET",
  });
}

export function useGetDashboardStatistik(
  role: DashboardRole,
  enabled = true,
) {
  return useFetch(
    () =>
      enabled
        ? getDashboardStatistik(role)
        : Promise.resolve(null as DashboardStatistik | null),
    [role, enabled],
  );
}
