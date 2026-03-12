import React from "react";
import { CalendarDays, Clock, GraduationCap, MapPin, User } from "lucide-react";
import type { JadwalUjianSiswaItem } from "@/types/Ujian/jadwalUjian";

const statusLabel: Record<string, string> = {
  belum_dimulai: "Mendatang",
  berlangsung: "Berlangsung",
  selesai: "Selesai",
};

type UjianCardProps = {
  item: JadwalUjianSiswaItem;
  actionLabel?: string;
  onAction?: (item: JadwalUjianSiswaItem) => void;
};

const UjianCard: React.FC<UjianCardProps> = ({
  item,
  actionLabel,
  onAction,
}) => {
  const kelasLabel =
    [item.nama_kelas?.trim(), item.tingkat_kelas ? `Tingkat ${item.tingkat_kelas}` : ""]
      .filter(Boolean)
      .join(" • ") || "Kelas belum ditentukan";
  const pengawasLabel =
    item.pengawas_nama_lengkap?.trim() || item.pengawas_ujian || "-";

  return (
    <article className="flex h-full flex-col justify-between rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition hover:shadow-lg hover:shadow-[#397e50]/5">
      <div className="space-y-4">
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="text-xs font-semibold uppercase tracking-wide text-gray-400">
              {statusLabel[item.status_ujian] ?? "Jadwal Ujian"}
            </p>
            <h3 className="text-lg font-bold text-[#37513d]">
              {item.nama_ujian}
            </h3>
            <p className="text-sm font-medium text-[#397e50]">
              {kelasLabel}
            </p>
          </div>
          <div className="rounded-full bg-[#397e50]/10 px-3 py-1 text-xs font-semibold text-[#397e50]">
            {item.sesi_ujian ? `Sesi ${item.sesi_ujian}` : "Sesi"}
          </div>
        </div>

        <div className="space-y-2 text-sm text-gray-500">
          <div className="flex items-center gap-2">
            <CalendarDays className="h-4 w-4 text-[#397e50]" />
            <span>{item.tgl_ujian}</span>
          </div>
          <div className="flex items-center gap-2">
            <Clock className="h-4 w-4 text-[#397e50]" />
            <span>
              {item.waktu_mulai} - {item.waktu_selesai ?? "-"}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <MapPin className="h-4 w-4 text-[#397e50]" />
            <span>{item.ruang_ujian ?? "Ruang belum ditentukan"}</span>
          </div>
          <div className="flex items-center gap-2">
            <GraduationCap className="h-4 w-4 text-[#397e50]" />
            <span>{kelasLabel}</span>
          </div>
          <div className="flex items-center gap-2">
            <User className="h-4 w-4 text-[#397e50]" />
            <span>{pengawasLabel}</span>
          </div>
        </div>
      </div>

      {actionLabel && (
        <button
          type="button"
          onClick={() => onAction?.(item)}
          className="mt-5 w-full rounded-full bg-[#397e50] px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
        >
          {actionLabel}
        </button>
      )}
    </article>
  );
};

export default UjianCard;
