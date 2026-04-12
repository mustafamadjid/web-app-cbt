import React from "react";
import toast from "react-hot-toast";

import {
  applySavedJawabanToDrafts,
  buildEssayAnswersMap,
  buildSaveJawabanPayload,
  buildSelectedOptionsMap,
  createJawabanDraft,
  mergeServerJawabanIntoDrafts,
  type JawabanDraftMap,
} from "@/helper/Ujian/jawabanDraft";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

import type { UseUjianMulaiSiswaControllerParams } from "./types";

type UseUjianMulaiSiswaDraftParams = Pick<
  UseUjianMulaiSiswaControllerParams,
  | "attemptId"
  | "clearSessionExitError"
  | "executeSaveJawaban"
  | "jawabanUjianData"
  | "resetKey"
  | "resetSaveJawabanState"
  | "saveJawabanError"
  | "savingJawaban"
  | "sessionExitError"
  | "soalPreview"
>;

type UseUjianMulaiSiswaDraftResult = {
  currentIndex: number;
  currentSoal: SoalPreviewItem | null;
  selectedOptions: Record<number, number>;
  essayAnswers: Record<number, string>;
  saveCurrentJawaban: (soalId: number | null) => Promise<boolean>;
  handleSelectOption: (soalId: number, optionId: number) => void;
  handleEssayAnswerChange: (soalId: number, value: string) => void;
  handleNavigateQuestion: (nextIndex: number) => void;
};

export function useUjianMulaiSiswaDraft({
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
}: UseUjianMulaiSiswaDraftParams): UseUjianMulaiSiswaDraftResult {
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [draftAnswers, setDraftAnswers] = React.useState<JawabanDraftMap>({});

  const currentSoal = soalPreview[currentIndex] ?? null;
  const selectedOptions = React.useMemo(
    () => buildSelectedOptionsMap(draftAnswers),
    [draftAnswers],
  );
  const essayAnswers = React.useMemo(
    () => buildEssayAnswersMap(draftAnswers),
    [draftAnswers],
  );

  React.useEffect(() => {
    setCurrentIndex(0);
    setDraftAnswers({});
    resetSaveJawabanState();
  }, [resetKey, resetSaveJawabanState]);

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
    saveCurrentJawaban,
    handleSelectOption,
    handleEssayAnswerChange,
    handleNavigateQuestion,
  };
}
