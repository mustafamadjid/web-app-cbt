import React from "react";
import { User,CalendarRange } from "lucide-react";
import type { Role } from "@/types/Sidebar/SidebarMenu";
import { paths } from "@/routes/paths";
import { useNavigate } from "react-router";

type HeaderProps = {
  title?: string;
  userName: string;
  roleLabel: Role;
  isOnline?: boolean;
  avatarUrl?: string | null;
  schoolName?: string;
  academicYear?: string;
  semester?: string;
  className?: string;
  onMenuClick?: () => void;
  onSettingsClick?: () => void;
};

const Header: React.FC<HeaderProps> = ({
  title = "Dashboard",
  userName,
  roleLabel,
  isOnline = true,
  avatarUrl = null,
  academicYear = "2025/2026",
  semester = "Ganjil",
  className,

}) => {
  const today = new Date().toLocaleDateString("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  const roleLabelMap = (role : Role) => {
    if(role === "ADMIN") return "Administrator";
    if(role === "GURU") return "Guru";
    if(role === "SISWA") return "Siswa";
    return "";
  }

  const navigate = useNavigate();

  const pathNavigationByRole = (role : Role) => {
    if(role === "ADMIN") return paths.dashboard.profil_admin;
    if(role === "GURU") return paths.dashboard.profil_guru;
    if(role === "SISWA") return paths.dashboard.profil_siswa;
    return "";
  }


  return (
    <header
      className={[
        " top-0 z-40 w-full flex-none transition-all duration-300",
        className ?? "",
      ].join(" ")}
    >
      <div className="h-1 w-full" />

      <div className="flex items-center justify-between border-b border-gray-200 bg-white/90 px-4 py-3 shadow-sm backdrop-blur-md sm:px-6 lg:px-8">
        {/* --- LEFT: Mobile Menu & Title --- */}
        <div className="flex items-center gap-4">
        

          <div className="flex flex-col">
            <h1 className="text-xl font-bold tracking-tight text-[#37513d] sm:text-2xl">
              {title}
            </h1>
            <p className="hidden text-xs font-medium text-gray-500 sm:block">
              {today}
            </p>
          </div>
        </div>

        {/* --- CENTER: REDESIGNED ACADEMIC BADGE --- */}
        <div className="hidden md:flex flex-1 justify-center px-4">
          {/* Container Utama */}
          <div className="group flex items-center overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm transition-all hover:border-[#397e50]/30 hover:shadow-md">
            {/* Bagian Tahun Ajaran (Kiri) */}
            <div className="hidden md:flex flex-1 justify-center px-4">
              <div className="flex items-center gap-3 rounded-lg px-4 py-2 transition-colors hover:bg-gray-100">
                <div className="text-right">
                  <p className="text-xs font-bold text-gray-900">
                    {academicYear}
                  </p>
                  <p className="text-2xs font-medium text-gray-500">
                    Semester {semester}
                  </p>
                </div>
                <div className="h-8 w-px] bg-gray-300"></div>
                <div className="rounded-full bg-white p-1.5 shadow-sm text-[#397e50]">
                  <CalendarRange className="h-5 w-5" />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* --- RIGHT: Profile --- */}
        <div className="flex items-center gap-3 sm:gap-6">
          <div className="hidden h-8 w-px bg-gray-200 sm:block" />

          <div className="flex items-center gap-3">
            <div className="hidden text-right lg:block">
              <div className="text-sm font-bold text-gray-800">{userName}</div>
              <div className="flex justify-end">
                <span className="inline-flex items-center rounded-md bg-[#397e50] px-2 py-0.5 text-2xs font-bold uppercase tracking-wider text-white shadow-sm">
                  {roleLabelMap(roleLabel)}
                </span>
              </div>
            </div>

            <div className="group relative cursor-pointer">
              <div 
              onClick={()=> navigate(pathNavigationByRole(roleLabel))}
              className="relative h-10 w-10 overflow-hidden rounded-full ring-2 ring-gray-100 transition-all group-hover:ring-[#397e50]/30 sm:h-11 sm:w-11">
                {avatarUrl ? (
                  <img
                    src={avatarUrl}
                    alt={userName}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center bg-slate-100 text-slate-500">
                    <User className="h-5 w-5" />
                  </div>
                )}
              </div>
              <span
                className={[
                  "absolute bottom-0.5 right-0.5 h-3 w-3 rounded-full ring-2 ring-white transition-colors",
                  isOnline ? "bg-emerald-500" : "bg-gray-400",
                ].join(" ")}
              />
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
