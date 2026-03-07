import React from "react";
import { useNavigate, useParams } from "react-router";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { useGetMockSoalUjianSiswa } from "@/services/Api/features-api/Ujian/soalUjianSiswa.mock";
import type { SoalPreviewData } from "@/types/Ujian/SoalPreview";

type SiswaSoalPreviewContentProps = {
  soalData: SoalPreviewData;
  onBack: () => void;
};

const SiswaSoalPreviewContent: React.FC<SiswaSoalPreviewContentProps> = ({
  soalData,
  onBack,
}) => {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [selectedOptions, setSelectedOptions] = React.useState<
    Record<number, number>
  >({});

  const totalSoal = soalData.soal.length;
  const currentSoal = soalData.soal[currentIndex];

  const handleSelectOption = (soalId: number, optionId: number) => {
    setSelectedOptions((prev) => ({ ...prev, [soalId]: optionId }));
  };

  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">Nomor Soal</p>
      <div className="flex flex-wrap gap-2">
        {soalData.soal.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id}
              type="button"
              onClick={() => setCurrentIndex(index)}
              className={[
                "flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
                isActive
                  ? "border-[#397e50] bg-[#397e50] text-white"
                  : "border-slate-200 text-slate-500 hover:border-[#397e50] hover:text-[#397e50]",
              ].join(" ")}
              aria-label={`Soal nomor ${soal.nomor}`}
            >
              {soal.nomor}
            </button>
          );
        })}
      </div>
    </div>
  );

  return (
    <SoalLayout
      title={soalData.title}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={soalData.sisa_waktu}
      soal={currentSoal}
      questionNavigator={questionNavigator}
      selectedOptionId={selectedOptions[currentSoal.id]}
      onSelectOption={(optionId) => handleSelectOption(currentSoal.id, optionId)}
      onPrev={
        currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={onBack}
    />
  );
};

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { bankSoalId } = useParams();

  const parsedBankSoalId = Number(bankSoalId);
  const isBankSoalIdValid =
    Number.isInteger(parsedBankSoalId) && parsedBankSoalId > 0;
  const {
    data: soalData,
    loading,
    error,
  } = useGetMockSoalUjianSiswa(parsedBankSoalId, isBankSoalIdValid);

  const errorMessage =
    !bankSoalId || !isBankSoalIdValid ? "Bank soal tidak ditemukan." : error;

  if (loading) {
    return (
      <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
        Memuat soal ujian...
      </div>
    );
  }

  if (errorMessage || !soalData || soalData.soal.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-red-200 bg-white p-6 text-center text-sm text-red-500">
        {errorMessage ?? "Soal ujian tidak tersedia."}
      </div>
    );
  }

  return (
    <SiswaSoalPreviewContent
      key={parsedBankSoalId}
      soalData={soalData}
      onBack={() => navigate(paths.dashboard.ujian_siswa)}
    />
  );
};

export default UjianMulaiSiswa;
