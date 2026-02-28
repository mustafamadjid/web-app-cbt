import { useCallback, useEffect, useRef, useState } from "react";
import type { DependencyList } from "react";
import { ApiError } from "@/services/Api/api";

/**
 * State yang dikembalikan oleh useFetch.
 */
export type UseFetchState<T> = {
  /** Data hasil fetch, null saat belum ada atau error. */
  data: T | null;
  /** Pesan error, null jika tidak ada error. */
  error: string | null;
  /** true selama proses fetching berlangsung. */
  loading: boolean;
  /** Panggil untuk memicu re-fetch manual. */
  refetch: () => Promise<void>;
};

/**
 * Generic hook untuk operasi GET (read) yang dipicu otomatis oleh useEffect.
 *
 * @param fetcher - Fungsi async yang mengembalikan data. Boleh memanggil
 *   service function apapun yang menggunakan `api()` wrapper dari `api.ts`.
 * @param deps - Dependency array; setiap kali berubah, fetcher dijalankan ulang.
 *
 * @example
 * ```ts
 * const { data, loading, error, refetch } = useFetch(
 *   () => getMapel({ search, tingkatKelas }),
 *   [search, tingkatKelas],
 * );
 * ```
 */
export default function useFetch<T>(
  fetcher: () => Promise<T>,
  deps: DependencyList = [],
): UseFetchState<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  // Untuk mencegah race condition saat dependency berubah cepat
  const requestIdRef = useRef(0);
  // Untuk mencegah state update setelah unmount
  const mountedRef = useRef(true);

  const run = useCallback(async () => {
    const requestId = ++requestIdRef.current;

    setLoading(true);
    setError(null);

    try {
      const result = await fetcher();

      // Abaikan response jika sudah ada request baru atau komponen unmount
      if (requestId !== requestIdRef.current || !mountedRef.current) return;

      setData(result);
    } catch (e: unknown) {
      if (requestId !== requestIdRef.current || !mountedRef.current) return;

      if (e instanceof ApiError) {
        setError(e.message);
      } else if (e instanceof Error) {
        setError(e.message);
      } else {
        setError("Terjadi kesalahan yang tidak diketahui.");
      }

      setData(null);
    } finally {
      if (requestId === requestIdRef.current && mountedRef.current) {
        setLoading(false);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps);

  useEffect(() => {
    mountedRef.current = true;
    void run();

    return () => {
      mountedRef.current = false;
    };
  }, [run]);

  return { data, error, loading, refetch: run };
}
