import React from "react";
import { ClipboardCheck, CheckCircle2 } from "lucide-react";

type UjianTerlaksanaWidgetProps = {
  totalSelesai: number;
  totalUjian: number;
  className?: string;
};

const UjianTerlaksanaWidget: React.FC<UjianTerlaksanaWidgetProps> = ({
  totalSelesai,
  totalUjian,
  className,
}) => {
  const progress = totalUjian > 0 ? (totalSelesai / totalUjian) * 100 : 0;

  return (
    <section
      className={[
        "flex h-full flex-col justify-between rounded-xl border border-gray-200 bg-white p-5 shadow-sm",
        "transition-all duration-300 hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      <header className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
          <ClipboardCheck className="h-5 w-5" />
        </div>
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-gray-500">
            Ujian Terlaksana
          </p>
          <h3 className="text-lg font-bold text-[#37513d]">
            {totalSelesai} dari {totalUjian} ujian
          </h3>
        </div>
      </header>

      <div className="mt-6 space-y-3">
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <CheckCircle2 className="h-4 w-4 text-[#397e50]" />
          <span>
            Progres penyelesaian{" "}
            <strong className="text-[#37513d]">
              {Math.round(progress)}%
            </strong>
          </span>
        </div>
        <div className="h-2 w-full rounded-full bg-gray-100">
          <div
            className="h-2 rounded-full bg-[#397e50]"
            style={{ width: `${progress}%` }}
          />
        </div>
        <p className="text-xs text-gray-500">
          Terus pertahankan performa belajar agar semua ujian tuntas tepat waktu.
        </p>
      </div>
    </section>
  );
};

export default UjianTerlaksanaWidget;
