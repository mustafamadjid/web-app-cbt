import React from "react";
import { useNavigate } from "react-router";
import UjianFilterBar from "@/components/features/Ujian/UjianFilterBar";
import UjianSection from "@/components/features/Ujian/UjianSection";
import UjianCard from "@/components/features/Ujian/UjianCard";
import UjianResultCard from "@/components/features/Ujian/UjianResultCard";
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
  const [activeCategory, setActiveCategory] = React.useState<
    "upcoming" | "ongoing" | "completed"
  >("upcoming");
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

  const handleStartExam = (id: number, bankSoalId: number) => {
    navigate(
      paths.dashboard.ujian_siswa_token
        .replace(":id", String(id))
        .replace(":bankSoalId", String(bankSoalId))
    );
  };

  const categories = [
    {
      key: "upcoming" as const,
      title: "Ujian Mendatang",
      subtitle: "Lihat jadwal ujian yang akan dilaksanakan.",
    },
    {
      key: "ongoing" as const,
      title: "Ujian Berlangsung",
      subtitle: "Mulai ujian yang sedang berlangsung.",
    },
    {
      key: "completed" as const,
      title: "Hasil Ujian",
      subtitle: "Cek hasil ujian yang sudah selesai.",
    },
  ];

  return (
    <div className="space-y-6 px-8 py-10">
      <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-[#37513d]">Ujian Siswa</h1>
        <p className="mt-2 text-sm text-gray-500">
          Pilih kategori ujian terlebih dahulu untuk melihat detailnya.
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
          <div className="grid gap-4 lg:grid-cols-3">
            {categories.map((category) => {
              const isActive = category.key === activeCategory;
              return (
                <button
                  key={category.key}
                  type="button"
                  onClick={() => setActiveCategory(category.key)}
                  className={`rounded-xl border p-5 text-left transition ${
                    isActive
                      ? "border-[#37513d] bg-[#f3f7f4] shadow-sm"
                      : "border-gray-200 bg-white hover:border-[#37513d]/60"
                  }`}
                >
                  <p className="text-base font-semibold text-[#37513d]">
                    {category.title}
                  </p>
                  <p className="mt-2 text-sm text-gray-500">
                    {category.subtitle}
                  </p>
                </button>
              );
            })}
          </div>

          {activeCategory === "upcoming" && (
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
          )}

          {activeCategory === "ongoing" && (
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
                    onAction={() => handleStartExam(item.id, item.id_bank_soal)}
                  />
                ))
              )}
            </UjianSection>
          )}

          {activeCategory === "completed" && (
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
          )}
        </div>
      )}
    </div>
  );
};

export default UjianSiswa;
