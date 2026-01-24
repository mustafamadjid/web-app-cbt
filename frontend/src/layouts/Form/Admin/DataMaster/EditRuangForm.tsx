import React, { useEffect, useState } from "react";

import InputField from "@/components/common/Input/InputField";

import type { RuangUjianFormValues } from "@/types/DataMaster/RuangUjian";
import { ApiError } from "@/services/Api/api";

import { createSetField } from "@/helper/setField/setField";
import { createValidator, requiredString } from "@/helper/validate/validateForm";

type EditRuangFormProps = {
  initialValues: RuangUjianFormValues;
  onSubmit: (values: RuangUjianFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const sectionTitle = "text-sm font-semibold text-slate-800";
const helperText = "text-xs text-slate-500";

const EditRuangForm = ({
  initialValues,
  onSubmit,
  loading = false,
  submitting = false,
}: EditRuangFormProps) => {
  const [values, setValues] = useState<RuangUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<
    Record<keyof RuangUjianFormValues, boolean>
  >({
    nama_ruangan_ujian: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    setValues(initialValues);
    setTouched({ nama_ruangan_ujian: false });
  }, [initialValues]);

  const setField = createSetField(setValues);
  const onBlur = (name: keyof RuangUjianFormValues) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  };

  const normalizeNama = (s: string) => s.trim().replace(/\s+/g, " ");

  const validate = createValidator<RuangUjianFormValues>({
    nama_ruangan_ujian: [
      requiredString("Nama ruangan ujian wajib diisi."),
      (value) => {
        if (!value.trim()) return null;
        return normalizeNama(value).length < 2
          ? "Nama ruangan ujian terlalu pendek."
          : null;
      },
    ],
  });

  const errors = validate(values);
  const hasError = (name: keyof RuangUjianFormValues) =>
    !!errors[name] && !!touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
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

    try {
      await onSubmit(payload);
    } catch (error) {
      if (error instanceof ApiError) {
        setSubmitError(error.message);
      }
    }
  };

  const isDisabled = loading || submitting;

  return (
    <div className="min-h-screen w-full py-8">
      <div className="mx-auto w-full max-w-5xl px-4">
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div>
            <h1 className="text-base font-semibold text-slate-900">
              Edit Data Ruangan Ujian
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Perbarui nama ruangan ujian.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
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
                  disabled={isDisabled}
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
              disabled={isDisabled}
            >
              Reset
            </button>

            <button
              type="submit"
              className="inline-flex cursor-pointer items-center justify-center rounded-lg bg-[#397e50] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-[#2f6a43] disabled:cursor-not-allowed disabled:opacity-70"
              disabled={isDisabled}
            >
              Simpan Perubahan
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditRuangForm;
