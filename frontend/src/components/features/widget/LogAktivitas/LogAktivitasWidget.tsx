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
  ChevronRight,
} from "lucide-react";

import { Link } from "react-router";

import type { UserRole,AktivitasLogItem } from "@/types/Log/LogAktivitas";

type AktivitasLogWidgetProps = {
  title?: string;
  items: AktivitasLogItem[];
  className?: string;

  /** tinggi maksimum list (opsional) */
  maxHeightClassName?: string; // default disediakan

  /** opsional: link lihat semua */
  lihatSemuaTo?: string;
};

function roleBadge(role: UserRole) {
  switch (role) {
    case "admin":
      return {
        label: "Admin",
        className:
          "bg-[#397e50] text-white ring-1 ring-inset ring-indigo-500/30",
        Icon: Shield,
      };
    case "guru":
      return {
        label: "Guru",
        className:
          "bg-emerald-50 text-emerald-700 ring-1 ring-inset ring-emerald-200",
        Icon: GraduationCap,
      };
    case "siswa":
      return {
        label: "Siswa",
        className: "bg-sky-50 text-sky-700 ring-1 ring-inset ring-sky-200",
        Icon: User,
      };
    default:
      return {
        label: role,
        className:
          "bg-slate-50 text-slate-700 ring-1 ring-inset ring-slate-200",
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
      iconClassName:
        "bg-emerald-50 text-emerald-700 ring-1 ring-inset ring-emerald-200",
      pillClassName:
        "bg-emerald-50 text-emerald-700 ring-1 ring-inset ring-emerald-200",
    };

  if (a.includes("LOGOUT"))
    return {
      Icon: LogOut,
      label: "Logout",
      iconClassName:
        "bg-amber-50 text-amber-800 ring-1 ring-inset ring-amber-200",
      pillClassName:
        "bg-amber-50 text-amber-800 ring-1 ring-inset ring-amber-200",
    };

  if (a.includes("CREATE") || a.includes("ADD"))
    return {
      Icon: FilePlus2,
      label: "Create",
      iconClassName: "bg-sky-50 text-sky-700 ring-1 ring-inset ring-sky-200",
      pillClassName: "bg-sky-50 text-sky-700 ring-1 ring-inset ring-sky-200",
    };

  if (a.includes("UPDATE") || a.includes("EDIT"))
    return {
      Icon: FilePenLine,
      label: "Update",
      iconClassName:
        "bg-green-100 text-green-700 ring-1 ring-inset ring-indigo-200",
      pillClassName:
        "bg-green-100 text-green-700 ring-1 ring-inset ring-indigo-200",
    };

  if (a.includes("DELETE") || a.includes("REMOVE"))
    return {
      Icon: Trash2,
      label: "Delete",
      iconClassName: "bg-rose-50 text-rose-700 ring-1 ring-inset ring-rose-200",
      pillClassName: "bg-rose-50 text-rose-700 ring-1 ring-inset ring-rose-200",
    };

  if (a.includes("SETTING") || a.includes("CONFIG"))
    return {
      Icon: Settings,
      label: "Setting",
      iconClassName:
        "bg-slate-100 text-slate-700 ring-1 ring-inset ring-slate-200",
      pillClassName:
        "bg-slate-100 text-slate-700 ring-1 ring-inset ring-slate-200",
    };

  return {
    Icon: RefreshCw,
    label: aksi,
    iconClassName:
      "bg-slate-100 text-slate-700 ring-1 ring-inset ring-slate-200",
    pillClassName:
      "bg-slate-100 text-slate-700 ring-1 ring-inset ring-slate-200",
  };
}


export const LogAktivitasWidget = ({
  title = "Log Aktivitas",
  items,
  className,
  maxHeightClassName = "max-h-[60vh] sm:max-h-[520px]",
  lihatSemuaTo,
}: AktivitasLogWidgetProps) => {
  return (
    <>
      <section
        className={[
          "rounded-2xl border border-slate-200 bg-white",
          "shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
          "p-4 sm:p-5",
          className ?? "",
        ].join(" ")}
      >
        <header className="flex items-start justify-between gap-3">
          <div className="flex items-center gap-2">
            <span className="grid h-8 w-8 place-items-center rounded-full bg-slate-50 text-slate-600">
              <Activity className="h-4 w-4" />
            </span>
            <div>
              <h2 className="text-sm font-semibold text-slate-900">{title}</h2>
              <p className="text-xs text-slate-500">{items.length} aktivitas</p>
            </div>
          </div>

          {lihatSemuaTo ? (
            <Link
              to={lihatSemuaTo}
              className={[
                "shrink-0 bg-white",
                "px-3 py-1.5 text-sm font-semibold text-[#397e50]",
                " hover:underline",
              ].join(" ")}
            >
              Lihat Semua
            </Link>
          ) : null}
        </header>

        <div className="mt-4 border-t border-slate-100 pt-4">
          {items.length === 0 ? (
            <div className="rounded-xl border border-dashed border-slate-200 p-4 text-sm text-slate-600">
              Belum ada aktivitas.
            </div>
          ) : (
            <div
              className={["overflow-y-auto pr-1", maxHeightClassName].join(" ")}
            >
              <ul className="space-y-2">
                {items.map((it) => {
                  const rb = roleBadge(it.role);
                  const am = aksiMeta(it.aksi);
                  const AksiIcon = am.Icon;
                  const RoleIcon = rb.Icon;

                  return (
                    <li
                      key={it.id}
                      className={[
                        "rounded-xl border border-slate-200 bg-white",
                        "transition hover:border-slate-300 hover:bg-slate-50",
                        "hover:shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
                      ].join(" ")}
                    >
                      <div className="flex items-start gap-3 p-3 sm:p-4">
                        {/* Icon aksi */}
                        <div
                          className={[
                            "mt-0.5 grid h-9 w-9 shrink-0 place-items-center rounded-lg",
                            am.iconClassName,
                          ].join(" ")}
                        >
                          <AksiIcon className="h-4 w-4" />
                        </div>

                        {/* Konten utama */}
                        <div className="min-w-0 flex-1">
                          {/* Baris atas: username + role badge */}
                          <div className="flex flex-wrap items-center gap-2">
                            <p className="text-sm font-semibold text-slate-900">
                              {it.username}
                            </p>

                            <span
                              className={[
                                "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-semibold",
                                rb.className,
                              ].join(" ")}
                            >
                              <RoleIcon className="h-3.5 w-3.5" />
                              {rb.label}
                            </span>

                            <span className="text-xs text-slate-400">•</span>

                            <span
                              className={[
                                "text-[11px] font-semibold rounded-full px-2 py-0.5",
                                am.pillClassName,
                              ].join(" ")}
                            >
                              {am.label}
                            </span>

                            {it.waktu ? (
                              <>
                                <span className="text-xs text-slate-400">
                                  •
                                </span>
                                <span className="text-xs text-slate-500">
                                  {it.waktu}
                                </span>
                              </>
                            ) : null}
                          </div>

                          {/* Deskripsi */}
                          <p className="mt-1 text-sm text-slate-600">
                            {it.deskripsi}
                          </p>

                          {/* Footer mini (mobile-friendly) */}
                          <div className="mt-2 flex items-center justify-between">
                            <span className="text-[11px] text-slate-400">
                              {it.aksi}
                            </span>
                            <ChevronRight className="h-4 w-4 text-slate-300" />
                          </div>
                        </div>
                      </div>
                    </li>
                  );
                })}
              </ul>
            </div>
          )}
        </div>
      </section>
    </>
  );
};