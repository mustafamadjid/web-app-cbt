import {
  Activity,
  Shield,
  GraduationCap,
  User,
  LogIn,
  LogOut,
  FilePlus2,
  FilePenLine,
  Trash2,
  Settings,
  RefreshCw,

  Clock,
} from "lucide-react";



import type { UserRole, AktivitasLogItem } from "@/types/Log/LogAktivitas";

type AktivitasLogWidgetProps = {
  title?: string;
  items: AktivitasLogItem[];
  className?: string;
  maxHeightClassName?: string;
  lihatSemuaTo?: string;
};

// --- Helpers ---

function roleBadge(role: UserRole) {
  switch (role) {
    case "ADMIN":
      return {
        label: "Admin",
        className: "bg-[#397e50] text-white",
        Icon: Shield,
      };
    case "GURU":
      return {
        label: "Guru",
        className: "bg-emerald-100 text-emerald-800",
        Icon: GraduationCap,
      };
    case "SISWA":
      return {
        label: "Siswa",
        className: "bg-sky-100 text-sky-800",
        Icon: User,
      };
    default:
      return {
        label: role,
        className: "bg-gray-100 text-gray-700",
        Icon: User,
      };
  }
}

function aksiMeta(aksi: string) {
  const a = aksi.trim().toUpperCase();

  if (a.includes("LOGIN"))
    return {
      Icon: LogIn,
      label: "Login",
      colorClass: "text-emerald-600 bg-emerald-50 border-emerald-100",
    };

  if (a.includes("LOGOUT"))
    return {
      Icon: LogOut,
      label: "Logout",
      colorClass: "text-amber-600 bg-amber-50 border-amber-100",
    };

  if (a.includes("CREATE") || a.includes("ADD"))
    return {
      Icon: FilePlus2,
      label: "Create",
      colorClass: "text-blue-600 bg-blue-50 border-blue-100",
    };

  if (a.includes("UPDATE") || a.includes("EDIT"))
    return {
      Icon: FilePenLine,
      label: "Update",
      colorClass: "text-indigo-600 bg-indigo-50 border-indigo-100",
    };

  if (a.includes("DELETE") || a.includes("REMOVE"))
    return {
      Icon: Trash2,
      label: "Delete",
      colorClass: "text-rose-600 bg-rose-50 border-rose-100",
    };

  if (a.includes("SETTING") || a.includes("CONFIG"))
    return {
      Icon: Settings,
      label: "Setting",
      colorClass: "text-slate-600 bg-slate-50 border-slate-100",
    };

  return {
    Icon: RefreshCw,
    label: aksi,
    colorClass: "text-gray-600 bg-gray-50 border-gray-100",
  };
}

const LogAktivitasWidget = ({
  title = "Log Aktivitas",
  items,
  className,
  maxHeightClassName = "max-h-[60vh] sm:max-h-[520px]",
}: AktivitasLogWidgetProps) => {
  return (
    <section
      className={[
        "relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      {/* Top Accent Line */}
      {/* <div className="h-1.5 w-full bg-linear-to-r from-[#397e50] to-[#37513d]" /> */}

      {/* Header */}
      <header className="flex items-center justify-between px-5 pt-5 pb-2">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
            <Activity className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
            <p className="text-xs font-medium text-gray-500">
              {items.length} aktivitas terakhir
            </p>
          </div>
        </div>

        
      </header>

      {/* Content */}
      <div className="mt-2 flex-1 border-t border-gray-100 bg-gray-50/30 p-5">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-gray-300 bg-white py-10 text-center">
            <Activity className="h-8 w-8 text-gray-300" />
            <p className="text-sm text-gray-500">
              Belum ada aktivitas tercatat.
            </p>
          </div>
        ) : (
          <div
            className={[
              "overflow-y-auto pr-2 relative",
              maxHeightClassName,
            ].join(" ")}
          >
            {/* Vertical Line for Timeline effect */}
            <div className="absolute left-6 top-4 bottom-4 w-px bg-gray-200" />

            <div className="space-y-6">
              {items.map((it) => {
                const rb = roleBadge(it.role);
                const am = aksiMeta(it.action);
                const AksiIcon = am.Icon;
                const RoleIcon = rb.Icon;

                return (
                  <div key={it.id_aktivitas} className="group relative flex gap-4">
                    {/* Timeline Node (Icon Aksi) */}
                    <div
                      className={[
                        "relative z-10 flex h-12 w-12 shrink-0 items-center justify-center rounded-full border-2 bg-white transition-all duration-300",
                        am.colorClass
                          .replace("bg-", "border-")
                          .replace("text-", "text-"), // logic kasar utk border color
                        "group-hover:scale-110 shadow-sm",
                      ].join(" ")}
                    >
                      <div
                        className={[
                          "flex h-8 w-8 items-center justify-center rounded-full",
                          am.colorClass,
                        ].join(" ")}
                      >
                        <AksiIcon className="h-4 w-4" />
                      </div>
                    </div>

                    {/* Content Card */}
                    <div className="flex-1 min-w-0 rounded-xl border border-gray-100 bg-white p-3.5 shadow-sm transition-all hover:border-[#397e50]/30 hover:shadow-md">
                      {/* Header: User & Role */}
                      <div className="mb-2 flex flex-wrap items-center justify-between gap-2 border-b border-gray-50 pb-2">
                        <div className="flex items-center gap-2">
                          <span className="text-sm font-bold text-gray-800">
                            {it.username}
                          </span>
                          <span
                            className={[
                              "inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-2xs font-bold uppercase tracking-wide",
                              rb.className,
                            ].join(" ")}
                          >
                            <RoleIcon className="h-3 w-3" />
                            {rb.label}
                          </span>
                        </div>

                        {/* Waktu */}
                        <div className="flex items-center gap-1 text-2xs font-medium text-gray-400">
                          <Clock className="h-3 w-3" />
                          {it.created_at}
                        </div>
                      </div>

                      {/* Body: Deskripsi */}
                      <div className="mb-2">
                        <p className="text-xs leading-relaxed text-gray-600">
                          <span className="font-semibold text-gray-800 mr-1">
                            {am.label}:
                          </span>
                          {it.description}
                        </p>
                      </div>

                      {/* Footer: Meta Aksi (Raw) */}
                      <div className="flex flex-wrap items-center gap-2 text-2xs text-gray-400">
                        <span className="font-mono bg-gray-50 px-1 py-0.5 rounded text-gray-500">
                          {it.action}
                        </span>
                        <span className="font-mono bg-gray-50 px-1 py-0.5 rounded text-gray-500">
                          {it.ip_address}
                        </span>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </section>
  );
};

export default LogAktivitasWidget;
