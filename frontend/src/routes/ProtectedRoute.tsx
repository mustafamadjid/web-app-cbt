import React from "react";
import { Navigate, Outlet, useLocation } from "react-router";

import Spinner from "@/components/ui/spinner";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import type { Role } from "@/types/Sidebar/SidebarMenu";

type ProtectedRouteProps = {
  allowedRoles?: Role[];
  children?: React.ReactNode;
};

const roleHomeMap: Record<Role, string> = {
  ADMIN: paths.dashboard.home_admin,
  GURU: paths.dashboard.home_guru,
  SISWA: paths.dashboard.home_siswa,
};

const ProtectedRoute = ({ allowedRoles, children }: ProtectedRouteProps) => {
  const { status, user } = useAuth();
  const location = useLocation();

  if (status === "loading") {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (status === "guest") {
    return (
      <Navigate
        to={paths.public.login}
        replace
        state={{ from: location }}
      />
    );
  }

  const role = user?.role as Role | undefined;

  if (allowedRoles && role && !allowedRoles.includes(role)) {
    const fallback = roleHomeMap[role] ?? paths.public.login;
    return <Navigate to={fallback} replace />;
  }

  return <>{children ?? <Outlet />}</>;
};

export default ProtectedRoute;
