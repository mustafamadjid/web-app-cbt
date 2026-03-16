type DeriveUjianMulaiSiswaStateParams = {
  hasIdJadwalUjianParam: boolean;
  isIdJadwalUjianValid: boolean;
  isIdSiswaValid: boolean;
  isAuthLoading: boolean;
  loadingActiveAttempt: boolean;
  loadingSoal: boolean;
  loadingWaktuSelesai: boolean;
  activeAttemptError: string | null;
  activeAttemptErrorCode: string | null;
  soalError: string | null;
  waktuSelesaiError: string | null;
  waktuSelesai: string;
  hasValidWaktuSelesai: boolean;
};

type UjianMulaiSiswaState = {
  loading: boolean;
  errorMessage: string | null;
};

const getWaktuSelesaiStateError = ({
  loadingWaktuSelesai,
  waktuSelesai,
  hasValidWaktuSelesai,
}: Pick<
  DeriveUjianMulaiSiswaStateParams,
  "loadingWaktuSelesai" | "waktuSelesai" | "hasValidWaktuSelesai"
>): string | null => {
  if (!loadingWaktuSelesai && !waktuSelesai) {
    return "Waktu selesai ujian tidak tersedia.";
  }

  if (waktuSelesai && !hasValidWaktuSelesai) {
    return "Format waktu selesai ujian tidak valid.";
  }

  return null;
};

const findFirstError = (...errors: Array<string | null>): string | null =>
  errors.find((error): error is string => Boolean(error)) ?? null;

export function deriveUjianMulaiSiswaState({
  hasIdJadwalUjianParam,
  isIdJadwalUjianValid,
  isIdSiswaValid,
  isAuthLoading,
  loadingActiveAttempt,
  loadingSoal,
  loadingWaktuSelesai,
  activeAttemptError,
  activeAttemptErrorCode,
  soalError,
  waktuSelesaiError,
  waktuSelesai,
  hasValidWaktuSelesai,
}: DeriveUjianMulaiSiswaStateParams): UjianMulaiSiswaState {
  const shouldRedirectToToken =
    isIdJadwalUjianValid && activeAttemptErrorCode === "NOT_FOUND";
  const loading =
    isAuthLoading ||
    loadingActiveAttempt ||
    loadingSoal ||
    loadingWaktuSelesai ||
    shouldRedirectToToken;
  const waktuSelesaiStateError = getWaktuSelesaiStateError({
    loadingWaktuSelesai,
    waktuSelesai,
    hasValidWaktuSelesai,
  });

  if (!hasIdJadwalUjianParam || !isIdJadwalUjianValid) {
    return {
      loading,
      errorMessage: "Jadwal ujian tidak ditemukan.",
    };
  }

  if (!isAuthLoading && !isIdSiswaValid) {
    return {
      loading,
      errorMessage: "Akun siswa tidak valid. Silakan login ulang.",
    };
  }

  if (activeAttemptErrorCode !== "NOT_FOUND") {
    return {
      loading,
      errorMessage: findFirstError(
        activeAttemptError,
        soalError,
        waktuSelesaiError,
        waktuSelesaiStateError,
      ),
    };
  }

  return {
    loading,
    errorMessage: findFirstError(
      soalError,
      waktuSelesaiError,
      waktuSelesaiStateError,
    ),
  };
}
