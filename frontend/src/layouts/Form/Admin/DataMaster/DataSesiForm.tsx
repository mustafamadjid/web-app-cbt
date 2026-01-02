import React, { useState } from "react";

import { InputField } from "@/components/common/Input/InputField";

import type { SesiFormValues } from "@/types/DataMaster/Sesi";

const initialValues: SesiFormValues = {
  kode_sesi: "",
  nama_sesi: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const DataSesiForm = () => {
  const [values, setValues] = useState<SesiFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Record<keyof SesiFormValues, boolean>
  >({
    kode_sesi: false,
    nama_sesi: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  const setField = <K extends keyof SesiFormValues>(
    key: K,
    value: SesiFormValues[K]
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
  };

  const onBlur = (name: keyof SesiFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const normalizeText = (text: string) => text.trim().replace(/\s+/g, " ");

  const validate = (v: SesiFormValues) => {
    const errors: Partial<Record<keyof SesiFormValues, string>> = {};

    if (!v.kode_sesi.trim()) errors.kode_sesi = "Kode sesi wajib diisi.";

    if (!v.nama_sesi.trim()) errors.nama_sesi = "Nama sesi wajib diisi.";

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof SesiFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      kode_sesi: true,
      nama_sesi: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError(
        "Periksa kembali input yang masih kosong atau tidak valid."
      );
      return;
    }

    const payload = {
      kode_sesi: normalizeText(values.kode_sesi),
      nama_sesi: normalizeText(values.nama_sesi),
    };

    console.log("READY_TO_SUBMIT", payload);
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Data Sesi Ujian
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Lengkapi informasi sesi ujian yang akan ditambahkan.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Sesi</h2>
              <p className={helperText}>Isi kode dan nama sesi ujian.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div>
                <InputField
                  id="kode_sesi"
                  label="Kode Sesi"
                  value={values.kode_sesi}
                  onChange={(v) => setField("kode_sesi", v)}
                  onBlur={() => onBlur("kode_sesi")}
                  placeholder="Contoh: SESI-01"
                  inputClassName={
                    hasError("kode_sesi")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("kode_sesi") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.kode_sesi}
                  </p>
                )}
              </div>

              <div>
                <InputField
                  id="nama_sesi"
                  label="Nama Sesi"
                  value={values.nama_sesi}
                  onChange={(v) => setField("nama_sesi", v)}
                  onBlur={() => onBlur("nama_sesi")}
                  placeholder="Contoh: Sesi Pagi"
                  inputClassName={
                    hasError("nama_sesi")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_sesi") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.nama_sesi}
                  </p>
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
              className="inline-flex items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({
                  kode_sesi: false,
                  nama_sesi: false,
                });
                setSubmitError(null);
              }}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
            >
              Simpan Sesi
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
