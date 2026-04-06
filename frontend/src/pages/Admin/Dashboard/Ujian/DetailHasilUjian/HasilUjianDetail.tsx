import { useMemo } from "react";
import { Link, useParams } from "react-router";
import {
  ArrowLeft,
  Users,
  Trophy,
  TrendingUp,
  TrendingDown,
  GraduationCap,
  Hash,
  Medal,
  type LucideIcon,
} from "lucide-react";
import {
  useGetHasilUjianDetail,
} from "@/services/Api/features-api/Ujian/hasilUjian.service";
import { useGetListPesertaUjianSubmitted } from "@/services/Api/features-api/Ujian/attemptUjian.service";
import type { PesertaUjianSubmittedItem } from "@/types/Ujian/AttemptUjian";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";

// --- KOMPONEN WIDGET FLAT (NO SHADOW, NO BG DECOR) ---
const StatWidget = ({
  title,
  value,
  icon: Icon,
  colorTheme, // 'emerald' | 'amber' | 'blue' | 'rose'
}: {
  title: string;
  value: number | string;
  icon: LucideIcon;
  colorTheme: string;
}) => {
  // Mapping warna simpel (Background soft + Text color)
  const colors: Record<string, string> = {
    emerald: "bg-emerald-50 text-emerald-600",
    amber: "bg-amber-50 text-amber-600",
    blue: "bg-blue-50 text-blue-600",
    rose: "bg-rose-50 text-rose-600",
  };

  const themeClass = colors[colorTheme] || colors.blue;

  return (
    <div className="flex items-center justify-between rounded-2xl border border-slate-200 bg-white p-6">
      <div>
        <p className="text-sm font-medium text-slate-500">{title}</p>
        <h3 className="mt-2 text-2xl font-bold text-slate-800">{value}</h3>
      </div>
      <div
        className={`flex h-12 w-12 items-center justify-center rounded-xl ${themeClass}`}
      >
        <Icon size={24} strokeWidth={2.5} />
      </div>
    </div>
  );
};

const HasilUjianDetail = () => {
  const { user } = useAuth();
  const { id } = useParams();
  const jadwalUjianId = useMemo(() => Number(id), [id]);
  const isJadwalUjianIdValid =
    Number.isInteger(jadwalUjianId) && jadwalUjianId > 0;
  const {
    data: hasilDetail,
    error: statistikError,
  } = useGetHasilUjianDetail(jadwalUjianId, isJadwalUjianIdValid);
  const {
    data: pesertaSubmitted,
    loading: pesertaLoading,
    error: pesertaError,
  } = useGetListPesertaUjianSubmitted(
    jadwalUjianId,
    isJadwalUjianIdValid,
  );

  const statistik = hasilDetail?.statistik ?? null;
  const daftarPeserta: PesertaUjianSubmittedItem[] = pesertaSubmitted ?? [];
  const invalidIdMessage = !isJadwalUjianIdValid
    ? "ID jadwal ujian tidak valid."
    : "";

    console.log(jadwalUjianId);

  const formatDateTime = (value: string | null) => {
    if (!value) return "-";

    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return "-";

    return parsed.toLocaleString("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatNilaiAkhir = (value: number | null) => {
    if (typeof value !== "number") return "-";

    return value.toLocaleString("id-ID", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  };

  const widgetData = [
    {
      title: "Nilai Tertinggi",
      value: statistik?.nilai_tertinggi ?? 0,
      icon: Trophy,
      colorTheme: "amber",
    },
    {
      title: "Rata-rata Kelas",
      value: statistik?.rata_rata ?? 0,
      icon: TrendingUp,
      colorTheme: "blue",
    },
    {
      title: "Nilai Terendah",
      value: statistik?.nilai_terendah ?? 0,
      icon: TrendingDown,
      colorTheme: "rose",
    },
    {
      title: "Total Peserta",
      value: statistik?.jumlah_peserta ?? 0,
      icon: Users,
      colorTheme: "emerald",
    },
  ];
  const backPath =
    user?.role === "GURU"
      ? paths.dashboard.hasil_ujian_guru
      : paths.dashboard.hasil_ujian;

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      {/* --- HEADER --- */}
      <div className="flex flex-col gap-2">
        <Link
          to={backPath}
          className="group inline-flex items-center gap-2 text-sm font-medium text-slate-500 transition-colors hover:text-[#397e50]"
        >
          <ArrowLeft
            size={16}
            className="transition-transform group-hover:-translate-x-1"
          />
          Kembali ke Daftar Hasil
        </Link>
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-[#397e50]/10 text-[#397e50]">
            <Medal size={20} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">
              Detail Hasil Ujian
            </h1>
            <p className="text-sm text-slate-500">
              Analisis statistik dan daftar peserta yang telah submit ujian.
            </p>
          </div>
        </div>
      </div>

      {invalidIdMessage && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm font-medium text-rose-700">
          {invalidIdMessage}
        </div>
      )}

      {statistikError && (
        <div className="rounded-xl border border-amber-200 bg-amber-50 p-4 text-sm font-medium text-amber-700">
          Gagal memuat statistik ujian: {statistikError}
        </div>
      )}

      {pesertaError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm font-medium text-rose-700">
          Gagal memuat daftar peserta: {pesertaError}
        </div>
      )}

      {/* --- WIDGET GRID (FLAT STYLE) --- */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {widgetData.map((widget, idx) => (
          <StatWidget
            key={idx}
            title={widget.title}
            value={widget.value}
            icon={widget.icon}
            colorTheme={widget.colorTheme}
          />
        ))}
      </div>

      {/* --- TABLE SECTION (SCROLLABLE) --- */}
      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5 bg-white">
          <h2 className="text-lg font-bold text-slate-800">
            Daftar Peserta Submit
          </h2>
          <p className="text-sm text-slate-500">
            Total {daftarPeserta.length} peserta.
          </p>
        </div>

        {/* CONTAINER SCROLLABLE - Max 600px height */}
        <div className="relative max-h-[600px] overflow-y-auto scrollbar-thin scrollbar-track-slate-50 scrollbar-thumb-slate-200">
          <table className="min-w-full divide-y divide-slate-100 text-left">
            <thead className="sticky top-0 z-10 bg-slate-50 shadow-sm">
              <tr>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Peserta
                </th>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Kelas & Absen
                </th>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Waktu Mulai
                </th>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Waktu Submit
                </th>
                <th className="px-6 py-4 text-center text-xs font-bold uppercase tracking-wider text-slate-500">
                  Nilai Akhir
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 bg-white">
              {pesertaLoading ? (
                <tr>
                  <td colSpan={5} className="py-20 text-center text-slate-500">
                    <div className="mx-auto h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-[#397e50]" />
                    <p className="mt-2 text-sm">Memuat data...</p>
                  </td>
                </tr>
              ) : !isJadwalUjianIdValid ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-slate-500">
                    ID jadwal ujian tidak valid.
                  </td>
                </tr>
              ) : daftarPeserta.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-slate-500">
                    Belum ada peserta yang submit attempt.
                  </td>
                </tr>
              ) : (
                daftarPeserta.map((peserta) => (
                  <tr
                    key={peserta.id_attempt}
                    className="group transition-colors hover:bg-slate-50"
                  >
                    <td className="px-6 py-4">
                      <div>
                        <p className="font-semibold text-slate-800">
                          {peserta.nama_lengkap}
                        </p>
                        <p className="text-xs text-slate-500">
                          ID Siswa: {peserta.id_siswa}
                        </p>
                      </div>
                    </td>

                    <td className="px-6 py-4">
                      <div className="flex flex-col gap-1 text-sm text-slate-600">
                        <div className="flex items-center gap-1.5 font-medium">
                          <GraduationCap size={14} className="text-slate-400" />
                          Kelas {peserta.tingkat_kelas} - {peserta.nama_kelas}
                        </div>
                        <div className="flex items-center gap-1 text-xs text-slate-400">
                          <span className="flex items-center gap-1">
                            <Hash size={12} /> No. {peserta.no_absen}
                          </span>
                        </div>
                      </div>
                    </td>

                    <td className="px-6 py-4 text-sm text-slate-600">
                      {formatDateTime(peserta.waktu_mulai)}
                    </td>

                    <td className="px-6 py-4 text-sm text-slate-600">
                      {formatDateTime(peserta.waktu_submit)}
                    </td>

                    <td className="px-6 py-4 text-center">
                      <span
                        className={`inline-block min-w-16 rounded-md px-2 py-1 text-center font-bold ${
                          typeof peserta.nilai_akhir === "number" &&
                          peserta.nilai_akhir >= 75
                            ? "bg-emerald-100 text-emerald-700"
                            : typeof peserta.nilai_akhir === "number"
                              ? "bg-slate-100 text-slate-700"
                              : "bg-slate-50 text-slate-400"
                        }`}
                      >
                        {formatNilaiAkhir(peserta.nilai_akhir)}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default HasilUjianDetail;
