import { useCallback, useEffect, useRef, useState } from "react";
import { ApiError } from "@/services/Api/api";

/**
 * State yang dikembalikan oleh usePut.
 */
export type UsePutState<TPayload, TResponse> = {
  /** Panggil untuk menjalankan PUT/PATCH request. */
  execute: (payload: TPayload) => Promise<TResponse>;
  /** true selama request berlangsung. */
  loading: boolean;
  /** Pesan error dari request terakhir, null jika sukses. */
  error: string | null;
  /** Reset state error dan loading ke default. */
  reset: () => void;
};

/**
 * Generic hook untuk operasi PUT/PATCH (update).
 *
 * Tidak meng-trigger useEffect. Caller harus memanggil `execute()` secara
 * manual, misalnya di dalam event handler (onSubmit, onClick, dsb).
 *
 * @param action - Fungsi async yang menerima payload dan mengembalikan response.
 *   Biasanya memanggil service function yang menggunakan `api()` wrapper.
 *
 * @example
 * ```ts
 * const { execute, loading, error } = usePut(
 *   (payload: { id: number; values: UpdatePayload }) =>
 *     updateMataPelajaran(payload.id, payload.values),
 * );
 * ```
 */
export default function usePut<TPayload, TResponse = unknown>(
  action: (payload: TPayload) => Promise<TResponse>,
): UsePutState<TPayload, TResponse> {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const mountedRef = useRef(true);

  useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
    };
  }, []);

  const execute = useCallback(
    async (payload: TPayload): Promise<TResponse> => {
      setLoading(true);
      setError(null);

      try {
        const result = await action(payload);
        if (mountedRef.current) {
          setLoading(false);
        }
        return result;
      } catch (e: unknown) {
        const message =
          e instanceof ApiError
            ? e.message
            : e instanceof Error
              ? e.message
              : "Terjadi kesalahan yang tidak diketahui.";

        if (mountedRef.current) {
          setError(message);
          setLoading(false);
        }

        throw e;
      }
    },
    [action],
  );

  const reset = useCallback(() => {
    setError(null);
    setLoading(false);
  }, []);

  return { execute, loading, error, reset };
}
