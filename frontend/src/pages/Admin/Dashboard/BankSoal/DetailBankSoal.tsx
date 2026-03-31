import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";

import SiswaEssayQuestionContent from "@/components/features/Ujian/SiswaEssayQuestionContent";
import { useAuth } from "@/contexts/AuthContext";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { isEssaySoal } from "@/helper/Ujian/soalType";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { useGetBankSoalById } from "@/services/Api/features-api/BankSoal/banksoal.service";
import { useGetSoalUjian } from "@/services/Api/features-api/Ujian/soalUjian.service";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

const OPTION_LABELS = ["A", "B", "C", "D", "E", "F", "G", "H"];

const selectBaseClass =
  "block w-full cursor-pointer appearance-none rounded-xl border border-slate-200 bg-white py-2.5 pl-3 pr-8 text-sm font-medium text-slate-700 outline-none transition hover:border-slate-300 focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50]/20";

type BankSoalPreviewContentProps = {
  title: string;
  soalPreview: SoalPreviewItem[];
  acakSoal: boolean;
  onChangeAcakSoal: (value: boolean) => void;
  onBack: () => void;
};

const BankSoalPreviewContent = ({
  title,
  soalPreview,
  acakSoal,
  onChangeAcakSoal,
  onBack,
}: BankSoalPreviewContentProps) => {
  const [selectedOptions, setSelectedOptions] = useState<Record<number, number>>(
    {},
  );
  const [currentIndex, setCurrentIndex] = useState(0);

  const totalSoal = soalPreview.length;
  const currentSoal = soalPreview[currentIndex];
  const isCurrentSoalEssay = isEssaySoal(currentSoal.tipe);

  const handleSelectOption = (soalId: number, optionId: number) => {
    setSelectedOptions((prev) => ({
      ...prev,
      [soalId]: optionId,
    }));
  };

  const questionNavigator = (
    <div className="space-y-4">
      <div className="space-y-2">
        <label
          htmlFor="acak_soal_preview"
          className="text-xs font-medium text-slate-600"
        >
          Acak Soal
        </label>
        <select
          id="acak_soal_preview"
          value={acakSoal ? "ya" : "tidak"}
          onChange={(event) => onChangeAcakSoal(event.target.value === "ya")}
          className={selectBaseClass}
        >
          <option value="tidak">Tidak, urutan tetap</option>
          <option value="ya">Ya, acak soal</option>
        </select>
      </div>

      <div className="space-y-3">
        <p className="text-xs font-semibold text-slate-400">Nomor Urut Soal</p>
        <div className="flex flex-wrap gap-1.5">
          {soalPreview.map((soal, index) => {
            const isActive = index === currentIndex;
            return (
              <button
                key={soal.id}
                type="button"
                onClick={() => setCurrentIndex(index)}
                className={[
                  "flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg border text-xs font-semibold transition",
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
    </div>
  );

  return (
    <SoalLayout
      title={title}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu="00:00:00"
      soal={currentSoal}
      questionNavigator={questionNavigator}
      questionContent={
        isCurrentSoalEssay ? (
          <SiswaEssayQuestionContent
            soal={currentSoal}
            value=""
            readOnly
          />
        ) : undefined
      }
      selectedOptionId={selectedOptions[currentSoal.id]}
      onSelectOption={(optionId) => handleSelectOption(currentSoal.id, optionId)}
      onPrev={currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined}
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={onBack}
    />
  );
};

const DetailBankSoal = () => {
  const params = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const bankSoalId = Number(params.id);
  const isBankSoalIdValid = Number.isInteger(bankSoalId) && bankSoalId > 0;
  const bankSoalListPath =
    user?.role === "GURU"
      ? paths.dashboard.bank_soal_guru
      : paths.dashboard.bank_soal;

  const [acakSoal, setAcakSoal] = useState(false);

  const {
    data: soalRows,
    loading: loadingSoal,
    error: soalError,
  } = useGetSoalUjian(bankSoalId, acakSoal, isBankSoalIdValid);

  const { data: bankSoalData, loading: loadingBankSoal } = useGetBankSoalById(
    bankSoalId,
    isBankSoalIdValid,
  );

  const soalPreview = useMemo<SoalPreviewItem[]>(
    () =>
      (soalRows ?? []).map((soal) => ({
        id: soal.id_soal,
        nomor: soal.no_urut_soal,
        tipe: soal.tipe_soal,
        pertanyaan: soal.pertanyaan,
        gambar_url: resolveImageUrl(soal.gambar) || undefined,
        opsi: soal.opsi_jawaban.map((opsi, index) => ({
          id: opsi.id_pilihan_ganda,
          label: OPTION_LABELS[index] ?? String(index + 1),
          text: opsi.isi_pilihan,
        })),
      })),
    [soalRows],
  );

  const errorMsg = !isBankSoalIdValid ? "ID bank soal tidak valid." : soalError;
  const title =
    bankSoalData?.nama_bank_soal ??
    (isBankSoalIdValid ? `Bank Soal #${bankSoalId}` : "Bank Soal");

  if (loadingSoal || (isBankSoalIdValid && loadingBankSoal && !bankSoalData)) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50 px-4 text-center text-sm text-slate-500">
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
          onClick={() => navigate(bankSoalListPath)}
          className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-500 transition hover:border-[#397e50] hover:text-[#397e50]"
        >
          Kembali ke Bank Soal
        </button>
      </div>
    );
  }

  if (soalPreview.length === 0) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-slate-50 px-4 text-center">
        <p className="text-sm font-semibold text-slate-700">
          Bank soal ini belum memiliki pertanyaan.
        </p>
        <button
          type="button"
          onClick={() => navigate(bankSoalListPath)}
          className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-500 transition hover:border-[#397e50] hover:text-[#397e50]"
        >
          Kembali ke Bank Soal
        </button>
      </div>
    );
  }

  return (
    <BankSoalPreviewContent
      key={`${bankSoalId}-${acakSoal ? "acak" : "tetap"}`}
      title={title}
      soalPreview={soalPreview}
      acakSoal={acakSoal}
      onChangeAcakSoal={setAcakSoal}
      onBack={() => navigate(bankSoalListPath)}
    />
  );
};

export default DetailBankSoal;
