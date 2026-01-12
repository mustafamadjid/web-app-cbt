import type { InfoItem } from "@/types/Ujian/DetailUjian";

const InfoCard = ({ title, items }: { title: string; items: InfoItem[] }) => (
  <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <h2 className="text-sm font-semibold text-slate-800">{title}</h2>
    <div className="mt-4 grid gap-4 sm:grid-cols-2">
      {items.map((item) => (
        <div key={item.label} className="flex gap-3">
          {item.icon && (
            <div className="mt-0.5 text-[#397e50]">{item.icon}</div>
          )}
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
              {item.label}
            </p>
            <p className="text-sm font-medium text-slate-700">{item.value}</p>
          </div>
        </div>
      ))}
    </div>
  </div>
);

export default InfoCard;