import React from "react";
import UjianFilterBar from "@/components/features/Ujian/UjianFilterBar";
import UjianSection from "@/components/features/Ujian/UjianSection";
import UjianResultCard from "@/components/features/Ujian/UjianResultCard";
import { useGetUjianBySiswa } from "@/services/Api/features-api/Ujian/ujianSiswa.service";
import type {
  UjianSiswaFilterParams,
  UjianSiswaResponse,
} from "@/types/Ujian/ujianSiswa";

const DEFAULT_SISWA_ID = 14;

const HasilUjianSiswa: React.FC = () => {
  const [filter, setFilter] = React.useState<UjianSiswaFilterParams>({});

  const { data, loading } = useGetUjianBySiswa({
    siswaId: DEFAULT_SISWA_ID,
    filter,
  });

  const ujianData: UjianSiswaResponse = React.useMemo(
    () => ({
      upcoming: data?.upcoming ?? [],
      ongoing: data?.ongoing ?? [],
      completed: data?.completed ?? [],
      mapelOptions: data?.mapelOptions ?? [],
    }),
    [data],
  );

  return (
    <div className="space-y-6 px-8 py-10">
      <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-[#37513d]">Hasil Ujian</h1>
        <p className="mt-2 text-sm text-gray-500">
          Daftar ujian yang sudah selesai kamu ikuti dan sudah dikumpulkan.
        </p>
      </header>

      <UjianFilterBar
        month={filter.bulan}
        year={filter.tahun}
        mapel={filter.mapel}
        mapelOptions={ujianData.mapelOptions}
        onMonthChange={(value) =>
          setFilter((prev) => ({ ...prev, bulan: value }))
        }
        onYearChange={(value) =>
          setFilter((prev) => ({ ...prev, tahun: value }))
        }
        onMapelChange={(value) =>
          setFilter((prev) => ({ ...prev, mapel: value }))
        }
        onReset={() => setFilter({})}
      />

      {loading ? (
        <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
          Memuat hasil ujian...
        </div>
      ) : (
        <UjianSection
          title="Daftar Hasil Ujian"
          description="Klik salah satu hasil ujian untuk melihat detail nilai."
        >
          {ujianData.completed.length === 0 ? (
            <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
              Belum ada hasil ujian yang tersedia.
            </div>
          ) : (
            ujianData.completed.map((item) => (
              <UjianResultCard key={item.id} item={item} />
            ))
          )}
        </UjianSection>
      )}
    </div>
  );
};

export default HasilUjianSiswa;
