import React from "react";
import { useNavigate } from "react-router";
import UjianFilterBar from "@/pages/Siswa/Ujian/UjianSiswa/components/UjianFilterBar";
import UjianSection from "@/pages/Siswa/Ujian/UjianSiswa/components/UjianSection";
import UjianCard from "@/pages/Siswa/Ujian/UjianSiswa/components/UjianCard";
import UjianResultCard from "@/pages/Siswa/Ujian/UjianSiswa/components/UjianResultCard";
import {
  getUjianSiswaOverview,
} from "@/services/Api/features-api/Ujian/ujianSiswa.service";
import type {
  UjianSiswaFilterParams,
  UjianSiswaResponse,
} from "@/types/Ujian/ujianSiswa";
import { paths } from "@/routes/paths";

const DEFAULT_SISWA_ID = 14;

const UjianSiswa: React.FC = () => {
  const navigate = useNavigate();
  const [filter, setFilter] = React.useState<UjianSiswaFilterParams>({});
  const [data, setData] = React.useState<UjianSiswaResponse>({
    upcoming: [],
    ongoing: [],
    completed: [],
    mapelOptions: [],
  });
  const [loading, setLoading] = React.useState(true);

  const fetchData = React.useCallback(async () => {
    setLoading(true);
    try {
      const response = await getUjianSiswaOverview({
        siswaId: DEFAULT_SISWA_ID,
        filter,
      });
      setData(response);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  React.useEffect(() => {
    void fetchData();
  }, [fetchData]);

  const handleStartExam = (id: number) => {
    navigate(paths.dashboard.ujian_siswa_token.replace(":id", String(id)));
  };

  return (
    <div className="space-y-6">
      <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-[#37513d]">Ujian Siswa</h1>
        <p className="mt-2 text-sm text-gray-500">
          Pantau jadwal ujian mendatang, ujian yang berlangsung hari ini, serta
          hasil ujian yang sudah selesai.
        </p>
      </header>

      <UjianFilterBar
        month={filter.bulan}
        year={filter.tahun}
        mapel={filter.mapel}
        mapelOptions={data.mapelOptions}
        onMonthChange={(value) =>
          setFilter((prev) => ({ ...prev, bulan: value }))
        }
        onYearChange={(value) => setFilter((prev) => ({ ...prev, tahun: value }))}
        onMapelChange={(value) =>
          setFilter((prev) => ({ ...prev, mapel: value }))
        }
        onReset={() => setFilter({})}
      />

      {loading ? (
        <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
          Memuat data ujian...
        </div>
      ) : (
        <div className="space-y-8">
          <UjianSection
            title="Jadwal Ujian Mendatang"
            description="Daftar ujian yang akan dilaksanakan dalam waktu dekat."
          >
            {data.upcoming.length === 0 ? (
              <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
                Tidak ada jadwal ujian mendatang.
              </div>
            ) : (
              data.upcoming.map((item) => (
                <UjianCard key={item.id} item={item} />
              ))
            )}
          </UjianSection>

          <UjianSection
            title="Ujian Berlangsung Hari Ini"
            description="Ujian yang sedang berlangsung dan bisa dikerjakan sekarang."
          >
            {data.ongoing.length === 0 ? (
              <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
                Tidak ada ujian yang sedang berlangsung.
              </div>
            ) : (
              data.ongoing.map((item) => (
                <UjianCard
                  key={item.id}
                  item={item}
                  actionLabel="Mulai Sekarang"
                  onAction={() => handleStartExam(item.id)}
                />
              ))
            )}
          </UjianSection>

          <UjianSection
            title="Hasil Ujian"
            description="Rekap ujian yang sudah kamu selesaikan beserta nilai."
          >
            {data.completed.length === 0 ? (
              <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
                Belum ada hasil ujian yang tersedia.
              </div>
            ) : (
              data.completed.map((item) => (
                <UjianResultCard key={item.id} item={item} />
              ))
            )}
          </UjianSection>
        </div>
      )}
    </div>
  );
};

export default UjianSiswa;
