import React from "react";
import toast from "react-hot-toast";

import type { JawabanUjianSiswaItem, JawabanUjianSiswaResponse } from "@/types/Ujian/ujianSiswa";

import type { UseUjianMulaiSiswaControllerParams } from "./types";

type UseUjianMulaiSiswaPreviewParams = Pick<
  UseUjianMulaiSiswaControllerParams,
  | "allowNavigation"
  | "attemptId"
  | "clearAttemptCache"
  | "clearSoalCache"
  | "executeSubmitUjian"
  | "expiringAttempt"
  | "jawabanUjianData"
  | "loadingJawabanUjian"
  | "navigateToResults"
  | "refetchJawabanUjian"
  | "resetKey"
  | "resetSubmitAttemptState"
  | "savingJawaban"
  | "submitAttemptError"
  | "submittingAttempt"
> & {
  currentSoalId: number | null;
  saveCurrentJawaban: (soalId: number | null) => Promise<boolean>;
};

type UseUjianMulaiSiswaPreviewResult = {
  previewJawabanItems: JawabanUjianSiswaItem[];
  preparingPreview: boolean;
  isSubmitPreviewMode: boolean;
  isManualSubmitBusy: boolean;
  handleOpenSubmitPreview: () => void;
  handleBackToQuestionMode: () => void;
  handleFinalSubmit: () => void;
};

export function useUjianMulaiSiswaPreview({
  allowNavigation,
  attemptId,
  clearAttemptCache,
  clearSoalCache,
  currentSoalId,
  executeSubmitUjian,
  expiringAttempt,
  jawabanUjianData,
  loadingJawabanUjian,
  navigateToResults,
  refetchJawabanUjian,
  resetKey,
  resetSubmitAttemptState,
  saveCurrentJawaban,
  savingJawaban,
  submitAttemptError,
  submittingAttempt,
}: UseUjianMulaiSiswaPreviewParams): UseUjianMulaiSiswaPreviewResult {
  const [viewMode, setViewMode] = React.useState<"question" | "submit_preview">(
    "question",
  );
  const [preparingPreview, setPreparingPreview] = React.useState(false);
  const [previewJawabanData, setPreviewJawabanData] =
    React.useState<JawabanUjianSiswaResponse | null>(null);

  const previewJawabanItems = React.useMemo(
    () => (previewJawabanData ?? jawabanUjianData)?.jawaban ?? [],
    [jawabanUjianData, previewJawabanData],
  );

  React.useEffect(() => {
    setViewMode("question");
    setPreparingPreview(false);
    setPreviewJawabanData(null);
    resetSubmitAttemptState();
  }, [resetKey, resetSubmitAttemptState]);

  const handleOpenSubmitPreview = React.useCallback(() => {
    if (
      preparingPreview ||
      savingJawaban ||
      submittingAttempt ||
      expiringAttempt
    ) {
      return;
    }

    void (async () => {
      setPreparingPreview(true);

      try {
        const saved = await saveCurrentJawaban(currentSoalId);
        if (!saved) {
          return;
        }

        const latestJawaban = await refetchJawabanUjian();
        if (latestJawaban === null) {
          toast.error("Gagal memuat preview jawaban ujian.");
          return;
        }

        if (submitAttemptError) {
          resetSubmitAttemptState();
        }

        setPreviewJawabanData(latestJawaban);
        setViewMode("submit_preview");
      } finally {
        setPreparingPreview(false);
      }
    })();
  }, [
    currentSoalId,
    expiringAttempt,
    preparingPreview,
    refetchJawabanUjian,
    resetSubmitAttemptState,
    saveCurrentJawaban,
    savingJawaban,
    submitAttemptError,
    submittingAttempt,
  ]);

  const handleBackToQuestionMode = React.useCallback(() => {
    if (preparingPreview || submittingAttempt || expiringAttempt) {
      return;
    }

    setPreviewJawabanData(null);
    setViewMode("question");
  }, [expiringAttempt, preparingPreview, submittingAttempt]);

  const handleFinalSubmit = React.useCallback(() => {
    if (
      attemptId === null ||
      attemptId <= 0 ||
      viewMode !== "submit_preview" ||
      submittingAttempt ||
      preparingPreview ||
      expiringAttempt
    ) {
      return;
    }

    void (async () => {
      try {
        await executeSubmitUjian(attemptId);
        toast.success("Ujian berhasil disubmit");
        allowNavigation();
        clearAttemptCache();
        clearSoalCache();
        navigateToResults();
      } catch {
        toast.error("Submit ujian gagal. Silakan coba lagi.");
      }
    })();
  }, [
    allowNavigation,
    attemptId,
    clearAttemptCache,
    clearSoalCache,
    executeSubmitUjian,
    expiringAttempt,
    navigateToResults,
    preparingPreview,
    submittingAttempt,
    viewMode,
  ]);

  return {
    previewJawabanItems,
    preparingPreview,
    isSubmitPreviewMode: viewMode === "submit_preview",
    isManualSubmitBusy:
      preparingPreview ||
      savingJawaban ||
      submittingAttempt ||
      expiringAttempt ||
      loadingJawabanUjian,
    handleOpenSubmitPreview,
    handleBackToQuestionMode,
    handleFinalSubmit,
  };
}
