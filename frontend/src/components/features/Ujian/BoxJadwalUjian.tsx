import type { JadwalUjianItem } from "../../../types/Ujian/jadwalUjian";
import {
  Calendar,
  Clock,
  MapPin,
  User,
  Hash,
  GraduationCap,
} from "lucide-react";
import { Link } from "react-router";

const statusConfig: Record<string, { label: string; color: string }> = {
  belum_dimulai: {
    label: "Belum Dimulai",
    color: "bg-amber-100 text-amber-700 border-amber-200",
  },
  berlangsung: {
    label: "Berlangsung",
    color: "bg-emerald-100 text-emerald-700 border-emerald-200",
  },
  selesai: {
    label: "Selesai",
    color: "bg-slate-100 text-slate-600 border-slate-200",
  },
};

type BoxJadwalUjianProps = JadwalUjianItem & {
  linkJadwal?: string;
};

export const BoxJadwalUjian = ({
  nama_ujian,
  pengawas_ujian,
  tgl_ujian,
  waktu_mulai,
  sesi_ujian,
  ruang_ujian,
  status_ujian,
  tingkat_kelas,
  nama_kelas,
  linkJadwal = "",
}: BoxJadwalUjianProps) => {
  const status = status_ujian
    ? statusConfig[status_ujian]
    : { label: "Unknown", color: "bg-gray-100 text-gray-500" };

  return (
    <Link to={linkJadwal} className="block">
      <div className="group relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 transition-all hover:border-[#397e50]/50 hover:shadow-lg">
        {/* Aksen Garis Samping Flat */}
        <div className="absolute left-0 top-0 h-full w-1.5 bg-[#397e50]" />

        {/* Header Section */}
        <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
          <div className="space-y-2">
            {/* Badge Kelas */}
            <div className="flex items-center gap-2">
              <span className="inline-flex items-center gap-1 rounded-md bg-slate-100 px-2 py-0.5 text-2xs font-bold uppercase tracking-wider text-slate-600 border border-slate-200">
                <GraduationCap size={12} />
                Kelas {tingkat_kelas} - {nama_kelas}
              </span>
            </div>

            <h3 className="text-xl font-bold text-slate-800 leading-tight group-hover:text-[#397e50] transition-colors">
              {nama_ujian}
            </h3>

            <div className="flex items-center gap-2 text-slate-500">
              <User size={14} className="text-slate-400" />
              <span className="text-sm font-medium">{pengawas_ujian}</span>
            </div>
          </div>

          <span
            className={`rounded-full border px-3 py-1 text-[11px] font-bold uppercase tracking-wider ${status.color}`}
          >
            {status.label}
          </span>
        </div>

        {/* Info Grid Section */}
        <div className="grid grid-cols-2 gap-y-5 gap-x-4 border-t border-slate-100 pt-5 sm:grid-cols-4">
          <div className="space-y-1">
            <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
              <Calendar size={12} className="text-[#397e50]" /> Tanggal
            </p>
            <p className="text-sm font-semibold text-slate-700">{tgl_ujian}</p>
          </div>

          <div className="space-y-1">
            <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
              <Clock size={12} className="text-[#397e50]" /> Waktu
            </p>
            <p className="text-sm font-semibold text-slate-700">
              {waktu_mulai}
            </p>
          </div>

          <div className="space-y-1">
            <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
              <Hash size={12} className="text-[#397e50]" /> Sesi
            </p>
            <p className="text-sm font-semibold text-slate-700">
              {sesi_ujian ?? "-"}
            </p>
          </div>

          <div className="space-y-1">
            <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
              <MapPin size={12} className="text-[#397e50]" /> Ruang
            </p>
            <p className="text-sm font-semibold text-slate-700">
              {ruang_ujian ?? "-"}
            </p>
          </div>
        </div>
      </div>
    </Link>
  );
};
