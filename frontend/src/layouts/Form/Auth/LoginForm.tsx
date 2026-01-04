import { useState } from "react";
import { useLogin } from "../../../hooks/Auth/useLogin";

// CSS
import "../../../index.css";

// UI
import { Spinner } from "@/components/ui/spinner";

// Components
import { LoginInputField } from "../../../components/common/Input/Auth/LoginInputField";

// Props Types
type LoginFormProps = {
  onSuccess: () => void;
};

export const LoginForm = ({ onSuccess }: LoginFormProps) => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const { login, loading, error } = useLogin();

  // Submit handling
  const submit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    await login({ username, password });
    onSuccess();
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
            {/* <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-gray-600 hover:text-black cursor-pointer"
              aria-label="Tampilkan password"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="currentColor"
                viewBox="0 0 16 16"
              >
                <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z" />
                <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0" />
              </svg>
            </button> */}
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
                "font-bold tracking-[0.05em]",
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
