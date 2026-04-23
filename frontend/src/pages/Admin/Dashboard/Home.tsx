// Widgets
import SimpleStatWidget from "@/components/features/widget/Soal/StatistikWidgetSimpler";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import JadwalUjianWidget from "@/components/features/widget/JadwalUjian/JadwalUjian";
import LogAktivitasWidget from "@/components/features/widget/LogAktivitas/LogAktivitasWidget";
import TotalSiswaGuruWidget from "@/components/features/widget/StatistikPengguna/TotalSiswaGuruWidget";

import type { Role } from "@/types/Sidebar/SidebarMenu";
import { useAuth } from "@/contexts/AuthContext";
import { useGetLogAktivitasEnabled } from "@/services/Api/features-api/LogAktivitas/log-aktivitas.service";
import { useGetPengumumanActive } from "@/services/Api/features-api/pengumuman/pengumuman.service";
import { useGetDashboardStatistik } from "@/services/Api/features-api/Dashboard/dashboard.service";
import { useGetJadwalUjian } from "@/services/Api/features-api/Ujian/jadwalujian.service";

export const Home = () => {
  const { user } = useAuth();
  const role = (user?.role ?? "SISWA") as Role;
  const shouldFetchDashboardStatistik = role === "ADMIN" || role === "GURU";

  // Hook: fetch log aktivitas (only for ADMIN)
  const { data: aktivitasItems } = useGetLogAktivitasEnabled(role === "ADMIN");

  // Hook: fetch pengumuman aktif
  const { data: pengumumanItems } = useGetPengumumanActive();
  const {
    data: jadwalUjianItems,
    loading: jadwalUjianLoading,
    error: jadwalUjianError,
  } = useGetJadwalUjian({
    kategoriUjian: "belum_dimulai",
  });
  const {
    data: dashboardStatistik,
    loading: dashboardStatistikLoading,
    error: dashboardStatistikError,
  } = useGetDashboardStatistik(
    shouldFetchDashboardStatistik ? role : "ADMIN",
    shouldFetchDashboardStatistik,
  );

  const statistik = dashboardStatistik ?? {
    total_siswa: 0,
    total_guru: 0,
    total_ujian_terlaksana: 0,
    total_bank_soal: 0,
    total_mapel_aktif: 0,
  };

  return (
    <div className="min-h-screen bg-[#ecf1ed]  pb-20">
      <div className="mx-auto max-w-[1920px] space-y-8 p-4 sm:p-6 lg:p-8">
        {(dashboardStatistikLoading || dashboardStatistikError) && (
          <div
            className={`rounded-xl border px-4 py-3 text-sm ${
              dashboardStatistikError
                ? "border-amber-200 bg-amber-50 text-amber-800"
                : "border-emerald-200 bg-emerald-50 text-emerald-700"
            }`}
          >
            {dashboardStatistikError
              ? `Statistik dashboard belum bisa dimuat: ${dashboardStatistikError}`
              : "Menyegarkan statistik dashboard..."}
          </div>
        )}

        {(jadwalUjianLoading || jadwalUjianError) && (
          <div
            className={`rounded-xl border px-4 py-3 text-sm ${
              jadwalUjianError
                ? "border-amber-200 bg-amber-50 text-amber-800"
                : "border-emerald-200 bg-emerald-50 text-emerald-700"
            }`}
          >
            {jadwalUjianError
              ? `Jadwal ujian belum bisa dimuat: ${jadwalUjianError}`
              : "Memuat jadwal ujian mendatang..."}
          </div>
        )}

        {/* === TOP STATS (GRID 4 KOLOM) === */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-5">
          {/* Card 1: Total User (Dibuat Lebih Lebar) */}
          {role === "ADMIN" && (
            <div className="h-full sm:col-span-2 lg:col-span-2">
              <TotalSiswaGuruWidget
                totalGuru={statistik.total_guru}
                totalSiswa={statistik.total_siswa}
                className="h-full"
              />
            </div>
          )}

          {/* Card 2: Statistik Ujian */}
          <div
            className={`h-full  ${
              role === "ADMIN" ? "lg:col-span-1" : "lg:col-span-3"
            }`}
          >
            <SimpleStatWidget
              title="Total Ujian Terlaksana"
              value={statistik.total_ujian_terlaksana}
              trend="neutral"
              trendText="Total ujian yang sudah terlaksana"
              className="h-full"
            />
          </div>

          {/* Card 3: Bank Soal */}
          <div className="h-full lg:col-span-1">
            <SimpleStatWidget
              title="Total Bank Soal"
              value={statistik.total_bank_soal}
              trend="neutral"
              trendText="Total bank soal tersedia"
              className="h-full"
            />
          </div>

          {/* Card 4: Mata Pelajaran */}
          <div className="h-full lg:col-span-1">
            <SimpleStatWidget
              title="Mata Pelajaran Aktif"
              value={statistik.total_mapel_aktif}
              trend="neutral"
              trendText="Total mata pelajaran aktif"
              className="h-full"
            />
          </div>
        </div>

        {/* === MAIN CONTENT === */}
        <div className="grid gap-6 lg:grid-cols-12 items-stretch">
          {/* LEFT COLUMN */}
          <div
            className={`flex flex-col gap-6 xl:col-span-8 h-full min-h-0 lg:col-span-6 ${
              role === "ADMIN" ? "" : "lg:h-300"
            }`}
          >
            <PengumumanWidget
              title="Papan Pengumuman"
              items={pengumumanItems ?? []}
              className="flex-1 min-h-0"
            />
          </div>

          {/* RIGHT COLUMN */}
          <div className="flex flex-col gap-6 lg:col-span-5 xl:col-span-4 h-full min-h-0">
            <JadwalUjianWidget
              items={jadwalUjianItems ?? []}
              className="flex-1 min-h-0"
              maxHeightClassName="h-full"
            />

            {role === "ADMIN" && (
              <LogAktivitasWidget
                items={aktivitasItems ?? []}
                className="flex-1 min-h-0"
                maxHeightClassName="max-h-[400px] overflow-y-auto "
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
