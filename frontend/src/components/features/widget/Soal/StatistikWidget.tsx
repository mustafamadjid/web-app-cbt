import React from "react";

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
  donutSize = 112,
  donutStroke = 10,
  trackColor = "#EEF2F6",
  gradientFrom = "#8dd19d",
  gradientTo = "#1ba23a",
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
        "flex flex-col sm:flex-row sm:items-center sm:justify-between",
        "gap-5 sm:gap-8",
        "rounded-2xl border border-slate-100 bg-white",
        "px-4 py-5 sm:px-6 sm:py-6",
        "shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
        className ?? "",
      ].join(" ")}
    >
      {/* Left */}
      <div className="min-w-0 flex flex-1 flex-col space-y-4 sm:space-y-6">
        <div className="truncate text-xs font-medium text-slate-400 sm:text-sm">
          {title}
        </div>

        <div className="text-[24px] font-bold leading-none tracking-[-0.02em] text-slate-900 sm:text-[28px]">
          {formattedValue}
        </div>

        <div className="flex items-start gap-3 text-xs font-medium text-slate-400 sm:items-center sm:text-sm">
          <span className="mt-0.5 h-4 w-[3px] rounded-full bg-slate-300 sm:mt-0" />
          <span className="leading-snug sm:truncate">{footerText}</span>
        </div>
      </div>

      {/* Right */}
      <div
        className="shrink-0 self-start sm:self-auto"
        style={{ width: donutSize, height: donutSize }}
        aria-label={`Progress ${p}%`}
      >
        <div
          className="relative grid place-items-center"
          style={{ width: donutSize, height: donutSize }}
        >
          <svg
            width={donutSize}
            height={donutSize}
            viewBox={`0 0 ${donutSize} ${donutSize}`}
          >
            <defs>
              <linearGradient id={`pw-grad-${gid}`} x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stopColor={gradientFrom} />
                <stop offset="100%" stopColor={gradientTo} />
              </linearGradient>
            </defs>

            <circle
              cx={donutSize / 2}
              cy={donutSize / 2}
              r={r}
              fill="none"
              stroke={trackColor}
              strokeWidth={donutStroke}
            />

            <circle
              cx={donutSize / 2}
              cy={donutSize / 2}
              r={r}
              fill="none"
              stroke={`url(#pw-grad-${gid})`}
              strokeWidth={donutStroke}
              strokeLinecap="round"
              strokeDasharray={`${dash} ${c - dash}`}
              transform={`rotate(-90 ${donutSize / 2} ${donutSize / 2})`}
            />
          </svg>

          <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
            <div className="text-[18px] font-bold leading-tight text-slate-700 sm:text-[22px]">
              {Math.round(p)}%
            </div>
            <div className="mt-1 text-[11px] font-semibold text-slate-400 sm:text-xs">
              {centerLabel}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

function clamp(n: number, min: number, max: number) {
  return Math.max(min, Math.min(max, n));
}
