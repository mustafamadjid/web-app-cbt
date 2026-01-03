import React from "react";
import { User } from "lucide-react";

type HeaderProps = {
  title?: string;
  userName: string;
  roleLabel: string;
  isOnline?: boolean;
  avatarUrl?: string | null;
  className?: string;
};

export const Header: React.FC<HeaderProps> = ({
  title = "Dashboard",
  userName,
  roleLabel,
  isOnline = true,
  avatarUrl = null,
  className,
}) => {
  return (
    <header className={["px-8 py-4 shrink-0", className ?? ""].join(" ")}>
      <div className="flex items-center justify-between gap-4 bg-white w-full p-4 sm:p-5 rounded-xl border border-slate-200 shadow-[0_10px_24px_rgba(15,23,42,0.06)]">
        <h1 className="font-semibold text-lg sm:text-xl text-slate-900">
          {title}
        </h1>

        {/* Profil */}
        <div className="flex items-center gap-3 min-w-0">
          <div className="min-w-0 text-right">
            <h2 className="font-semibold text-sm sm:text-base text-slate-900 truncate">
              {userName}
            </h2>

            <span className="inline-flex items-center rounded-2xl bg-[#fff5d5] px-2.5 py-1 text-xs font-semibold text-[#724b00] ">
              {roleLabel}
            </span>
          </div>

          <div className="relative shrink-0">
            <div className="bg-slate-100 text-slate-700 w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center overflow-hidden ring-1 ring-slate-200">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={`Foto profil ${userName}`}
                  className="h-full w-full object-cover"
                />
              ) : (
                <User
                  className="h-5 w-5 sm:h-6 sm:w-6"
                  aria-label="Profil pengguna"
                />
              )}
            </div>

            <span
              className={[
                "absolute bottom-0 right-0 h-3 w-3 rounded-full ring-2 ring-white",
                isOnline ? "bg-green-600" : "bg-slate-300",
              ].join(" ")}
              aria-label={isOnline ? "Online" : "Offline"}
              title={isOnline ? "Online" : "Offline"}
            />
          </div>
        </div>
      </div>
    </header>
  );
};
