// src/routes/ProtectedRoute.tsx
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

  // Jika belum autentik, redirect ke login (simpan lokasi asal untuk redirect setelah login)
  if (status === "guest") {
    return (
      <Navigate
        to={paths.public.login}
        replace
        state={{
          from: { pathname: location.pathname, search: location.search },
        }}
      />
    );
  }

  const role = user?.role as Role | undefined;

  // Jika ada allowedRoles, pastikan user punya role dan role termasuk di allowedRoles.
  if (allowedRoles) {
    if (!role) {
      // user tidak punya role => treat as unauthorized
      return (
        <Navigate
          to={paths.public.login}
          replace
          state={{ from: { pathname: location.pathname } }}
        />
      );
    }

    if (!allowedRoles.includes(role)) {
      const fallback = roleHomeMap[role] ?? paths.public.login;
      return <Navigate to={fallback} replace />;
    }
  }

  // Render children jika ada, kalau tidak render nested routes via Outlet
  return <>{children ?? <Outlet />}</>;
};

export default ProtectedRoute;
