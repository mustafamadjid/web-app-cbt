import { useCallback, useEffect, useRef, useState } from "react";
import { ApiError } from "@/services/Api/api";

/**
 * State yang dikembalikan oleh usePost.
 */
export type UsePostState<TPayload, TResponse> = {
  /** Panggil untuk menjalankan POST request. */
  execute: (payload: TPayload) => Promise<TResponse>;
  /** true selama request berlangsung. */
  loading: boolean;
  /** Pesan error dari request terakhir, null jika sukses. */
  error: string | null;
  /** Reset state error dan loading ke default. */
  reset: () => void;
};

/**
 * Generic hook untuk operasi POST (create).
 *
 * Tidak meng-trigger useEffect. Caller harus memanggil `execute()` secara
 * manual, misalnya di dalam event handler (onSubmit, onClick, dsb).
 *
 * @param action - Fungsi async yang menerima payload dan mengembalikan response.
 *   Biasanya memanggil service function yang menggunakan `api()` wrapper.
 *
 * @example
 * ```ts
 * const { execute, loading, error } = usePost(
 *   (values: CreateMapelPayload) => createMataPelajaran(values),
 * );
 *
 * const handleSubmit = async (values) => {
 *   try {
 *     await execute(values);
 *     toast.success("Berhasil!");
 *   } catch {
 *     // error sudah di-set di hook state
 *   }
 * };
 * ```
 */
export default function usePost<TPayload, TResponse = unknown>(
  action: (payload: TPayload) => Promise<TResponse>,
): UsePostState<TPayload, TResponse> {
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

        throw e; // Re-throw agar caller bisa handle di try/catch
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
