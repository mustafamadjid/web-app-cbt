import { useMemo, useState } from "react";
import { CalendarRange, Edit3, Megaphone, Trash2 } from "lucide-react";
import { useNavigate } from "react-router";

import AddButton from "@/components/common/Button/AddButton";
import ConfirmAlert from "@/components/ui/ConfirmAlert/ConfirmAlert";
import { useAuth } from "@/contexts/AuthContext";
import { resolveDocumentUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import {
  deletePengumuman,
  useGetPengumumanActive,
  useGetPengumumanIncoming,
  useGetPengumumanNonActive,
} from "@/services/Api/features-api/pengumuman/pengumuman.service";
import type {
  PengumumanGetResponse,
  PengumumanStatusKey,
} from "@/types/Widget/Pengumuman";

const statusOptions: Array<{ key: PengumumanStatusKey; label: string }> = [
  { key: "incoming", label: "Akan Rilis" },
  { key: "active", label: "Sedang Rilis" },
  { key: "non-active", label: "Sudah Rilis" },
];

const toLocalDate = (value: string) => {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;
  return parsed.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

const getDocumentName = (path: string) => {
  const raw = path.split("/").filter(Boolean).pop() ?? "Dokumen";
  return decodeURIComponent(raw);
};

const PengumumanManagement = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [statusKey, setStatusKey] = useState<PengumumanStatusKey>("incoming");
  const [actionErrorMsg, setActionErrorMsg] = useState("");

  const [deleteTarget, setDeleteTarget] = useState<PengumumanGetResponse | null>(
    null,
  );
  const [deleting, setDeleting] = useState(false);

  const isGuru = user?.role === "GURU";
  const createPath = isGuru
    ? paths.dashboard.tambah_pengumuman_guru
    : paths.dashboard.tambah_pengumuman_admin;

  const getEditPath = (id: number) =>
    (isGuru ? paths.dashboard.edit_pengumuman_guru : paths.dashboard.edit_pengumuman_admin).replace(
      ":id",
      String(id),
    );

  const incomingState = useGetPengumumanIncoming();
  const activeState = useGetPengumumanActive();
  const nonActiveState = useGetPengumumanNonActive();

  const activeStatusState = useMemo(() => {
    switch (statusKey) {
      case "incoming":
        return incomingState;
      case "active":
        return activeState;
      case "non-active":
        return nonActiveState;
      default:
        return incomingState;
    }
  }, [activeState, incomingState, nonActiveState, statusKey]);

  const items = activeStatusState.data ?? [];
  const loading = activeStatusState.loading;
  const fetchErrorMsg = activeStatusState.error ?? "";
  const errorMsg = actionErrorMsg || fetchErrorMsg;

  const canEdit = statusKey !== "non-active";

  const onConfirmDelete = async () => {
    if (!deleteTarget) return;

    setDeleting(true);
    setActionErrorMsg("");
    try {
      await deletePengumuman(deleteTarget.id_pengumuman);
      setDeleteTarget(null);
      await activeStatusState.refetch();
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "data not found"
            ? "Data pengumuman tidak ditemukan."
            : e.message === "delete restricted : constraint violation"
              ? "Pengumuman tidak bisa dihapus karena masih dipakai."
              : "Pengumuman gagal dihapus."
          : "Pengumuman gagal dihapus.";
      setActionErrorMsg(message);
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="w-full space-y-6 px-8 py-13">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Manajemen Pengumuman
          </h2>
          <p className="mt-1 text-sm text-slate-500">
            Kelola pengumuman berdasarkan status rilis.
          </p>
        </div>
        <AddButton label="Tambah Pengumuman" onClick={() => navigate(createPath)} />
      </div>

      <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <div className="flex flex-wrap items-center gap-2">
          {statusOptions.map((status) => (
            <button
              key={status.key}
              type="button"
              onClick={() => {
                setStatusKey(status.key);
                setActionErrorMsg("");
              }}
              className={`cursor-pointer rounded-lg px-4 py-2 text-sm font-medium transition ${
                statusKey === status.key
                  ? "bg-[#397e50] text-white"
                  : "border border-slate-200 bg-white text-slate-700 hover:bg-slate-50"
              }`}
            >
              {status.label}
            </button>
          ))}
        </div>
      </div>

      {errorMsg && (
        <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
          {errorMsg}
        </div>
      )}

      <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        {loading ? (
          <div className="px-6 py-12 text-center text-sm text-slate-500">
            Memuat data...
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-slate-200 text-sm">
              <thead className="bg-slate-50">
                <tr>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Judul
                  </th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Periode
                  </th>
                  <th className="px-6 py-3 text-left font-semibold text-slate-600">
                    Dokumen
                  </th>
                  <th className="px-6 py-3 text-right font-semibold text-slate-600">
                    Aksi
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {items.length > 0 ? (
                  items.map((item) => {
                    const dokumenPath = item.dokumen_pengumuman ?? "";
                    const dokumenURL = resolveDocumentUrl(dokumenPath);

                    return (
                      <tr
                        key={item.id_pengumuman}
                        className="transition-colors hover:bg-slate-50"
                      >
                        <td className="px-6 py-4">
                          <div className="flex flex-col">
                            <span className="font-semibold text-slate-900">
                              {item.judul_pengumuman}
                            </span>
                            <span className="line-clamp-2 text-xs text-slate-500">
                              {item.isi_pengumuman}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-slate-700">
                          <div className="flex items-center gap-2">
                            <CalendarRange className="h-4 w-4 text-slate-400" />
                            <span>
                              {toLocalDate(item.tanggal_rilis_pengumuman)} -{" "}
                              {toLocalDate(item.tanggal_selesai_pengumuman)}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          {dokumenURL ? (
                            <a
                              href={dokumenURL}
                              target="_blank"
                              rel="noreferrer"
                              className="inline-flex cursor-pointer items-center rounded-lg border border-slate-200 bg-slate-50 px-3 py-1.5 text-xs text-slate-700 hover:bg-slate-100"
                            >
                              {getDocumentName(dokumenPath)}
                            </a>
                          ) : (
                            <span className="text-xs text-slate-400">Tidak ada</span>
                          )}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <div className="flex items-center justify-end gap-2">
                            <button
                              type="button"
                              onClick={() =>
                                navigate(getEditPath(item.id_pengumuman))
                              }
                              disabled={!canEdit}
                              className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-green-600 disabled:cursor-not-allowed disabled:opacity-40"
                              title={
                                canEdit
                                  ? "Edit"
                                  : "Pengumuman yang sudah rilis tidak bisa diedit"
                              }
                            >
                              <Edit3 className="h-4 w-4" />
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteTarget(item)}
                              className="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-100 hover:text-red-600"
                              title="Hapus"
                            >
                              <Trash2 className="h-4 w-4" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    );
                  })
                ) : (
                  <tr>
                    <td colSpan={4} className="px-6 py-12 text-center">
                      <div className="flex flex-col items-center gap-2">
                        <Megaphone className="h-10 w-10 text-slate-300" />
                        <p className="text-base font-medium text-slate-900">
                          Tidak ada pengumuman
                        </p>
                        <p className="text-sm text-slate-500">
                          Pilih kategori lain atau tambah pengumuman baru.
                        </p>
                      </div>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <ConfirmAlert
        isOpen={Boolean(deleteTarget)}
        title="Konfirmasi Hapus Pengumuman"
        message="Pengumuman yang dihapus tidak bisa dikembalikan. Lanjutkan?"
        onClose={() => {
          if (deleting) return;
          setDeleteTarget(null);
        }}
        onConfirm={onConfirmDelete}
        isLoading={deleting}
        confirmLabel="Ya, Hapus"
        loadingLabel="Menghapus..."
      />
    </div>
  );
};

export default PengumumanManagement;
