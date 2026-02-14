import React, { useEffect, useState } from "react";

import InputField from "@/components/common/Input/InputField";
import { createSetField } from "@/helper/setField/setField";
import { createValidator, requiredString } from "@/helper/validate/validateForm";
import { ApiError } from "@/services/Api/api";
import type { RuangUjianFormValues } from "@/types/DataMaster/RuangUjian";

type EditRuangFormProps = {
  initialValues: RuangUjianFormValues;
  onSubmit: (values: RuangUjianFormValues) => Promise<void>;
  loading?: boolean;
  submitting?: boolean;
};

const EditRuangForm = ({
  initialValues,
  onSubmit,
  loading = false,
  submitting = false,
}: EditRuangFormProps) => {
  const [values, setValues] = useState<RuangUjianFormValues>(initialValues);
  const [touched, setTouched] = useState<Record<keyof RuangUjianFormValues, boolean>>({
    kode_ruang: false,
    nama_ruangan: false,
  });
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    setValues(initialValues);
    setTouched({ kode_ruang: false, nama_ruangan: false });
  }, [initialValues]);

  const setField = createSetField(setValues);
  const normalize = (s: string) => s.trim().replace(/\s+/g, " ");

  const validate = createValidator<RuangUjianFormValues>({
    kode_ruang: [requiredString("Kode ruang wajib diisi.")],
    nama_ruangan: [requiredString("Nama ruangan wajib diisi.")],
  });

  const errors = validate(values);
  const hasError = (name: keyof RuangUjianFormValues) => !!errors[name] && touched[name];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);
    setTouched({ kode_ruang: true, nama_ruangan: true });

    const currentErrors = validate(values);
    if (Object.keys(currentErrors).length > 0) {
      setSubmitError("Periksa kembali input yang masih kosong atau tidak valid.");
      return;
    }

    const payload: RuangUjianFormValues = {
      kode_ruang: normalize(values.kode_ruang),
      nama_ruangan: normalize(values.nama_ruangan),
    };

    try {
      await onSubmit(payload);
    } catch (e) {
      const message =
        e instanceof ApiError
          ? e.message === "bad request: kode ruang ujian already exist"
            ? "Kode ruang ujian sudah ada."
            : "Ruang ujian gagal diperbarui."
          : "Ruang ujian gagal diperbarui.";
      setSubmitError(message);
    }
  };

  const isDisabled = loading || submitting;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <InputField
            id="kode_ruang"
            label="Kode Ruang"
            value={values.kode_ruang}
            onChange={(v) => setField("kode_ruang", v)}
            onBlur={() => setTouched((p) => ({ ...p, kode_ruang: true }))}
            placeholder="Contoh: RU-01"
            inputClassName={hasError("kode_ruang") ? "border-rose-300 ring-rose-100" : ""}
            disabled={isDisabled}
            required
          />
          <InputField
            id="nama_ruangan"
            label="Nama Ruangan"
            value={values.nama_ruangan}
            onChange={(v) => setField("nama_ruangan", v)}
            onBlur={() => setTouched((p) => ({ ...p, nama_ruangan: true }))}
            placeholder="Contoh: Lab Komputer"
            inputClassName={hasError("nama_ruangan") ? "border-rose-300 ring-rose-100" : ""}
            disabled={isDisabled}
            required
          />
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
          className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
          onClick={() => {
            setValues(initialValues);
            setTouched({ kode_ruang: false, nama_ruangan: false });
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
  );
};

export default EditRuangForm;
