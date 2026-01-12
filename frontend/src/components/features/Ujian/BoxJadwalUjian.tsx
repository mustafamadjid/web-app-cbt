import type { JadwalUjianItem } from "../../../types/Ujian/jadwalUjian";
import {
  Calendar,
  Clock,
  MapPin,
  User,
  Hash,
  GraduationCap,
  Play,
  XCircle,
} from "lucide-react";
import { Link } from "react-router";

const statusConfig: Record<string, { label: string; color: string }> = {
  belum_dimulai: {
    label: "Belum Dimulai",
    color: "bg-amber-50 text-amber-700 border-amber-200",
  },
  berlangsung: {
    label: "Berlangsung",
    color: "bg-emerald-50 text-emerald-700 border-emerald-200",
  },
  selesai: {
    label: "Selesai",
    color: "bg-slate-50 text-slate-600 border-slate-200",
  },
  dibatalkan: {
    label: "Dibatalkan",
    color: "bg-rose-50 text-rose-700 border-rose-200",
  },
};

type BoxJadwalUjianProps = JadwalUjianItem & {
  linkJadwal?: string;
  onStart?: (id: number) => void;
  onCancel?: (id: number) => void;
  canControl?: boolean;
};

const BoxJadwalUjian = ({
  nama_ujian,
  pengawas_ujian,
  tgl_ujian,
  waktu_mulai,
  sesi_ujian,
  ruang_ujian,
  status_ujian,
  tingkat_kelas,
  nama_kelas,
  started,
  id,
  linkJadwal = "",
  onStart,
  onCancel,
  canControl = false,
}: BoxJadwalUjianProps) => {
  const status = status_ujian
    ? statusConfig[status_ujian]
    : { label: "Unknown", color: "bg-gray-100 text-gray-500" };

  return (
    <Link to={linkJadwal} className="block group">
      <div className="relative overflow-hidden rounded-2xl border-2 border-slate-100 bg-white transition-all duration-300 hover:border-[#397e50]/30 hover:shadow-xl hover:shadow-slate-200/50">
        {/* Accent Bar */}
        <div className="absolute left-0 top-0 h-full w-1.5 bg-[#397e50]" />

        <div className="p-6">
          {/* Header: Title & Status */}
          <div className="mb-6 flex items-start justify-between gap-4">
            <div className="space-y-2">
              <span className="inline-flex items-center gap-1.5 rounded-lg bg-slate-100 px-2.5 py-1 text-2xs font-bold uppercase tracking-wider text-slate-600 border border-slate-200">
                <GraduationCap size={14} />
                Kelas {tingkat_kelas} • {nama_kelas}
              </span>
              <h3 className="text-xl font-extrabold text-slate-800 leading-tight transition-colors group-hover:text-[#397e50]">
                {nama_ujian}
              </h3>
              <div className="flex items-center gap-2 text-slate-500">
                <User size={14} className="text-[#397e50]" />
                <span className="text-sm font-semibold">{pengawas_ujian}</span>
              </div>
            </div>

            <span
              className={`shrink-0 rounded-full border px-3 py-1 text-2xs font-black uppercase tracking-widest ${status.color}`}
            >
              {status.label}
            </span>
          </div>

          {/* Info Grid: Flat & Clean */}
          <div className="grid grid-cols-2 gap-4 rounded-xl border border-slate-50 bg-slate-50/50 p-4 sm:grid-cols-4">
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <Calendar size={12} className="text-[#397e50]" /> Tanggal
              </p>
              <p className="text-sm font-bold text-slate-700">{tgl_ujian}</p>
            </div>
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <Clock size={12} className="text-[#397e50]" /> Waktu
              </p>
              <p className="text-sm font-bold text-slate-700">{waktu_mulai}</p>
            </div>
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <Hash size={12} className="text-[#397e50]" /> Sesi
              </p>
              <p className="text-sm font-bold text-slate-700">
                {sesi_ujian ?? "-"}
              </p>
            </div>
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <MapPin size={12} className="text-[#397e50]" /> Ruangan
              </p>
              <p className="text-sm font-bold text-slate-700">
                {ruang_ujian ?? "-"}
              </p>
            </div>
          </div>

          {/* Action Footer: New Buttons Design */}
          {canControl && (
            <div className="mt-6 flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:items-center sm:justify-end">
              <button
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  onCancel?.(id);
                }}
                disabled={started === 0}
                className="cursor-pointer flex items-center justify-center gap-2 rounded-xl border-2 border-slate-200 bg-white px-5 py-2.5 text-xs font-bold uppercase tracking-widest text-slate-500 transition-all hover:border-rose-200 hover:bg-rose-50 hover:text-rose-600 disabled:opacity-30 disabled:hover:bg-transparent disabled:hover:border-slate-200 disabled:hover:text-slate-500"
              >
                <XCircle size={16} />
                Batalkan
              </button>

              <button
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  onStart?.(id);
                }}
                disabled={started === 1}
                className="cursor-pointer flex items-center justify-center gap-2 rounded-xl bg-[#397e50] px-6 py-2.5 text-xs font-bold uppercase tracking-widest text-white transition-all hover:bg-[#2d633f] hover:shadow-lg hover:shadow-emerald-900/20 disabled:bg-slate-200 disabled:text-slate-400 disabled:shadow-none"
              >
                <Play size={16} fill="currentColor" />
                Mulai Ujian
              </button>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
};

export default BoxJadwalUjian;
