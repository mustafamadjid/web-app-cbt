import React from "react";
import { TrendingUp, TrendingDown, Minus, BarChart3 } from "lucide-react";

type SimpleStatWidgetProps = {
  title: string;
  value: number | string;
  trendText: string;
  trend?: "up" | "down" | "neutral";
  valueFormatter?: (n: number) => string;
  className?: string;
  trendColorClassName?: string;
};

const SimpleStatWidget = ({
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

  // Tentukan Style berdasarkan Trend
  const trendConfig = React.useMemo(() => {
    if (trend === "up") {
      return {
        icon: TrendingUp,
        baseClass: "text-emerald-600 bg-emerald-50 border-emerald-100",
        customClass: trendColorClassName,
      };
    }
    if (trend === "down") {
      return {
        icon: TrendingDown,
        baseClass: "text-rose-600 bg-rose-50 border-rose-100",
        customClass: trendColorClassName,
      };
    }
    return {
      icon: Minus,
      baseClass: "text-slate-600 bg-slate-50 border-slate-100",
      customClass: trendColorClassName,
    };
  }, [trend, trendColorClassName]);

  const TrendIcon = trendConfig.icon;
  // Jika ada custom class, pakai itu. Jika tidak, pakai default style pill.
  const trendStyleClass = trendConfig.customClass
    ? trendConfig.customClass
    : `inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-bold border ${trendConfig.baseClass}`;

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

      <div className="relative z-10 flex flex-1 flex-col justify-between p-5">
        {/* Header Title */}
        <div className="flex items-center gap-2">
          <h3 className="text-[11px] font-bold uppercase tracking-wider text-gray-500">
            {title}
          </h3>
        </div>

        {/* Main Value */}
        <div className="mt-2 text-3xl font-black text-[#37513d] sm:text-4xl">
          {formattedValue}
        </div>

        {/* Footer: Trend */}
        <div className="mt-4">
          <div className={trendStyleClass}>
            <TrendIcon className="h-3.5 w-3.5" />
            <span className="truncate">{trendText}</span>
          </div>
        </div>
      </div>

      {/* Decorative Watermark */}
      <BarChart3
        className="absolute -bottom-3 -right-3 h-24 w-24 text-gray-100/50 transition-transform duration-500 group-hover:scale-110 group-hover:text-[#397e50]/5"
        strokeWidth={1}
      />
    </div>
  );
};

export default SimpleStatWidget;
