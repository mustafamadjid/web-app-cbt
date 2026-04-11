import type { LucideIcon } from "lucide-react";

export type StatWidgetColorTheme = "emerald" | "amber" | "blue" | "rose";

export type StatWidgetProps = {
  title: string;
  value: number | string;
  icon: LucideIcon;
  colorTheme: StatWidgetColorTheme;
};

const colors: Record<StatWidgetColorTheme, string> = {
  emerald: "bg-emerald-50 text-emerald-600",
  amber: "bg-amber-50 text-amber-600",
  blue: "bg-blue-50 text-blue-600",
  rose: "bg-rose-50 text-rose-600",
};

const StatWidget = ({
  title,
  value,
  icon: Icon,
  colorTheme,
}: StatWidgetProps) => {
  return (
    <div className="flex items-center justify-between rounded-2xl border border-slate-200 bg-white p-6">
      <div>
        <p className="text-sm font-medium text-slate-500">{title}</p>
        <h3 className="mt-2 text-2xl font-bold text-slate-800">{value}</h3>
      </div>
      <div
        className={`flex h-12 w-12 items-center justify-center rounded-xl ${colors[colorTheme]}`}
      >
        <Icon size={24} strokeWidth={2.5} />
      </div>
    </div>
  );
};

export default StatWidget;
