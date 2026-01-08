import React, { useState } from "react";

import { InputField } from "@/components/common/Input/InputField";

import type { RuangUjianFormValues } from "@/types/DataMaster/RuangUjian";

import { createSetField } from "@/helper/setField/setField";

const initialValues: RuangUjianFormValues = {
  nama_ruangan_ujian: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

export const DataRuangForm = () => {
  const [values, setValues] = useState<RuangUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Record<keyof RuangUjianFormValues, boolean>
  >({
    nama_ruangan_ujian: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  const setField = createSetField(setValues);
  const onBlur = (name: keyof RuangUjianFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const normalizeNama = (s: string) => s.trim().replace(/\s+/g, " ");

  const validate = (v: RuangUjianFormValues) => {
    const errors: Partial<Record<keyof RuangUjianFormValues, string>> = {};

    if (!v.nama_ruangan_ujian.trim())
      errors.nama_ruangan_ujian = "Nama ruangan ujian wajib diisi.";
    else if (normalizeNama(v.nama_ruangan_ujian).length < 2)
      errors.nama_ruangan_ujian = "Nama ruangan ujian terlalu pendek.";

    return errors;
  };

  const errors = validate(values);
  const hasError = (name: keyof RuangUjianFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({ nama_ruangan_ujian: true });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError(
        "Periksa kembali input yang masih kosong atau tidak valid."
      );
      return;
    }

    const payload = {
      nama_ruangan_ujian: normalizeNama(values.nama_ruangan_ujian),
    };

    console.log("READY_TO_SUBMIT", payload);
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Data Ruangan Ujian
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Masukkan nama ruangan ujian yang akan ditambahkan.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Ruangan</h2>
              <p className={helperText}>Isi nama ruangan ujian.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="md:col-span-1">
                <InputField
                  id="nama_ruangan_ujian"
                  label="Nama Ruangan Ujian"
                  value={values.nama_ruangan_ujian}
                  onChange={(v) => setField("nama_ruangan_ujian", v)}
                  onBlur={() => onBlur("nama_ruangan_ujian")}
                  placeholder="Contoh: Ruang 01 / Lab Komputer / Aula"
                  inputClassName={
                    hasError("nama_ruangan_ujian")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_ruangan_ujian") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.nama_ruangan_ujian}
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
              className="inline-flex items-center justify-center cursor-pointer rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
              onClick={() => {
                setValues(initialValues);
                setTouched({ nama_ruangan_ujian: false });
                setSubmitError(null);
              }}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43]"
            >
              Simpan Ruangan Ujian
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
