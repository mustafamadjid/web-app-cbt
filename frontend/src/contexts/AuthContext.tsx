// src/auth/AuthContext.tsx
import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { api, ApiError } from "@/services/Api/api";
import { authToken } from "@/services/auth/token";

type User = {
  username: string;
  role: string;
};

type AuthStatus = "loading" | "authenticated" | "guest";

type LoginPayload = {
  username: string;
  password: string;
};

type AuthContextValue = {
  user: User | null;
  status: AuthStatus;
  login: (payload: LoginPayload) => Promise<void>;
  logout: () => Promise<void>;
  refetchMe: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

// Hanya untuk mode development
function getDebugAuthUser(): { username: string; role: string } | null {
  if (!import.meta.env.DEV) return null;
  const raw = localStorage.getItem("debug:auth");
  if (!raw) return null;

  try {
    const parsed = JSON.parse(raw) as { username?: string; role?: string };
    if (!parsed.username || !parsed.role) return null;
    return { username: parsed.username, role: parsed.role };
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [status, setStatus] = useState<AuthStatus>("loading");

  const refetchMe = async () => {
    const me = await api<User>("/auth/me", { method: "GET" });
    setUser({ username: me.username, role: me.role });
    setStatus("authenticated");
  };

  const boot = async () => {
    setStatus("loading");

    // Hanya untuk development
    const debugUser = getDebugAuthUser();
    if (debugUser) {
      setUser(debugUser);
      setStatus("authenticated");
      return; // stop: ini mock total

      // Jalankan di console browser
      // localStorage.setItem("debug:auth", JSON.stringify({ username: "dev", role: "ADMIN" }))
    }
    try {
      // Catatan: kalau access token belum ada, /me kemungkinan 401.
      // Wrapper api() akan mencoba refresh otomatis jika 401 (karena interceptor-like logic ada di wrapper).
      await refetchMe();
    } catch (e) {
      // Jika refresh gagal / tidak ada sesi refresh cookie, jatuh ke guest.
      authToken.clear();
      setUser(null);
      setStatus("guest");
    }
  };

  useEffect(() => {
    void boot();
  }, []);

  const login = async (payload: LoginPayload) => {
    // Login: dapat accessToken, simpan di memory
    const res = await api<{ accessToken: string }>("/auth/login", {
      method: "POST",
      data: payload,
    });

    authToken.set(res.accessToken);

    // Ambil user (source of truth)
    await refetchMe();
  };

  const logout = async () => {
    try {
      await api<void>("/auth/logout", { method: "POST" });
    } catch (e) {
      const err = e as ApiError;
      console.warn("Logout error:", err.status, err.message);
    } finally {
      authToken.clear();
      setUser(null);
      setStatus("guest");
    }
  };

  const value = useMemo(
    () => ({
      user,
      status,
      login,
      logout,
      refetchMe,
    }),
    [user, status]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
