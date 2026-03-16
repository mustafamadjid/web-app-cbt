import React from "react";

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
}) => {
  const totalSoal = soalPreview.length;
  const currentSoal = soalPreview[currentIndex];

  const questionNavigator = (
    <div className="space-y-3">
      <p className="text-xs font-semibold text-slate-400">Nomor Soal</p>
      <div className="flex flex-wrap gap-2">
        {soalPreview.map((soal, index) => {
          const isActive = index === currentIndex;
          return (
            <button
              key={soal.id}
              type="button"
              onClick={() => onNavigateQuestion(index)}
              className={[
                "flex cursor-pointer h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold transition",
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
      selectedOptionId={selectedOptions[currentSoal.id]}
      onSelectOption={(optionId) => onSelectOption(currentSoal.id, optionId)}
      essayAnswer={essayAnswers[currentSoal.id] ?? ""}
      onEssayAnswerChange={(value) =>
        onEssayAnswerChange(currentSoal.id, value)
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
