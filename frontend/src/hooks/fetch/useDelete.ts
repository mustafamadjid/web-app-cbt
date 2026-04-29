import { useCallback, useEffect, useRef, useState } from "react";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";

/**
 * State yang dikembalikan oleh useDelete.
 */
export type UseDeleteState<TPayload, TResponse> = {
  /** Panggil untuk menjalankan DELETE request. */
  execute: (payload: TPayload) => Promise<TResponse>;
  /** true selama request berlangsung. */
  loading: boolean;
  /** Pesan error dari request terakhir, null jika sukses. */
  error: string | null;
  /** Reset state error dan loading ke default. */
  reset: () => void;
};

/**
 * Generic hook untuk operasi DELETE.
 *
 * Tidak meng-trigger useEffect. Caller harus memanggil `execute()` secara
 * manual, misalnya di dalam event handler (onClick, onConfirm, dsb).
 *
 * @param action - Fungsi async yang menerima payload dan mengembalikan response.
 *   Biasanya memanggil service function yang menggunakan `api()` wrapper.
 *
 * @example
 * ```ts
 * const { execute, loading, error } = useDelete(
 *   (id: number) => deleteMataPelajaran(id),
 * );
 *
 * const handleDelete = async (id: number) => {
 *   try {
 *     await execute(id);
 *     toast.success("Berhasil dihapus!");
 *   } catch {
 *     // error sudah di-set di hook state
 *   }
 * };
 * ```
 */
export default function useDelete<TPayload, TResponse = unknown>(
  action: (payload: TPayload) => Promise<TResponse>,
): UseDeleteState<TPayload, TResponse> {
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
        const message = getUserFriendlyErrorMessage(e, {
          action: "delete",
          entity: "data",
        });

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
