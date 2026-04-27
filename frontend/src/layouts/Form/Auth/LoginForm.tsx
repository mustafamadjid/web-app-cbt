import { useState } from "react";
import { ApiError } from "@/services/Api/api";
import { useAuth } from "@/contexts/AuthContext";

// CSS
import "../../../index.css";

// UI
import Spinner from "@/components/ui/spinner";
import {
  USERNAME_LENGTH_INVALID_MESSAGE,
  USERNAME_MAX_LENGTH,
} from "@/constants/username";

// Components
import LoginInputField from "../../../components/common/Input/Auth/LoginInputField";

// Props Types
type LoginFormProps = {
  onSuccess?: () => void;
};

const LoginForm = ({ onSuccess }: LoginFormProps) => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Submit handling
  const submit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await login({ username, password });
      onSuccess?.();
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.code === "USERNAME_LENGTH_INVALID"
            ? USERNAME_LENGTH_INVALID_MESSAGE
            : e.code === "HAS_SESSION"
              ? "Login gagal. Silakan logout terlebih dahulu pada device sebelumnya"
              : "Login Gagal. Silakan Coba lagi"
          : "Login gagal.";
      setError(message);
    } finally {
      setLoading(false);
    }
  };
  return (
    <>
      {error && <p className="text-red-500">{error}</p>}
      {/* Form */}
      <form onSubmit={submit} className="flex flex-col gap-10 ">
        <div className="flex flex-col gap-6">
          <LoginInputField
            id="username"
            label="Username"
            type="text"
            value={username}
            onChange={setUsername}
            autoComplete="username"
            maxLength={USERNAME_MAX_LENGTH}
          />

          <div className="relative">
            <LoginInputField
              id="password"
              label="Password"
              type="password"
              value={password}
              onChange={setPassword}
              autoComplete="current-password"
            />
            
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className={[
            "group relative inline-flex items-center justify-center cursor-pointer",
            "rounded-xl border border-gray-200 bg-white px-4 py-2",
            "focus:outline-none focus:ring-4 focus:ring-black/10 focus:border-black",
            "transition active:scale-95 disabled:cursor-not-allowed disabled:opacity-60",
            "button-box-shadow",
          ].join(" ")}
        >
          {/* Background pill yang melebar saat hover (replacement untuk :before) */}
          <span
            className={[
              "absolute left-0 top-1/2 -translate-y-1/2",
              "h-11 w-8 rounded-xl bg-[#397e50]",
              "transition-all duration-300 ease-out",
              "group-hover:w-[calc(101%-0px)]",
              "z-0",
            ].join(" ")}
          />

          {/* Content */}
          <span className="relative z-10 inline-flex items-center gap-2">
            {loading && <Spinner />}

            <span
              className={[
                "font-bold tracking-wider",
                "text-[#397e50] transition-colors duration-300",
                "group-hover:text-white",
              ].join(" ")}
            >
              Masuk
            </span>

            {/* Icon panah (replacement svg animasi) */}
            <svg
              viewBox="0 0 24 24"
              className={[
                "h-5 w-5",
                "translate-x-[-5px] transition-all duration-300 ease-out",
                "group-hover:translate-x-0",
                "stroke-[#397e50] group-hover:stroke-white",
              ].join(" ")}
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              aria-hidden="true"
            >
              <path d="M5 12h14" />
              <path d="M13 5l7 7-7 7" />
            </svg>
          </span>
        </button>
      </form>
    </>
  );
};

export default LoginForm;
