import { useState } from "react";

// TODO : Tambahkan ID Device juga saat login

type LoginPayload = {
  username: string;
  password: string;
};

type LoginResult = {
  token: string;
};

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (payload: LoginPayload): Promise<LoginResult> => {
    setLoading(true);
    setError(null);

    try {
      // API Call
      const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        if (res.status === 401) throw new Error("Kombinasi Username dan password salah.");
        throw new Error("Terjadi kesalahan server.");
      }

      const data = (await res.json()) as LoginResult;

    //   TODO : Buat mekanisme simpan token

      return data;
    } catch (e) {
      const message = e instanceof Error ? e.message : "Login gagal.";
      setError(message);
      throw e; 
    } finally {
      setLoading(false);
    }
  };

  return { login, loading, error, setError };
}
