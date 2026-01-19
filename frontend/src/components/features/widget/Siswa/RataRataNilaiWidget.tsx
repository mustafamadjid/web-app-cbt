import React from "react";
import { Award, Target } from "lucide-react";
import type { SiswaSemesterAverage } from "@/types/Widget/SiswaDashboard";

type RataRataNilaiWidgetProps = {
  items: SiswaSemesterAverage[];
  className?: string;
};

const RataRataNilaiWidget: React.FC<RataRataNilaiWidgetProps> = ({
  items,
  className,
}) => {
  return (
    <section
      className={[
        "flex h-full flex-col rounded-xl border border-gray-200 bg-white p-5 shadow-sm",
        "transition-all duration-300 hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      <header className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
          <Award className="h-5 w-5" />
        </div>
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
            Rata-rata Nilai
          </p>
          <h3 className="text-lg font-bold text-[#37513d]">
            Performa Semester
          </h3>
        </div>
      </header>

      <div className="mt-4 space-y-3">
        {items.length === 0 ? (
          <div className="rounded-lg border border-dashed border-gray-200 bg-gray-50 px-3 py-4 text-center text-xs text-gray-500">
            Data rata-rata nilai belum tersedia.
          </div>
        ) : (
          items.map((item) => (
            <div
              key={item.semester}
              className="rounded-lg border border-gray-100 bg-gray-50 px-3 py-3"
            >
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs text-gray-500">{item.semester}</p>
                  <p className="text-lg font-bold text-[#37513d]">
                    {item.rata_rata.toFixed(1)}
                  </p>
                </div>
                {item.target != null && (
                  <div className="flex items-center gap-1 text-xs font-semibold text-[#397e50]">
                    <Target className="h-4 w-4" />
                    Target {item.target}
                  </div>
                )}
              </div>
              <div className="mt-2 h-2 w-full rounded-full bg-gray-100">
                <div
                  className="h-2 rounded-full bg-[#397e50]"
                  style={{
                    width: `${Math.min(item.rata_rata, 100)}%`,
                  }}
                />
              </div>
            </div>
          ))
        )}
      </div>
    </section>
  );
};

export default RataRataNilaiWidget;
