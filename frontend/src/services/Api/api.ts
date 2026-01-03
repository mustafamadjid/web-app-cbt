// services/api.ts
import axios from "axios";
import type { AxiosError, AxiosRequestConfig } from "axios";

import { authToken } from "../auth/token";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:3000";

export type RequestOptions = Omit<
  AxiosRequestConfig,
  "url" | "baseURL" | "headers"
> & {
  /** Optional: override token untuk request ini saja */
  token?: string | null;
  headers?: Record<string, string>;
  /** Internal */
  _retry?: boolean;
};

export class ApiError extends Error {
  status?: number;
  data?: unknown;

  constructor(message: string, status?: number, data?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}



/** ===== Axios clients ===== */
const client = axios.create({
  baseURL: API_URL,
  withCredentials: true, // penting: kirim refresh cookie
  headers: {
    "Content-Type": "application/json",
  },
});

// client khusus refresh supaya tidak terjebak interceptor/loop
const refreshClient = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

/** ===== Refresh concurrency control ===== */
let isRefreshing = false;
let pending: Array<(t: string | null) => void> = [];

function subscribeRefresh(cb: (t: string | null) => void) {
  pending.push(cb);
}
function flushSubscribers(t: string | null) {
  pending.forEach((cb) => cb(t));
  pending = [];
}

async function refreshAccessToken(): Promise<string> {
  const r = await refreshClient.post<{ accessToken: string }>("/auth/refresh");
  const newToken = r.data?.accessToken;

  if (!newToken || typeof newToken !== "string") {
    throw new ApiError(
      "Refresh succeeded but no accessToken returned",
      500,
      r.data
    );
  }

  authToken.set(newToken);
  return newToken;
}

/** helper untuk convert AxiosError -> ApiError */
function toApiError(err: unknown): ApiError {
  const e = err as AxiosError;
  const status = e.response?.status;
  const data = e.response?.data;

  const backendMsg =
    typeof data === "object" && data && "message" in (data as any)
      ? String((data as any).message)
      : undefined;

  return new ApiError(backendMsg ?? e.message, status, data);
}

/** ===== helper utama ===== */
export async function api<T>(
  path: string,
  opts: RequestOptions = {}
): Promise<T> {
  const { token, headers, _retry, ...rest } = opts;

  // token precedence: opts.token > global authToken
  const usedToken = token ?? authToken.get();

  try {
    const res = await client.request<T>({
      url: path,
      ...rest,
      headers: {
        ...(usedToken ? { Authorization: `Bearer ${usedToken}` } : {}),
        ...(headers ?? {}),
      },
    });

    return res.data;
  } catch (err) {
    const apiErr = toApiError(err);

    // hanya tangani 401 dan hanya sekali retry
    if (apiErr.status !== 401 || _retry) {
      throw apiErr;
    }

    // kalau refresh sedang berjalan, tunggu hasilnya lalu retry
    if (isRefreshing) {
      return new Promise<T>((resolve, reject) => {
        subscribeRefresh(async (t) => {
          if (!t) return reject(apiErr);
          try {
            const data = await api<T>(path, {
              ...opts,
              token: t,
              _retry: true,
            });
            resolve(data);
          } catch (e2) {
            reject(e2);
          }
        });
      });
    }

    // mulai refresh
    isRefreshing = true;
    try {
      const newToken = await refreshAccessToken();
      flushSubscribers(newToken);

      // retry request original dengan token baru
      return await api<T>(path, { ...opts, token: newToken, _retry: true });
    } catch (refreshErr) {
      authToken.clear();
      flushSubscribers(null);
      throw toApiError(refreshErr);
    } finally {
      isRefreshing = false;
    }
  }
}
