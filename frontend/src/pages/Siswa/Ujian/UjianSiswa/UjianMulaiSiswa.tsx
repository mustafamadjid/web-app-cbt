import React from "react";
import { useBeforeUnload, useBlocker, useNavigate, useParams } from "react-router";
import type { BlockerFunction } from "react-router";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { useAuth } from "@/contexts/AuthContext";
import { getRemainingExamTime } from "@/helper/Countdown/getRemainingExamTime";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { useCachedActiveAttemptId } from "@/hooks/Ujian/useCachedActiveAttemptId";
import { useCachedSoalUjianForSiswa } from "@/hooks/Ujian/useCachedSoalUjianForSiswa";
import SoalLayout from "@/layouts/BankSoalLayout/SoalLayout";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  expireAttemptUjianSiswaOnPageLeave,
  useExpireAttemptUjianSiswa,
  useGetWaktuSelesaiUjian,
} from "@/services/Api/features-api/Ujian/ujian.service";
import type { SoalPreviewItem } from "@/types/Ujian/SoalPreview";

const OPTION_LABELS = ["A", "B", "C", "D", "E", "F", "G", "H"];
const DEFAULT_TIMER_LABEL = "00:00:00";
const EXPIRE_ATTEMPT_ERROR_MESSAGE =
  "Gagal mengakhiri sesi ujian. Silakan coba lagi.";

type SelectedOptionsMap = Record<number, number>;

type SiswaSoalPreviewContentProps = {
  title: string;
  sisaWaktu: string;
  soalPreview: SoalPreviewItem[];
  selectedOptions: SelectedOptionsMap;
  onSelectOption: (soalId: number, optionId: number) => void;
  onBack: () => void;
};

const SiswaSoalPreviewContent: React.FC<SiswaSoalPreviewContentProps> = ({
  title,
  sisaWaktu,
  soalPreview,
  selectedOptions,
  onSelectOption,
  onBack,
}) => {
  const [currentIndex, setCurrentIndex] = React.useState(0);

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
              onClick={() => setCurrentIndex(index)}
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
      onPrev={
        currentIndex > 0 ? () => setCurrentIndex((prev) => prev - 1) : undefined
      }
      onNext={
        currentIndex < totalSoal - 1
          ? () => setCurrentIndex((prev) => prev + 1)
          : undefined
      }
      onBack={onBack}
    />
  );
};

const mapExpireAttemptErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError && error.message) {
    return error.message;
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return EXPIRE_ATTEMPT_ERROR_MESSAGE;
};

const UjianMulaiSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { user, status } = useAuth();
  const { idJadwalUjian } = useParams();
  const [selectedOptions, setSelectedOptions] =
    React.useState<SelectedOptionsMap>({});
  const [sisaWaktu, setSisaWaktu] = React.useState(DEFAULT_TIMER_LABEL);
  const [isTimeExpired, setIsTimeExpired] = React.useState(false);
  const [isLeaveConfirmOpen, setIsLeaveConfirmOpen] = React.useState(false);
  const [sessionExitError, setSessionExitError] = React.useState<string | null>(null);
  const allowNavigationRef = React.useRef(false);
  const expireTriggeredRef = React.useRef(false);
  const { execute: executeExpireAttempt, loading: expiringAttempt } =
    useExpireAttemptUjianSiswa();

  const parsedIdJadwalUjian = Number(idJadwalUjian);
  const parsedIdSiswa = user?.id_pengguna ?? 0;
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

  const soalPreview = React.useMemo<SoalPreviewItem[]>(
    () =>
      (soalRows ?? []).map((soal, index) => ({
        id: soal.id_soal,
        nomor: index + 1,
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

  const waktuSelesai = waktuSelesaiData?.waktu_selesai ?? "";
  const initialRemainingTime = React.useMemo(() => {
    if (!waktuSelesai) return null;
    return getRemainingExamTime(waktuSelesai);
  }, [waktuSelesai]);

  React.useEffect(() => {
    setSelectedOptions({});
    setIsLeaveConfirmOpen(false);
    setSessionExitError(null);
    allowNavigationRef.current = false;
    expireTriggeredRef.current = false;
  }, [parsedIdJadwalUjian]);

  React.useEffect(() => {
    if (!isIdJadwalUjianValid || activeAttemptErrorCode !== "NOT_FOUND") {
      return;
    }

    allowNavigationRef.current = true;
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
    clearAttemptCache,
    isIdJadwalUjianValid,
    navigate,
    parsedIdJadwalUjian,
  ]);

  React.useEffect(() => {
    if (!waktuSelesai) {
      setSisaWaktu(DEFAULT_TIMER_LABEL);
      setIsTimeExpired(false);
      return;
    }

    const currentRemainingTime = getRemainingExamTime(waktuSelesai);
    if (!currentRemainingTime) {
      setSisaWaktu(DEFAULT_TIMER_LABEL);
      setIsTimeExpired(false);
      return;
    }

    setSisaWaktu(currentRemainingTime.formattedTime);
    setIsTimeExpired(currentRemainingTime.isExpired);

    if (currentRemainingTime.isExpired) return;

    const timerId = window.setInterval(() => {
      const nextRemainingTime = getRemainingExamTime(waktuSelesai);
      if (!nextRemainingTime) {
        window.clearInterval(timerId);
        setSisaWaktu(DEFAULT_TIMER_LABEL);
        setIsTimeExpired(false);
        return;
      }

      setSisaWaktu(nextRemainingTime.formattedTime);

      if (nextRemainingTime.isExpired) {
        setIsTimeExpired(true);
        window.clearInterval(timerId);
      }
    }, 1000);

    return () => {
      window.clearInterval(timerId);
    };
  }, [waktuSelesai]);

  const handleSelectOption = React.useCallback(
    (soalId: number, optionId: number) => {
      if (sessionExitError) {
        setSessionExitError(null);
      }
      setSelectedOptions((prev) => ({ ...prev, [soalId]: optionId }));
    },
    [sessionExitError],
  );

  const shouldRedirectToToken =
    isIdJadwalUjianValid && activeAttemptErrorCode === "NOT_FOUND";
  const loading =
    status === "loading" ||
    loadingActiveAttempt ||
    loadingSoal ||
    loadingWaktuSelesai ||
    shouldRedirectToToken;
  const waktuSelesaiStateError =
    !loadingWaktuSelesai && !waktuSelesai
      ? "Waktu selesai ujian tidak tersedia."
      : waktuSelesai && !initialRemainingTime
        ? "Format waktu selesai ujian tidak valid."
        : null;
  const errorMessage =
    !idJadwalUjian || !isIdJadwalUjianValid
      ? "Jadwal ujian tidak ditemukan."
      : status !== "loading" && !isIdSiswaValid
        ? "Akun siswa tidak valid. Silakan login ulang."
        : activeAttemptErrorCode !== "NOT_FOUND"
          ? activeAttemptError ?? soalError ?? waktuSelesaiError ?? waktuSelesaiStateError
          : soalError ?? waktuSelesaiError ?? waktuSelesaiStateError;
  const title = isIdJadwalUjianValid ? `Ujian ` : "Ujian";
  const hasActiveExamSession =
    !loading &&
    !errorMessage &&
    attemptId !== null &&
    soalPreview.length > 0 &&
    !isTimeExpired;
  const shouldBlockNavigation = React.useCallback<BlockerFunction>(
    ({ currentLocation, nextLocation }) => {
      if (allowNavigationRef.current || !hasActiveExamSession) {
        return false;
      }

      return (
        currentLocation.pathname !== nextLocation.pathname ||
        currentLocation.search !== nextLocation.search ||
        currentLocation.hash !== nextLocation.hash
      );
    },
    [hasActiveExamSession],
  );
  const navigationBlocker = useBlocker(shouldBlockNavigation);

  const expireAttemptBeforeLeave = React.useCallback(async () => {
    if (attemptId === null || expireTriggeredRef.current) {
      return true;
    }

    setSessionExitError(null);
    expireTriggeredRef.current = true;

    try {
      await executeExpireAttempt(attemptId);
      clearAttemptCache();
      clearSoalCache();
      return true;
    } catch (error) {
      expireTriggeredRef.current = false;
      setSessionExitError(mapExpireAttemptErrorMessage(error));
      return false;
    }
  }, [
    attemptId,
    clearAttemptCache,
    clearSoalCache,
    executeExpireAttempt,
  ]);

  useBeforeUnload(
    React.useCallback(
      (event) => {
        if (allowNavigationRef.current || !hasActiveExamSession) {
          return;
        }

        event.preventDefault();
        event.returnValue = "";
      },
      [hasActiveExamSession],
    ),
  );

  React.useEffect(() => {
    if (attemptId === null) {
      return;
    }

    const handlePageHide = () => {
      if (allowNavigationRef.current || expireTriggeredRef.current) {
        return;
      }

      expireTriggeredRef.current = true;
      clearAttemptCache();
      clearSoalCache();
      expireAttemptUjianSiswaOnPageLeave(attemptId);
    };

    window.addEventListener("pagehide", handlePageHide);
    return () => {
      window.removeEventListener("pagehide", handlePageHide);
    };
  }, [attemptId, clearAttemptCache, clearSoalCache]);

  React.useEffect(() => {
    if (isTimeExpired && navigationBlocker.state === "blocked") {
      navigationBlocker.reset();
    }
  }, [isTimeExpired, navigationBlocker]);

  React.useEffect(() => {
    if (navigationBlocker.state === "blocked") {
      setIsLeaveConfirmOpen(true);
      return;
    }

    setIsLeaveConfirmOpen(false);
  }, [navigationBlocker.state]);

  const handleExpiredSubmit = React.useCallback(async () => {
    const expired = await expireAttemptBeforeLeave();
    if (!expired) {
      return;
    }

    allowNavigationRef.current = true;
    navigate(paths.dashboard.ujian_siswa);
  }, [expireAttemptBeforeLeave, navigate]);

  const handleLeaveConfirm = React.useCallback(async () => {
    const expired = await expireAttemptBeforeLeave();
    if (!expired) {
      return;
    }

    allowNavigationRef.current = true;
    setIsLeaveConfirmOpen(false);
    if (navigationBlocker.state === "blocked") {
      navigationBlocker.proceed();
      return;
    }

    navigate(paths.dashboard.ujian_siswa);
  }, [expireAttemptBeforeLeave, navigate, navigationBlocker]);

  const handleLeaveCancel = React.useCallback(() => {
    setIsLeaveConfirmOpen(false);

    if (navigationBlocker.state === "blocked") {
      navigationBlocker.reset();
    }
  }, [navigationBlocker]);

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

      <SiswaSoalPreviewContent
        key={parsedIdJadwalUjian}
        title={title}
        sisaWaktu={sisaWaktu}
        soalPreview={soalPreview}
        selectedOptions={selectedOptions}
        onSelectOption={handleSelectOption}
        onBack={() => navigate(paths.dashboard.ujian_siswa)}
      />

      <ConfirmAlert
        isOpen={isLeaveConfirmOpen}
        title="Keluar dari Ujian?"
        message="Jawaban yang belum dikirim bisa hilang. Yakin ingin meninggalkan halaman ujian ini?"
        onClose={handleLeaveCancel}
        onConfirm={() => {
          void handleLeaveConfirm();
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
          void handleExpiredSubmit();
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
