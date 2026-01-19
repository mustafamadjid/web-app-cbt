import React from "react";
import {
  CalendarCheck2,
  CalendarDays,
  ChevronLeft,
  ChevronRight,
  Clock,
  MapPin,
  User,
} from "lucide-react";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import { formatTanggalToIso } from "@/helper/dateFormatting/formatToIso";

type JadwalUjianSiswaWidgetProps = {
  title?: string;
  items: JadwalUjianItem[];
  className?: string;
  defaultView?: "list" | "calendar";
};

type CalendarItem = JadwalUjianItem & {
  isoDate?: string;
  displayDate: string;
};

const MONTHS_ID = [
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

const WEEKDAYS_ID = ["Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"];

const buildDisplayDate = (item: JadwalUjianItem) => {
  if (item.tgl_ujian) return item.tgl_ujian;
  if (item.tanggal_ujian) {
    const dt = new Date(item.tanggal_ujian);
    return dt.toLocaleDateString("id-ID", {
      weekday: "long",
      day: "2-digit",
      month: "long",
      year: "numeric",
    });
  }
  return "-";
};

const toIsoDate = (item: JadwalUjianItem) =>
  item.tanggal_ujian ?? formatTanggalToIso(item.tgl_ujian ?? "");

const getMonthLabel = (date: Date) =>
  `${MONTHS_ID[date.getMonth()]}, ${date.getFullYear()}`;

const getInitialMonth = (items: CalendarItem[]) => {
  const today = new Date();
  const upcoming = items
    .map((item) => item.isoDate)
    .filter(Boolean)
    .map((iso) => new Date(iso as string))
    .sort((a, b) => a.getTime() - b.getTime())
    .find((date) => date >= new Date(today.getFullYear(), today.getMonth(), 1));

  return upcoming ?? today;
};

const JadwalUjianSiswaWidget: React.FC<JadwalUjianSiswaWidgetProps> = ({
  title = "Jadwal Ujian Mendatang",
  items,
  className,
  defaultView = "list",
}) => {
  const [view, setView] = React.useState<"list" | "calendar">(defaultView);

  const normalizedItems = React.useMemo<CalendarItem[]>(
    () =>
      items
        .map((item) => ({
          ...item,
          isoDate: toIsoDate(item) ?? undefined,
          displayDate: buildDisplayDate(item),
        }))
        .sort((a, b) =>
          (a.isoDate ?? "").localeCompare(b.isoDate ?? "")
        ),
    [items]
  );

  const [currentMonth, setCurrentMonth] = React.useState<Date>(() =>
    getInitialMonth(normalizedItems)
  );

  const scheduleByDate = React.useMemo(() => {
    const map = new Map<string, CalendarItem[]>();
    normalizedItems.forEach((item) => {
      if (!item.isoDate) return;
      const list = map.get(item.isoDate) ?? [];
      list.push(item);
      map.set(item.isoDate, list);
    });
    return map;
  }, [normalizedItems]);

  const handlePrevMonth = () => {
    setCurrentMonth(
      (prev) => new Date(prev.getFullYear(), prev.getMonth() - 1, 1)
    );
  };

  const handleNextMonth = () => {
    setCurrentMonth(
      (prev) => new Date(prev.getFullYear(), prev.getMonth() + 1, 1)
    );
  };

  const daysInMonth = new Date(
    currentMonth.getFullYear(),
    currentMonth.getMonth() + 1,
    0
  ).getDate();
  const firstDay = new Date(
    currentMonth.getFullYear(),
    currentMonth.getMonth(),
    1
  ).getDay();
  const startOffset = (firstDay + 6) % 7;
  const totalCells = Math.ceil((startOffset + daysInMonth) / 7) * 7;

  return (
    <section
      className={[
        "relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      <header className="flex flex-wrap items-center justify-between gap-3 px-5 pt-5 pb-3">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
            <CalendarCheck2 className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-[#37513d]">{title}</h2>
            <p className="text-xs font-medium text-gray-500">
              {items.length} jadwal terdekat
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 rounded-full border border-gray-200 bg-gray-50 p-1 text-xs">
          <button
            type="button"
            onClick={() => setView("list")}
            className={[
              "rounded-full px-3 py-1 font-semibold transition",
              view === "list"
                ? "bg-[#397e50] text-white"
                : "text-gray-600 hover:text-[#397e50]",
            ].join(" ")}
          >
            List
          </button>
          <button
            type="button"
            onClick={() => setView("calendar")}
            className={[
              "rounded-full px-3 py-1 font-semibold transition",
              view === "calendar"
                ? "bg-[#397e50] text-white"
                : "text-gray-600 hover:text-[#397e50]",
            ].join(" ")}
          >
            Kalender
          </button>
        </div>
      </header>

      <div className="border-t border-gray-100 bg-gray-50/50 p-4">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-gray-300 bg-white py-12 text-center">
            <CalendarDays className="h-10 w-10 text-gray-300" />
            <p className="text-sm font-medium text-gray-500">
              Tidak ada jadwal ujian dalam waktu dekat.
            </p>
          </div>
        ) : view === "list" ? (
          <div className="space-y-3">
            {normalizedItems.map((item) => (
              <div
                key={item.id}
                className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm transition hover:border-[#397e50]/30"
              >
                <div className="flex flex-wrap items-start gap-4">
                  <div className="flex shrink-0 flex-col items-start rounded-lg bg-[#397e50]/5 px-3 py-2 text-xs font-bold text-[#37513d]">
                    {item.displayDate}
                    <span className="mt-1 inline-flex items-center gap-1 text-gray-500">
                      <Clock className="h-3.5 w-3.5" />
                      {item.waktu_mulai}
                    </span>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-base font-bold text-gray-800">
                      {item.nama_ujian}
                    </h3>
                    <div className="mt-2 flex flex-wrap items-center gap-4 text-xs text-gray-500">
                      <span className="inline-flex items-center gap-1">
                        <User className="h-3.5 w-3.5" />
                        {item.pengawas_ujian}
                      </span>
                      <span className="inline-flex items-center gap-1">
                        <MapPin className="h-3.5 w-3.5" />
                        {item.ruang_ujian ?? "Belum ditentukan"}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-gray-200 bg-white p-4">
            <div className="flex items-center justify-between gap-3">
              <button
                type="button"
                onClick={handlePrevMonth}
                className="rounded-full border border-gray-200 p-2 text-gray-600 transition hover:border-[#397e50]/40 hover:text-[#397e50]"
                aria-label="Bulan sebelumnya"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
              <div className="text-sm font-bold text-[#37513d]">
                {getMonthLabel(currentMonth)}
              </div>
              <button
                type="button"
                onClick={handleNextMonth}
                className="rounded-full border border-gray-200 p-2 text-gray-600 transition hover:border-[#397e50]/40 hover:text-[#397e50]"
                aria-label="Bulan berikutnya"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>

            <div className="mt-4 grid grid-cols-7 gap-2 text-xs text-gray-500">
              {WEEKDAYS_ID.map((day) => (
                <div key={day} className="text-center font-semibold">
                  {day}
                </div>
              ))}
            </div>

            <div className="mt-2 grid grid-cols-7 gap-2 text-xs">
              {Array.from({ length: totalCells }).map((_, idx) => {
                const dayNumber = idx - startOffset + 1;
                const isCurrentMonth =
                  dayNumber >= 1 && dayNumber <= daysInMonth;
                const isoDate = isCurrentMonth
                  ? `${currentMonth.getFullYear()}-${String(
                      currentMonth.getMonth() + 1
                    ).padStart(2, "0")}-${String(dayNumber).padStart(2, "0")}`
                  : null;
                const dayItems = isoDate
                  ? scheduleByDate.get(isoDate) ?? []
                  : [];
                const hasEvents = dayItems.length > 0;

                return (
                  <div
                    key={`day-${idx}`}
                    className={[
                      "min-h-[84px] rounded-lg border p-2",
                      isCurrentMonth
                        ? "border-gray-200 bg-white"
                        : "border-transparent bg-transparent text-gray-300",
                      hasEvents
                        ? "border-[#397e50]/40 bg-[#397e50]/10"
                        : "",
                    ].join(" ")}
                  >
                    <div className="flex items-center justify-between">
                      <span
                        className={[
                          "text-xs font-semibold",
                          hasEvents ? "text-[#37513d]" : "text-gray-500",
                        ].join(" ")}
                      >
                        {isCurrentMonth ? dayNumber : ""}
                      </span>
                      {hasEvents && (
                        <span className="h-1.5 w-1.5 rounded-full bg-[#397e50]" />
                      )}
                    </div>
                    {hasEvents && (
                      <ul className="mt-1 space-y-1 text-[10px] text-[#37513d]">
                        {dayItems.slice(0, 2).map((event) => (
                          <li
                            key={event.id}
                            className="line-clamp-2 rounded bg-white/70 px-1 py-0.5 shadow-sm"
                            title={event.nama_ujian}
                          >
                            {event.nama_ujian}
                          </li>
                        ))}
                        {dayItems.length > 2 && (
                          <li className="text-[10px] text-gray-500">
                            +{dayItems.length - 2} ujian lainnya
                          </li>
                        )}
                      </ul>
                    )}
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

export default JadwalUjianSiswaWidget;
