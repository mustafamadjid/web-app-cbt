import { useMemo, useState } from "react";
import { Clock3, LogOut, ShieldCheck, UserRound } from "lucide-react";
import toast from "react-hot-toast";

import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { useAuth } from "@/contexts/AuthContext";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import {
  useAdminRevokeSession,
  useGetActiveLoginSessions,
} from "@/services/Api/features-api/Session/session.service";
import type { ActiveSessionRow } from "@/types/Session/Session";

const roleLabels: Record<string, string> = {
  ADMIN: "Administrator",
  GURU: "Guru",
  SISWA: "Siswa",
};

const statusLabels: Record<string, string> = {
  AKTIF: "Aktif",
  NONAKTIF: "Nonaktif",
};

const getInitials = (name?: string) => {
  if (!name) return "U";

  const parts = name.trim().split(/\s+/);
  const first = parts[0]?.[0] ?? "";
  const last = parts.length > 1 ? parts[parts.length - 1]?.[0] ?? "" : "";

  return `${first}${last}`.toUpperCase();
};

const formatOptionalValue = (value?: string | null) => {
  const trimmed = value?.trim();
  return trimmed ? trimmed : "-";
};

const formatDateTime = (value: string) => {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;

  return parsed.toLocaleString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const KelolaSesi = () => {
  const { user } = useAuth();
  const {
    data: sessionRows,
    loading,
    error,
    refetch,
  } = useGetActiveLoginSessions();
  const {
    execute: executeRevokeSession,
    loading: revokeLoading,
  } = useAdminRevokeSession();

  const [selectedSession, setSelectedSession] = useState<ActiveSessionRow | null>(null);

  const rows = useMemo(() => sessionRows ?? [], [sessionRows]);

  const handleOpenConfirm = (item: ActiveSessionRow) => {
    setSelectedSession(item);
  };

  const handleCloseConfirm = () => {
    if (revokeLoading) return;
    setSelectedSession(null);
  };

  const handleConfirmRevoke = async () => {
    if (!selectedSession) return;

    try {
      await executeRevokeSession({ session_id: selectedSession.session.session_id });
      toast.success(
        `Sesi ${selectedSession.pengguna.nama_lengkap} berhasil di-logout.`,
      );
      setSelectedSession(null);
      await refetch();
    } catch (e) {
      toast.error(
        getUserFriendlyErrorMessage(e, {
          action: "delete",
          entity: "sesi login",
          fallbackMessage: "Gagal melakukan logout user.",
        }),
      );
    }
  };

  return (
    <div className="w-full space-y-6 px-8 py-13">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Kelola Sesi
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Pantau sesi login aktif dan logout user lain saat diperlukan.
          </p>
        </div>

        <div className="inline-flex items-center gap-2 rounded-full border border-[#397e50]/20 bg-[#397e50]/8 px-4 py-2 text-sm font-medium text-[#397e50]">
          <ShieldCheck className="h-4 w-4" />
          Hanya admin yang dapat merevoke sesi
        </div>
      </div>

      {error && (
        <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
          {error}
        </div>
      )}

      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        {loading ? (
          <div className="px-6 py-12 text-center text-sm text-slate-500">
            Memuat data sesi aktif...
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-slate-200 text-sm">
              <thead className="bg-slate-50">
                <tr>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Pengguna
                  </th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Kontak
                  </th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Role / Status
                  </th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Informasi Sesi
                  </th>
                  <th className="px-6 py-3 text-right font-semibold text-slate-600">
                    Aksi
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {rows.length > 0 ? (
                  rows.map((item) => {
                    const isSelfSession =
                      item.session.id_pengguna === user?.id_pengguna;
                    const avatarUrl = item.pengguna.foto_profil
                      ? resolveImageUrl(item.pengguna.foto_profil)
                      : "";

                    return (
                      <tr
                        key={item.session.session_id}
                        className="transition-colors hover:bg-slate-50"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="relative h-11 w-11 overflow-hidden rounded-full border border-slate-200 bg-slate-100">
                              {avatarUrl ? (
                                <img
                                  src={avatarUrl}
                                  alt={item.pengguna.nama_lengkap}
                                  className="h-full w-full object-cover"
                                />
                              ) : (
                                <div className="flex h-full w-full items-center justify-center bg-[#397e50]/10 text-sm font-semibold text-[#397e50]">
                                  {getInitials(item.pengguna.nama_lengkap)}
                                </div>
                              )}
                            </div>
                            <div className="flex flex-col">
                              <span className="font-semibold text-slate-900">
                                {item.pengguna.nama_lengkap}
                              </span>
                              <span className="text-xs text-slate-500">
                                @{item.pengguna.username}
                              </span>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-slate-700">
                          <div className="flex flex-col gap-1">
                            <span>{formatOptionalValue(item.pengguna.email)}</span>
                            <span className="text-xs text-slate-500">
                              {formatOptionalValue(item.pengguna.no_hp)}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex flex-col gap-2">
                            <span className="inline-flex w-fit rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-700">
                              {roleLabels[item.pengguna.role] ?? item.pengguna.role}
                            </span>
                            <span
                              className={`inline-flex w-fit rounded-full px-2.5 py-1 text-xs font-medium ${
                                item.pengguna.status_akun === "AKTIF"
                                  ? "bg-emerald-50 text-emerald-700"
                                  : "bg-slate-100 text-slate-600"
                              }`}
                            >
                              {statusLabels[item.pengguna.status_akun] ??
                                item.pengguna.status_akun}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-slate-700">
                          <div className="flex flex-col gap-1">
                            <span className="text-xs font-medium text-slate-900">
                              {isSelfSession ? "Sesi Anda" : "Sesi aktif"}
                            </span>
                            <span className="text-xs text-slate-500">
                              Berakhir: {formatDateTime(item.session.expires_at)}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-right">
                          {isSelfSession ? (
                            <button
                              type="button"
                              disabled
                              className="cursor-not-allowed rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium text-slate-400 opacity-80"
                            >
                              Sesi Anda
                            </button>
                          ) : (
                            <button
                              type="button"
                              onClick={() => handleOpenConfirm(item)}
                              className="inline-flex cursor-pointer items-center gap-2 rounded-lg bg-rose-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-rose-700"
                            >
                              <LogOut className="h-4 w-4" />
                              Logout
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })
                ) : (
                  <tr>
                    <td colSpan={5} className="px-6 py-12 text-center">
                      <div className="flex flex-col items-center gap-3">
                        <div className="rounded-full bg-slate-100 p-4 text-slate-400">
                          <UserRound className="h-8 w-8" />
                        </div>
                        <div>
                          <p className="text-base font-medium text-slate-900">
                            Tidak ada sesi aktif
                          </p>
                          <p className="mt-1 text-sm text-slate-500">
                            Semua user saat ini sedang tidak memiliki sesi login aktif.
                          </p>
                        </div>
                      </div>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}

        {!loading && rows.length > 0 && (
          <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
            <p className="text-sm text-slate-700">
              Total sesi aktif: <span className="font-medium">{rows.length}</span>
            </p>
            <div className="flex items-center gap-2 text-xs text-slate-500">
              <Clock3 className="h-4 w-4" />
              Data diperbarui setiap kali halaman dimuat ulang atau setelah revoke.
            </div>
          </div>
        )}
      </div>

      <ConfirmAlert
        isOpen={Boolean(selectedSession)}
        title="Konfirmasi Logout User"
        message={
          selectedSession
            ? `Anda akan melogout ${selectedSession.pengguna.nama_lengkap} (@${selectedSession.pengguna.username}) dari sesi aktifnya. Lanjutkan?`
            : ""
        }
        onClose={handleCloseConfirm}
        onConfirm={() => void handleConfirmRevoke()}
        isLoading={revokeLoading}
        confirmLabel="Ya, Logout"
        loadingLabel="Memproses..."
        confirmClassName="bg-rose-600 hover:bg-rose-700"
      />
    </div>
  );
};

export default KelolaSesi;
