import { useAuth } from "@/contexts/AuthContext";
import { PengumumanWidget } from "@/components/features/widget/Pengumuman/PengumumanWidget";
import { useGetPengumumanActive } from "@/services/Api/features-api/pengumuman/pengumuman.service";

const PengumumanSiswa = () => {
  const { user } = useAuth();
  const { data: pengumumanItems, loading, error } = useGetPengumumanActive();

  return (
    <div className="min-h-screen bg-[#ecf1ed] pb-20">
      <div className="mx-auto flex max-w-[1920px] flex-col gap-6 p-4 sm:p-6 lg:p-8">
        <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <p className="text-sm font-semibold text-[#397e50]">
            Pengumuman Siswa
          </p>
          <h1 className="mt-2 text-2xl font-bold text-[#37513d]">
            Halo, {user?.username ?? "Siswa"}
          </h1>
          <p className="mt-2 text-sm text-gray-600">
            Semua pengumuman aktif ditampilkan dari service yang sama dengan
            admin dan guru.
          </p>
        </header>

        {loading ? (
          <div className="rounded-xl border border-dashed border-gray-200 bg-white px-4 py-3 text-sm text-gray-500">
            Memuat pengumuman...
          </div>
        ) : error ? (
          <div className="rounded-xl border border-dashed border-red-200 bg-white px-4 py-3 text-sm text-red-500">
            Gagal memuat pengumuman: {error}
          </div>
        ) : (
          <PengumumanWidget
            title="Papan Pengumuman"
            items={pengumumanItems ?? []}
            allowMultipleOpen
            className="min-h-[60vh]"
          />
        )}
      </div>
    </div>
  );
};

export default PengumumanSiswa;
