import React from "react";
import { Sparkles } from "lucide-react";
import type { SiswaTip } from "@/types/Widget/SiswaDashboard";

type TipsBelajarWidgetProps = {
  tips: SiswaTip[];
  className?: string;
};

const TipsBelajarWidget: React.FC<TipsBelajarWidgetProps> = ({
  tips,
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
          <Sparkles className="h-5 w-5" />
        </div>
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
            Tips Belajar
          </p>
          <h3 className="text-lg font-bold text-[#37513d]">
            Fokus Persiapan Ujian
          </h3>
        </div>
      </header>

      <ul className="mt-4 space-y-3 text-sm text-gray-600">
        {tips.length === 0 ? (
          <li className="rounded-lg border border-dashed border-gray-200 bg-gray-50 px-3 py-4 text-center text-xs text-gray-500">
            Tips belajar akan muncul di sini.
          </li>
        ) : (
          tips.map((tip) => (
            <li
              key={tip.id}
              className="rounded-lg border border-gray-100 bg-gray-50 px-3 py-3"
            >
              <p className="font-semibold text-[#37513d]">{tip.judul}</p>
              <p className="mt-1 text-xs text-gray-500">{tip.deskripsi}</p>
            </li>
          ))
        )}
      </ul>
    </section>
  );
};

export default TipsBelajarWidget;
