import React from "react";
import { useAuth } from "@/contexts/AuthContext";
import UjianTerlaksanaWidget from "@/components/features/widget/Siswa/UjianTerlaksanaWidget";
import RataRataNilaiWidget from "@/components/features/widget/Siswa/RataRataNilaiWidget";
import JadwalUjianSiswaWidget from "@/components/features/widget/JadwalUjian/JadwalUjianSiswaWidget";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import type { PengumumanItem } from "@/types/Widget/Pengumuman";
import type { JadwalUjianItem } from "@/types/Ujian/jadwalUjian";
import type {
  SiswaDashboardSummary,
  SiswaSemesterAverage,
} from "@/types/Widget/SiswaDashboard";
import {
  getSiswaDashboardSummary,
  getSiswaJadwalUjian,
  getSiswaPengumuman,
  getSiswaRataRataSemester,
} from "@/services/Api/features-api/Siswa/homeSiswa.service";

const HomeSiswa = () => {
  const { user } = useAuth();
  const [summary, setSummary] = React.useState<SiswaDashboardSummary | null>(
    null
  );
  const [jadwal, setJadwal] = React.useState<JadwalUjianItem[]>([]);
  const [pengumuman, setPengumuman] = React.useState<PengumumanItem[]>([]);
  const [rataRata, setRataRata] = React.useState<SiswaSemesterAverage[]>([]);
  const [loading, setLoading] = React.useState(true);

  React.useEffect(() => {
    let mounted = true;
    const loadData = async () => {
      setLoading(true);
      try {
        const [summaryData, jadwalData, pengumumanData, rataRataData] =
          await Promise.all([
            getSiswaDashboardSummary(),
            getSiswaJadwalUjian(),
            getSiswaPengumuman(),
            getSiswaRataRataSemester(),
          ]);
        if (!mounted) return;
        setSummary(summaryData);
        setJadwal(jadwalData);
        setPengumuman(pengumumanData);
        setRataRata(rataRataData);
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    void loadData();
    return () => {
      mounted = false;
    };
  }, []);

  return (
    <div className="min-h-screen bg-[#ecf1ed] pb-20">
      <div className="mx-auto flex max-w-[1920px] flex-col gap-6 p-4 sm:p-6 lg:p-8">
        <header className="flex flex-col gap-2">
          <p className="text-sm font-semibold text-[#397e50]">
            Selamat datang kembali
          </p>
          <h1 className="text-2xl font-bold text-[#37513d]">
            Halo, {user?.username ?? "Siswa"}
          </h1>
          <p className="text-sm text-gray-600">
            Pantau jadwal ujian, pengumuman terbaru, dan progress belajarmu di
            sini.
          </p>
        </header>

        <div className="grid gap-5 lg:grid-cols-12">
          <div className="lg:col-span-6">
            <UjianTerlaksanaWidget
              totalSelesai={summary?.ujian_selesai ?? 0}
              totalUjian={summary?.total_ujian ?? 0}
              className="h-full"
            />
          </div>
          <div className="lg:col-span-6">
            <RataRataNilaiWidget items={rataRata} className="h-full" />
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-12">
          <div className="lg:col-span-7">
            <PengumumanWidget items={pengumuman} />
          </div>
          <div className="lg:col-span-5">
            <JadwalUjianSiswaWidget items={jadwal} />
          </div>
        </div>

        {loading && (
          <div className="rounded-xl border border-dashed border-gray-200 bg-white px-4 py-3 text-xs text-gray-500">
            Menyegarkan data dashboard siswa...
          </div>
        )}
      </div>
    </div>
  );
};

export default HomeSiswa;
