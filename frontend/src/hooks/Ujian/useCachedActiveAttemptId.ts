import React from "react";
import { ApiError } from "@/services/Api/api";
import { getActiveAttemptUjian } from "@/services/Api/features-api/Ujian/ujian.service";

const CACHE_VERSION = 1;
const STORAGE_PREFIX = `ujian_active_attempt:v${CACHE_VERSION}`;

type CachedActiveAttemptPayload = {
  version: number;
  siswaId: number;
  idJadwalUjian: number;
  idAttempt: number;
};

export type UseCachedActiveAttemptIdState = {
  attemptId: number | null;
  error: string | null;
  errorCode: string | null;
  loading: boolean;
  clearCache: () => void;
};

const buildCacheKey = (siswaId: number, idJadwalUjian: number) =>
  `${STORAGE_PREFIX}:${siswaId}:${idJadwalUjian}`;

const clearCachedActiveAttemptId = (siswaId: number, idJadwalUjian: number) => {
  if (typeof window === "undefined") return;

  try {
    window.localStorage.removeItem(buildCacheKey(siswaId, idJadwalUjian));
  } catch {
    // Ignore storage access errors.
  }
};

const isValidCachePayload = (
  value: unknown,
  siswaId: number,
  idJadwalUjian: number,
): value is CachedActiveAttemptPayload => {
  if (!value || typeof value !== "object") return false;

  const cache = value as Partial<CachedActiveAttemptPayload>;
  return (
    cache.version === CACHE_VERSION &&
    cache.siswaId === siswaId &&
    cache.idJadwalUjian === idJadwalUjian &&
    typeof cache.idAttempt === "number" &&
    Number.isInteger(cache.idAttempt) &&
    cache.idAttempt > 0
  );
};

const readCachedActiveAttemptId = (
  siswaId: number,
  idJadwalUjian: number,
): number | null => {
  if (typeof window === "undefined") return null;

  const cacheKey = buildCacheKey(siswaId, idJadwalUjian);

  try {
    const raw = window.localStorage.getItem(cacheKey);
    if (!raw) return null;

    const parsed = JSON.parse(raw) as unknown;
    if (!isValidCachePayload(parsed, siswaId, idJadwalUjian)) {
      window.localStorage.removeItem(cacheKey);
      return null;
    }

    return parsed.idAttempt;
  } catch {
    try {
      window.localStorage.removeItem(cacheKey);
    } catch {
      // Ignore storage cleanup errors.
    }

    return null;
  }
};

const writeCachedActiveAttemptId = (
  siswaId: number,
  idJadwalUjian: number,
  idAttempt: number,
) => {
  if (typeof window === "undefined") return;

  const payload: CachedActiveAttemptPayload = {
    version: CACHE_VERSION,
    siswaId,
    idJadwalUjian,
    idAttempt,
  };

  try {
    window.localStorage.setItem(
      buildCacheKey(siswaId, idJadwalUjian),
      JSON.stringify(payload),
    );
  } catch {
    // Ignore storage access errors.
  }
};

const mapErrorState = (error: unknown) => {
  if (error instanceof ApiError) {
    return {
      error: error.message,
      errorCode: error.code ?? null,
    };
  }

  if (error instanceof Error) {
    return {
      error: error.message,
      errorCode: null,
    };
  }

  return {
    error: "Terjadi kesalahan yang tidak diketahui.",
    errorCode: null,
  };
};

export function useCachedActiveAttemptId(
  idJadwalUjian: number,
  siswaId: number,
  enabled = true,
): UseCachedActiveAttemptIdState {
  const hasValidIds =
    Number.isInteger(idJadwalUjian) &&
    idJadwalUjian > 0 &&
    Number.isInteger(siswaId) &&
    siswaId > 0;

  const hasValidContext = enabled && hasValidIds;

  const initialCachedAttemptId = React.useMemo(() => {
    if (!hasValidContext) return null;
    return readCachedActiveAttemptId(siswaId, idJadwalUjian);
  }, [hasValidContext, siswaId, idJadwalUjian]);

  const [attemptId, setAttemptId] = React.useState<number | null>(
    initialCachedAttemptId,
  );
  const [error, setError] = React.useState<string | null>(null);
  const [errorCode, setErrorCode] = React.useState<string | null>(null);
  const [loading, setLoading] = React.useState<boolean>(
    hasValidContext && initialCachedAttemptId === null,
  );

  const clearCache = React.useCallback(() => {
    if (!hasValidIds) return;
    clearCachedActiveAttemptId(siswaId, idJadwalUjian);
  }, [hasValidIds, siswaId, idJadwalUjian]);

  React.useEffect(() => {
    if (!hasValidContext) {
      setAttemptId(null);
      setError(null);
      setErrorCode(null);
      setLoading(false);
      return;
    }

    const cachedAttemptId = readCachedActiveAttemptId(siswaId, idJadwalUjian);
    if (cachedAttemptId !== null) {
      setAttemptId(cachedAttemptId);
      setError(null);
      setErrorCode(null);
      setLoading(false);
      return;
    }

    let cancelled = false;

    setAttemptId(null);
    setError(null);
    setErrorCode(null);
    setLoading(true);

    void (async () => {
      try {
        const result = await getActiveAttemptUjian(idJadwalUjian);
        if (cancelled) return;

        if (
          !Number.isInteger(result.id_attempt) ||
          result.id_attempt <= 0
        ) {
          setAttemptId(null);
          setError("ID attempt aktif tidak valid.");
          setErrorCode(null);
          return;
        }

        setAttemptId(result.id_attempt);
        writeCachedActiveAttemptId(siswaId, idJadwalUjian, result.id_attempt);
      } catch (fetchError) {
        if (cancelled) return;

        const nextErrorState = mapErrorState(fetchError);
        setAttemptId(null);
        setError(nextErrorState.error);
        setErrorCode(nextErrorState.errorCode);
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [hasValidContext, siswaId, idJadwalUjian]);

  return { attemptId, error, errorCode, loading, clearCache };
}
