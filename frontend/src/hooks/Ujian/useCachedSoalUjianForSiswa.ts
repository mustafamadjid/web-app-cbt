import React from "react";

import { hasRichContent } from "@/types/Content/RichContent";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import { GetSoalUjianForSiswa } from "@/services/Api/features-api/Ujian/soalUjian.service";
import type {
  OpsiJawabanSoalUjianSiswa,
  SoalUjianSiswa,
} from "@/types/Ujian/SoalUjian";

const CACHE_VERSION = 2;
const STORAGE_PREFIX = `ujian_soal_siswa:v${CACHE_VERSION}`;

type CachedSoalUjianPayload = {
  version: number;
  siswaId: number;
  idJadwalUjian: number;
  soalRows: SoalUjianSiswa[];
};

export type UseCachedSoalUjianForSiswaState = {
  data: SoalUjianSiswa[] | null;
  error: string | null;
  loading: boolean;
  clearCache: () => void;
};

const isValidOption = (value: unknown): value is OpsiJawabanSoalUjianSiswa => {
  if (!value || typeof value !== "object") return false;

  const option = value as Partial<OpsiJawabanSoalUjianSiswa>;
  return (
    typeof option.id_pilihan_ganda === "number" &&
    typeof option.isi_pilihan === "string" &&
    (option.isi_pilihan_content == null ||
      (typeof option.isi_pilihan_content === "object" &&
        hasRichContent(option.isi_pilihan_content as never)))
  );
};

const isValidSoalRow = (value: unknown): value is SoalUjianSiswa => {
  if (!value || typeof value !== "object") return false;

  const soal = value as Partial<SoalUjianSiswa>;
  return (
    typeof soal.id_soal === "number" &&
    typeof soal.tipe_soal === "string" &&
    typeof soal.pertanyaan === "string" &&
    (soal.pertanyaan_content == null ||
      (typeof soal.pertanyaan_content === "object" &&
        hasRichContent(soal.pertanyaan_content as never))) &&
    typeof soal.gambar === "string" &&
    typeof soal.bobot_soal === "number" &&
    typeof soal.no_urut_soal === "number" &&
    Array.isArray(soal.opsi_jawaban) &&
    soal.opsi_jawaban.every(isValidOption)
  );
};

const buildCacheKey = (siswaId: number, idJadwalUjian: number) =>
  `${STORAGE_PREFIX}:${siswaId}:${idJadwalUjian}`;

const clearCachedSoalUjian = (siswaId: number, idJadwalUjian: number) => {
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
): value is CachedSoalUjianPayload => {
  if (!value || typeof value !== "object") return false;

  const cache = value as Partial<CachedSoalUjianPayload>;
  return (
    cache.version === CACHE_VERSION &&
    cache.siswaId === siswaId &&
    cache.idJadwalUjian === idJadwalUjian &&
    Array.isArray(cache.soalRows) &&
    cache.soalRows.every(isValidSoalRow)
  );
};

const readCachedSoalUjian = (
  siswaId: number,
  idJadwalUjian: number,
): SoalUjianSiswa[] | null => {
  if (typeof window === "undefined") {
    return null;
  }

  const cacheKey = buildCacheKey(siswaId, idJadwalUjian);

  try {
    const raw = window.localStorage.getItem(cacheKey);

    if (!raw) {
      return null;
    }

    const parsed = JSON.parse(raw) as unknown;
    const valid = isValidCachePayload(parsed, siswaId, idJadwalUjian);

    if (!valid) {
      window.localStorage.removeItem(cacheKey);
      return null;
    }

    return parsed.soalRows;
  } catch {
    try {
      window.localStorage.removeItem(cacheKey);
    } catch {
      // Ignore storage cleanup errors.
    }
    return null;
  }
};

const writeCachedSoalUjian = (
  siswaId: number,
  idJadwalUjian: number,
  soalRows: SoalUjianSiswa[],
) => {
  if (typeof window === "undefined") return;

  const payload: CachedSoalUjianPayload = {
    version: CACHE_VERSION,
    siswaId,
    idJadwalUjian,
    soalRows,
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

const mapErrorMessage = (error: unknown) => {
  return getUserFriendlyErrorMessage(error, {
    action: "fetch",
    entity: "soal ujian",
  });
};

export function useCachedSoalUjianForSiswa(
  idJadwalUjian: number,
  siswaId: number,
  enabled = true,
): UseCachedSoalUjianForSiswaState {
  const hasValidIds =
    Number.isInteger(idJadwalUjian) &&
    idJadwalUjian > 0 &&
    Number.isInteger(siswaId) &&
    siswaId > 0;

  const hasValidContext = enabled && hasValidIds;

  const initialCachedData = React.useMemo(() => {
    if (!hasValidContext) return null;
    return readCachedSoalUjian(siswaId, idJadwalUjian);
  }, [hasValidContext, siswaId, idJadwalUjian]);

  const [data, setData] = React.useState<SoalUjianSiswa[] | null>(
    initialCachedData,
  );
  const [error, setError] = React.useState<string | null>(null);
  const [loading, setLoading] = React.useState<boolean>(
    hasValidContext && initialCachedData === null,
  );

  const clearCache = React.useCallback(() => {
    if (!hasValidIds) return;
    clearCachedSoalUjian(siswaId, idJadwalUjian);
  }, [hasValidIds, siswaId, idJadwalUjian]);

  React.useEffect(() => {
    if (!hasValidContext) {
      setData(null);
      setError(null);
      setLoading(false);
      return;
    }

    const cachedRows = readCachedSoalUjian(siswaId, idJadwalUjian);
    if (cachedRows !== null) {
      setData(cachedRows);
      setError(null);
      setLoading(false);
      return;
    }

    let cancelled = false;

    setData(null);
    setError(null);
    setLoading(true);

    void (async () => {
      try {
        const result = await GetSoalUjianForSiswa(idJadwalUjian);
        if (cancelled) return;

        setData(result);
        writeCachedSoalUjian(siswaId, idJadwalUjian, result);
      } catch (fetchError) {
        if (cancelled) return;

        setData(null);
        setError(mapErrorMessage(fetchError));
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

  return { data, error, loading, clearCache };
}
