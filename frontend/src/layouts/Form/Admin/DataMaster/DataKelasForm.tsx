import React, { useState } from "react";

import InputField from "@/components/common/Input/InputField";

import type { KelasFormValues } from "@/types/DataMaster/Kelas";
import { submitKelasResponse } from "@/services/Api/features-api/DataMaster/kelas.service";
import { ApiError } from "@/services/Api/api";

import { createSetField } from "@/helper/setField/setField";
import {
  createValidator,
  integerNumber,
  minNumber,
  requiredString,
  requiredValue,
} from "@/helper/validate/validateForm";

const initialValues: KelasFormValues = {
  tingkat_kelas: "",
  nama_kelas: "",
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const DataKelasForm = () => {
  const [values, setValues] = useState<KelasFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Record<keyof KelasFormValues, boolean>
  >({
    tingkat_kelas: false,
    nama_kelas: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  const setField = createSetField(setValues);
  const onBlur = (name: keyof KelasFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const validate = createValidator<KelasFormValues>({
    tingkat_kelas: [
      requiredValue("Tingkat kelas wajib diisi."),
      integerNumber("Tingkat kelas harus bilangan bulat."),
      minNumber(1, "Tingkat kelas tidak valid."),
    ],
    nama_kelas: [requiredString("Nama kelas wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof KelasFormValues) =>
    !!errors[name] && !!touched[name];

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    setTouched({
      tingkat_kelas: true,
      nama_kelas: true,
    });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError(
        "Periksa kembali input yang masih kosong atau tidak valid."
      );
      return;
    }

    try {
      await submitKelasResponse(values);
      alert("Kelas berhasil ditambahkan.");
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    }
  };

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Data Kelas
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Lengkapi informasi kelas yang akan ditambahkan.
            </p>
          </div>
        </div>

        <form onSubmit={onSubmit} className="space-y-6">
          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="mb-4">
              <h2 className={sectionTitle}>Informasi Kelas</h2>
              <p className={helperText}>Isi tingkat dan nama kelas.</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              {/* Tingkat Kelas (number) */}
              <div>
                <label
                  htmlFor="tingkat_kelas"
                  className="text-xs font-medium text-slate-600"
                >
                  Tingkat Kelas
                </label>

                <input
                  id="tingkat_kelas"
                  type="number"
                  inputMode="numeric"
                  className={`w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-[#397e50] focus:ring-2 focus:ring-[#397e50] disabled:bg-slate-50 disabled:text-slate-500 ${
                    hasError("tingkat_kelas")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }`}
                  placeholder="Contoh: 10"
                  value={values.tingkat_kelas}
                  onChange={(e) => {
                    const raw = e.target.value;
                    setField("tingkat_kelas", raw === "" ? "" : Number(raw));
                  }}
                  onBlur={() => onBlur("tingkat_kelas")}
                  min={1}
                  step={1}
                  required
                />

                {hasError("tingkat_kelas") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.tingkat_kelas}
                  </p>
                )}
              </div>

              {/* Nama Kelas */}
              <div>
                <InputField
                  id="nama_kelas"
                  label="Nama Kelas"
                  value={values.nama_kelas}
                  onChange={(v) => setField("nama_kelas", v)}
                  onBlur={() => onBlur("nama_kelas")}
                  placeholder="Contoh: X IPA 1"
                  inputClassName={
                    hasError("nama_kelas")
                      ? "border-rose-300 ring-rose-100"
                      : ""
                  }
                  required
                />
                {hasError("nama_kelas") && (
                  <p className="mt-1 text-xs text-rose-500">
                    {errors.nama_kelas}
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
                  tingkat_kelas: false,
                  nama_kelas: false,
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
              Simpan Kelas
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default DataKelasForm;
