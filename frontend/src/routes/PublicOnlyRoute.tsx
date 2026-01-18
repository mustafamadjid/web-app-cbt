import React from "react";
import { Navigate, Outlet } from "react-router";

import Spinner from "@/components/ui/spinner";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import type { Role } from "@/types/Sidebar/SidebarMenu";

type PublicOnlyRouteProps = {
  children?: React.ReactNode;
};

const roleHomeMap: Partial<Record<Role, string>> = {
  ADMIN: paths.dashboard.home_admin,
  GURU: paths.dashboard.home_guru,
  SISWA: paths.dashboard.home_siswa,
};

const PublicOnlyRoute = ({ children }: PublicOnlyRouteProps) => {
  const { status, user } = useAuth();

  if (status === "loading") {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (status === "authenticated") {
    const role = user?.role as Role | undefined;

   
    const fallback =
      (role && roleHomeMap[role]) ?? paths.dashboard.home_admin ?? "/";

    return <Navigate to={fallback} replace />;
  }

 
  return <>{children ?? <Outlet />}</>;
};

export default PublicOnlyRoute;
