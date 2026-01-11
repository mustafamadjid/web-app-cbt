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
          <span className="inline-flex items-center gap-2 rounded-full border border-[#397e50]/20 bg-[#397e50]/10  px-3 py-1 text-2xs font-semibold uppercase tracking-wide text-[#397e50]">
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

// import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
// import {
//   Calendar,
//   Clock,
//   MapPin,
//   User,
//   Layout,
// } from "lucide-react";
// import type { ReactNode } from "react";

// type BoxCetakUjianProps = JadwalUjianItem & {
//   actions?: ReactNode;
// };

// export const BoxCetakUjian = ({
//   nama_ujian,
//   pengawas_ujian,
//   tgl_ujian,
//   waktu_mulai,
//   ruang_ujian,
//   tingkat_kelas,
//   nama_kelas,
//   actions,
// }: BoxCetakUjianProps) => {
//   return (
//     <div className="relative overflow-hidden rounded-xl border-2 border-slate-100 bg-white transition-all duration-200 hover:border-[#397e50] hover:shadow-lg">
//       {/* Indikator Samping (Solid Color) */}
//       <div className="absolute left-0 top-0 h-full w-1.5 bg-[#397e50]" />

//       <div className="p-6">
//         {/* Header: Badge & Title */}
//         <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
//           <div className="space-y-1">
//             <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-[#397e50]">
//               <div className="flex h-6 w-6 items-center justify-center rounded bg-[#397e50] text-white">
//                 <Layout size={14} />
//               </div>
//               <span>
//                 Kelas {tingkat_kelas ?? "-"}{" "}
//                 {nama_kelas ? ` • ${nama_kelas}` : ""}
//               </span>
//             </div>
//             <h3 className="text-xl font-bold text-slate-900 leading-tight">
//               {nama_ujian}
//             </h3>
//           </div>

//           <div className="flex items-center gap-2 rounded-lg border border-slate-100 bg-slate-50 px-3 py-1.5">
//             <User size={14} className="text-slate-400" />
//             <span className="text-xs font-medium text-slate-600">
//               {pengawas_ujian}
//             </span>
//           </div>
//         </div>

//         {/* Info Grid: Solid Background */}
//         <div className="mt-6 flex flex-wrap items-center rounded-xl border border-slate-200 bg-white py-3">
//           {/* Tanggal */}
//           <div className="flex flex-1 items-center gap-3 px-4 min-w-[150px]">
//             <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-[#397e50]">
//               <Calendar size={18} />
//             </div>
//             <div className="flex flex-col">
//               <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
//                 Tanggal
//               </span>
//               <span className="text-sm font-bold text-slate-700">
//                 {tgl_ujian}
//               </span>
//             </div>
//           </div>

//           {/* Divider - Hilang di mobile */}
//           <div className="hidden h-8 w-px bg-slate-200 sm:block" />

//           {/* Waktu */}
//           <div className="flex flex-1 items-center gap-3 px-4 min-w-[150px]">
//             <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-[#397e50]">
//               <Clock size={18} />
//             </div>
//             <div className="flex flex-col">
//               <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
//                 Waktu
//               </span>
//               <span className="text-sm font-bold text-slate-700">
//                 {waktu_mulai}
//               </span>
//             </div>
//           </div>

//           {/* Divider - Hilang di mobile */}
//           <div className="hidden h-8 w-px bg-slate-200 sm:block" />

//           {/* Ruangan */}
//           <div className="flex flex-1 items-center gap-3 px-4 min-w-[150px]">
//             <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-[#397e50]">
//               <MapPin size={18} />
//             </div>
//             <div className="flex flex-col">
//               <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
//                 Ruangan
//               </span>
//               <span className="text-sm font-bold text-[#397e50]">
//                 {ruang_ujian ?? "-"}
//               </span>
//             </div>
//           </div>
//         </div>

//         {/* Footer Actions */}
//         {actions && (
//           <div className="mt-5 flex items-center justify-end gap-3 border-t border-slate-100 pt-4">
//             {actions}
//           </div>
//         )}
//       </div>
//     </div>
//   );
// };
