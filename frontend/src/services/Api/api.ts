// services/api.ts
import axios from "axios";
import type { AxiosError, AxiosRequestConfig } from "axios";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:3000";

// Omit -> menghilangkan property dari type
export type RequestOptions = Omit<
  AxiosRequestConfig,
  "url" | "baseURL" | "headers"
> & {
  token?: string | null;
  headers?: Record<string, string>;
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

const client = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// helper utama
export async function api<T>(
  path: string,
  opts: RequestOptions = {}
): Promise<T> {
  const { token, headers, ...rest } = opts;

  try {
    const res = await client.request<T>({
      url: path,
      ...rest,
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(headers ?? {}),
      },
    });

   
    return res.data;
  } catch (err) {
    const e = err as AxiosError;

    
    const status = e.response?.status;
    const data = e.response?.data;

  
    const backendMsg =
      typeof data === "object" && data && "message" in (data as any)
        ? String((data as any).message)
        : undefined;

    throw new ApiError(backendMsg ?? e.message, status, data);
  }
}
