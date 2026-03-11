import React from "react";
import { useNavigate } from "react-router";
import UjianSection from "@/components/features/Ujian/UjianSection";
import UjianCard from "@/components/features/Ujian/UjianCard";
import { useGetUjianBySiswa } from "@/services/Api/features-api/Ujian/ujianSiswa.service";
import { paths } from "@/routes/paths";

const DEFAULT_SISWA_ID = 14;

const UjianSiswa: React.FC = () => {
  const navigate = useNavigate();
  const [activeCategory, setActiveCategory] = React.useState<"upcoming" | "ongoing">(
    "upcoming",
  );

  const { data, loading } = useGetUjianBySiswa({ siswaId: DEFAULT_SISWA_ID });
  const upcomingExams = data?.upcoming ?? [];
  const ongoingExams = data?.ongoing ?? [];

  const handleStartExam = (id: number, bankSoalId: number) => {
    navigate(
      paths.dashboard.ujian_siswa_token
        .replace(":id", String(id))
        .replace(":bankSoalId", String(bankSoalId)),
    );
  };

  const categories = [
    {
      key: "upcoming" as const,
      title: "Ujian Mendatang",
      subtitle: "Lihat jadwal ujian yang akan dilaksanakan.",
      badge: "Segera",
    },
    {
      key: "ongoing" as const,
      title: "Ujian Berlangsung",
      subtitle: "Mulai ujian yang sedang berlangsung.",
      badge: "Aktif",
    },
  ];

  return (
    <div className="space-y-6 px-8 py-10">
      <header className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-[#37513d]">Ujian Siswa</h1>
        <p className="mt-2 text-sm text-gray-500">
          Lihat jadwal ujian yang akan datang atau mulai ujian yang sedang
          berlangsung.
        </p>
      </header>

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
                  className={`group relative cursor-pointer overflow-hidden rounded-2xl border p-5 text-left transition duration-200 ${
                    isActive
                      ? "border-[#397e50] bg-[#f3f8f5] shadow-md"
                      : "border-gray-200 bg-white hover:-translate-y-0.5 hover:border-[#397e50]/60 hover:shadow"
                  }`}
                >
                  <div
                    className="absolute inset-x-0 top-0 h-1 bg-[#397e50]"
                  />
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex items-center gap-3">
                      <div
                        className="flex h-11 w-11 items-center justify-center rounded-xl bg-[#397e50] text-base font-semibold text-white shadow-sm"
                      >
                        {category.title.charAt(0)}
                      </div>
                      <div>
                        <p className="text-base font-semibold text-[#397e50]">
                          {category.title}
                        </p>
                        <p className="mt-1 text-sm text-gray-500">
                          {category.subtitle}
                        </p>
                      </div>
                    </div>
                    <span
                      className={`rounded-full px-3 py-1 text-xs font-semibold ${
                        isActive
                          ? "bg-[#397e50] text-white"
                          : "bg-gray-100 text-gray-600 group-hover:bg-[#397e50]/10"
                      }`}
                    >
                      {category.badge}
                    </span>
                  </div>
                  {isActive && (
                    <div className="mt-4 flex items-center gap-2 text-xs font-semibold uppercase tracking-wide text-[#397e50]">
                      <span className="h-2 w-2 rounded-full bg-[#397e50]" />
                      Aktif dipilih
                    </div>
                  )}
                </button>
              );
            })}
          </div>

          {activeCategory === "upcoming" && (
            <UjianSection
              title="Jadwal Ujian Mendatang"
              description="Daftar ujian yang akan dilaksanakan dalam waktu dekat."
            >
              {upcomingExams.length === 0 ? (
                <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
                  Tidak ada jadwal ujian mendatang.
                </div>
              ) : (
                upcomingExams.map((item) => (
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
              {ongoingExams.length === 0 ? (
                <div className="col-span-full rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
                  Tidak ada ujian yang sedang berlangsung.
                </div>
              ) : (
                ongoingExams.map((item) => (
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
        </div>
      )}
    </div>
  );
};

export default UjianSiswa;
