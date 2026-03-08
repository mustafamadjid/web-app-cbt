import { useMemo } from "react";
import { Link, useParams } from "react-router";
import {
  AlertCircle,
  ArrowLeft,
  BookOpen,
  Calendar,
  GraduationCap,
  Hash,
  MapPin,
} from "lucide-react";

import PrintButton from "@/components/common/Input/PrintButton";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import { useGetDetailUjianById } from "@/services/Api/features-api/Ujian/ujian.service";

const statusLabelMap: Record<string, string> = {
  belum_dimulai: "Belum Dimulai",
  berlangsung: "Berlangsung",
  selesai: "Selesai",
  dibatalkan: "Dibatalkan",
};

const statusColorMap: Record<string, string> = {
  belum_dimulai: "bg-amber-50 text-amber-700 border-amber-200",
  berlangsung: "bg-emerald-50 text-emerald-700 border-emerald-200",
  selesai: "bg-slate-50 text-slate-600 border-slate-200",
  dibatalkan: "bg-rose-50 text-rose-700 border-rose-200",
};

const calculateDuration = (mulai: string, selesai: string) => {
  const [mulaiHour = 0, mulaiMinute = 0] = mulai.split(":").map(Number);
  const [selesaiHour = 0, selesaiMinute = 0] = selesai.split(":").map(Number);
  const mulaiMinuteTotal = mulaiHour * 60 + mulaiMinute;
  const selesaiMinuteTotal = selesaiHour * 60 + selesaiMinute;
  return Math.max(0, selesaiMinuteTotal - mulaiMinuteTotal);
};

const DetailUjian = () => {
  const params = useParams();
  const { user } = useAuth();
  const ujianId = Number(params.id);
  const isUjianIdValid = Number.isFinite(ujianId) && ujianId > 0;
  const backPath =
    user?.role === "GURU"
      ? paths.dashboard.jadwal_ujian_guru
      : paths.dashboard.jadwal_ujian;

  const {
    data: detail,
    loading,
    error,
  } = useGetDetailUjianById(ujianId, isUjianIdValid);

  const errorMsg = !isUjianIdValid ? "ID ujian tidak valid." : (error ?? "");

  const durasiMenit = useMemo(() => {
    if (!detail) return 0;
    return calculateDuration(detail.waktu_mulai, detail.waktu_selesai);
  }, [detail]);

  const tingkatLabel = useMemo(() => {
    if (!detail) return "-";
    return `Kelas ${detail.tingkat_kelas}`;
  }, [detail]);

  const namaKelasLabel = useMemo(() => {
    if (!detail) return "-";
    return detail.nama_kelas.trim() || "Semua kelas di tingkat ini";
  }, [detail]);

  const statusLabel = detail?.status_ujian
    ? (statusLabelMap[detail.status_ujian] ?? "Tidak Diketahui")
    : "Tidak Diketahui";

  const handlePrint = () => {
    if (typeof window !== "undefined") window.print();
  };

  return (
    <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-8 flex flex-col gap-4 border-b border-slate-100 pb-6 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <Link
            to={backPath}
            className="group mb-2 inline-flex items-center gap-2 text-sm font-medium text-slate-500 transition hover:text-[#397e50]"
          >
            <ArrowLeft
              size={16}
              className="transition-transform group-hover:-translate-x-1"
            />
            Kembali ke Jadwal
          </Link>
          <h1 className="text-3xl font-extrabold tracking-tight text-slate-900">
            {detail?.nama_ujian ?? "Memuat Detail..."}
          </h1>
          <span className="font-semibold text-[#397e50] ml-2">{tingkatLabel}</span>
        </div>

        {!loading && detail && (
          <div
            className={`inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-bold uppercase tracking-wider ${
              statusColorMap[detail.status_ujian] || statusColorMap.selesai
            }`}
          >
            <span
              className={`h-2.5 w-2.5 rounded-full ${
                detail.status_ujian === "berlangsung"
                  ? "animate-pulse bg-emerald-500"
                  : "bg-current"
              }`}
            />
            {statusLabel}
          </div>
        )}
      </div>

      {loading ? (
        <div className="flex min-h-[400px] flex-col items-center justify-center rounded-3xl border border-slate-200 bg-white p-12 text-slate-400">
          <div className="h-10 w-10 animate-spin rounded-full border-4 border-[#397e50] border-t-transparent" />
          <p className="mt-4 font-medium">Sinkronisasi data ujian...</p>
        </div>
      ) : errorMsg ? (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-red-100 bg-red-50 p-10 text-center">
          <AlertCircle size={40} className="mb-3 text-red-500" />
          <p className="text-lg font-semibold text-red-700">{errorMsg}</p>
          <Link
            to={backPath}
            className="mt-4 text-sm font-bold text-red-600 underline"
          >
            Kembali
          </Link>
        </div>
      ) : (
        detail && (
          <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
            <div className="space-y-6 lg:col-span-2">
              <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
                <div className="flex items-center gap-2 border-b border-slate-100 bg-slate-50/50 px-6 py-4 text-sm font-bold text-slate-800">
                  <BookOpen size={16} className="text-[#397e50]" />
                  Deskripsi Ujian
                </div>
                <div className="p-6">
                  <p className="text-sm leading-relaxed text-slate-600">
                    {detail.deskripsi_ujian ||
                      "Tidak ada deskripsi tambahan untuk ujian ini."}
                  </p>
                </div>
              </div>

              <div className="grid gap-6 sm:grid-cols-2">
                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 text-sm font-bold uppercase tracking-wide text-slate-800">
                    <Calendar size={16} className="text-[#397e50]" />
                    Waktu
                  </div>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Tanggal</span>
                      <span className="font-bold text-slate-700">
                        {detail.tgl_ujian}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Jam</span>
                      <span className="font-bold text-slate-700">
                        {detail.waktu_mulai} - {detail.waktu_selesai}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Durasi</span>
                      <span className="font-bold text-[#397e50]">
                        {durasiMenit} Menit
                      </span>
                    </div>
                  </div>
                </div>

                <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                  <div className="mb-4 flex items-center gap-2 border-b border-slate-50 pb-3 text-sm font-bold uppercase tracking-wide text-slate-800">
                    <MapPin size={16} className="text-[#397e50]" />
                    Lokasi
                  </div>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Ruangan</span>
                      <span className="font-bold text-[#397e50]">
                        {detail.ruang_ujian || "-"}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Sesi</span>
                      <span className="font-bold text-slate-700">
                        {detail.sesi_ujian > 0
                          ? `Sesi ${detail.sesi_ujian}`
                          : "-"}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-500">Pengawas</span>
                      <span className="font-bold text-slate-700">
                        {detail.pengawas_ujian || "-"}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="rounded-2xl p-5">
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <PrintButton
                    label="Daftar Hadir"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg bg-[#397e50] text-xs font-bold text-white transition hover:bg-[#2d633f]"
                    onClick={handlePrint}
                  />
                  <PrintButton
                    label="Berita Acara"
                    variant="outline"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg border border-slate-300 bg-white text-xs font-bold text-slate-600 transition hover:bg-slate-50"
                    onClick={handlePrint}
                  />
                  <PrintButton
                    label="Kartu Peserta"
                    variant="outline"
                    className="flex h-10 w-full items-center justify-center gap-2 rounded-lg border border-slate-300 bg-white text-xs font-bold text-slate-600 transition hover:bg-slate-50"
                    onClick={handlePrint}
                  />
                </div>
              </div>
            </div>

            <div className="space-y-6">
              <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                <h3 className="mb-5 text-2xs font-bold uppercase tracking-[0.15em] text-slate-400">
                  Target Akademik
                </h3>
                <div className="space-y-6">
                  <div className="flex items-start gap-3">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-[#397e50]">
                      <GraduationCap size={18} />
                    </div>
                    <div>
                      <p className="text-2xs font-bold uppercase text-slate-400">
                        Kelas
                      </p>
                      <p className="text-xs font-bold leading-tight text-slate-800">
                        {tingkatLabel}
                      </p>
                      <p className="text-2xs text-slate-500">
                        {namaKelasLabel}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <div className="rounded-2xl border border-emerald-100 bg-emerald-50/20 p-6">
                <h3 className="mb-5 text-2xs font-bold uppercase tracking-[0.15em] text-emerald-800/60">
                  Akses Ujian
                </h3>
                <div className="space-y-5">
                  <div className="rounded-xl border border-white bg-white p-5 text-center shadow-sm">
                    <p className="mb-2 text-[9px] font-bold uppercase tracking-widest text-slate-400">
                      Token Akses
                    </p>
                    <p className="text-2xl font-mono font-black tracking-[0.2em] text-[#397e50]">
                      {detail.token || "------"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )
      )}
    </div>
  );
};

export default DetailUjian;
