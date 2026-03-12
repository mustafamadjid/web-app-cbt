import React from "react";
import { useNavigate, useParams } from "react-router";
import { KeyRound } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { useAttemptUjian } from "@/services/Api/features-api/Ujian/attemptUjian.service";
import type { AttemptUjianRequest } from "@/types/Ujian/AttemptUjian";

const ATTEMPT_ERROR_MESSAGES: Record<string, string> = {
  INVALID_TOKEN_UJIAN: "Token ujian salah",
  SISWA_NOT_ALLOWED: "Anda tidak diizinkan mengikuti ujian ini",
  UJIAN_ATTEMPT_TIME_EXPIRED: "Ujian telah selesai",
};

const mapAttemptErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError) {
    if (error.code && ATTEMPT_ERROR_MESSAGES[error.code]) {
      return ATTEMPT_ERROR_MESSAGES[error.code];
    }
  }

  return "Gagal memulai ujian. Silakan coba lagi.";
};

const UjianTokenSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { idJadwalUjian, bankSoalId } = useParams();
  const { execute: executeAttempt, loading: submittingAttempt } = useAttemptUjian();
  const [token, setToken] = React.useState("");
  const [submitError, setSubmitError] = React.useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    setSubmitError(null);

    const trimmedToken = token.trim();
    if (!trimmedToken) {
      setSubmitError("Token ujian wajib diisi.");
      return;
    }

    const parsedIdJadwalUjian = Number(idJadwalUjian);
    if (!bankSoalId || !Number.isInteger(parsedIdJadwalUjian) || parsedIdJadwalUjian <= 0) {
      setSubmitError("Data ujian tidak valid.");
      return;
    }

    const idSiswa = user?.id_pengguna ?? 0;
    if (!Number.isInteger(idSiswa) || idSiswa <= 0) {
      setSubmitError("Akun siswa tidak valid. Silakan login ulang.");
      return;
    }

    const payload: AttemptUjianRequest = {
      id_siswa: idSiswa,
      id_jadwal_ujian: parsedIdJadwalUjian,
      token_ujian: trimmedToken,
      waktu_mulai: new Date().toISOString(),
    };

    try {
      await executeAttempt(payload);
      navigate(
        paths.dashboard.ujian_siswa_mulai
          .replace(":idJadwalUjian", String(parsedIdJadwalUjian))
          .replace(":bankSoalId", bankSoalId),
      );
    } catch (error) {
      setSubmitError(mapAttemptErrorMessage(error));
    }
  };

  return (
    <div className="flex min-h-[70vh] items-center justify-center px-4">
      <div className="w-full max-w-md rounded-2xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#397e50]/10 text-[#397e50]">
            <KeyRound className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-lg font-bold text-[#37513d]">
              Masukkan Token Ujian
            </h1>
            <p className="text-sm text-gray-500">
              Token dibagikan oleh pengawas sebelum ujian dimulai.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <label className="flex flex-col gap-2 text-xs font-semibold text-gray-500">
            Token Ujian
            <input
              value={token}
              onChange={(event) => {
                setToken(event.target.value);
                if (submitError) setSubmitError(null);
              }}
              placeholder="Contoh: MAT-UTS-2026"
              className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#397e50] focus:outline-none"
            />
          </label>

          {submitError && (
            <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-600">
              {submitError}
            </div>
          )}

          <button
            type="submit"
            disabled={submittingAttempt}
            className="w-full rounded-full bg-[#397e50] px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-70"
          >
            {submittingAttempt ? "Memproses..." : "Mulai Ujian"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default UjianTokenSiswa;
