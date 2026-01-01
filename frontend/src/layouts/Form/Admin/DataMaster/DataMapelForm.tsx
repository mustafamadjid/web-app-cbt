import { useState } from "react";

import type {
  KelasOption,
  MataPelajaranFormValues,
} from "@/types/DataMaster/MataPelajaran";

const kelasOptions: KelasOption[] = [
  { id: "kelas-10", label: "Kelas 10" },
  { id: "kelas-11", label: "Kelas 11" },
  { id: "kelas-12", label: "Kelas 12" },
];

const initialValues: MataPelajaranFormValues = {
  kelasId: "",
  namaMapel: "",
  deskripsiMapel: "",
};

const inputBase =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500";
const labelBase = "text-xs font-medium text-slate-600";
const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const DataMapelForm = () => {
  const [values, setValues] = useState<MataPelajaranFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const setField = <K extends keyof MataPelajaranFormValues>(
    key: K,
    value: MataPelajaranFormValues[K]
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const onBlur = (name: keyof MataPelajaranFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = (v: MataPelajaranFormValues) => {
    const errors: Partial<Record<keyof MataPelajaranFormValues, string>> = {};

    if (!v.kelasId) errors.kelasId = "Kelas wajib dipilih.";
    if (!v.namaMapel.trim())
      errors.namaMapel = "Nama mata pelajaran wajib diisi.";
    if (!v.deskripsiMapel.trim())
      errors.deskripsiMapel = "Deskripsi mata pelajaran wajib diisi.";

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof MataPelajaranFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      kelasId: true,
      namaMapel: true,
      deskripsiMapel: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong.");
      return;
    }

    console.log("READY_TO_SUBMIT", values);
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Data Mata Pelajaran
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Lengkapi informasi mata pelajaran yang akan ditambahkan.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Mata Pelajaran</h2>
              <p className={helperText}>
                Pilih kelas dan isi detail mata pelajaran.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <label className={labelBase} htmlFor="kelasId">
                  Kelas
                </label>
                <select
                  id="kelasId"
                  className={`${inputBase} ${
                    hasError("kelasId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelasId}
                  onChange={(e) => setField("kelasId", e.target.value)}
                  onBlur={() => onBlur("kelasId")}
                >
                  <option value="">Pilih kelas</option>
                  {kelasOptions.map((option) => (
                    <option key={option.id} value={option.id}>
                      {option.label}
                    </option>
                  ))}
                </select>
                {hasError("kelasId") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kelasId}
                  </p>
                )}
              </div>

              <div>
                <label className={labelBase} htmlFor="namaMapel">
                  Nama Mata Pelajaran
                </label>
                <input
                  id="namaMapel"
                  className={`${inputBase} ${
                    hasError("namaMapel")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Contoh: Matematika"
                  value={values.namaMapel}
                  onChange={(e) => setField("namaMapel", e.target.value)}
                  onBlur={() => onBlur("namaMapel")}
                />
                {hasError("namaMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.namaMapel}
                  </p>
                )}
              </div>
            </div>

            <div className="mt-4">
              <label className={labelBase} htmlFor="deskripsiMapel">
                Deskripsi Mata Pelajaran
              </label>
              <textarea
                id="deskripsiMapel"
                rows={4}
                className={`${inputBase} ${
                  hasError("deskripsiMapel")
                    ? "border-rose-300 ring-rose-100"
                    : ""
                }`}
                placeholder="Tuliskan ringkasan materi atau fokus pelajaran."
                value={values.deskripsiMapel}
                onChange={(e) => setField("deskripsiMapel", e.target.value)}
                onBlur={() => onBlur("deskripsiMapel")}
              />
              {hasError("deskripsiMapel") && (
                <p className="mt-1 text-xs text-rose-500">
                  {errors.deskripsiMapel}
                </p>
              )}
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
              className="inline-flex items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({});
                setSubmitError(null);
              }}
            >
              Reset
            </button>
            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
            >
              Simpan Mata Pelajaran
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
