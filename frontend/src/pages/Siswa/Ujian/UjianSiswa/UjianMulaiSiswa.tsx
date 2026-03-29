import React from "react";
import toast from "react-hot-toast";
import { useNavigate, useParams } from "react-router";

import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import SiswaSubmitPreviewContent from "@/components/features/Ujian/SiswaSubmitPreviewContent";
import SiswaSoalPreviewContent, {
  type EssayAnswersMap,
  type SelectedOptionsMap,
} from "@/components/features/Ujian/SiswaSoalPreviewContent";
import { useAuth } from "@/contexts/AuthContext";
import { buildSoalPreview } from "@/helper/Ujian/buildSoalPreview";
import { deriveUjianMulaiSiswaState } from "@/helper/Ujian/deriveUjianMulaiSiswaState";
import { useCachedActiveAttemptId } from "@/hooks/Ujian/useCachedActiveAttemptId";
import { useCachedSoalUjianForSiswa } from "@/hooks/Ujian/useCachedSoalUjianForSiswa";
import { useExamCountdown } from "@/hooks/Ujian/useExamCountdown";
import { useExamSessionExit } from "@/hooks/Ujian/useExamSessionExit";
import { paths } from "@/routes/paths";
import {
  useGetJawabanUjianSiswaByAttemptId,
  useSubmitAttemptUjianSiswa,
  useGetWaktuSelesaiUjian,
  useSaveJawabanUjianSiswa,
} from "@/services/Api/features-api/Ujian/ujian.service";
import type {
  JawabanUjianSiswaResponse,
  SaveJawabanUjianSiswaRequest,
} from "@/types/Ujian/ujianSiswa";

const EXAM_LIST_PATH = paths.dashboard.ujian_siswa;

type JawabanDraft = {
  id_soal: number;
  id_pilihan: number | null;
  jawaban_essay: string;
  waktu_jawab: string | null;
  isDirty: boolean;
};

type JawabanDraftMap = Record<number, JawabanDraft>;

const createJawabanDraft = (idSoal: number): JawabanDraft => ({
  id_soal: idSoal,
  id_pilihan: null,
  jawaban_essay: "",
  waktu_jawab: null,
  isDirty: false,
});

const buildSaveJawabanPayload = (
  idAttempt: number,
  jawaban: JawabanDraft,
): SaveJawabanUjianSiswaRequest => {
  const normalizedEssay = jawaban.jawaban_essay.trim();
  return {
    id_attempt: idAttempt,
    jawaban: [
      {
        id_soal: jawaban.id_soal,
        id_pilihan: jawaban.id_pilihan,
        jawaban_essay: jawaban.id_pilihan === null ? normalizedEssay || null : null,
        waktu_jawab: jawaban.waktu_jawab ?? new Date().toISOString(),
      },
    ],
  };
};

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { user, status } = useAuth();
  const { idJadwalUjian } = useParams();
  const pendingBrowserExitRef = React.useRef(false);
  const [currentIndex, setCurrentIndex] = React.useState(0);
  const [draftAnswers, setDraftAnswers] = React.useState<JawabanDraftMap>({});
  const [viewMode, setViewMode] = React.useState<"question" | "submit_preview">(
    "question",
  );
  const [preparingPreview, setPreparingPreview] = React.useState(false);
  const [previewJawabanData, setPreviewJawabanData] =
    React.useState<JawabanUjianSiswaResponse | null>(null);

  const parsedIdJadwalUjian = Number(idJadwalUjian);
  const parsedIdSiswa = user?.id_pengguna ?? 0;
  const hasIdJadwalUjianParam = Boolean(idJadwalUjian);
  const isIdJadwalUjianValid =
    Number.isInteger(parsedIdJadwalUjian) && parsedIdJadwalUjian > 0;
  const isIdSiswaValid = Number.isInteger(parsedIdSiswa) && parsedIdSiswa > 0;
  const {
    attemptId,
    loading: loadingActiveAttempt,
    error: activeAttemptError,
    errorCode: activeAttemptErrorCode,
    clearCache: clearAttemptCache,
  } = useCachedActiveAttemptId(
    parsedIdJadwalUjian,
    parsedIdSiswa,
    isIdJadwalUjianValid && isIdSiswaValid,
  );
  const {
    data: soalRows,
    loading: loadingSoal,
    error: soalError,
    clearCache: clearSoalCache,
  } = useCachedSoalUjianForSiswa(
    parsedIdJadwalUjian,
    parsedIdSiswa,
    isIdJadwalUjianValid && isIdSiswaValid,
  );
  const {
    data: waktuSelesaiData,
    loading: loadingWaktuSelesai,
    error: waktuSelesaiError,
  } = useGetWaktuSelesaiUjian(parsedIdJadwalUjian, isIdJadwalUjianValid);
  const shouldSyncJawaban =
    attemptId !== null &&
    attemptId > 0 &&
    isIdJadwalUjianValid &&
    isIdSiswaValid;
  const {
    data: jawabanUjianData,
    loading: loadingJawabanUjian,
    error: jawabanUjianError,
    refetch: refetchJawabanUjian,
  } = useGetJawabanUjianSiswaByAttemptId(
    attemptId ?? 0,
    shouldSyncJawaban,
    [currentIndex],
  );
  const {
    execute: executeSaveJawaban,
    loading: savingJawaban,
    error: saveJawabanError,
    reset: resetSaveJawabanState,
  } = useSaveJawabanUjianSiswa();
  const {
    execute: executeSubmitUjian,
    loading: submittingAttempt,
    error: submitAttemptError,
    reset: resetSubmitAttemptState,
  } = useSubmitAttemptUjianSiswa();

  const soalPreview = React.useMemo(() => buildSoalPreview(soalRows), [soalRows]);
  const currentSoal = soalPreview[currentIndex] ?? null;
  const waktuSelesai = waktuSelesaiData?.waktu_selesai ?? "";
  const { sisaWaktu, isTimeExpired, hasValidWaktuSelesai } =
    useExamCountdown(waktuSelesai);
  const selectedOptions = React.useMemo<SelectedOptionsMap>(() => {
    const mapped: SelectedOptionsMap = {};

    for (const jawaban of Object.values(draftAnswers)) {
      if (jawaban.id_pilihan !== null) {
        mapped[jawaban.id_soal] = jawaban.id_pilihan;
      }
    }

    return mapped;
  }, [draftAnswers]);
  const essayAnswers = React.useMemo<EssayAnswersMap>(() => {
    const mapped: EssayAnswersMap = {};

    for (const jawaban of Object.values(draftAnswers)) {
      if (jawaban.jawaban_essay !== "") {
        mapped[jawaban.id_soal] = jawaban.jawaban_essay;
      }
    }

    return mapped;
  }, [draftAnswers]);
  const previewJawabanItems = React.useMemo(
    () => (previewJawabanData ?? jawabanUjianData)?.jawaban ?? [],
    [jawabanUjianData, previewJawabanData],
  );
  const jawabanSyncError = saveJawabanError ?? jawabanUjianError;

  React.useEffect(() => {
    setCurrentIndex(0);
    setDraftAnswers({});
    setViewMode("question");
    setPreparingPreview(false);
    setPreviewJawabanData(null);
    resetSaveJawabanState();
    resetSubmitAttemptState();
  }, [
    parsedIdJadwalUjian,
    resetSaveJawabanState,
    resetSubmitAttemptState,
  ]);

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

    setDraftAnswers((prev) => {
      const next = { ...prev };
      const serverJawabanBySoal = new Map(
        jawabanUjianData.jawaban.map((item) => [item.id_soal, item]),
      );

      for (const soal of soalPreview) {
        const existing = prev[soal.id];
        if (existing?.isDirty) {
          continue;
        }

        const serverJawaban = serverJawabanBySoal.get(soal.id);
        if (!serverJawaban) {
          delete next[soal.id];
          continue;
        }

        next[soal.id] = {
          id_soal: serverJawaban.id_soal,
          id_pilihan: serverJawaban.id_pilihan,
          jawaban_essay: serverJawaban.jawaban_essay ?? "",
          waktu_jawab: serverJawaban.waktu_jawab,
          isDirty: false,
        };
      }

      return next;
    });
  }, [jawabanUjianData, soalPreview]);

  const handleBack = React.useCallback(() => {
    navigate(EXAM_LIST_PATH);
  }, [navigate]);

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

  const { loading, errorMessage } = React.useMemo(
    () =>
      deriveUjianMulaiSiswaState({
        hasIdJadwalUjianParam,
        isIdJadwalUjianValid,
        isIdSiswaValid,
        isAuthLoading: status === "loading",
        loadingActiveAttempt,
        loadingSoal,
        loadingWaktuSelesai,
        activeAttemptError,
        activeAttemptErrorCode,
        soalError,
        waktuSelesaiError,
        waktuSelesai,
        hasValidWaktuSelesai,
      }),
    [
      activeAttemptError,
      activeAttemptErrorCode,
      hasIdJadwalUjianParam,
      hasValidWaktuSelesai,
      isIdJadwalUjianValid,
      isIdSiswaValid,
      loadingActiveAttempt,
      loadingSoal,
      loadingWaktuSelesai,
      soalError,
      status,
      waktuSelesai,
      waktuSelesaiError,
    ],
  );
  const title = isIdJadwalUjianValid ? `Ujian ` : "Ujian";
  const hasActiveExamSession =
    !loading &&
    !errorMessage &&
    attemptId !== null &&
    soalPreview.length > 0 &&
    !isTimeExpired;

  const {
    isLeaveConfirmOpen,
    sessionExitError,
    expiringAttempt,
    allowNavigation,
    clearSessionExitError,
    handleExpiredSubmit,
    handleLeaveConfirm,
    handleLeaveCancel,
  } = useExamSessionExit({
    attemptId,
    hasActiveExamSession,
    isTimeExpired,
    resetKey: parsedIdJadwalUjian,
    clearAttemptCache,
    clearSoalCache,
    onFallbackLeave: handleBack,
  });

  React.useEffect(() => {
    if (!hasActiveExamSession) {
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

  React.useEffect(() => {
    if (!isIdJadwalUjianValid || activeAttemptErrorCode !== "NOT_FOUND") {
      return;
    }

    allowNavigation();
    clearAttemptCache();
    navigate(
      paths.dashboard.ujian_siswa_token.replace(
        ":idJadwalUjian",
        String(parsedIdJadwalUjian),
      ),
      { replace: true },
    );
  }, [
    activeAttemptErrorCode,
    allowNavigation,
    clearAttemptCache,
    isIdJadwalUjianValid,
    navigate,
    parsedIdJadwalUjian,
  ]);

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

        setDraftAnswers((prev) => {
          const latestJawaban = prev[soalId];
          if (!latestJawaban) {
            return prev;
          }

          const [savedItem] = payload.jawaban;
          const shouldClearAnswer =
            savedItem.id_pilihan === null && savedItem.jawaban_essay === null;

          if (shouldClearAnswer) {
            const rest = { ...prev };
            delete rest[soalId];
            return rest;
          }

          return {
            ...prev,
            [soalId]: {
              ...latestJawaban,
              id_pilihan: savedItem.id_pilihan,
              jawaban_essay: savedItem.jawaban_essay ?? "",
              waktu_jawab: savedItem.waktu_jawab,
              isDirty: false,
            },
          };
        });

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
        navigate(paths.dashboard.hasil_ujian_siswa, { replace: true });
      } catch {
        toast.error("Submit ujian gagal. Silakan coba lagi.");
        return;
      }
    })();
  }, [
    allowNavigation,
    attemptId,
    clearAttemptCache,
    clearSoalCache,
    executeSubmitUjian,
    expiringAttempt,
    navigate,
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

  if (loading) {
    return (
      <div className="rounded-xl border border-dashed border-gray-200 bg-white p-6 text-center text-sm text-gray-500">
        Memuat soal ujian...
      </div>
    );
  }

  if (errorMessage || soalPreview.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-red-200 bg-white p-6 text-center text-sm text-red-500">
        {errorMessage ?? "Soal ujian tidak tersedia."}
      </div>
    );
  }

  const isSubmitPreviewMode = viewMode === "submit_preview";
  const isManualSubmitBusy =
    preparingPreview ||
    savingJawaban ||
    submittingAttempt ||
    expiringAttempt ||
    loadingJawabanUjian;

  return (
    <>
      {sessionExitError && (
        <div className="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600">
          {sessionExitError}
        </div>
      )}

      {jawabanSyncError && (
        <div className="mb-4 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-700">
          Sinkronisasi jawaban gagal. {jawabanSyncError}
        </div>
      )}

      {submitAttemptError && (
        <div className="mb-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600">
          Submit ujian gagal. {submitAttemptError}
        </div>
      )}

      {isSubmitPreviewMode ? (
        <SiswaSubmitPreviewContent
          title={title}
          sisaWaktu={sisaWaktu}
          soalPreview={soalPreview}
          jawabanItems={previewJawabanItems}
          onBackToQuestionMode={handleBackToQuestionMode}
          onSubmitFinal={handleFinalSubmit}
          backDisabled={isManualSubmitBusy}
          submitDisabled={isManualSubmitBusy}
          submitLoading={submittingAttempt}
        />
      ) : (
        <SiswaSoalPreviewContent
          key={parsedIdJadwalUjian}
          title={title}
          sisaWaktu={sisaWaktu}
          soalPreview={soalPreview}
          currentIndex={currentIndex}
          selectedOptions={selectedOptions}
          essayAnswers={essayAnswers}
          onSelectOption={handleSelectOption}
          onEssayAnswerChange={handleEssayAnswerChange}
          onNavigateQuestion={handleNavigateQuestion}
          onBack={handleBack}
          onSubmitExam={handleOpenSubmitPreview}
          submitDisabled={isManualSubmitBusy}
          submitLoading={preparingPreview}
        />
      )}

      <ConfirmAlert
        isOpen={isLeaveConfirmOpen}
        title="Keluar dari Ujian?"
        message="Jawaban yang belum dikirim bisa hilang. Yakin ingin meninggalkan halaman ujian ini?"
        onClose={handleLeaveCancel}
        onConfirm={() => {
          void (async () => {
            await saveCurrentJawaban(currentSoal?.id ?? null);
            await handleLeaveConfirm();
          })();
        }}
        confirmLabel="Ya, Keluar"
        cancelLabel="Lanjut Ujian"
        loadingLabel="Keluar..."
        isLoading={expiringAttempt}
      />

      <ConfirmAlert
        isOpen={isTimeExpired}
        title="Waktu Ujian Habis"
        message="Waktu pengerjaan ujian telah mencapai batas akhir. Tekan Submit untuk mengakhiri sesi ujian ini."
        onClose={() => undefined}
        onConfirm={() => {
          void (async () => {
            await saveCurrentJawaban(currentSoal?.id ?? null);
            await handleExpiredSubmit();
          })();
        }}
        confirmLabel="Submit"
        loadingLabel="Submit..."
        isLoading={expiringAttempt}
        hideCancel
        dismissible={false}
        confirmClassName="bg-[#397e50] hover:bg-[#326f45]"
      />
    </>
  );
};

export default UjianMulaiSiswa;
