import React from "react";
import { useNavigate, useParams } from "react-router";
import { KeyRound } from "lucide-react";
import { paths } from "@/routes/paths";

const UjianTokenSiswa: React.FC = () => {
  const navigate = useNavigate();
  const { id, bankSoalId } = useParams();
  const [token, setToken] = React.useState("");

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!id || !bankSoalId) return;
    if (!token.trim()) return;
    navigate(
      paths.dashboard.ujian_siswa_mulai
        .replace(":id", id)
        .replace(":bankSoalId", bankSoalId)
    );
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
              onChange={(event) => setToken(event.target.value)}
              placeholder="Contoh: MAT-UTS-2026"
              className="rounded-lg border border-gray-200 px-3 py-2 text-sm text-gray-600 focus:border-[#397e50] focus:outline-none"
            />
          </label>

          <button
            type="submit"
            className="w-full rounded-full bg-[#397e50] px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
          >
            Mulai Ujian
          </button>
        </form>
      </div>
    </div>
  );
};

export default UjianTokenSiswa;
