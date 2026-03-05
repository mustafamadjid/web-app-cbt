import { useEffect, useRef, useState } from "react";

type TimePickerProps = {
  label: string;
  value: string;
  onChange: (time: string) => void;
  error?: boolean;
  onBlur?: () => void;
  id?: string;
};

const ClockIcon = ({ className }: { className?: string }) => (
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
    <circle cx="12" cy="12" r="9" />
    <polyline points="12 7 12 12 15 14" />
  </svg>
);

const parseTime = (value: string) => {
  const [hourRaw = "00", minuteRaw = "00"] = value.split(":");
  const hour = Number.parseInt(hourRaw, 10);
  const minute = Number.parseInt(minuteRaw, 10);

  if (!Number.isFinite(hour) || hour < 0 || hour > 23) {
    return { hour: 0, minute: Number.isFinite(minute) ? minute : 0 };
  }
  if (!Number.isFinite(minute) || minute < 0 || minute > 59) {
    return { hour, minute: 0 };
  }

  return { hour, minute };
};

const pad = (value: number) => String(value).padStart(2, "0");

const formatTime = (hour: number, minute: number) => `${pad(hour)}:${pad(minute)}`;

const HOURS = Array.from({ length: 24 }, (_, idx) => idx);
const MINUTES = Array.from({ length: 60 }, (_, idx) => idx);

const TimePicker = ({
  label,
  value,
  onChange,
  error = false,
  onBlur,
  id,
}: TimePickerProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [draftHour, setDraftHour] = useState(0);
  const [draftMinute, setDraftMinute] = useState(0);

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

  const display = value ? `${value} WIB` : "Pilih Waktu (WIB)";

  const openPicker = () => {
    const parsed = parseTime(value);
    setDraftHour(parsed.hour);
    setDraftMinute(parsed.minute);
    setIsOpen(true);
  };

  const closePicker = () => {
    setIsOpen(false);
    onBlur?.();
  };

  const applyPicker = () => {
    onChange(formatTime(draftHour, draftMinute));
    closePicker();
  };

  return (
    <div className="relative w-full" ref={containerRef}>
      <label htmlFor={id} className="mb-1 block text-xs font-medium text-slate-600">
        {label}
      </label>

      <button
        id={id}
        type="button"
        onClick={() => (isOpen ? closePicker() : openPicker())}
        className={`flex w-full cursor-pointer items-center justify-between rounded-lg border px-3 py-2 text-sm transition hover:bg-slate-50 ${
          isOpen
            ? "border-[#397e50] ring-1 ring-[#397e50]"
            : error
              ? "border-rose-300 ring-rose-100"
              : "border-slate-200"
        }`}
        aria-expanded={isOpen}
        aria-haspopup="dialog"
      >
        <span className={value ? "text-slate-900" : "text-slate-400"}>{display}</span>
        <ClockIcon className="h-4 w-4 text-slate-400" />
      </button>

      {isOpen && (
        <div className="absolute z-50 mt-2 w-full rounded-xl border border-slate-200 bg-white p-4 shadow-xl">
          <div className="mb-3 text-xs font-medium text-slate-500">
            Format 24 jam (WIB), tanpa AM/PM.
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="mb-1 block text-xs text-slate-500">Jam</label>
              <select
                value={draftHour}
                onChange={(event) => setDraftHour(Number(event.target.value))}
                className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50]"
              >
                {HOURS.map((hour) => (
                  <option key={hour} value={hour}>
                    {pad(hour)}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="mb-1 block text-xs text-slate-500">Menit</label>
              <select
                value={draftMinute}
                onChange={(event) => setDraftMinute(Number(event.target.value))}
                className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50]"
              >
                {MINUTES.map((minute) => (
                  <option key={minute} value={minute}>
                    {pad(minute)}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <button
              type="button"
              onClick={closePicker}
              className="cursor-pointer rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:bg-slate-50"
            >
              Batal
            </button>
            <button
              type="button"
              onClick={applyPicker}
              className="cursor-pointer rounded-lg bg-[#397e50] px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-[#2f6a43]"
            >
              Pakai Waktu
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default TimePicker;
