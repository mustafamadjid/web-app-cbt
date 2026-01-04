import React from "react";
import { TrendingUp } from "lucide-react";

type StatistikWidgetProps = {
  title: string;
  value: number | string;
  footerText: string;
  percent: number;
  centerLabel: string;
  valueFormatter?: (n: number) => string;
  className?: string;
  donutSize?: number;
  donutStroke?: number;
  trackColor?: string;
  gradientFrom?: string;
  gradientTo?: string;
};

export const StatistikWidget = ({
  title,
  value,
  footerText,
  percent,
  centerLabel,
  valueFormatter,
  className,
  donutSize = 100, // Sedikit diperkecil agar proporsional
  donutStroke = 8,
  trackColor = "#f1f5f9", // Slate-100
  gradientFrom = "#397e50", // Warna branding 1
  gradientTo = "#37513d", // Warna branding 2
}: StatistikWidgetProps) => {
  const p = clamp(percent, 0, 100);

  const formattedValue =
    typeof value === "number"
      ? valueFormatter?.(value) ?? new Intl.NumberFormat("id-ID").format(value)
      : value;

  const r = (donutSize - donutStroke) / 2;
  const c = 2 * Math.PI * r;
  const dash = (p / 100) * c;

  const gid = React.useId().replace(/:/g, "");

  return (
    <div
      className={[
        "group relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:-translate-y-1 hover:border-[#397e50]/30 hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      {/* Top Accent Line */}
      <div className="h-1 w-full bg-linear-to-r from-[#397e50] to-[#37513d]" />

      <div className="flex flex-1 items-center justify-between p-5 sm:p-6">
        {/* Left Content */}
        <div className="flex min-w-0 flex-1 flex-col justify-between self-stretch">
          <div className="flex flex-col items-star gap-2 md:gap-6 ">
            {/* Title */}
            <h3 className="text-[11px] font-bold uppercase tracking-wider text-gray-500">
              {title}
            </h3>

            {/* Main Value */}
            <div className="mt-1 text-3xl font-black text-[#37513d] sm:text-4xl">
              {formattedValue}
            </div>
          </div>

          {/* Footer Text */}
          <div className="mt-4 flex items-center gap-2 text-xs font-medium text-gray-500">
            <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-[#397e50]">
              <TrendingUp className="h-3 w-3" />
            </div>
            <span className="truncate leading-snug">{footerText}</span>
          </div>
        </div>

        {/* Right Content (Chart) */}
        <div
          className="relative ml-4 shrink-0"
          style={{ width: donutSize, height: donutSize }}
          aria-label={`Progress ${p}%`}
        >
          {/* SVG Chart */}
          <svg
            width={donutSize}
            height={donutSize}
            viewBox={`0 0 ${donutSize} ${donutSize}`}
            className="-rotate-90 transition-all duration-500 group-hover:scale-105"
          >
            <defs>
              <linearGradient id={`pw-grad-${gid}`} x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stopColor={gradientFrom} />
                <stop offset="100%" stopColor={gradientTo} />
              </linearGradient>
            </defs>

            {/* Track Circle */}
            <circle
              cx={donutSize / 2}
              cy={donutSize / 2}
              r={r}
              fill="none"
              stroke={trackColor}
              strokeWidth={donutStroke}
              className="opacity-80"
            />

            {/* Progress Circle */}
            <circle
              cx={donutSize / 2}
              cy={donutSize / 2}
              r={r}
              fill="none"
              stroke={`url(#pw-grad-${gid})`}
              strokeWidth={donutStroke}
              strokeLinecap="round"
              strokeDasharray={`${dash} ${c - dash}`}
              className="transition-all duration-1000 ease-out"
            />
          </svg>

          {/* Center Label */}
          <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
            <span className="text-xl font-bold text-[#37513d]">
              {Math.round(p)}%
            </span>
            <span className="text-2xs font-semibold uppercase tracking-wide text-gray-400">
              {centerLabel}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

function clamp(n: number, min: number, max: number) {
  return Math.max(min, Math.min(max, n));
}
