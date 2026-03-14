import type { JadwalUjianItem } from "../../../types/Ujian/jadwalUjian";
import {
  Calendar,
  Clock,
  GraduationCap,
  Hash,
  MapPin,
  Trash2,
  User,
} from "lucide-react";
import { Link } from "react-router";


type BoxJadwalUjianProps = JadwalUjianItem & {
  linkJadwal?: string;
  onDelete?: (idUjian: number) => void;
  canControl?: boolean;
  deleting?: boolean;
};

const BoxJadwalUjian = ({
  nama_ujian,
  pengawas_ujian,
  tgl_ujian,
  waktu_mulai,
  sesi_ujian,
  ruang_ujian,
  tingkat_kelas,
  nama_kelas,
  id_ujian,
  linkJadwal = "",
  onDelete,
  canControl = false,
  deleting = false,
}: BoxJadwalUjianProps) => {
  const idUjian = id_ujian ?? 0;
  const disableDelete = deleting || idUjian <= 0;

  return (
    <Link to={linkJadwal} className="group block">
      <div className="relative overflow-hidden rounded-2xl border-2 border-slate-100 bg-white transition-all duration-300 hover:border-[#397e50]/30 hover:shadow-xl hover:shadow-slate-200/50">
        <div className="absolute left-0 top-0 h-full w-1.5 bg-[#397e50]" />

        <div className="p-6">
          <div className="mb-6 flex items-start justify-between gap-4">
            <div className="space-y-2">
              <span className="inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-slate-100 px-2.5 py-1 text-2xs font-bold uppercase tracking-wider text-slate-600">
                <GraduationCap size={14} />
                Kelas {tingkat_kelas ?? "-"} - {nama_kelas ?? "-"}
              </span>
              <h3 className="text-xl font-extrabold leading-tight text-slate-800 transition-colors group-hover:text-[#397e50]">
                {nama_ujian}
              </h3>
            </div>

            
          </div>

          <div className="grid grid-cols-2 gap-4 rounded-xl border border-slate-50 bg-slate-50/50 p-4 sm:grid-cols-5">
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <Calendar size={12} className="text-[#397e50]" /> Tanggal
              </p>
              <p className="text-sm font-bold text-slate-700">{tgl_ujian}</p>
            </div>
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <User size={12} className="text-[#397e50]" /> Pengawas Ujian
              </p>
              <p className="text-sm font-bold text-slate-700">{pengawas_ujian}</p>
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
              <p className="text-sm font-bold text-slate-700">{sesi_ujian ?? "-"}</p>
            </div>
            <div className="space-y-1">
              <p className="flex items-center gap-1.5 text-2xs font-bold uppercase tracking-wider text-slate-400">
                <MapPin size={12} className="text-[#397e50]" /> Ruangan
              </p>
              <p className="text-sm font-bold text-slate-700">{ruang_ujian ?? "-"}</p>
            </div>
          </div>

          {canControl && (
            <div className="mt-6 flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:flex-wrap sm:items-center sm:justify-end">
              <button
                type="button"
                onClick={(event) => {
                  event.preventDefault();
                  event.stopPropagation();
                  if (idUjian > 0) onDelete?.(idUjian);
                }}
                disabled={disableDelete}
                className="flex cursor-pointer items-center justify-center gap-2 rounded-xl border-2 border-slate-200 bg-white px-5 py-2.5 text-xs font-bold uppercase tracking-widest text-slate-500 transition-all hover:border-rose-200 hover:bg-rose-50 hover:text-rose-600 disabled:opacity-30 disabled:hover:border-slate-200 disabled:hover:bg-transparent disabled:hover:text-slate-500"
              >
                <Trash2 size={16} />
                Hapus
              </button>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
};

export default BoxJadwalUjian;
