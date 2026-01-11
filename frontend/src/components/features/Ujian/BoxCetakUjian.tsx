import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import {
  Calendar,
  Clock,
  GraduationCap,
  MapPin,
  User,
} from "lucide-react";
import type { ReactNode } from "react";

type BoxCetakUjianProps = JadwalUjianItem & {
  actions?: ReactNode;
};

export const BoxCetakUjian = ({
  nama_ujian,
  pengawas_ujian,
  tgl_ujian,
  waktu_mulai,
  ruang_ujian,
  tingkat_kelas,
  nama_kelas,
  actions,
}: BoxCetakUjianProps) => {
  return (
    <div className="group rounded-2xl border border-[#397e50]/20 bg-white p-6 shadow-sm transition hover:border-[#397e50]/50 hover:shadow-md">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-2">
          <span className="inline-flex items-center gap-2 rounded-full border border-[#397e50]/20 bg-[#397e50]/10 px-3 py-1 text-2xs font-semibold uppercase tracking-wide text-[#397e50]">
            <GraduationCap size={12} />
            Kelas {tingkat_kelas ?? "-"} {nama_kelas ? `• ${nama_kelas}` : ""}
          </span>
          <h3 className="text-lg font-bold text-slate-800 transition group-hover:text-[#397e50]">
            {nama_ujian}
          </h3>
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <User size={14} className="text-[#397e50]" />
            <span>{pengawas_ujian}</span>
          </div>
        </div>
      </div>

      <div className="mt-5 grid gap-4 rounded-xl border border-slate-100 bg-slate-50/60 p-4 sm:grid-cols-3">
        <div className="space-y-1">
          <p className="flex items-center gap-2 text-2xs font-semibold uppercase tracking-wide text-slate-400">
            <Calendar size={12} className="text-[#397e50]" />
            Tanggal
          </p>
          <p className="text-sm font-semibold text-slate-700">{tgl_ujian}</p>
        </div>
        <div className="space-y-1">
          <p className="flex items-center gap-2 text-2xs font-semibold uppercase tracking-wide text-slate-400">
            <Clock size={12} className="text-[#397e50]" />
            Waktu
          </p>
          <p className="text-sm font-semibold text-slate-700">{waktu_mulai}</p>
        </div>
        <div className="space-y-1">
          <p className="flex items-center gap-2 text-2xs font-semibold uppercase tracking-wide text-slate-400">
            <MapPin size={12} className="text-[#397e50]" />
            Ruangan
          </p>
          <p className="text-sm font-semibold text-slate-700">
            {ruang_ujian ?? "-"}
          </p>
        </div>
      </div>

      {actions ? (
        <div className="mt-5 flex flex-wrap gap-3 border-t border-[#397e50]/10 pt-4">
          {actions}
        </div>
      ) : null}
    </div>
  );
};
