import React from "react";

import {
  QuestionList,
  QuestionNavigator,
  SubmitActions,
  SubmitPreviewHeader,
} from "@/components/features/Ujian/SiswaSubmitPreviewParts";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type { JawabanUjianSiswaItem } from "@/types/Ujian/ujianSiswa";

type SiswaSubmitPreviewContentProps = {
  title: string;
  sisaWaktu: string;
  soalPreview: SoalPreviewItem[];
  jawabanItems: JawabanUjianSiswaItem[];
  onBackToQuestionMode: () => void;
  onSubmitFinal: () => void;
  backDisabled?: boolean;
  submitDisabled?: boolean;
  submitLoading?: boolean;
};

const hasAnswer = (jawaban?: JawabanUjianSiswaItem) =>
  Boolean(
    jawaban &&
      (jawaban.id_pilihan !== null ||
        (jawaban.jawaban_essay?.trim() ?? "") !== ""),
  );

const SiswaSubmitPreviewContent: React.FC<SiswaSubmitPreviewContentProps> = ({
  title,
  sisaWaktu,
  soalPreview,
  jawabanItems,
  onBackToQuestionMode,
  onSubmitFinal,
  backDisabled = false,
  submitDisabled = false,
  submitLoading = false,
}) => {
  const getJawabanBySoalId = React.useCallback(
    (soalId: number) =>
      jawabanItems.find((item) => String(item.id_soal) === String(soalId)),
    [jawabanItems],
  );

  const isAnswered = React.useCallback(
    (soalId: number) => hasAnswer(getJawabanBySoalId(soalId)),
    [getJawabanBySoalId],
  );

  const answeredCount = soalPreview.filter((soal) =>
    isAnswered(soal.id),
  ).length;

  const scrollToQuestion = React.useCallback((soalId: number) => {
    const element = document.getElementById(`submit-preview-soal-${soalId}`);
    element?.scrollIntoView({ behavior: "smooth", block: "start" });
  }, []);

  return (
    <div className="min-h-screen w-full bg-slate-50">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-4 py-6 md:px-6">
        <SubmitPreviewHeader
          title={title}
          sisaWaktu={sisaWaktu}
          answeredCount={answeredCount}
          totalSoal={soalPreview.length}
        />

        <div className="grid gap-6 lg:grid-cols-[minmax(0,820px)_240px] lg:items-start lg:justify-center">
          <QuestionList
            soalPreview={soalPreview}
            getJawabanBySoalId={getJawabanBySoalId}
            isAnswered={isAnswered}
          />

          <QuestionNavigator
            soalPreview={soalPreview}
            isAnswered={isAnswered}
            onNavigate={scrollToQuestion}
          />
        </div>

        <SubmitActions
          onBackToQuestionMode={onBackToQuestionMode}
          onSubmitFinal={onSubmitFinal}
          backDisabled={backDisabled}
          submitDisabled={submitDisabled}
          submitLoading={submitLoading}
        />
      </div>
    </div>
  );
};

export default SiswaSubmitPreviewContent;
