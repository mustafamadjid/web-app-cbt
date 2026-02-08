// src/contexts/AuthContext.tsx
import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  useCallback,
} from "react";
import { ApiError } from "@/services/Api/api";
import {
  getAuthMe,
  login as loginService,
  logout as logoutService,
} from "@/services/Api/auth-api/auth.service";

import type {
  User,
  AuthContextValue,
  AuthStatus,
  LoginPayload,
} from "@/types/Auth/Auth";
import type { Role } from "@/types/Sidebar/SidebarMenu";

const AuthContext = createContext<AuthContextValue | null>(null);

const validRoles = new Set<Role>(["ADMIN", "GURU", "SISWA"]);

const isValidRole = (role: unknown): role is Role =>
  typeof role === "string" && validRoles.has(role as Role);

// function getDebugAuthUser(): User | null {
//   if (!import.meta.env.DEV) return null;
//   const raw = localStorage.getItem("debug:auth");
//   if (!raw) return null;

//   try {
//     const parsed = JSON.parse(raw) as Partial<User>;
//     if (!parsed.username || !isValidRole(parsed.role)) return null;
//     return {
//       id_pengguna: typeof parsed.id_pengguna === "number" ? parsed.id_pengguna : 0,
//       username: parsed.username,
//       role: parsed.role,
//     };
//   } catch {
//     return null;
//   }
// }

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [status, setStatus] = useState<AuthStatus>("loading");

  const forceGuest = useCallback(() => {
    setUser(null);
    setStatus("guest");
  }, []);

  const handleAuthError = useCallback(
    (e: unknown) => {
      const err = e as ApiError;


      if (err.code === "SESSION_EXPIRED") {
        forceGuest();
        return true;
      }

      if (err.status === 401) {
        forceGuest();
        return true;
      }

      return false;
    },
    [forceGuest],
  );

  const refetchMe = useCallback(async () => {
    try {
      const me = await getAuthMe();
      if (!isValidRole(me.role)) {
        forceGuest();
        return null;
      }
      setUser({ id_pengguna: me.id_pengguna, username: me.username, role: me.role });
      setStatus("authenticated");
      return me;
    } catch (e) {
      if (!handleAuthError(e)) throw e;
      return null;
    }
  }, [forceGuest, handleAuthError]);

  const boot = useCallback(async () => {
    setStatus("loading");

    // const debugUser = getDebugAuthUser();
    // if (debugUser) {
    //   setUser(debugUser);
    //   setStatus("authenticated");
    //   return;
    // }

    try {
      await refetchMe();
    } catch {
      // kalau error non-auth, kamu bisa log kalau mau
      forceGuest();
    }
  }, [refetchMe, forceGuest]);

  useEffect(() => {
    void boot();
  }, [boot]);

  const login = useCallback(
    async (payload: LoginPayload) => {
      await loginService(payload);
      await refetchMe();
    },
    [refetchMe],
  );

  const logout = useCallback(async () => {
    try {
      await logoutService();
    } catch (e) {
      const err = e as ApiError;
      console.warn("Logout error:", err.status, err.message);
    } finally {
      forceGuest();
    }
  }, [forceGuest]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      status,
      login,
      logout,
      refetchMe,
    }),
    [user, status, login, logout, refetchMe],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
