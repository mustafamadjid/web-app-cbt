import React from "react";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { Link } from "react-router";
import {
  CalendarCheck2,
  ChevronRight,
  Clock,
  MapPin,
  User,
  CalendarDays,
} from "lucide-react";

type JadwalUjianProps = {
  title?: string;
  items: JadwalUjianItem[];
  className?: string;
  /** tinggi maksimum list (opsional) */
  maxHeightClassName?: string;
  /** tujuan halaman list ujian keseluruhan */
  lihatSemuaTo?: string;
};

const JadwalUjianWidget: React.FC<JadwalUjianProps> = ({
  title = "Jadwal Ujian",
  items,
  className,
  maxHeightClassName = "max-h-[60vh] sm:max-h-[520px]",
  lihatSemuaTo = "/jadwal-ujian",
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
      {/* Top Accent Line
      <div className="h-1.5 w-full bg-linear-to-r from-[#397e50] to-[#37513d]" /> */}

      {/* Header Section */}
      <header className="flex items-center justify-between px-5 pt-5 pb-3">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
            <CalendarCheck2 className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
            <p className="text-xs font-medium text-gray-500">
              {items.length} ujian mendatang
            </p>
          </div>
        </div>

        <Link
          to={lihatSemuaTo}
          className="group flex items-center gap-1 text-xs font-bold text-[#397e50] transition-colors hover:text-[#2f5c3f]"
        >
          Lihat Semua
          <ChevronRight className="h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" />
        </Link>
      </header>

      {/* Content Section */}
      <div className="flex-1 border-t border-gray-100 bg-gray-50/50 p-4 overflow-y-scroll">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-gray-300 bg-white py-12 text-center">
            <CalendarDays className="h-10 w-10 text-gray-300" />
            <p className="text-sm font-medium text-gray-500">
              Tidak ada ujian terjadwal.
            </p>
          </div>
        ) : (
          <div
            className={[
              "overflow-y-auto pr-1 space-y-3",
              maxHeightClassName,
            ].join(" ")}
          >
            {items.map((item) => (
              <Link
                key={item.id}
                to={`/jadwal-ujian/${item.id}`}
                className={[
                  "group relative block overflow-hidden rounded-xl border border-gray-200 bg-white p-4 shadow-sm transition-all duration-300",
                  "hover:-translate-y-0.5 hover:border-[#397e50]/30 hover:shadow-md",
                  "cursor-pointer focus:outline-none focus:ring-2 focus:ring-[#397e50]/20",
                ].join(" ")}
              >
                {/* Decorative side bar on hover */}
                <div className="absolute inset-y-0 left-0 w-1 bg-[#397e50] opacity-0 transition-opacity group-hover:opacity-100" />

                <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:gap-4">
                  {/* Date Box (Left Side) */}
                  <div className="flex shrink-0 flex-row items-center gap-3 rounded-lg bg-[#397e50]/5 p-2.5 sm:flex-col sm:items-center sm:justify-center sm:gap-1 sm:p-3 sm:w-20">
                    <CalendarDays className="h-4 w-4 text-[#397e50] sm:h-5 sm:w-5" />
                    <span className="text-xs font-bold text-[#37513d] text-center leading-tight">
                      {item.tgl_ujian}
                    </span>
                  </div>

                  {/* Main Content */}
                  <div className="flex-1 min-w-0">
                    {/* Top Row: Session Badge & Time */}
                    <div className="mb-1.5 flex flex-wrap items-center gap-2 text-xs">
                      {typeof item.sesi_ujian === "number" && (
                        // PERUBAHAN DISINI: Menggunakan warna solid bg-[#397e50]
                        <span className="inline-flex items-center rounded-md bg-[#397e50] px-2 py-0.5 font-bold text-white shadow-sm">
                          Sesi {item.sesi_ujian}
                        </span>
                      )}
                      <div className="flex items-center gap-1 font-medium text-gray-500">
                        <Clock className="h-3.5 w-3.5" />
                        {item.waktu_mulai}
                      </div>
                    </div>

                    {/* Title */}
                    <h3 className="mb-2 text-base font-bold text-gray-800 transition-colors group-hover:text-[#397e50] line-clamp-2">
                      {item.nama_ujian}
                    </h3>

                    {/* Footer Details: Pengawas & Ruang */}
                    <div className="flex flex-wrap items-center gap-y-2 gap-x-4 border-t border-gray-100 pt-2.5 text-xs text-gray-600">
                      <div className="flex items-center gap-1.5 min-w-0 max-w-[50%]">
                        <User className="h-3.5 w-3.5 shrink-0 text-gray-400" />
                        <span className="truncate" title={item.pengawas_ujian}>
                          {item.pengawas_ujian}
                        </span>
                      </div>

                      <div className="flex items-center gap-1.5 min-w-0">
                        <MapPin className="h-3.5 w-3.5 shrink-0 text-gray-400" />
                        {item.ruang_ujian ? (
                          <span className="font-medium text-[#37513d]">
                            {item.ruang_ujian}
                          </span>
                        ) : (
                          <span className="italic text-gray-400">
                            Belum ditentukan
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </section>
  );
};

export default JadwalUjianWidget;
