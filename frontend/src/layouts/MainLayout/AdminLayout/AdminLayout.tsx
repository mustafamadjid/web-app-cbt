// AdminLayout.tsx
import { Outlet } from "react-router";
import { SidebarAdmin } from "@/layouts/MainLayout/Sidebar/SidebarAdmin";

export function AdminLayout() {
  return (
    <div className="min-h-screen flex">
      {/* Sidebar */}
      <aside className="w-64 shrink-0">
        <SidebarAdmin />
      </aside>

      {/* Main content */}
      <main className=" flex-1 min-h-screen bg-gray-100">
        <Outlet />
      </main>
    </div>
  );
}
