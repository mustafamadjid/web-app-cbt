import React from "react";
import { useNavigate, useParams } from "react-router";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import type { SoalUjianResponse } from "@/types/BankSoal/BankSoal";
import { paths } from "@/routes/paths";
import { getSoalUjian } from "@/services/Api/features-api/BankSoal/soalUjian.service";

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { id, bankSoalId } = useParams();
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [soalData, setSoalData] = React.useState<SoalUjianResponse | null>(
    null
  );
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [selectedOptions, setSelectedOptions] = React.useState<
    Record<number, string>
  >({});

  React.useEffect(() => {
    const loadSoal = async () => {
      if (!bankSoalId) {
        setError("Bank soal tidak ditemukan.");
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        setError(null);
        const response = await getSoalUjian(Number(bankSoalId));
        setSoalData(response);
        setCurrentIndex(0);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Gagal memuat soal.");
      } finally {
        setLoading(false);
      }
    };

    void loadSoal();
  }, [bankSoalId]);

  const totalSoal = soalData?.soal.length ?? 0;
  const currentSoal = soalData?.soal[currentIndex];

  const handleSelectOption = (soalId: number, value: string) => {
    setSelectedOptions((prev) => ({ ...prev, [soalId]: value }));
  };

  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">Nomor Soal</p>
      <div className="flex flex-wrap gap-2">
        {soalData?.soal.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id_soal}
              type="button"
              onClick={() => setCurrentIndex(index)}
              className={[
                "flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
                isActive
                  ? "border-[#397e50] bg-[#397e50] text-white"
                  : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
              ].join(" ")}
              aria-label={`Soal nomor ${soal.nomor_urut_soal}`}
            >
              {soal.nomor_urut_soal}
            </button>
          );
        })}
      </div>
    </div>
  );

  if (loading) {
    return (
      <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
        Memuat soal ujian...
      </div>
    );
  }

  if (error || !currentSoal) {
    return (
      <div className="rounded-xl border border-dashed border-red-200 bg-white p-6 text-center text-sm text-red-500">
        {error ?? "Soal ujian tidak tersedia."}
      </div>
    );
  }

  return (
    <SoalLayout
      title={soalData?.nama_ujian ?? `Ujian ${id ?? ""}`}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={soalData?.sisa_waktu ?? "00:00:00"}
      soal={currentSoal}
      questionNavigator={questionNavigator}
      selectedOption={selectedOptions[currentSoal.id_soal]}
      onSelectOption={(value) => handleSelectOption(currentSoal.id_soal, value)}
      onPrev={
        currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={() => navigate(paths.dashboard.ujian_siswa)}
    />
  );
};

export default UjianMulaiSiswa;
