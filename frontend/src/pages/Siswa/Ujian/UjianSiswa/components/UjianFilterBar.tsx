import React from "react";
import { Filter } from "lucide-react";

const MONTHS = [
  "Januari",
  "Februari",
  "Maret",
  "April",
  "Mei",
  "Juni",
  "Juli",
  "Agustus",
  "September",
  "Oktober",
  "November",
  "Desember",
];

type UjianFilterBarProps = {
  month?: number;
  year?: number;
  mapel?: string;
  mapelOptions: string[];
  onMonthChange: (value?: number) => void;
  onYearChange: (value?: number) => void;
  onMapelChange: (value?: string) => void;
  onReset: () => void;
};

const UjianFilterBar: React.FC<UjianFilterBarProps> = ({
  month,
  year,
  mapel,
  mapelOptions,
  onMonthChange,
  onYearChange,
  onMapelChange,
  onReset,
}) => {
  const years = React.useMemo(() => {
    const currentYear = new Date().getFullYear();
    return Array.from({ length: 4 }, (_, idx) => currentYear + idx);
  }, []);

  return (
    <section className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
            <Filter className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-[#37513d]">
              Filter Ujian
            </h2>
            <p className="text-xs text-gray-500">
              Pilih bulan, tahun, atau mata pelajaran.
            </p>
          </div>
        </div>
        <button
          type="button"
          onClick={onReset}
          className="rounded-full border border-gray-200 px-4 py-2 text-xs font-semibold text-gray-500 transition hover:border-[#397e50] hover:text-[#397e50]"
        >
          Reset Filter
        </button>
      </div>

      <div className="mt-4 grid gap-3 md:grid-cols-3">
        <label className="flex flex-col gap-2 text-xs font-semibold text-gray-500">
          Bulan
          <select
            value={month ?? ""}
            onChange={(event) =>
              onMonthChange(
                event.target.value
                  ? Number(event.target.value)
                  : undefined
              )
            }
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#397e50] focus:outline-none"
          >
            <option value="">Semua Bulan</option>
            {MONTHS.map((label, index) => (
              <option key={label} value={index + 1}>
                {label}
              </option>
            ))}
          </select>
        </label>

        <label className="flex flex-col gap-2 text-xs font-semibold text-gray-500">
          Tahun
          <select
            value={year ?? ""}
            onChange={(event) =>
              onYearChange(
                event.target.value
                  ? Number(event.target.value)
                  : undefined
              )
            }
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#397e50] focus:outline-none"
          >
            <option value="">Semua Tahun</option>
            {years.map((item) => (
              <option key={item} value={item}>
                {item}
              </option>
            ))}
          </select>
        </label>

        <label className="flex flex-col gap-2 text-xs font-semibold text-gray-500">
          Mata Pelajaran
          <select
            value={mapel ?? ""}
            onChange={(event) =>
              onMapelChange(event.target.value || undefined)
            }
            className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#397e50] focus:outline-none"
          >
            <option value="">Semua Mapel</option>
            {mapelOptions.map((option) => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
        </label>
      </div>
    </section>
  );
};

export default UjianFilterBar;
