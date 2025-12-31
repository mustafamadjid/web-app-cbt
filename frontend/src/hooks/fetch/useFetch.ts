import { useCallback, useEffect, useRef, useState } from "react";
import axios from "axios";
import { api } from "@/services/Api/api";

type UseFetchState<T> = {
  data: T | null;
  error: string | null;
  loading: boolean;
  refetch: () => Promise<void>;
};

export default function useFetch<T>(url: string): UseFetchState<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const controllerRef = useRef<AbortController | null>(null);
  const requestIdRef = useRef(0);

  const run = useCallback(async () => {
    // batalkan request sebelumnya kalau masih jalan
    controllerRef.current?.abort();

    const controller = new AbortController();
    controllerRef.current = controller;

    const requestId = ++requestIdRef.current;

    setLoading(true);
    setError(null);

    try {
      const res = await api.get<T>(url, { signal: controller.signal });

      // abaikan kalau ada request baru yang sudah dimulai
      if (requestId !== requestIdRef.current) return;

      setData(res.data);
    } catch (e: unknown) {
      // Abort bukan error yang perlu ditampilkan
      if (e instanceof DOMException && e.name === "AbortError") return;

      if (requestId !== requestIdRef.current) return;

      if (axios.isAxiosError(e)) {
        setError((e.response?.data as any)?.message ?? e.message);
      } else if (e instanceof Error) {
        setError(e.message);
      } else {
        setError("Unknown error");
      }

      setData(null);
    } finally {
      if (requestId === requestIdRef.current) {
        setLoading(false);
      }
    }
  }, [url]);

  useEffect(() => {
    void run();

    return () => {
      controllerRef.current?.abort();
    };
  }, [run]);

  return { data, error, loading, refetch: run };
}
