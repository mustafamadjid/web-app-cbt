import { useExamBrowserExitCleanup } from "./useExamBrowserExitCleanup";
import { useUjianMulaiSiswaDraft } from "./useUjianMulaiSiswaDraft";
import { useUjianMulaiSiswaPreview } from "./useUjianMulaiSiswaPreview";
import type {
  UseUjianMulaiSiswaControllerParams,
  UseUjianMulaiSiswaControllerResult,
} from "./types";

export function useUjianMulaiSiswaController({
  resetKey,
  attemptId,
  soalPreview,
  jawabanUjianData,
  loadingJawabanUjian,
  refetchJawabanUjian,
  executeSaveJawaban,
  savingJawaban,
  saveJawabanError,
  resetSaveJawabanState,
  executeSubmitUjian,
  submittingAttempt,
  submitAttemptError,
  resetSubmitAttemptState,
  expiringAttempt,
  sessionExitError,
  clearSessionExitError,
  hasActiveExamSession,
  allowNavigation,
  handleBrowserViolation,
  closeExamWindowOrLeave,
  clearAttemptCache,
  clearSoalCache,
  navigateToResults,
}: UseUjianMulaiSiswaControllerParams): UseUjianMulaiSiswaControllerResult {
  const {
    currentIndex,
    currentSoal,
    selectedOptions,
    essayAnswers,
    saveCurrentJawaban,
    handleSelectOption,
    handleEssayAnswerChange,
    handleNavigateQuestion,
  } = useUjianMulaiSiswaDraft({
    resetKey,
    attemptId,
    soalPreview,
    jawabanUjianData,
    executeSaveJawaban,
    savingJawaban,
    saveJawabanError,
    resetSaveJawabanState,
    sessionExitError,
    clearSessionExitError,
  });

  const handleBrowserExitViolation = async () => {
    if (!hasActiveExamSession) {
      return false;
    }

    clearSessionExitError();
    await saveCurrentJawaban(currentSoal?.id ?? null);

    const expired = await handleBrowserViolation();
    if (!expired) {
      return false;
    }

    closeExamWindowOrLeave();
    return true;
  };

  useExamBrowserExitCleanup({
    hasActiveExamSession,
    onViolation: handleBrowserExitViolation,
  });

  const {
    previewJawabanItems,
    preparingPreview,
    isSubmitPreviewMode,
    isManualSubmitBusy,
    handleOpenSubmitPreview,
    handleBackToQuestionMode,
    handleFinalSubmit,
  } = useUjianMulaiSiswaPreview({
    resetKey,
    attemptId,
    jawabanUjianData,
    loadingJawabanUjian,
    refetchJawabanUjian,
    savingJawaban,
    saveCurrentJawaban,
    currentSoalId: currentSoal?.id ?? null,
    executeSubmitUjian,
    submittingAttempt,
    submitAttemptError,
    resetSubmitAttemptState,
    expiringAttempt,
    allowNavigation,
    clearAttemptCache,
    clearSoalCache,
    navigateToResults,
  });

  return {
    currentIndex,
    currentSoal,
    selectedOptions,
    essayAnswers,
    previewJawabanItems,
    preparingPreview,
    isSubmitPreviewMode,
    isManualSubmitBusy,
    saveCurrentJawaban,
    handleSelectOption,
    handleEssayAnswerChange,
    handleNavigateQuestion,
    handleOpenSubmitPreview,
    handleBackToQuestionMode,
    handleFinalSubmit,
  };
}
