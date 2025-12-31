

type SimpleStatWidgetProps = {
  title: string;
  value: number | string;
  trendText: string;
  trend?: "up" | "down" | "neutral";
  valueFormatter?: (n: number) => string;
  className?: string;
  trendColorClassName?: string;
};

export const SimpleStatWidget = ({
  title,
  value,
  trendText,
  trend = "up",
  valueFormatter,
  className,
  trendColorClassName,
}: SimpleStatWidgetProps) => {
  const formattedValue =
    typeof value === "number"
      ? valueFormatter?.(value) ?? new Intl.NumberFormat("id-ID").format(value)
      : value;

  const trendColor =
    trendColorClassName ??
    (trend === "up"
      ? "text-emerald-600"
      : trend === "down"
      ? "text-rose-600"
      : "text-slate-500");

  return (
    <div
      className={[
        "rounded-2xl border border-slate-100 bg-white",
        "shadow-[0_10px_24px_rgba(15,23,42,0.06)]",
        "px-4 py-5 sm:px-6 sm:py-8",
        className ?? "",
      ].join(" ")}
    >
      {/* pakai space-y untuk jarak vertikal yang rapi */}
      <div className="flex flex-col space-y-5 sm:space-y-6">
        <div className="text-xs font-medium text-slate-400 sm:text-sm">
          {title}
        </div>

        <div className="text-[28px] font-bold leading-none tracking-[-0.02em] text-slate-900 sm:text-[35px]">
          {formattedValue}
        </div>

        <div
          className={[
            "flex items-start gap-2",
            "text-xs font-medium sm:text-sm",
            trendColor,
          ].join(" ")}
        >
          <TrendIcon direction={trend} className="mt-px sm:mt-0" />
          <span className="leading-snug sm:truncate">{trendText}</span>
        </div>
      </div>
    </div>
  );
};

function TrendIcon({
  direction,
  className,
}: {
  direction: "up" | "down" | "neutral";
  className?: string;
}) {
  if (direction === "neutral") {
    return (
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        className={["shrink-0", className ?? ""].join(" ")}
        aria-hidden="true"
      >
        <path
          d="M4 12h16"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
        />
      </svg>
    );
  }

  if (direction === "down") {
    return (
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        className={["shrink-0", className ?? ""].join(" ")}
        aria-hidden="true"
      >
        <path
          d="M4 7l7 7 4-4 5 5"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path
          d="M20 20v-6h-6"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    );
  }

  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 24 24"
      fill="none"
      className={["shrink-0", className ?? ""].join(" ")}
      aria-hidden="true"
    >
      <path
        d="M4 17l7-7 4 4 5-5"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M20 7v6h-6"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
