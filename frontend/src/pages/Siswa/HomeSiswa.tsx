import { useAuth } from "@/contexts/AuthContext";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import {
  useGetSiswaDashboardSummary,
  useGetSiswaPengumuman,
} from "@/services/Api/features-api/Siswa/homeSiswa.service";

const HomeSiswa = () => {
  const { user } = useAuth();

  useGetSiswaDashboardSummary();
  const { data: pengumuman } = useGetSiswaPengumuman();

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
        <div className="lg:col-span-7">
          <PengumumanWidget items={pengumuman ?? []} />
        </div>
      </div>
    </div>
  );
};

export default HomeSiswa;