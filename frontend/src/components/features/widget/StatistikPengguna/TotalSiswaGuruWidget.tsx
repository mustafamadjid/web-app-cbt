import React from "react";
import { GraduationCap, Users } from "lucide-react";

import type { TotalSiswaGuru } from "@/types/Widget/TotalSiswaGuru";

type TotalSiswaGuruWidgetProps = TotalSiswaGuru & {
  title?: string;
  className?: string;
};

export const TotalSiswaGuruWidget: React.FC<TotalSiswaGuruWidgetProps> = ({
  title = "Statistik Pengguna",
  totalSiswa,
  totalGuru,
  className,
}) => {
  return (
    <section
      className={[
        "rounded-xl border border-gray-200 bg-white p-5 shadow-sm",
        "transition-all duration-300 hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      <header className="mb-4">
        <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
        <p className="text-xs font-medium text-gray-500">
          Ringkasan total pengguna aktif
        </p>
      </header>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="flex items-center gap-4 rounded-xl border border-emerald-100 bg-emerald-50/60 p-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-700">
            <GraduationCap className="h-6 w-6" />
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-emerald-700">
              Total Siswa
            </p>
            <p className="text-2xl font-bold text-emerald-900">
              {totalSiswa.toLocaleString("id-ID")}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-4 rounded-xl border border-sky-100 bg-sky-50/60 p-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-sky-100 text-sky-700">
            <Users className="h-6 w-6" />
          </div>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-sky-700">
              Total Guru
            </p>
            <p className="text-2xl font-bold text-sky-900">
              {totalGuru.toLocaleString("id-ID")}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
