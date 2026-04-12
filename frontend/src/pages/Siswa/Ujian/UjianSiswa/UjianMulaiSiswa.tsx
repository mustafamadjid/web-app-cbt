import React from "react";
import { useNavigate, useParams } from "react-router";

import SiswaSubmitPreviewContent from "@/components/features/Ujian/SiswaSubmitPreviewContent";
import SiswaSoalPreviewContent from "@/components/features/Ujian/SiswaSoalPreviewContent";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { useAuth } from "@/contexts/AuthContext";
import { buildSoalPreview } from "@/helper/Ujian/buildSoalPreview";
import { deriveUjianMulaiSiswaState } from "@/helper/Ujian/deriveUjianMulaiSiswaState";
import { useCachedActiveAttemptId } from "@/hooks/Ujian/useCachedActiveAttemptId";
import { useCachedSoalUjianForSiswa } from "@/hooks/Ujian/useCachedSoalUjianForSiswa";
import { useExamCountdown } from "@/hooks/Ujian/useExamCountdown";
import { useExamSessionExit } from "@/hooks/Ujian/useExamSessionExit";
import { useUjianMulaiSiswaController } from "@/hooks/Ujian/UjianMulaiSiswa";
import { paths } from "@/routes/paths";
import {
  useGetJawabanUjianSiswaByAttemptId,
  useGetWaktuSelesaiUjian,
  useSaveJawabanUjianSiswa,
  useSubmitAttemptUjianSiswa,
} from "@/services/Api/features-api/Ujian/ujian.service";

const EXAM_LIST_PATH = paths.dashboard.ujian_siswa;

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { user, status } = useAuth();
  const { idJadwalUjian } = useParams();

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
  } = useGetJawabanUjianSiswaByAttemptId(attemptId ?? 0, shouldSyncJawaban);
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
  const waktuSelesai = waktuSelesaiData?.waktu_selesai ?? "";
  const { sisaWaktu, isTimeExpired, hasValidWaktuSelesai } =
    useExamCountdown(waktuSelesai);
  const jawabanSyncError = saveJawabanError ?? jawabanUjianError;

  const handleBack = React.useCallback(() => {
    navigate(EXAM_LIST_PATH);
  }, [navigate]);

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

  const navigateToResults = React.useCallback(() => {
    navigate(paths.dashboard.hasil_ujian_siswa, { replace: true });
  }, [navigate]);

  const {
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
  } = useUjianMulaiSiswaController({
    resetKey: parsedIdJadwalUjian,
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
  });

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
