import React from "react";

import SiswaEssayQuestionContent from "@/components/features/Ujian/SiswaEssayQuestionContent";
import { isEssaySoal } from "@/helper/Ujian/soalType";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

export type SelectedOptionsMap = Record<number, number>;
export type EssayAnswersMap = Record<number, string>;

type SiswaSoalPreviewContentProps = {
  title: string;
  sisaWaktu: string;
  soalPreview: SoalPreviewItem[];
  currentIndex: number;
  selectedOptions: SelectedOptionsMap;
  essayAnswers: EssayAnswersMap;
  onSelectOption: (soalId: number, optionId: number) => void;
  onEssayAnswerChange: (soalId: number, value: string) => void;
  onNavigateQuestion: (nextIndex: number) => void;
  onBack: () => void;
  onSubmitExam: () => void;
  submitDisabled?: boolean;
  submitLoading?: boolean;
};

const SiswaSoalPreviewContent: React.FC<SiswaSoalPreviewContentProps> = ({
  title,
  sisaWaktu,
  soalPreview,
  currentIndex,
  selectedOptions,
  essayAnswers,
  onSelectOption,
  onEssayAnswerChange,
  onNavigateQuestion,
  onBack,
  onSubmitExam,
  submitDisabled = false,
  submitLoading = false,
}) => {
  const totalSoal = soalPreview.length;
  const currentSoal = soalPreview[currentIndex];
  const isCurrentSoalEssay = isEssaySoal(currentSoal.tipe);

  const questionNavigator = (
    <div className="space-y-4">
      <div className="space-y-1 border-b border-slate-100 pb-3">
        <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
          Navigasi
        </p>
        <p className="text-sm font-semibold text-slate-700">Nomor Soal</p>
        <p className="text-xs text-slate-500">
          Pilih nomor soal dari panel kanan.
        </p>
      </div>

      <div className="grid grid-cols-4 gap-2 sm:grid-cols-5 lg:grid-cols-4">
        {soalPreview.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id}
              type="button"
              onClick={() => onNavigateQuestion(index)}
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
  );

  return (
    <SoalLayout
      title={title}
      currentNumber={currentIndex + 1}
      totalSoal={totalSoal}
      sisaWaktu={sisaWaktu}
      soal={currentSoal}
      questionNavigator={questionNavigator}
      questionContent={
        isCurrentSoalEssay ? (
          <SiswaEssayQuestionContent
            soal={currentSoal}
            value={essayAnswers[currentSoal.id] ?? ""}
            onChange={(value) => onEssayAnswerChange(currentSoal.id, value)}
          />
        ) : undefined
      }
      compactSplitLayout
      selectedOptionId={selectedOptions[currentSoal.id]}
      onSelectOption={(optionId) => onSelectOption(currentSoal.id, optionId)}
      footerActions={
        <button
          type="button"
          onClick={onSubmitExam}
          disabled={submitDisabled}
          className={[
            "rounded-lg px-4 py-2 text-sm font-semibold text-white transition",
            submitDisabled
              ? "cursor-not-allowed bg-slate-300"
              : "cursor-pointer bg-amber-500 hover:bg-amber-600",
          ].join(" ")}
        >
          {submitLoading ? "Menyiapkan..." : "Submit Ujian"}
        </button>
      }
      onPrev={
        currentIndex > 0
          ? () => onNavigateQuestion(currentIndex - 1)
          : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => onNavigateQuestion(currentIndex + 1)
          : undefined
      }
      onBack={onBack}
    />
  );
};

export default SiswaSoalPreviewContent;
