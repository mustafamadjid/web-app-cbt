import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { useGetSoalUjian } from "@/services/Api/features-api/BankSoal/soalUjian.service";

const DetailBankSoal = () => {
  const params = useParams();
  const navigate = useNavigate();
  const bankSoalId = Number(params.id);
  const isBankSoalIdValid = Number.isFinite(bankSoalId);

  const [selectedOptions, setSelectedOptions] = useState<
    Record<number, string>
  >({});
  const [currentIndex, setCurrentIndex] = useState(0);

  const {
    data: rawSoalUjian,
    loading,
    error,
  } = useGetSoalUjian(isBankSoalIdValid ? bankSoalId : -1);

  const soalUjian = useMemo(
    () =>
      rawSoalUjian
        ? {
            namaUjian: rawSoalUjian.nama_ujian,
            sisaWaktu: rawSoalUjian.sisa_waktu ?? "00:00:00",
            soal: rawSoalUjian.soal,
          }
        : null,
    [rawSoalUjian],
  );

  useEffect(() => {
    setCurrentIndex(0);
  }, [soalUjian?.namaUjian]);

  const errorMsg = !isBankSoalIdValid ? "ID bank soal tidak valid." : (error ?? "");

  const totalSoal = useMemo(
    () => soalUjian?.soal.length ?? 0,
    [soalUjian?.soal.length]
  );

  const handleSelectOption = (soalId: number, value: string) => {
    setSelectedOptions((prev) => ({
      ...prev,
      [soalId]: value,
    }));
  };

  const handlePrev = () => {
    setCurrentIndex((prev) => Math.max(prev - 1, 0));
  };

  const handleNext = () => {
    setCurrentIndex((prev) =>
      Math.min(prev + 1, Math.max(totalSoal - 1, 0))
    );
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50 text-sm text-slate-500">
        Memuat soal bank soal...
      </div>
    );
  }

  if (errorMsg) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-slate-50 px-4 text-center">
        <p className="text-sm font-semibold text-slate-700">{errorMsg}</p>
        <button
          type="button"
          onClick={() => navigate(paths.dashboard.bank_soal)}
          className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-500 transition hover:border-[#397e50] hover:text-[#397e50]"
        >
          Kembali ke Bank Soal
        </button>
      </div>
    );
  }

  if (!soalUjian || totalSoal === 0) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-slate-50 px-4 text-center">
        <p className="text-sm font-semibold text-slate-700">
          Bank soal ini belum memiliki pertanyaan.
        </p>
        <button
          type="button"
          onClick={() => navigate(paths.dashboard.bank_soal)}
          className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-500 transition hover:border-[#397e50] hover:text-[#397e50]"
        >
          Kembali ke Bank Soal
        </button>
      </div>
    );
  }

  const currentSoal = soalUjian.soal[currentIndex];
  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">
        Nomor Urut Soal
      </p>
      <div className="flex flex-wrap gap-2">
        {soalUjian.soal.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id_soal}
              type="button"
              onClick={() => setCurrentIndex(index)}
              className={[
                "flex h-10 cursor-pointer w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
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

  return (
    <SoalLayout
      title={soalUjian.namaUjian}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={soalUjian.sisaWaktu}
      soal={currentSoal}
      questionNavigator={questionNavigator}
      selectedOption={selectedOptions[currentSoal.id_soal]}
      onSelectOption={(value) => handleSelectOption(currentSoal.id_soal, value)}
      onPrev={handlePrev}
      onNext={handleNext}
      onBack={() => navigate(paths.dashboard.bank_soal)}
    />
  );
};

export default DetailBankSoal;
