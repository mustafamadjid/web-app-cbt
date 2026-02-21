import { useState } from "react";
import { useNavigate } from "react-router";

import { useAuth } from "@/contexts/AuthContext";
import { createSetField } from "@/helper/setField/setField";
import { createValidator, requiredString } from "@/helper/validate/validateForm";
import { paths } from "@/routes/paths";
import { ApiError } from "@/services/Api/api";
import { createPengumuman } from "@/services/Api/features-api/pengumuman/pengumuman.service";
import type { PengumumanFormValues } from "@/types/Widget/Pengumuman";

const MAX_DOCUMENT_SIZE = 10 * 1024 * 1024;

const initialValues: PengumumanFormValues = {
  judul_pengumuman: "",
  isi_pengumuman: "",
  tanggal_rilis_pengumuman: "",
  tanggal_selesai_pengumuman: "",
  dokumen_pengumuman: null,
};

const isAllowedDocument = (file: File) => {
  const name = file.name.toLowerCase();
  return name.endsWith(".pdf") || name.endsWith(".docx");
};

const CreatePengumumanForm = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [values, setValues] = useState<PengumumanFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const setField = createSetField(setValues);

  const validate = createValidator<PengumumanFormValues>({
    judul_pengumuman: [requiredString("Judul pengumuman wajib diisi.")],
    isi_pengumuman: [requiredString("Isi pengumuman wajib diisi.")],
    tanggal_rilis_pengumuman: [
      requiredString("Tanggal rilis pengumuman wajib diisi."),
    ],
    tanggal_selesai_pengumuman: [
      requiredString("Tanggal selesai pengumuman wajib diisi."),
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof PengumumanFormValues) =>
    !!errors[name] && !!touched[name];

  const onBlur = (name: keyof PengumumanFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const goBackPath =
    user?.role === "GURU"
      ? paths.dashboard.pengumuman_guru
      : paths.dashboard.pengumuman_admin;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouched({
      judul_pengumuman: true,
      isi_pengumuman: true,
      tanggal_rilis_pengumuman: true,
      tanggal_selesai_pengumuman: true,
      dokumen_pengumuman: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
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

    setSubmitting(true);
    try {
      await createPengumuman({
        judul_pengumuman: values.judul_pengumuman.trim(),
        isi_pengumuman: values.isi_pengumuman.trim(),
        tanggal_rilis_pengumuman: values.tanggal_rilis_pengumuman.trim(),
        tanggal_selesai_pengumuman: values.tanggal_selesai_pengumuman.trim(),
        dokumen_pengumuman: values.dokumen_pengumuman,
      });
      navigate(goBackPath);
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "bad request: invalid date format"
            ? "Format tanggal pengumuman tidak valid."
            : e.message === "bad request: missing fields"
              ? "Semua data pengumuman wajib diisi."
              : e.message === "bad request: invalid dokumen_pengumuman"
                ? "Dokumen pengumuman harus berupa PDF atau DOCX."
                : e.message === "file too large"
                  ? "Ukuran dokumen pengumuman maksimal 10MB."
                  : "Pengumuman gagal ditambahkan."
          : "Pengumuman gagal ditambahkan.";
      setSubmitError(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h1 className="text-base font-semibold text-slate-900">
            Tambah Pengumuman
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Isi data pengumuman dan unggah dokumen jika diperlukan.
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
                    setField("judul_pengumuman", event.target.value)
                  }
                  onBlur={() => onBlur("judul_pengumuman")}
                  placeholder="Masukkan judul pengumuman"
                  className={`mt-1 w-full cursor-pointer rounded-lg border px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] ${
                    hasError("judul_pengumuman")
                      ? "border-rose-300 ring-rose-100"
                      : "border-slate-200"
                  }`}
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
                    setField("isi_pengumuman", event.target.value)
                  }
                  onBlur={() => onBlur("isi_pengumuman")}
                  placeholder="Masukkan isi pengumuman"
                  className={`mt-1 w-full rounded-lg border px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] ${
                    hasError("isi_pengumuman")
                      ? "border-rose-300 ring-rose-100"
                      : "border-slate-200"
                  }`}
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
                    setField("tanggal_rilis_pengumuman", event.target.value)
                  }
                  onBlur={() => onBlur("tanggal_rilis_pengumuman")}
                  className={`mt-1 w-full cursor-pointer rounded-lg border px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] ${
                    hasError("tanggal_rilis_pengumuman")
                      ? "border-rose-300 ring-rose-100"
                      : "border-slate-200"
                  }`}
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
                    setField("tanggal_selesai_pengumuman", event.target.value)
                  }
                  onBlur={() => onBlur("tanggal_selesai_pengumuman")}
                  className={`mt-1 w-full cursor-pointer rounded-lg border px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-1 focus:ring-[#397e50] ${
                    hasError("tanggal_selesai_pengumuman")
                      ? "border-rose-300 ring-rose-100"
                      : "border-slate-200"
                  }`}
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
                    setField("dokumen_pengumuman", file);
                    onBlur("dokumen_pengumuman");
                  }}
                  className="mt-1 block w-full cursor-pointer rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-slate-100 file:px-3 file:py-1.5 file:text-sm file:font-medium"
                />
                {values.dokumen_pengumuman && (
                  <div className="mt-2 flex items-center justify-between rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-600">
                    <span className="truncate">{values.dokumen_pengumuman.name}</span>
                    <button
                      type="button"
                      onClick={() => setField("dokumen_pengumuman", null)}
                      className="cursor-pointer rounded-md border border-slate-200 bg-white px-2 py-1 text-xs text-slate-600 hover:bg-slate-100"
                    >
                      Hapus
                    </button>
                  </div>
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
                setTouched({});
                setSubmitError(null);
              }}
              disabled={submitting}
              className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-70"
            >
              Reset
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
            >
              {submitting ? "Menyimpan..." : "Simpan Pengumuman"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreatePengumumanForm;
