import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { getSoalUjian } from "@/services/Api/features-api/BankSoal/soalUjian.service";
import type { SoalUjianItem } from "@/types/BankSoal/BankSoal";

const DetailBankSoal = () => {
  const params = useParams();
  const navigate = useNavigate();
  const bankSoalId = Number(params.id);

  const [soalUjian, setSoalUjian] = useState<{
    namaUjian: string;
    sisaWaktu: string;
    soal: SoalUjianItem[];
  } | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<
    Record<number, string>
  >({});
  const [currentIndex, setCurrentIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState("");

  useEffect(() => {
    let active = true;
    const loadSoal = async () => {
      if (!Number.isFinite(bankSoalId)) {
        setErrorMsg("ID bank soal tidak valid.");
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        setErrorMsg("");
        const data = await getSoalUjian(bankSoalId);
        if (!active) return;
        setSoalUjian({
          namaUjian: data.nama_ujian,
          sisaWaktu: data.sisa_waktu ?? "00:00:00",
          soal: data.soal,
        });
        setCurrentIndex(0);
      } catch {
        if (!active) return;
        setErrorMsg("Soal bank soal tidak ditemukan.");
        setSoalUjian(null);
      } finally {
        if (active) setLoading(false);
      }
    };
    loadSoal();
    return () => {
      active = false;
    };
  }, [bankSoalId]);

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

  return (
    <SoalLayout
      title={soalUjian.namaUjian}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={soalUjian.sisaWaktu}
      soal={currentSoal}
      selectedOption={selectedOptions[currentSoal.id_soal]}
      onSelectOption={(value) => handleSelectOption(currentSoal.id_soal, value)}
      onPrev={handlePrev}
      onNext={handleNext}
      onBack={() => navigate(paths.dashboard.bank_soal)}
    />
  );
};

export default DetailBankSoal;
