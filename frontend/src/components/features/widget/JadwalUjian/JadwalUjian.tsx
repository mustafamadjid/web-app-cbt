import React from "react";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { Link } from "react-router";
import { CalendarCheck2, ChevronRight } from "lucide-react";

type JadwalUjianProps = {
  title?: string;
  items: JadwalUjianItem[];
  className?: string;

  /** tinggi maksimum list (opsional) */
  maxHeightClassName?: string;

  /** tujuan halaman list ujian keseluruhan */
  lihatSemuaTo?: string;
};

export const JadwalUjianWidget: React.FC<JadwalUjianProps> = ({
  title = "Jadwal Ujian",
  items,
  className,
  maxHeightClassName = "max-h-[60vh] sm:max-h-[520px]",
  lihatSemuaTo = "/jadwal-ujian",
}) => {
  return (
    <section
      className={[
        "rounded-2xl border border-slate-200 bg-white",
        "shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
        "p-4 sm:p-5",
        className ?? "",
      ].join(" ")}
    >
      {/* Header */}
      <header className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-full bg-slate-50 text-slate-600">
            <CalendarCheck2 className="h-4 w-4" />
          </span>
          <div>
            <h2 className="text-sm font-semibold text-slate-900">{title}</h2>
            <p className="text-xs text-slate-500">{items.length} jadwal</p>
          </div>
        </div>

        {/* Link pojok kanan atas */}
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
      </header>

      <div className="mt-4 border-t border-slate-100 pt-4">
        {items.length === 0 ? (
          <div className="rounded-xl border border-dashed border-slate-200 p-4 text-sm text-slate-600">
            Tidak ada ujian terjadwal.
          </div>
        ) : (
          <div
            className={["overflow-y-auto pr-1", maxHeightClassName].join(" ")}
          >
            <ul className="space-y-2">
              {items.map((item) => (
                <li key={item.id}>
                  <Link
                    to={`/jadwal-ujian/${item.id}`}
                    className={[
                      "group block rounded-xl border border-slate-200 bg-white",
                      "transition hover:border-slate-300 hover:bg-slate-50",
                      "hover:shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
                      "focus:outline-none focus:ring-2 focus:ring-slate-300",
                    ].join(" ")}
                  >
                    <div className="flex items-center justify-between gap-3 border-b border-slate-100 px-3 py-2 sm:px-4">
                      <div className="text-xs font-medium text-slate-600">
                        <span className="text-slate-900">{item.tgl_ujian}</span>
                        <span className="mx-2 text-slate-300">•</span>
                        <span>{item.waktu_mulai}</span>
                      </div>

                      <div className="flex items-center gap-2 text-xs text-slate-500">
                        {typeof item.sesi_ujian === "number" ? (
                          <span className="rounded-full bg-slate-100 px-2 py-0.5">
                            Sesi {item.sesi_ujian}
                          </span>
                        ) : null}
                        <ChevronRight className="h-4 w-4 text-slate-400 transition group-hover:translate-x-0.5" />
                      </div>
                    </div>

                    <div className="px-3 py-3 sm:px-4 sm:py-4">
                      <div className="grid gap-2 sm:grid-cols-[1fr_auto] sm:items-start">
                        <div>
                          <p className="text-sm font-semibold text-slate-900">
                            {item.nama_ujian}
                          </p>
                          <p className="mt-1 text-xs text-slate-600">
                            <span className="font-medium text-slate-700">
                              Pengawas:
                            </span>{" "}
                            {item.pengawas_ujian}
                          </p>
                        </div>

                        <div className="sm:text-right">
                          {item.ruang_ujian ? (
                            <p className="text-xs font-medium text-slate-700">
                              {item.ruang_ujian}
                            </p>
                          ) : (
                            <p className="text-xs text-slate-400">
                              Ruang belum ditentukan
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </section>
  );
};
