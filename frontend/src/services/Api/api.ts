// services/api.ts
import axios from "axios";
import type { AxiosError, AxiosRequestConfig } from "axios";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";

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
  rawMessage?: string;

  constructor(
    message: string,
    status?: number,
    data?: unknown,
    code?: string,
    rawMessage?: string,
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
    this.code = code;
    this.rawMessage = rawMessage ?? message;
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
  const data = e.response?.data as
    | {
        error?: { message?: unknown; code?: unknown } | null;
        message?: unknown;
      }
    | undefined;

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

  const rawMessage = backendErrorMsg ?? backendTopMsg ?? e.message;
  const apiError = new ApiError(
    rawMessage,
    status,
    data,
    backendErrorCode,
    rawMessage,
  );
  apiError.message = getUserFriendlyErrorMessage(apiError);
  return apiError;
}

/**
 * Refresh session
 */
async function refreshSession(): Promise<void> {
  await refreshClient.post("/auth/refresh");
}

/** unwrap envelope sukses */
function unwrapEnvelope<T>(env: ApiEnvelope<T>, path?: string): T {
  if (env?.error)
    throw new ApiError(
      getUserFriendlyErrorMessage(
        new ApiError(env.error.message, 400, env, env.error.code, env.error.message),
      ),
      400,
      env,
      env.error.code,
      env.error.message,
    );

  if (!env || typeof env !== "object" || !("data" in env)) {
    throw new ApiError(
      "Terjadi kendala pada sistem. Silakan coba lagi.",
      500,
      env,
      undefined,
      "Invalid API response envelope",
    );
  }

  if (env.data === null) {
    if (path === "/auth/me") {
      throw new ApiError(
        "Sesi Anda telah berakhir. Silakan login kembali.",
        401,
        env,
        "UNAUTHENTICATED",
        "Unauthenticated",
      );
    }
    throw new ApiError(
      "Terjadi kendala pada sistem. Silakan coba lagi.",
      500,
      env,
      undefined,
      "API returned null data",
    );
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

    if (apiErr.status !== 401 || _retry) throw apiErr;

   
    const p = String(path);
    if (p.startsWith("/auth/login") || p.startsWith("/auth/refresh"))
      throw apiErr;

  
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
      if (rErr.status === 401) {
        rErr.code = rErr.code ?? "SESSION_EXPIRED";
        rErr.message = getUserFriendlyErrorMessage(rErr);
      }
      flushSubscribers(false, rErr);
      throw rErr;
    } finally {
      isRefreshing = false;
    }
  }
}
