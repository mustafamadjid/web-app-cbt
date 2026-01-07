import { useMemo, useState } from "react";

import { InputField } from "@/components/common/Input/InputField";

import type {
  KelasOption,
  MataPelajaranFormValues,
} from "@/types/DataMaster/MataPelajaran";

const kelasOptions: KelasOption[] = [
  { id: "kelas-10", tingkat_kelas: 10, label: "Kelas 10" },
  { id: "kelas-11", tingkat_kelas: 11, label: "Kelas 11" },
  { id: "kelas-12", tingkat_kelas: 12, label: "Kelas 12" },
];

const initialValues: MataPelajaranFormValues = {
  kelasId: "",
  kodeMapel: "",
  namaMapel: "",
  deskripsiMapel: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const DataMapelForm = () => {
  const [values, setValues] = useState<MataPelajaranFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);

  const kelasTingkatOptions = useMemo(
    () => [...kelasOptions].sort((a, b) => a.tingkat_kelas - b.tingkat_kelas),
    []
  );

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

    if (!v.kelasId) errors.kelasId = "Tingkat kelas wajib dipilih.";
    if (!v.kodeMapel.trim()) errors.kodeMapel = "Kode mapel wajib diisi.";
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
      kodeMapel: true,
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
                Pilih tingkat kelas dan isi detail mata pelajaran.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              {/* Tingkat Kelas (select native) */}
              <div>
                <label
                  htmlFor="kelasId"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>
                <select
                  id="kelasId"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("kelasId") ? "border-rose-300 ring-rose-100" : ""
                  }`}
                  value={values.kelasId}
                  onChange={(e) => setField("kelasId", e.target.value)}
                  onBlur={() => onBlur("kelasId")}
                  required
                >
                  <option value="">Pilih tingkat kelas</option>

                  {/* 2) Render hanya tingkat kelas unik */}
                  {kelasTingkatOptions.map((option) => (
                    <option key={option.id} value={option.id}>
                      {option.label}
                    </option>
                  ))}
                </select>

                {hasError("kelasId") && (
                  <p className="mt-1 text-xs text-rose-500">{errors.kelasId}</p>
                )}
              </div>

              <div>
                <InputField
                  id="kodeMapel"
                  label="Kode Mapel"
                  value={values.kodeMapel}
                  onChange={(v) => setField("kodeMapel", v)}
                  onBlur={() => onBlur("kodeMapel")}
                  placeholder="Contoh: MAT-10-01"
                  inputClassName={
                    hasError("kodeMapel") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("kodeMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kodeMapel}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="namaMapel"
                  label="Nama Mata Pelajaran"
                  value={values.namaMapel}
                  onChange={(v) => setField("namaMapel", v)}
                  onBlur={() => onBlur("namaMapel")}
                  placeholder="Contoh: Matematika"
                  inputClassName={
                    hasError("namaMapel") ? "border-rose-300 ring-rose-100" : ""
                  }
                  required
                />
                {hasError("namaMapel") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.namaMapel}
                  </p>
                )}
              </div>
            </div>

            <div className="mt-4">
              <label
                htmlFor="deskripsiMapel"
                className="text-xs font-medium text-slate-600"
              >
                Deskripsi Mata Pelajaran
              </label>
              <textarea
                id="deskripsiMapel"
                rows={4}
                className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                  hasError("deskripsiMapel")
                    ? "border-rose-300 ring-rose-100"
                    : ""
                }`}
                placeholder="Tuliskan ringkasan materi atau fokus pelajaran."
                value={values.deskripsiMapel}
                onChange={(e) => setField("deskripsiMapel", e.target.value)}
                onBlur={() => onBlur("deskripsiMapel")}
                required
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
