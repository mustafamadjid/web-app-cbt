// AdminLayout.tsx
import { useEffect, useState } from "react";
import { Outlet } from "react-router";
import { SidebarAdmin } from "@/layouts/MainLayout/Sidebar/SidebarAdmin";

const STORAGE_KEY = "admin_sidebar_open";

export function AdminLayout() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  // restore preference
  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw === null) return;
      setIsSidebarOpen(raw === "1");
    } catch {
      // ignore
    }
  }, []);

  // persist preference
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, isSidebarOpen ? "1" : "0");
    } catch {
      // ignore
    }
  }, [isSidebarOpen]);

  return (
    <div className="min-h-screen bg-[#ecf1ed]">
      <SidebarAdmin
        isOpen={isSidebarOpen}
        onToggle={() => setIsSidebarOpen((v) => !v)}
        onClose={() => setIsSidebarOpen(false)}
        onOpen={() => setIsSidebarOpen(true)}
      />

      
      <main
        className={[
          "min-h-screen min-w-0 transition-[margin] duration-200 ease-out",
          isSidebarOpen ? "sm:ml-64" : "sm:ml-0",
        ].join(" ")}
      >
        <Outlet />
      </main>
    </div>
  );
}
