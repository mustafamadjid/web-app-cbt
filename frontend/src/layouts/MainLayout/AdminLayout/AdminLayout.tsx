// AdminLayout.tsx
import { Outlet } from "react-router";
import { SidebarAdmin } from "@/layouts/MainLayout/Sidebar/SidebarAdmin";

export function AdminLayout() {
  return (
    <div className="min-h-screen bg-[#ecf1ed] ">
      {/* Sidebar is fixed inside component; do NOT reserve width on mobile */}
      <SidebarAdmin />

      <main className="min-h-screen min-w-0 sm:ml-64">
        <Outlet />
      </main>
    </div>
  );
}
