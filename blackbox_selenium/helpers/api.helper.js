import { env } from "../config/env.js";
import { accountFor } from "../data/accounts.js";

let tokenByRole = new Map();

async function request(path, options = {}) {
  const response = await fetch(`${env.apiUrl}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : null;
  if (!response.ok) {
    throw new Error(`API ${options.method || "GET"} ${path} gagal: ${response.status} ${text}`);
  }
  return data;
}

export async function apiLogin(role = "admin") {
  const account = accountFor(role);
  const data = await request("/auth/login", {
    method: "POST",
    body: JSON.stringify({ username: account.username, password: account.password }),
  });
  const token = data?.data?.token || data?.token || data?.access_token;
  if (token) tokenByRole.set(role, token);
  return data;
}

function authHeader(role) {
  const token = tokenByRole.get(role);
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiGet(path, role = "admin") {
  return request(path, { headers: authHeader(role) });
}

export async function apiPost(path, body, role = "admin") {
  return request(path, {
    method: "POST",
    headers: authHeader(role),
    body: JSON.stringify(body),
  });
}

export async function apiDelete(path, role = "admin") {
  return request(path, {
    method: "DELETE",
    headers: authHeader(role),
  });
}

export async function ensureBaseData() {
  throw new Error("ensureBaseData belum diikat ke endpoint seed E2E. Siapkan data melalui UI atau endpoint khusus test.");
}

export async function cleanupE2EData() {
  throw new Error("cleanupE2EData belum diikat ke endpoint cleanup E2E.");
}
