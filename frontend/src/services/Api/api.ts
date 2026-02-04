// services/api.ts
import axios from "axios";
import type { AxiosError, AxiosRequestConfig } from "axios";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export type ApiErrorBody = {
  code: string;
  message: string;
};

export type ApiEnvelope<T> = {
  data: T | null;
  message?: string;
  meta?: unknown;
  error?: ApiErrorBody | null;
};

export type RequestOptions = Omit<
  AxiosRequestConfig,
  "url" | "baseURL" | "headers"
> & {
  headers?: Record<string, string>;

  _retry?: boolean;
};

export class ApiError extends Error {
  status?: number;
  data?: unknown;
  code?: string;

  constructor(message: string, status?: number, data?: unknown, code?: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
    this.code = code;
  }
}

/** ===== Axios clients (cookie-based auth) ===== */
const client = axios.create({
  baseURL: API_URL,
  withCredentials: true,
});


const refreshClient = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

/** ===== Refresh concurrency control (single-flight) ===== */
let isRefreshing = false;
let pending: Array<(ok: boolean, err?: ApiError) => void> = [];

function subscribeRefresh(cb: (ok: boolean, err?: ApiError) => void) {
  pending.push(cb);
}
function flushSubscribers(ok: boolean, err?: ApiError) {
  pending.forEach((cb) => cb(ok, err));
  pending = [];
}

/** Convert AxiosError -> ApiError, paham envelope backend kamu */
function toApiError(err: unknown): ApiError {
  const e = err as AxiosError;
  const status = e.response?.status;
  const data = e.response?.data as any;

  // backend error: { data:null, error:{code,message} }
  const backendErrorMsg =
    data?.error?.message && typeof data.error.message === "string"
      ? data.error.message
      : undefined;

  const backendErrorCode =
    data?.error?.code && typeof data.error.code === "string"
      ? data.error.code
      : undefined;

  // backend sukses biasanya punya "message" top-level (optional)
  const backendTopMsg =
    data?.message && typeof data.message === "string"
      ? data.message
      : undefined;

  const msg = backendErrorMsg ?? backendTopMsg ?? e.message;
  return new ApiError(msg, status, data, backendErrorCode);
}

/**
 * Refresh session
 */
async function refreshSession(): Promise<void> {
  await refreshClient.post("/auth/login/refresh");
}

/** unwrap envelope sukses */
function unwrapEnvelope<T>(env: ApiEnvelope<T>, path?: string): T {
  if (env?.error)
    throw new ApiError(env.error.message, 400, env, env.error.code);

  if (!env || typeof env !== "object" || !("data" in env)) {
    throw new ApiError("Invalid API response envelope", 500, env);
  }

  if (env.data === null) {
    // treat null as unauthenticated ONLY for auth-check endpoints
    if (path && (path === "/auth/me" || path.startsWith("/auth/"))) {
      throw new ApiError("Unauthenticated", 401, env, "UNAUTHENTICATED");
    }
    throw new ApiError("API returned null data", 500, env);
  }

  return env.data;
}


/** helper utama */
export async function api<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const { headers, _retry, ...rest } = opts;

  try {
    const res = await client.request<ApiEnvelope<T>>({
      url: path,
      ...rest,
      headers: { ...(headers ?? {}) },
    });

    return unwrapEnvelope(res.data,path);
  } catch (err) {
    const apiErr = toApiError(err);

    // hanya tangani 401 dan hanya sekali retry
    if (apiErr.status !== 401 || _retry) throw apiErr;

    // jangan refresh kalau request ini sendiri adalah refresh endpoint
    if (String(path).includes("/auth/login/refresh")) throw apiErr;

    // kalau refresh sedang berjalan, tunggu hasilnya lalu retry
    if (isRefreshing) {
      return new Promise<T>((resolve, reject) => {
        subscribeRefresh(async (ok, refreshErr) => {
          if (!ok) return reject(refreshErr ?? apiErr);
          try {
            const data = await api<T>(path, { ...opts, _retry: true });
            resolve(data);
          } catch (e2) {
            reject(e2);
          }
        });
      });
    }

    // mulai refresh (single-flight)
    isRefreshing = true;
    try {
      await refreshSession();
      flushSubscribers(true);

      // retry request original setelah refresh sukses
      return await api<T>(path, { ...opts, _retry: true });
    } catch (refreshErr) {
      const rErr = toApiError(refreshErr);
       if (rErr.status === 401 ) {
         rErr.code = rErr.code ?? "SESSION_EXPIRED";
       }
      flushSubscribers(false, rErr);
      throw rErr;
    } finally {
      isRefreshing = false;
    }
  }
}
