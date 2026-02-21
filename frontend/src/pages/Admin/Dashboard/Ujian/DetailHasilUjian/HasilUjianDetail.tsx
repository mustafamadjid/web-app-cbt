import { useEffect, useMemo, useRef, useState } from "react";
import { Link, useParams } from "react-router";
import {
  ArrowLeft,
  Users,
  Trophy,
  TrendingUp,
  TrendingDown,
  CheckCircle2,
  XCircle,
  MinusCircle,
  ChevronRight,
  GraduationCap,
  Hash,
  Calendar,
  Medal,
} from "lucide-react";
import {
  getHasilUjianDetail,
  type HasilUjianDetailResponse,
} from "@/services/Api/features-api/Ujian/hasilUjian.service";
import type { HasilUjianSiswa } from "@/types/Ujian/HasilUjian";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { paths } from "@/routes/paths";
import { useAuth } from "@/contexts/AuthContext";

// --- KOMPONEN WIDGET FLAT (NO SHADOW, NO BG DECOR) ---
const StatWidget = ({
  title,
  value,
  icon: Icon,
  colorTheme, // 'emerald' | 'amber' | 'blue' | 'rose'
}: {
  title: string;
  value: number | string;
  icon: any;
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
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [statistik, setStatistik] = useState<
    HasilUjianDetailResponse["statistik"] | null
  >(null);
  const [daftarSiswa, setDaftarSiswa] = useState<HasilUjianSiswa[]>([]);

  const requestSeq = useRef(0);
  const ujianId = useMemo(() => Number(id), [id]);

  useEffect(() => {
    if (!ujianId || Number.isNaN(ujianId)) return;
    const seq = ++requestSeq.current;
    (async () => {
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getHasilUjianDetail(ujianId);
        if (seq !== requestSeq.current) return;
        setStatistik(data.statistik);
        setDaftarSiswa(data.siswa);
      } catch {
        if (seq !== requestSeq.current) return;
        setErrorMsg("Gagal memuat detail hasil ujian.");
        setStatistik(null);
        setDaftarSiswa([]);
      } finally {
        if (seq !== requestSeq.current) return;
        setLoading(false);
      }
    })();
  }, [ujianId]);

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

  return (
    <div className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      {/* --- HEADER --- */}
      <div className="flex flex-col gap-2">
        <Link
          to={paths.dashboard.hasil_ujian}
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
              Analisis statistik dan peringkat siswa.
            </p>
          </div>
        </div>
      </div>

      {errorMsg && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm font-medium text-rose-700">
          {errorMsg}
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
            Daftar Peserta & Nilai
          </h2>
          <p className="text-sm text-slate-500">
            Total {daftarSiswa.length} siswa.
          </p>
        </div>

        {/* CONTAINER SCROLLABLE - Max 600px height */}
        <div className="relative max-h-[600px] overflow-y-auto scrollbar-thin scrollbar-track-slate-50 scrollbar-thumb-slate-200">
          <table className="min-w-full divide-y divide-slate-100 text-left">
            <thead className="sticky top-0 z-10 bg-slate-50 shadow-sm">
              <tr>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Siswa
                </th>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Info Akademik
                </th>
                <th className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-slate-500">
                  Analisis Jawaban
                </th>
                <th className="px-6 py-4 text-center text-xs font-bold uppercase tracking-wider text-slate-500">
                  Nilai
                </th>
                <th className="px-6 py-4 text-right text-xs font-bold uppercase tracking-wider text-slate-500">
                  Opsi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 bg-white">
              {loading ? (
                <tr>
                  <td colSpan={5} className="py-20 text-center text-slate-500">
                    <div className="mx-auto h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-[#397e50]" />
                    <p className="mt-2 text-sm">Memuat data...</p>
                  </td>
                </tr>
              ) : daftarSiswa.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-slate-500">
                    Belum ada data siswa.
                  </td>
                </tr>
              ) : (
                daftarSiswa.map((siswa) => (
                  <tr
                    key={siswa.id_pengguna}
                    className="group transition-colors hover:bg-slate-50"
                  >
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <img
                          src={
                            resolveImageUrl(siswa.foto_profil) ||
                            `https://ui-avatars.com/api/?name=${siswa.nama_lengkap}&background=random`
                          }
                          alt={siswa.nama_lengkap}
                          className="h-9 w-9 rounded-full bg-slate-200 object-cover ring-2 ring-white"
                        />
                        <div>
                          <p className="font-semibold text-slate-800">
                            {siswa.nama_lengkap}
                          </p>
                          <p className="text-xs text-slate-500 capitalize">
                            {siswa.role}
                          </p>
                        </div>
                      </div>
                    </td>

                    <td className="px-6 py-4">
                      <div className="flex flex-col gap-1 text-sm text-slate-600">
                        <div className="flex items-center gap-1.5 font-medium">
                          <GraduationCap size={14} className="text-slate-400" />
                          {siswa.tingkat_kelas} - {siswa.nama_kelas}
                        </div>
                        <div className="flex items-center gap-3 text-xs text-slate-400">
                          <span className="flex items-center gap-1">
                            <Hash size={12} /> {siswa.no_absen}
                          </span>
                          <span className="flex items-center gap-1">
                            <Calendar size={12} /> {siswa.angkatan}
                          </span>
                        </div>
                      </div>
                    </td>

                    <td className="px-6 py-4">
                      <div className="flex w-fit items-center gap-3 rounded-lg border border-slate-100 bg-slate-50 px-3 py-1.5">
                        <div
                          className="flex items-center gap-1 text-emerald-600"
                          title="Benar"
                        >
                          <CheckCircle2 size={14} />
                          <span className="text-sm font-semibold">
                            {siswa.jumlah_benar}
                          </span>
                        </div>
                        <div className="h-4 w-px bg-slate-200"></div>
                        <div
                          className="flex items-center gap-1 text-rose-500"
                          title="Salah"
                        >
                          <XCircle size={14} />
                          <span className="text-sm font-semibold">
                            {siswa.jumlah_salah}
                          </span>
                        </div>
                        <div className="h-4 w-px bg-slate-200"></div>
                        <div
                          className="flex items-center gap-1 text-slate-400"
                          title="Kosong"
                        >
                          <MinusCircle size={14} />
                          <span className="text-sm font-semibold">
                            {siswa.jumlah_kosong}
                          </span>
                        </div>
                      </div>
                    </td>

                    <td className="px-6 py-4 text-center">
                      <span
                        className={`inline-block min-w-12 rounded-md px-2 py-1 text-center font-bold ${
                          Number(siswa.nilai) >= 75
                            ? "bg-emerald-100 text-emerald-700"
                            : "bg-slate-100 text-slate-700"
                        }`}
                      >
                        {siswa.nilai ?? 0}
                      </span>
                    </td>

                    <td className="px-6 py-4 text-right">
                      <Link
                        to={
                          user?.role === "ADMIN"
                            ? paths.dashboard.hasil_ujian_detail_admin
                            : paths.dashboard.hasil_ujian_detail_admin_guru
                                .replace(":id", String(ujianId))
                                .replace(":siswaId", String(siswa.id_pengguna))
                        }
                        className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-200 text-slate-400 transition-colors hover:border-[#397e50] hover:bg-[#397e50] hover:text-white"
                        title="Lihat Detail"
                      >
                        <ChevronRight size={16} />
                      </Link>
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
