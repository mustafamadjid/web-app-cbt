import { useEffect, useRef, useState } from "react";

const DAYS = ["Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"];
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

interface DatePickerProps {
  label: string;
  value: string;
  onChange: (date: string) => void;
  error?: boolean;
  onBlur?: () => void;
  id?: string;
  disabled?: boolean;
}

const CalendarIcon = ({ className }: { className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
    <line x1="16" y1="2" x2="16" y2="6" />
    <line x1="8" y1="2" x2="8" y2="6" />
    <line x1="3" y1="10" x2="21" y2="10" />
  </svg>
);

const ChevronLeft = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="15 18 9 12 15 6" />
  </svg>
);

const ChevronRight = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <polyline points="9 18 15 12 9 6" />
  </svg>
);

const parseIsoDate = (value: string) => {
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !month || !day) return null;

  const parsed = new Date(year, month - 1, day);
  if (
    parsed.getFullYear() !== year ||
    parsed.getMonth() !== month - 1 ||
    parsed.getDate() !== day
  ) {
    return null;
  }

  return parsed;
};

const formatIsoDate = (date: Date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

const formatDisplayDate = (value: string) => {
  const parsed = parseIsoDate(value);
  if (!parsed) return "";

  return new Intl.DateTimeFormat("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(parsed);
};

const DatePicker = ({
  label,
  value,
  onChange,
  error = false,
  onBlur,
  id,
  disabled = false,
}: DatePickerProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [viewDate, setViewDate] = useState(() => parseIsoDate(value) ?? new Date());
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) return;

    const handleOutsideClick = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
        onBlur?.();
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsOpen(false);
        onBlur?.();
      }
    };

    document.addEventListener("mousedown", handleOutsideClick);
    document.addEventListener("keydown", handleEscape);

    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onBlur]);

  useEffect(() => {
    if (disabled) setIsOpen(false);
  }, [disabled]);

  const year = viewDate.getFullYear();
  const month = viewDate.getMonth();
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const firstDay = new Date(year, month, 1).getDay();
  const emptyDays = Array.from({ length: firstDay });
  const currentDays = Array.from({ length: daysInMonth }, (_, index) => index + 1);
  const today = formatIsoDate(new Date());
  const displayValue = formatDisplayDate(value);
  const currentYear = new Date().getFullYear();
  const startYear = Math.min(currentYear - 20, year - 10);
  const endYear = Math.max(currentYear + 20, year + 10);
  const yearOptions = Array.from(
    { length: endYear - startYear + 1 },
    (_, index) => startYear + index,
  );

  const closeCalendar = () => {
    setIsOpen(false);
    onBlur?.();
  };

  const handleDateClick = (day: number) => {
    const selectedDate = new Date(year, month, day);
    onChange(formatIsoDate(selectedDate));
    closeCalendar();
  };

  const handleToggleOpen = () => {
    if (disabled) return;

    if (isOpen) {
      setIsOpen(false);
      onBlur?.();
      return;
    }

    const selectedDate = parseIsoDate(value);
    setViewDate(selectedDate ?? new Date());
    setIsOpen(true);
  };

  return (
    <div className="relative w-full" ref={containerRef}>
      <label htmlFor={id} className="mb-1 block text-xs font-medium text-slate-600">
        {label}
      </label>

      <button
        id={id}
        type="button"
        onClick={handleToggleOpen}
        disabled={disabled}
        className={`flex w-full items-center justify-between rounded-lg border px-3 py-2 text-sm transition ${
          disabled
            ? "cursor-not-allowed bg-slate-50 text-slate-500 opacity-70"
            : "cursor-pointer hover:bg-slate-50"
        } ${
          isOpen
            ? "border-[#397e50] ring-1 ring-[#397e50]"
            : error
              ? "border-rose-300 ring-rose-100"
              : "border-slate-200"
        }`}
        aria-expanded={isOpen}
        aria-haspopup="dialog"
        aria-disabled={disabled}
      >
        <span
          className={
            disabled
              ? "text-slate-500"
              : displayValue
                ? "text-slate-900"
                : "text-slate-400"
          }
        >
          {displayValue || "Pilih Tanggal"}
        </span>
        <CalendarIcon className="h-4 w-4 text-slate-400" />
      </button>

      {isOpen && (
        <div className="absolute z-50 mt-2 w-72 rounded-xl border border-slate-200 bg-white p-4 shadow-xl animate-in fade-in zoom-in-95 duration-200">
          <div className="mb-4 flex items-center justify-between">
            <button
              type="button"
              onClick={() => setViewDate(new Date(year, month - 1, 1))}
              className="cursor-pointer rounded-full p-1 text-slate-500 transition hover:bg-slate-100"
              aria-label="Bulan sebelumnya"
            >
              <ChevronLeft />
            </button>
            <div className="flex items-center gap-2">
              <select
                value={month}
                onChange={(event) =>
                  setViewDate(new Date(year, Number(event.target.value), 1))
                }
                className="cursor-pointer rounded-md border border-slate-200 bg-white px-4 py-1 text-xs font-medium text-slate-700 outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50]"
                aria-label="Pilih bulan"
              >
                {MONTHS.map((monthName, monthIndex) => (
                  <option key={monthName} value={monthIndex}>
                    {monthName}
                  </option>
                ))}
              </select>
              <select
                value={year}
                onChange={(event) =>
                  setViewDate(new Date(Number(event.target.value), month, 1))
                }
                className="cursor-pointer rounded-md border border-slate-200 bg-white px-7 py-1 text-xs font-medium text-slate-700 outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50]"
                aria-label="Pilih tahun"
              >
                {yearOptions.map((yearOption) => (
                  <option key={yearOption} value={yearOption}>
                    {yearOption}
                  </option>
                ))}
              </select>
            </div>
            <button
              type="button"
              onClick={() => setViewDate(new Date(year, month + 1, 1))}
              className="cursor-pointer rounded-full p-1 text-slate-500 transition hover:bg-slate-100"
              aria-label="Bulan berikutnya"
            >
              <ChevronRight />
            </button>
          </div>

          <div className="mb-2 grid grid-cols-7 text-center">
            {DAYS.map((dayName) => (
              <span key={dayName} className="text-[10px] font-medium text-slate-400">
                {dayName}
              </span>
            ))}
          </div>

          <div className="grid grid-cols-7 gap-1 text-center">
            {emptyDays.map((_, index) => (
              <div key={`empty-${index}`} />
            ))}

            {currentDays.map((day) => {
              const isoDate = `${year}-${String(month + 1).padStart(2, "0")}-${String(
                day,
              ).padStart(2, "0")}`;
              const isSelected = value === isoDate;
              const isToday = today === isoDate;

              return (
                <button
                  key={day}
                  type="button"
                  onClick={() => handleDateClick(day)}
                  className={`flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-xs transition ${
                    isSelected
                      ? "bg-[#397e50] font-semibold text-white shadow-md shadow-green-200"
                      : isToday
                        ? "bg-slate-100 font-semibold text-[#397e50]"
                        : "text-slate-600 hover:bg-slate-100"
                  }`}
                >
                  {day}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};

export default DatePicker;
