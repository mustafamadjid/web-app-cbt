import { useEffect, useMemo, useState } from "react";

import { resolveDocumentUrl } from "@/helper/MediaUrl/resolveMediaUrl";
import { getUserFriendlyErrorMessage } from "@/services/Api/errorMessage";
import type { PengumumanFormValues } from "@/types/Widget/Pengumuman";

const MAX_DOCUMENT_SIZE = 10 * 1024 * 1024;

const isAllowedDocument = (file: File) => {
  const name = file.name.toLowerCase();
  return name.endsWith(".pdf") || name.endsWith(".docx");
};

type EditPengumumanFormProps = {
  initialValues: PengumumanFormValues;
  dokumenLama?: string;
  onSubmit: (values: PengumumanFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const EditPengumumanForm = ({
  initialValues,
  dokumenLama = "",
  onSubmit,
  loading = false,
  submitting = false,
}: EditPengumumanFormProps) => {
  const [values, setValues] = useState<PengumumanFormValues>(initialValues);
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    setValues(initialValues);
  }, [initialValues]);

  const dokumenLamaURL = useMemo(() => {
    if (!dokumenLama) return "";
    return resolveDocumentUrl(dokumenLama);
  }, [dokumenLama]);

  const isDisabled = loading || submitting;

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitError(null);

    if (
      !values.judul_pengumuman.trim() ||
      !values.isi_pengumuman.trim() ||
      !values.tanggal_rilis_pengumuman.trim() ||
      !values.tanggal_selesai_pengumuman.trim()
    ) {
      setSubmitError("Periksa kembali input yang masih kosong atau tidak valid.");
      return;
    }

    if (
      values.tanggal_selesai_pengumuman < values.tanggal_rilis_pengumuman
    ) {
      setSubmitError(
        "Tanggal selesai pengumuman tidak boleh lebih awal dari tanggal rilis.",
      );
      return;
    }

    if (
      values.dokumen_pengumuman &&
      values.dokumen_pengumuman.size > MAX_DOCUMENT_SIZE
    ) {
      setSubmitError("Ukuran dokumen pengumuman maksimal 10MB.");
      return;
    }

    if (values.dokumen_pengumuman && !isAllowedDocument(values.dokumen_pengumuman)) {
      setSubmitError("Dokumen pengumuman harus berupa PDF atau DOCX.");
      return;
    }

    try {
      await onSubmit({
        judul_pengumuman: values.judul_pengumuman.trim(),
        isi_pengumuman: values.isi_pengumuman.trim(),
        tanggal_rilis_pengumuman: values.tanggal_rilis_pengumuman.trim(),
        tanggal_selesai_pengumuman: values.tanggal_selesai_pengumuman.trim(),
        dokumen_pengumuman: values.dokumen_pengumuman,
      });
    } catch (e) {
      setSubmitError(
        getUserFriendlyErrorMessage(e, {
          action: "update",
          entity: "pengumuman",
        }),
      );
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">
            Edit Pengumuman
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Perbarui data pengumuman yang dipilih.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <label
                  htmlFor="judul_pengumuman"
                  className="text-xs font-medium text-slate-600"
                >
                  Judul Pengumuman
                </label>
                <input
                  id="judul_pengumuman"
                  value={values.judul_pengumuman}
                  onChange={(event) =>
                    setValues((prev) => ({
                      ...prev,
                      judul_pengumuman: event.target.value,
                    }))
                  }
                  placeholder="Masukkan judul pengumuman"
                  disabled={isDisabled}
                  className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                />
              </div>

              <div className="md:col-span-2">
                <label
                  htmlFor="isi_pengumuman"
                  className="text-xs font-medium text-slate-600"
                >
                  Isi Pengumuman
                </label>
                <textarea
                  id="isi_pengumuman"
                  rows={5}
                  value={values.isi_pengumuman}
                  onChange={(event) =>
                    setValues((prev) => ({
                      ...prev,
                      isi_pengumuman: event.target.value,
                    }))
                  }
                  placeholder="Masukkan isi pengumuman"
                  disabled={isDisabled}
                  className="mt-1 w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                />
              </div>

              <div>
                <label
                  htmlFor="tanggal_rilis_pengumuman"
                  className="text-xs font-medium text-slate-600"
                >
                  Tanggal Rilis
                </label>
                <input
                  id="tanggal_rilis_pengumuman"
                  type="date"
                  value={values.tanggal_rilis_pengumuman}
                  onChange={(event) =>
                    setValues((prev) => ({
                      ...prev,
                      tanggal_rilis_pengumuman: event.target.value,
                    }))
                  }
                  disabled={isDisabled}
                  className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                />
              </div>

              <div>
                <label
                  htmlFor="tanggal_selesai_pengumuman"
                  className="text-xs font-medium text-slate-600"
                >
                  Tanggal Selesai
                </label>
                <input
                  id="tanggal_selesai_pengumuman"
                  type="date"
                  value={values.tanggal_selesai_pengumuman}
                  onChange={(event) =>
                    setValues((prev) => ({
                      ...prev,
                      tanggal_selesai_pengumuman: event.target.value,
                    }))
                  }
                  disabled={isDisabled}
                  className="mt-1 w-full cursor-pointer rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500"
                />
              </div>

              <div className="md:col-span-2">
                <label
                  htmlFor="dokumen_pengumuman"
                  className="text-xs font-medium text-slate-600"
                >
                  Dokumen Pengumuman (Opsional)
                </label>
                <input
                  id="dokumen_pengumuman"
                  type="file"
                  accept=".pdf,.docx"
                  onChange={(event) => {
                    const file = event.target.files?.[0] ?? null;
                    setValues((prev) => ({ ...prev, dokumen_pengumuman: file }));
                  }}
                  disabled={isDisabled}
                  className="mt-1 block w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-slate-100 file:px-3 file:py-1.5 file:text-sm file:font-medium disabled:bg-slate-50 disabled:text-slate-500"
                />

                {values.dokumen_pengumuman && (
                  <div className="mt-2 flex items-center justify-between rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-600">
                    <span className="truncate">{values.dokumen_pengumuman.name}</span>
                    <button
                      type="button"
                      onClick={() =>
                        setValues((prev) => ({ ...prev, dokumen_pengumuman: null }))
                      }
                      className="cursor-pointer rounded-md border border-slate-200 bg-white px-2 py-1 text-xs text-slate-600 hover:bg-slate-100"
                    >
                      Hapus
                    </button>
                  </div>
                )}

                {!values.dokumen_pengumuman && dokumenLamaURL && (
                  <a
                    href={dokumenLamaURL}
                    target="_blank"
                    rel="noreferrer"
                    className="mt-2 inline-flex cursor-pointer items-center rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-700 hover:bg-slate-100"
                  >
                    Lihat dokumen saat ini
                  </a>
                )}
              </div>
            </div>
          </div>

          {submitError && (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600">
              {submitError}
            </div>
          )}

          <div className="flex flex-col gap-3 sm:flex-row sm:justify-end">
            <button
              type="button"
              onClick={() => {
                setValues(initialValues);
                setSubmitError(null);
              }}
              disabled={isDisabled}
              className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-70"
            >
              Reset
            </button>
            <button
              type="submit"
              disabled={isDisabled}
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
            >
              {submitting ? "Menyimpan..." : "Simpan Perubahan"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditPengumumanForm;
