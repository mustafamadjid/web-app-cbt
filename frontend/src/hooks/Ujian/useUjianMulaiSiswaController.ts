import React from "react";
import toast from "react-hot-toast";

import {
  applySavedJawabanToDrafts,
  buildEssayAnswersMap,
  buildSaveJawabanPayload,
  buildSelectedOptionsMap,
  createJawabanDraft,
  type JawabanDraftMap,
  mergeServerJawabanIntoDrafts,
} from "@/helper/Ujian/jawabanDraft";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";
import type {
  JawabanUjianSiswaItem,
  JawabanUjianSiswaResponse,
  SaveJawabanUjianSiswaRequest,
} from "@/types/Ujian/ujianSiswa";

type UseUjianMulaiSiswaControllerParams = {
  resetKey: number;
  attemptId: number | null;
  soalPreview: SoalPreviewItem[];
  jawabanUjianData: JawabanUjianSiswaResponse | null;
  loadingJawabanUjian: boolean;
  refetchJawabanUjian: () => Promise<JawabanUjianSiswaResponse | null>;
  executeSaveJawaban: (
    payload: SaveJawabanUjianSiswaRequest,
  ) => Promise<boolean>;
  savingJawaban: boolean;
  saveJawabanError: string | null;
  resetSaveJawabanState: () => void;
  executeSubmitUjian: (attemptId: number) => Promise<boolean>;
  submittingAttempt: boolean;
  submitAttemptError: string | null;
  resetSubmitAttemptState: () => void;
  expiringAttempt: boolean;
  sessionExitError: string | null;
  clearSessionExitError: () => void;
  hasActiveExamSession: boolean;
  allowNavigation: () => void;
  clearAttemptCache: () => void;
  clearSoalCache: () => void;
  navigateToResults: () => void;
};

type UseUjianMulaiSiswaControllerResult = {
  currentIndex: number;
  currentSoal: SoalPreviewItem | null;
  selectedOptions: Record<number, number>;
  essayAnswers: Record<number, string>;
  previewJawabanItems: JawabanUjianSiswaItem[];
  preparingPreview: boolean;
  isSubmitPreviewMode: boolean;
  isManualSubmitBusy: boolean;
  saveCurrentJawaban: (soalId: number | null) => Promise<boolean>;
  handleSelectOption: (soalId: number, optionId: number) => void;
  handleEssayAnswerChange: (soalId: number, value: string) => void;
  handleNavigateQuestion: (nextIndex: number) => void;
  handleOpenSubmitPreview: () => void;
  handleBackToQuestionMode: () => void;
  handleFinalSubmit: () => void;
};

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
  clearAttemptCache,
  clearSoalCache,
  navigateToResults,
}: UseUjianMulaiSiswaControllerParams): UseUjianMulaiSiswaControllerResult {
  const pendingBrowserExitRef = React.useRef(false);
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [draftAnswers, setDraftAnswers] = React.useState<JawabanDraftMap>({});
  const [viewMode, setViewMode] = React.useState<"question" | "submit_preview">(
    "question",
  );
  const [preparingPreview, setPreparingPreview] = React.useState(false);
  const [previewJawabanData, setPreviewJawabanData] =
    React.useState<JawabanUjianSiswaResponse | null>(null);

  const currentSoal = soalPreview[currentIndex] ?? null;
  const selectedOptions = React.useMemo(
    () => buildSelectedOptionsMap(draftAnswers),
    [draftAnswers],
  );
  const essayAnswers = React.useMemo(
    () => buildEssayAnswersMap(draftAnswers),
    [draftAnswers],
  );
  const previewJawabanItems = React.useMemo(
    () => (previewJawabanData ?? jawabanUjianData)?.jawaban ?? [],
    [jawabanUjianData, previewJawabanData],
  );

  React.useEffect(() => {
    setCurrentIndex(0);
    setDraftAnswers({});
    setViewMode("question");
    setPreparingPreview(false);
    setPreviewJawabanData(null);
    resetSaveJawabanState();
    resetSubmitAttemptState();
  }, [resetKey, resetSaveJawabanState, resetSubmitAttemptState]);

  React.useEffect(() => {
    if (currentIndex < soalPreview.length) {
      return;
    }

    setCurrentIndex(0);
  }, [currentIndex, soalPreview.length]);

  React.useEffect(() => {
    if (!jawabanUjianData) {
      return;
    }

    setDraftAnswers((prev) =>
      mergeServerJawabanIntoDrafts({
        draftAnswers: prev,
        jawabanUjianData,
        soalPreview,
      }),
    );
  }, [jawabanUjianData, soalPreview]);

  const clearBrowserStorageCache = React.useCallback(() => {
    if (typeof window === "undefined") {
      return;
    }

    try {
      window.localStorage.clear();
    } catch {
      // Ignore storage cleanup errors.
    }
  }, []);

  React.useEffect(() => {
    if (!hasActiveExamSession || typeof window === "undefined") {
      pendingBrowserExitRef.current = false;
      return;
    }

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      pendingBrowserExitRef.current = true;
      event.preventDefault();
      event.returnValue = "";
    };

    const handlePageHide = () => {
      if (!pendingBrowserExitRef.current) {
        return;
      }

      clearBrowserStorageCache();
    };

    const resetPendingBrowserExit = () => {
      pendingBrowserExitRef.current = false;
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    window.addEventListener("pagehide", handlePageHide);
    window.addEventListener("focus", resetPendingBrowserExit);
    window.addEventListener("pageshow", resetPendingBrowserExit);

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      window.removeEventListener("pagehide", handlePageHide);
      window.removeEventListener("focus", resetPendingBrowserExit);
      window.removeEventListener("pageshow", resetPendingBrowserExit);
      pendingBrowserExitRef.current = false;
    };
  }, [clearBrowserStorageCache, hasActiveExamSession]);

  const handleSelectOption = React.useCallback(
    (soalId: number, optionId: number) => {
      if (sessionExitError) {
        clearSessionExitError();
      }

      if (saveJawabanError) {
        resetSaveJawabanState();
      }

      setDraftAnswers((prev) => ({
        ...prev,
        [soalId]: {
          ...(prev[soalId] ?? createJawabanDraft(soalId)),
          id_soal: soalId,
          id_pilihan: optionId,
          jawaban_essay: "",
          waktu_jawab: new Date().toISOString(),
          isDirty: true,
        },
      }));
    },
    [
      clearSessionExitError,
      resetSaveJawabanState,
      saveJawabanError,
      sessionExitError,
    ],
  );

  const handleEssayAnswerChange = React.useCallback(
    (soalId: number, value: string) => {
      if (sessionExitError) {
        clearSessionExitError();
      }

      if (saveJawabanError) {
        resetSaveJawabanState();
      }

      setDraftAnswers((prev) => ({
        ...prev,
        [soalId]: {
          ...(prev[soalId] ?? createJawabanDraft(soalId)),
          id_soal: soalId,
          id_pilihan: null,
          jawaban_essay: value,
          waktu_jawab: new Date().toISOString(),
          isDirty: true,
        },
      }));
    },
    [
      clearSessionExitError,
      resetSaveJawabanState,
      saveJawabanError,
      sessionExitError,
    ],
  );

  const saveCurrentJawaban = React.useCallback(
    async (soalId: number | null): Promise<boolean> => {
      if (soalId === null || attemptId === null || attemptId <= 0) {
        return true;
      }

      const currentJawaban = draftAnswers[soalId];
      if (!currentJawaban?.isDirty) {
        return true;
      }

      const payload = buildSaveJawabanPayload(attemptId, currentJawaban);

      try {
        await executeSaveJawaban(payload);
        toast.success("Jawaban berhasil disimpan");

        setDraftAnswers((prev) =>
          applySavedJawabanToDrafts({
            draftAnswers: prev,
            soalId,
            payload,
          }),
        );

        return true;
      } catch {
        return false;
      }
    },
    [attemptId, draftAnswers, executeSaveJawaban],
  );

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
        const saved = await saveCurrentJawaban(currentSoal?.id ?? null);
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
    currentSoal?.id,
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

  const handleNavigateQuestion = React.useCallback(
    (nextIndex: number) => {
      if (
        savingJawaban ||
        nextIndex === currentIndex ||
        nextIndex < 0 ||
        nextIndex >= soalPreview.length
      ) {
        return;
      }

      void (async () => {
        await saveCurrentJawaban(currentSoal?.id ?? null);
        setCurrentIndex(nextIndex);
      })();
    },
    [
      currentIndex,
      currentSoal?.id,
      saveCurrentJawaban,
      savingJawaban,
      soalPreview.length,
    ],
  );

  return {
    currentIndex,
    currentSoal,
    selectedOptions,
    essayAnswers,
    previewJawabanItems,
    preparingPreview,
    isSubmitPreviewMode: viewMode === "submit_preview",
    isManualSubmitBusy:
      preparingPreview ||
      savingJawaban ||
      submittingAttempt ||
      expiringAttempt ||
      loadingJawabanUjian,
    saveCurrentJawaban,
    handleSelectOption,
    handleEssayAnswerChange,
    handleNavigateQuestion,
    handleOpenSubmitPreview,
    handleBackToQuestionMode,
    handleFinalSubmit,
  };
}
