import React from "react";
import { GraduationCap, Users, School } from "lucide-react";

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
        "relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      {/* Top Accent Line */}
      <div className="h-1.5 w-full bg-gradient-to-r from-[#397e50] to-[#37513d]" />

      {/* Header */}
      <header className="flex items-center gap-3 px-5 pt-5 pb-2">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
          <School className="h-5 w-5" />
        </div>
        <div>
          <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
          <p className="text-xs font-medium text-gray-500">
            Ringkasan data akademik
          </p>
        </div>
      </header>

      {/* Content Grid */}
      <div className="grid gap-4 p-5 sm:grid-cols-2">
        {/* Card Siswa */}
        <div className="group relative overflow-hidden rounded-xl border border-gray-100 bg-white p-4 shadow-sm transition-all hover:-translate-y-1 hover:border-[#397e50]/30 hover:shadow-md">
          {/* Side Accent */}
          <div className="absolute left-0 top-0 bottom-0 w-1 bg-[#397e50]" />

          <div className="relative z-10 flex items-center gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-[#397e50]/10 text-[#397e50] group-hover:bg-[#397e50] group-hover:text-white transition-colors">
              <GraduationCap className="h-6 w-6" />
            </div>
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Siswa
              </p>
              <p className="mt-0.5 text-2xl font-black text-[#37513d]">
                {totalSiswa.toLocaleString("id-ID")}
              </p>
            </div>
          </div>

          {/* Watermark Icon */}
          <GraduationCap
            className="absolute -bottom-2 -right-2 h-20 w-20 text-[#397e50]/5 transition-transform group-hover:scale-110 group-hover:text-[#397e50]/10"
            strokeWidth={1}
          />
        </div>

        {/* Card Guru */}
        <div className="group relative overflow-hidden rounded-xl border border-gray-100 bg-white p-4 shadow-sm transition-all hover:-translate-y-1 hover:border-sky-500/30 hover:shadow-md">
          {/* Side Accent */}
          <div className="absolute left-0 top-0 bottom-0 w-1 bg-sky-600" />

          <div className="relative z-10 flex items-center gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-sky-50 text-sky-600 group-hover:bg-sky-600 group-hover:text-white transition-colors">
              <Users className="h-6 w-6" />
            </div>
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Guru
              </p>
              <p className="mt-0.5 text-2xl font-black text-[#37513d]">
                {totalGuru.toLocaleString("id-ID")}
              </p>
            </div>
          </div>

          {/* Watermark Icon */}
          <Users
            className="absolute -bottom-2 -right-2 h-20 w-20 text-sky-600/5 transition-transform group-hover:scale-110 group-hover:text-sky-600/10"
            strokeWidth={1}
          />
        </div>
      </div>
    </section>
  );
};
