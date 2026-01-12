import React from "react";
import { Link } from "react-router";
import {
  Timer,
  Users,
  ArrowRight,
  BookOpen,
  Clock,

} from "lucide-react";

import type { UjianBerlangsungItem } from "@/types/Widget/UjianBerlangsung";


type UjianBerlangsungWidgetProps = {
  items: UjianBerlangsungItem[];
  className?: string;
  maxHeightClassName?: string;
};

const UjianBerlangsungWidget: React.FC<UjianBerlangsungWidgetProps> = ({
  items,
  className,
  maxHeightClassName = "max-h-[500px]",
}) => {
  return (
    <section
      className={[
        "relative flex flex-col overflow-hidden rounded-xl bg-white",
        "border border-gray-200 shadow-sm transition-all duration-300",
        "hover:shadow-lg hover:shadow-[#397e50]/5",
        className ?? "",
      ].join(" ")}
    >
      {/* Top Accent Line */}
      <div className="h-1.5 w-full bg-linear-to-r from-[#397e50] to-[#37513d]" />

      {/* Header */}
      <header className="flex items-center justify-between px-5 pt-5 pb-3">
        <div className="flex items-center gap-3">
          {/* Icon Timer dengan efek Pulse agar terasa 'Live' */}
          <div className="relative flex h-10 w-10 items-center justify-center rounded-full bg-rose-50 text-rose-600">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-rose-400 opacity-20"></span>
            <Timer className="relative h-5 w-5" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-[#37513d]">
              Sedang Berlangsung
            </h2>
            <div className="flex items-center gap-2">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
              </span>
              <p className="text-xs font-medium text-gray-500">
                {items.length} ujian aktif sekarang
              </p>
            </div>
          </div>
        </div>
      </header>

      {/* Content List */}
      <div className="flex-1 border-t border-gray-100 bg-gray-50/30 p-4">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 py-10 text-center rounded-xl border border-dashed border-gray-300 bg-white">
            <div className="rounded-full bg-gray-100 p-3 text-gray-400">
              <Clock className="h-6 w-6" />
            </div>
            <div>
              <p className="text-sm font-bold text-gray-700">
                Tidak ada ujian aktif
              </p>
              <p className="text-xs text-gray-500">
                Semua jadwal ujian telah selesai atau belum dimulai.
              </p>
            </div>
          </div>
        ) : (
          <div
            className={[
              "space-y-4 overflow-y-auto pr-1",
              maxHeightClassName,
            ].join(" ")}
          >
            {items.map((item) => (
              <ActiveExamCard key={item.id} item={item} />
            ))}
          </div>
        )}
      </div>
    </section>
  );
};

// --- Sub-Component: Card Item ---
const ActiveExamCard = ({ item }: { item: UjianBerlangsungItem }) => {
  // Hitung persentase progress (Selesai / Total)
  const progressPercent = Math.round(
    (item.siswa_selesai / item.total_siswa) * 100
  );

  // Hitung sisa (yang belum selesai / sedang mengerjakan)
  const activePercent = Math.round(
    (item.siswa_mengerjakan / item.total_siswa) * 100
  );

  return (
    <div className="group relative overflow-hidden rounded-xl border border-gray-200 bg-white p-4 shadow-sm transition-all hover:border-[#397e50]/40 hover:shadow-md">
      {/* Header Card: Mapel & Waktu */}
      <div className="mb-3 flex items-start justify-between">
        <div className="space-y-1">
          <div className="flex flex-wrap items-center gap-2">
            <span className="inline-flex items-center gap-1 rounded bg-[#397e50]/10 px-2 py-0.5 text-2xs font-bold uppercase tracking-wider text-[#397e50]">
              <BookOpen className="h-3 w-3" />
              {item.mata_pelajaran}
            </span>
            <span className="text-2xs font-medium text-gray-400">
              • {item.kelas.join(", ")}
            </span>
          </div>
          <h3 className="line-clamp-1 text-base font-bold text-gray-800 group-hover:text-[#37513d]">
            {item.nama_ujian}
          </h3>
        </div>

        {/* Badge Waktu */}
        <div className="shrink-0 rounded-lg bg-gray-50 px-2 py-1 text-center border border-gray-100">
          <div className="text-2xs uppercase text-gray-400 font-bold">
            Waktu
          </div>
          <div className="text-xs font-bold text-gray-700">
            {item.waktu_mulai} - {item.waktu_selesai}
          </div>
        </div>
      </div>

      {/* Progress Bar Section */}
      <div className="mb-4">
        <div className="mb-1.5 flex items-end justify-between text-xs">
          <span className="font-medium text-gray-500">Progress Siswa</span>
          <span className="font-bold text-[#37513d]">
            {item.siswa_selesai}{" "}
            <span className="text-gray-400">/ {item.total_siswa} Selesai</span>
          </span>
        </div>

        {/* Multi-segment Progress Bar */}
        <div className="flex h-2.5 w-full overflow-hidden rounded-full bg-gray-100">
          {/* Segment: Selesai (Hijau) */}
          <div
            className="bg-linear-to-r from-[#397e50] to-[#37513d] transition-all duration-500"
            style={{ width: `${progressPercent}%` }}
            title={`${progressPercent}% Selesai`}
          />
          {/* Segment: Sedang Mengerjakan (Kuning/Amber) */}
          <div
            className="bg-amber-400 transition-all duration-500"
            style={{ width: `${activePercent}%` }}
            title={`${activePercent}% Sedang Mengerjakan`}
          />
        </div>

        {/* Legend Kecil */}
        <div className="mt-1.5 flex justify-end gap-3 text-2xs text-gray-400">
          <div className="flex items-center gap-1">
            <span className="h-1.5 w-1.5 rounded-full bg-[#397e50]" />
            Selesai
          </div>
          <div className="flex items-center gap-1">
            <span className="h-1.5 w-1.5 rounded-full bg-amber-400" />
            Mengerjakan
          </div>
        </div>
      </div>

      {/* Footer Action */}
      <div className="mt-3 flex items-center justify-between border-t border-gray-50 pt-3">
        <div className="flex items-center gap-2 text-xs text-gray-500">
          <Users className="h-3.5 w-3.5" />
          <span>
            <strong className="text-gray-700">{item.siswa_mengerjakan}</strong>{" "}
            sedang online
          </span>
        </div>

        <Link
          to={`/proctoring/${item.id}`} // Sesuaikan dengan route monitoring Anda
          className="inline-flex items-center gap-1 rounded-lg bg-[#397e50] px-3 py-1.5 text-xs font-bold text-white transition-all hover:bg-[#2f5c3f] active:scale-95"
        >
          Pantau
          <ArrowRight className="h-3 w-3" />
        </Link>
      </div>
    </div>
  );
};

export default UjianBerlangsungWidget;
