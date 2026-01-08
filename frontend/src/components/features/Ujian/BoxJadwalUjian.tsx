import type { JadwalUjianItem } from "../../../types/Ujian/jadwalUjian";

const statusLabels: Record<string, string> = {
  belum_dimulai: "Belum dimulai",
  berlangsung: "Berlangsung",
  selesai: "Selesai",
};

const formatStatus = (status?: string) => {
  if (!status) return "Status belum tersedia";
  return statusLabels[status] ?? "Status belum tersedia";
};

type BoxJadwalUjianProps = JadwalUjianItem;

export const BoxJadwalUjian = ({
  nama_ujian,
  pengawas_ujian,
  tgl_ujian,
  waktu_mulai,
  sesi_ujian,
  ruang_ujian,
  status_ujian,
}: BoxJadwalUjianProps) => {
  return (
    <div className="flex flex-col gap-4 rounded-xl border border-[#397e50] bg-white p-5 shadow-sm">
      <div className="flex flex-wrap items-center justify-between gap-3 rounded-lg bg-[#397e50] px-4 py-2 text-white">
        <h3 className="text-lg font-semibold">{nama_ujian}</h3>
        <span className="rounded-full border border-white/60 px-3 py-1 text-xs font-semibold uppercase tracking-wide">
          {formatStatus(status_ujian)}
        </span>
      </div>

      <div className="grid gap-3 text-sm text-slate-700 sm:grid-cols-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-[#397e50]">Pengawas</p>
          <p className="font-medium">{pengawas_ujian}</p>
        </div>
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-[#397e50]">Tanggal</p>
          <p className="font-medium">{tgl_ujian}</p>
        </div>
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-[#397e50]">Waktu mulai</p>
          <p className="font-medium">{waktu_mulai}</p>
        </div>
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-[#397e50]">Sesi</p>
          <p className="font-medium">{sesi_ujian ?? "-"}</p>
        </div>
        <div>
          <p className="text-xs font-semibold uppercase tracking-wide text-[#397e50]">Ruang</p>
          <p className="font-medium">{ruang_ujian ?? "-"}</p>
        </div>
      </div>
    </div>
  );
};
